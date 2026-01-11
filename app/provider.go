package app

import (
	"context"
	"os"
	"path/filepath"

	"cosmossdk.io/log"
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

func ProvideEngineClient(logger log.Logger, appOpts servertypes.AppOptions) (*ethrpc.Client, *params.ChainConfig) {
	endpoint := cast.ToString(appOpts.Get(goatflags.GoatGeth))
	client, conf, err := ConnectEngineClient(context.Background(), logger, endpoint)
	if err != nil {
		panic(err)
	}
	return client, conf
}

func ProvideValidatorPrvKey(appOpts servertypes.AppOptions) cryptotypes.PrivKey {
	prvkey := cast.ToString(appOpts.Get("priv_validator_key_file"))
	if !filepath.IsAbs(prvkey) {
		prvkey = filepath.Join(cast.ToString(appOpts.Get("home")), prvkey)
	}

	keyJSONBytes, err := os.ReadFile(prvkey)
	if err != nil {
		panic(err)
	}

	var pvKey privval.FilePVKey
	err = cmtjson.Unmarshal(keyJSONBytes, &pvKey)
	if err != nil {
		panic(err)
	}

	if pvKey.PrivKey.Type() != bitcintypes.Secp256K1Name {
		panic(prvkey + " is not an secp256k1 key")
	}
	return &secp256k1.PrivKey{Key: pvKey.PrivKey.Bytes()}
}
