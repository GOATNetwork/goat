package app

import (
	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	"cosmossdk.io/core/appconfig"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	bitcoinmodulev1 "github.com/goatnetwork/goat/api/goat/bitcoin/module/v1"
	relayermodulev1 "github.com/goatnetwork/goat/api/goat/relayer/module/v1"
	_ "github.com/goatnetwork/goat/x/relayer/module"
	relayermoduletypes "github.com/goatnetwork/goat/x/relayer/types"

	goatmodulev1 "github.com/goatnetwork/goat/api/goat/goat/module/v1"
	lockingmodulev1 "github.com/goatnetwork/goat/api/goat/locking/module/v1"
	_ "github.com/goatnetwork/goat/x/bitcoin/module"
	bitcoinmoduletypes "github.com/goatnetwork/goat/x/bitcoin/types"
	_ "github.com/goatnetwork/goat/x/goat/module"
	goatmoduletypes "github.com/goatnetwork/goat/x/goat/types"
	_ "github.com/goatnetwork/goat/x/locking/module"
	lockingmoduletypes "github.com/goatnetwork/goat/x/locking/types"
)

var (
	genesisModuleOrder = []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		relayermoduletypes.ModuleName,
		bitcoinmoduletypes.ModuleName,
		lockingmoduletypes.ModuleName,
		goatmoduletypes.ModuleName,
	}

	beginBlockers = []string{}

	endBlockers = []string{
		relayermoduletypes.ModuleName,
		goatmoduletypes.ModuleName,
	}

	preBlockers = []string{}

	// module account permissions
	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: goatmoduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
		// this line is used by starport scaffolding # stargate/app/maccPerms
	}

	// blocked account addresses
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
	}

	// appConfig application configuration (used by depinject)
	appConfig = appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: runtime.ModuleName,
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName:       Name,
					PreBlockers:   preBlockers,
					BeginBlockers: beginBlockers,
					EndBlockers:   endBlockers,
					InitGenesis:   genesisModuleOrder,
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: appconfig.WrapAny(&authmodulev1.Module{
					Bech32Prefix:             AccountAddressPrefix,
					ModuleAccountPermissions: moduleAccPerms,
				}),
			},
			{
				Name:   "tx",
				Config: appconfig.WrapAny(&txconfigv1.Config{SkipAnteHandler: true}),
			},
			{
				Name: banktypes.ModuleName,
				Config: appconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name:   consensustypes.ModuleName,
				Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
			},
			{
				Name:   relayermoduletypes.ModuleName,
				Config: appconfig.WrapAny(&relayermodulev1.Module{}),
			},
			{
				Name:   bitcoinmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&bitcoinmodulev1.Module{}),
			},
			{
				Name:   lockingmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&lockingmodulev1.Module{}),
			},
			{
				Name:   goatmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&goatmodulev1.Module{}),
			},
			// this line is used by starport scaffolding # stargate/app/moduleConfig
		},
	})
)
