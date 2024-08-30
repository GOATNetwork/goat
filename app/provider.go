package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/privval"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/goatnetwork/goat/pkg/ethrpc"
)

func ProvideBitcoinNetworkConfig(appOpts servertypes.AppOptions) *chaincfg.Params {
	name, ok := appOpts.Get("goat.btc-network").(string)
	if !ok {
		panic("No bitcoin network config")
	}

	switch name {
	case chaincfg.MainNetParams.Name:
		return &chaincfg.MainNetParams
	case chaincfg.TestNet3Params.Name:
		return &chaincfg.TestNet3Params
	case chaincfg.SigNetParams.Name:
		return &chaincfg.SigNetParams
	case chaincfg.RegressionNetParams.Name:
		return &chaincfg.RegressionNetParams
	default:
		panic("Undefined bitcoin network name: " + name)
	}
}

func ProvideEngineClient(appOpts servertypes.AppOptions) *ethrpc.Client {
	endpoint, ok := appOpts.Get("goat.geth").(string)
	if !ok || endpoint == "" {
		panic("goat execution node endpoint not found")
	}

	jwtSecret := func() []byte {
		jwtPath, ok := appOpts.Get("goat.jwt-path").(string)
		if ok && jwtPath != "" {
			if data, err := os.ReadFile(jwtPath); err == nil {
				jwtSecret := common.FromHex(strings.TrimSpace(string(data)))
				if len(jwtSecret) == 32 {
					return jwtSecret
				}
			}
		}
		return nil
	}()

	var ethclient *ethrpc.Client
	for i := 0; i < 10; i++ {
		var err error
		ethclient, err = ethrpc.DialContext(context.Background(), endpoint, jwtSecret)
		if err == nil {
			break
		}
		<-time.After(time.Second / 2)
	}

	return ethclient
}

func ProvideValidatorPrvKey(appOpts servertypes.AppOptions) cryptotypes.PrivKey {
	prvkey := appOpts.Get("priv_validator_key_file").(string)
	if !filepath.IsAbs(prvkey) {
		prvkey = filepath.Join(appOpts.Get("home").(string), prvkey)
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

	if pvKey.PrivKey.Type() != "secp256k1" {
		panic(prvkey + " is not an secp256k1 key")
	}
	return &secp256k1.PrivKey{Key: pvKey.PrivKey.Bytes()}
}
