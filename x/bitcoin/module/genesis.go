package bitcoin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/goatnetwork/goat/x/bitcoin/keeper"
	"github.com/goatnetwork/goat/x/bitcoin/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := genState.Validate(); err != nil {
		panic(err)
	}

	// this line is used by starport scaffolding # genesis/module/init
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		panic(err)
	}

	if err := k.NewPubkey(ctx, genState.Pubkey); err != nil {
		panic(err)
	}

	blockNumber := genState.StartBlockNumber
	for _, v := range genState.BlockHash {
		if err := k.BlockHashes.Set(ctx, blockNumber, v); err != nil {
			panic(err)
		}
		blockNumber++
	}
	if err := k.BlockTip.Set(ctx, blockNumber); err != nil {
		panic(err)
	}

	queue := types.ExecuableQueue{
		BlockNumber: blockNumber,
	}
	if err := k.ExecuableQueue.Set(ctx, queue); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	pubkey, err := k.Pubkey.Get(ctx)
	if err != nil {
		panic(err)
	}

	genesis.Pubkey = &pubkey

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
