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
	"github.com/spf13/cast"
)

func ProvideBitcoinNetworkConfig(appOpts servertypes.AppOptions) *chaincfg.Params {
	name := cast.ToString(appOpts.Get("goat.btc-network"))
	if name == "" {
		panic("no bitcoin network name provided")
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
	endpoint := cast.ToString(appOpts.Get("goat.geth"))
	if endpoint == "" {
		panic("goat execution node endpoint not found")
	}

	jwtSecret := func() []byte {
		if jwtPath := cast.ToString(appOpts.Get("goat.jwt-path")); jwtPath != "" {
			data, err := os.ReadFile(jwtPath)
			if err != nil {
				panic("cannot open jwt secret file: " + jwtPath)
			}

			jwtSecret := common.FromHex(strings.TrimSpace(string(data)))
			if len(jwtSecret) != 32 {
				panic("jwt secret is not a 32 bytes hex string")
			}
			return jwtSecret
		}
		return nil
	}()

	var ethclient *ethrpc.Client

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	for i := 0; i < 15; i++ {
		var err error
		ethclient, err = ethrpc.DialContext(ctx, endpoint, jwtSecret)
		if err == nil {
			conf, err := ethclient.GetChainConfig(ctx)
			if err == nil {
				if conf.Goat == nil {
					panic("No goat config found in the goat-geth, please verify if you're using correct steup")
				}
				break
			}
		}
		<-time.After(time.Second / 2)
	}

	if ethclient == nil {
		panic("can not connect to goat-geth via " + endpoint)
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
