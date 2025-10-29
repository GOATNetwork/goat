package cmd

import (
	"time"

	cmtcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

const (
	FlagGoatGeth        = "goat.geth"
	FlagGoatPreset      = "goat.preset"
	FlagPersistentPeers = "p2p.persistent_peers"
	// FlagExternalIP      = "p2p.external_address"
)

type GoatConfig struct {
	Geth   string `mapstructure:"geth"`
	Preset string `mapstructure:"preset"`
}

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	cfg.Mempool.Size = 10
	cfg.Consensus.TimeoutPropose = 1500 * time.Millisecond
	cfg.Consensus.TimeoutPrevote = 1500 * time.Millisecond
	cfg.Consensus.TimeoutPrecommit = time.Second
	cfg.Consensus.TimeoutCommit = time.Second * 3

	return cfg
}

func initRegtestCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	cfg.Consensus.TimeoutPropose = 1500 * time.Millisecond
	cfg.Consensus.TimeoutPrevote = 500 * time.Millisecond
	cfg.Consensus.TimeoutPrecommit = 500 * time.Millisecond
	// the geth can't handle the duration which is less than 1s
	// and we can't use 1s due to cosmos-sdk updates it to 5s by default
	cfg.Consensus.TimeoutCommit = 1500 * time.Millisecond
	cfg.P2P.PexReactor = false
	cfg.Moniker = "regtest"

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, any) {
	// The following code snippet is just for reference.
	type GoatAppConfig struct {
		serverconfig.Config `mapstructure:",squash"`
		Goat                GoatConfig `mapstructure:"goat"`
	}

	srvCfg := serverconfig.DefaultConfig()
	srvCfg.MinGasPrices = "0gas"
	srvCfg.Mempool.MaxTxs = 10
	srvCfg.GRPCWeb.Enable = false

	customAppConfig := GoatAppConfig{
		Config: *srvCfg,
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate + `
[goat]
# the goat-geth node ipc path
# we don't use http due to the node server has body limit for it
geth = "{{ .Goat.Geth }}"
# the node preset configuration, e.g. rpc,bootnode
preset = "{{ .Goat.Preset }}"
`

	return customAppTemplate, customAppConfig
}
