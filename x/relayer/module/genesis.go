package relayer

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/goatnetwork/goat/x/relayer/keeper"
	"github.com/goatnetwork/goat/x/relayer/types"
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

	if err := k.ProposalSeq.Set(ctx, 0); err != nil {
		panic(err)
	}

	if err := k.Epoch.Set(ctx, 0); err != nil {
		panic(err)
	}

	sdkctx := sdk.UnwrapSDKContext(ctx)

	relayer := types.Relayer{
		Threshold:   genState.Threshold,
		LastElected: ctx.BlockTime(),
	}

	for idx, v := range genState.Voters {
		key := &secp256k1.PubKey{Key: v.TxKey}
		addr, err := k.AddrCodec.BytesToString(key.Address())
		if err != nil {
			panic(err)
		}

		if len(v.VoteKey) != 96 {
			panic("invalid vote key")
		}

		if idx == 0 {
			relayer.Proposer = addr
		} else {
			relayer.Voters = append(relayer.Voters, addr)
		}

		v.Status = types.Activated
		v.Height = sdkctx.BlockHeight()
		if err := k.Voters.Set(ctx, addr, *v); err != nil {
			panic(err)
		}
	}

	if err := k.Relayer.Set(ctx, relayer); err != nil {
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

	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		panic(err)
	}

	genesis.Threshold = relayer.Threshold

	iter, err := k.Voters.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		kv, err := iter.Value()
		if err != nil {
			panic(err)
		}
		if kv.Status == types.Activated {
			genesis.Voters = append(genesis.Voters, &kv)
		}
	}

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
