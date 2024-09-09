package modgen

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/ethereum/go-ethereum/triedb/hashdb"
)

var BitcoinNetworks = map[string]*chaincfg.Params{
	chaincfg.MainNetParams.Name:       &chaincfg.MainNetParams,
	chaincfg.TestNet3Params.Name:      &chaincfg.TestNet3Params,
	chaincfg.SigNetParams.Name:        &chaincfg.SigNetParams,
	chaincfg.RegressionNetParams.Name: &chaincfg.RegressionNetParams,
}

var DepositMagicPreifxs = map[string]string{
	chaincfg.MainNetParams.Name:       "GTV2",
	chaincfg.TestNet3Params.Name:      "GTV1",
	chaincfg.SigNetParams.Name:        "GTV1",
	chaincfg.RegressionNetParams.Name: "GTT0",
}

func DecodeHexOrBase64String(str string) ([]byte, error) {
	pubkeyRaw, err := hex.DecodeString(str)
	if err != nil {
		pubkeyRaw, err = base64.StdEncoding.DecodeString(str)
		if err != nil {
			pubkeyRaw, err = hex.DecodeString(strings.TrimPrefix(str, "0x"))
			if err != nil {
				return nil, fmt.Errorf("pubkey %s doesn't use base64 or hex encoding", str)
			}
		}
	}
	return pubkeyRaw, nil
}

func IsValidSecp256Pubkey(key []byte) error {
	if len(key) != secp256k1.PubKeySize {
		return errors.New("invalid secp256k1 pubkey length")
	}
	if key[0] != 2 && key[0] != 3 {
		return errors.New("invalid secp256k1 pubkey prefix")
	}
	return nil
}

func GetEthGenesisHeaderByFile(genesisPath string) (*ethtypes.Header, error) {
	file, err := os.Open(genesisPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	genesis := new(core.Genesis)
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		return nil, err
	}

	memdb := rawdb.NewMemoryDatabase()
	triedb := triedb.NewDatabase(memdb, &triedb.Config{Preimages: true, HashDB: hashdb.Defaults})
	defer triedb.Close()

	block, err := genesis.Commit(memdb, triedb)
	if err != nil {
		return nil, err
	}
	header := block.Header()

	if header.BaseFee == nil || header.WithdrawalsHash == nil {
		return nil, errors.New("shanghai upgrade should be activated")
	}

	if *header.WithdrawalsHash != ethtypes.EmptyWithdrawalsHash {
		return nil, errors.New("No withdrawals required")
	}

	if header.GasUsed != 0 || header.TxHash != ethtypes.EmptyTxsHash {
		return nil, errors.New("No txs required")
	}

	if header.BlobGasUsed == nil || header.ExcessBlobGas == nil || header.ParentBeaconRoot == nil {
		return nil, errors.New("cancun upgrade should be activated")
	}

	if *header.BlobGasUsed != 0 {
		return nil, errors.New("No blob txes required")
	}

	return header, nil
}
