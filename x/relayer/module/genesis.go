package relayer

import (
	"bytes"
	"slices"

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

	for addr, v := range genState.Voters {
		addrByte, err := k.AddrCodec.StringToBytes(addr)
		if err != nil {
			panic(err)
		}

		if len(v.VoteKey) != 96 {
			panic("invalid vote key")
		}

		relayer.Voters = append(relayer.Voters, addr)
		v.Status = types.Activated
		v.Height = sdkctx.BlockHeight()
		if err := k.Voters.Set(ctx, addrByte, *v); err != nil {
			panic(err)
		}
	}

	slices.Sort(relayer.Voters)
	relayer.Proposer = relayer.Voters[0]
	relayer.Voters = relayer.Voters[1:]
	if err := k.Relayer.Set(ctx, relayer); err != nil {
		panic(err)
	}

	if err := k.Randao.Set(ctx, bytes.Repeat([]byte{0}, 32)); err != nil {
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
		kv, err := iter.KeyValue()
		if err != nil {
			panic(err)
		}

		if kv.Value.Status == types.Activated {
			genesis.Voters[kv.Key.String()] = &kv.Value
		}
	}

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
