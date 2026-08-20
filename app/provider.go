package app

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"cosmossdk.io/log/v2"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/privval"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/goatnetwork/goat/cmd/goatd/cmd/goatflags"
	"github.com/goatnetwork/goat/pkg/ethrpc"
	bitcintypes "github.com/goatnetwork/goat/x/bitcoin/types"
	"github.com/spf13/cast"
)

// TxSigningKeyFile is the optional key file, kept beside the CometBFT privval
// key, that a validator signs MsgNewEthBlock with.
//
// The two keys are the same one for a validator that has never rotated, which
// is why the privval file is still the fallback. They stop being the same the
// moment a consensus key is rotated: the account was created from the original
// secp256k1 key (x/locking/keeper/msg_create.go), rotation deliberately leaves
// that account alone, and the ante handler keeps verifying MsgNewEthBlock
// against it. Rotating without splitting the files first replaces the account
// key as a side effect, which either fails the type check below or trips the
// pubkey comparison in createEthBlockProposal.
//
// The format is the privval one so that splitting is a copy:
//
//	cp config/priv_validator_key.json config/priv_tx_key.json
const TxSigningKeyFile = "priv_tx_key.json"

func ProvideEngineClient(logger log.Logger, appOpts servertypes.AppOptions) (*ethrpc.Client, *params.ChainConfig) {
	endpoint := cast.ToString(appOpts.Get(goatflags.GoatGeth))
	client, conf, err := ConnectEngineClient(context.Background(), logger, endpoint)
	if err != nil {
		panic(err)
	}
	return client, conf
}

// ProvideValidatorPrvKey loads the key this node signs MsgNewEthBlock with. It
// is the account key, not the consensus key, and it has to stay secp256k1 for
// as long as the account it belongs to does.
func ProvideValidatorPrvKey(logger log.Logger, appOpts servertypes.AppOptions) cryptotypes.PrivKey {
	consensusKeyFile := cast.ToString(appOpts.Get("priv_validator_key_file"))
	if !filepath.IsAbs(consensusKeyFile) {
		consensusKeyFile = filepath.Join(cast.ToString(appOpts.Get("home")), consensusKeyFile)
	}
	txKeyFile := filepath.Join(filepath.Dir(consensusKeyFile), TxSigningKeyFile)

	keyFile, split := txKeyFile, true
	if _, err := os.Stat(txKeyFile); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			panic(err)
		}
		keyFile, split = consensusKeyFile, false
	}

	pvKey, err := loadFilePVKey(keyFile)
	if err != nil {
		panic(err)
	}

	if pvKey.PrivKey.Type() != bitcintypes.Secp256K1Name {
		if split {
			panic(keyFile + " is not an secp256k1 key")
		}
		// The likely story: the operator rotated the consensus key and copied
		// the new one over the privval file without splitting first. Say what
		// to do rather than only what is wrong.
		panic(consensusKeyFile + " is not an secp256k1 key. A validator that has rotated its" +
			" consensus key has to keep the key its account was created with in " + txKeyFile)
	}

	logger.Info("Loaded the transaction signing key", "file", keyFile, "split_from_privval", split)
	return &secp256k1.PrivKey{Key: pvKey.PrivKey.Bytes()}
}

func loadFilePVKey(path string) (*privval.FilePVKey, error) {
	keyJSONBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pvKey := new(privval.FilePVKey)
	if err := cmtjson.Unmarshal(keyJSONBytes, pvKey); err != nil {
		return nil, err
	}
	if pvKey.PrivKey == nil {
		return nil, errors.New(path + " carries no private key")
	}
	return pvKey, nil
}
