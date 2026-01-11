package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/goatnetwork/goat/pkg/ethrpc"
)

func ConnectEngineClient(ctx context.Context, logger log.Logger, endpoint string) (*ethrpc.Client, *params.ChainConfig, error) {
	if endpoint == "" {
		return nil, nil, errors.New("goat-geth node endpoint not found")
	}
	logger.Info("try to connect goat-geth", "endpoint", endpoint)
	ctx, cancel := context.WithTimeout(ctx, time.Second*7)
	defer cancel()
	for range 10 {
		ethclient, err := ethrpc.DialContext(ctx, endpoint)
		if err == nil {
			var conf *params.ChainConfig
			conf, err = ethclient.GetChainConfig(ctx)
			if err == nil {
				if conf.Goat == nil {
					return nil, nil, errors.New("invalid goat-geth node, please verify if you're using correct setup")
				}
				// goat-geth upgrade check
				if conf.OsakaTime == nil {
					return nil, nil, errors.New("osaka time is undefined in goat-geth node, please upgrade your goat-geth node")
				}
				return ethclient, conf, nil
			}
		}
		logger.Warn("retry to connect goat-geth", "err", err.Error())
		<-time.After(time.Second / 2)
	}

	return nil, nil, fmt.Errorf("can not connect to goat-geth via %s", endpoint)
}
