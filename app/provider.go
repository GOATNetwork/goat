package app

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
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
