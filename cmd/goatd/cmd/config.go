package cmd

import (
	"time"

	cmtcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

const (
	FlagGoatGeth = "goat.geth"
)

type GoatConfig struct {
	Geth string `mapstructure:"geth"`
}

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	cfg.P2P.MaxNumInboundPeers = 150
	cfg.P2P.MaxNumOutboundPeers = 100
	cfg.Mempool.Size = 10
	cfg.Consensus.TimeoutPropose = 1500 * time.Millisecond
	cfg.Consensus.TimeoutPrevote = 1500 * time.Millisecond
	cfg.Consensus.TimeoutPrecommit = time.Second
	cfg.Consensus.TimeoutCommit = time.Second * 3

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// The following code snippet is just for reference.
	type GoatAppConfig struct {
		serverconfig.Config `mapstructure:",squash"`
		Goat                GoatConfig `mapstructure:"goat"`
	}

	srvCfg := serverconfig.DefaultConfig()
	srvCfg.MinGasPrices = "0gas"
	srvCfg.Mempool.MaxTxs = 10

	customAppConfig := GoatAppConfig{
		Config: *srvCfg,
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate + `
[goat]
# the goat-geth node endpoint, using ipc is recommended
geth = "{{ .Goat.Geth }}"
`

	return customAppTemplate, customAppConfig
}
