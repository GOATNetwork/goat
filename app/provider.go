package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cosmossdk.io/log"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/privval"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/goatnetwork/goat/pkg/ethrpc"
	bitcintypes "github.com/goatnetwork/goat/x/bitcoin/types"
	"github.com/spf13/cast"
)

func ProvideEngineClient(logger log.Logger, appOpts servertypes.AppOptions) *ethrpc.Client {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	for i := 0; i < 10; i++ {
		var err error
		logger.Info("try to connect goat-geth", "endpoint", endpoint)
		ethclient, err = ethrpc.DialContext(ctx, endpoint, jwtSecret)
		if err == nil {
			var conf *params.ChainConfig
			conf, err = ethclient.GetChainConfig(ctx)
			if err == nil {
				if conf.Goat == nil {
					panic("No goat config found in the goat-geth, please verify if you're using correct setup")
				}
				break
			}
		}
		logger.Error("retry to connect goat-geth", "err", err.Error())
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

	if pvKey.PrivKey.Type() != bitcintypes.Secp256K1Name {
		panic(prvkey + " is not an secp256k1 key")
	}
	return &secp256k1.PrivKey{Key: pvKey.PrivKey.Bytes()}
}
