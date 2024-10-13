package cmd

import (
	cmtcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

const (
	FlagGoatGeth = "goat.geth"
	FlagJwtPath  = "goat.jwt"
)

type GoatConfig struct {
	Geth string `mapstructure:"geth"`
	JWT  string `mapstructure:"jwt"`
}

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

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
# the jwt secret file for engine api, it's only required if connecting to an execution node via HTTP.
jwt = "{{ .Goat.jwt }}"
`

	return customAppTemplate, customAppConfig
}
