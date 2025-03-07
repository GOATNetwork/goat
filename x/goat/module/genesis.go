package goat

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/goat/keeper"
	"github.com/goatnetwork/goat/x/goat/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		panic(err)
	}

	if err := k.Block.Set(ctx, genState.EthBlock); err != nil {
		panic(err)
	}

	if err := k.BeaconRoot.Set(ctx, genState.BeaconRoot); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var err error

	genesis := new(types.GenesisState)
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	genesis.EthBlock, err = k.Block.Get(ctx)
	if err != nil {
		panic(err)
	}

	genesis.BeaconRoot, err = k.BeaconRoot.Get(ctx)
	if err != nil {
		panic(err)
	}

	return genesis
}
