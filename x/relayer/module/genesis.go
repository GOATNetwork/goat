package relayer

import (
	"bytes"
	"fmt"
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

	if err := k.Params.Set(ctx, genState.Params); err != nil {
		panic(err)
	}

	if err := k.Sequence.Set(ctx, 0); err != nil {
		panic(err)
	}

	relayer := types.Relayer{
		Epoch:            genState.Epoch,
		ProposerAccepted: true,
		LastElected:      ctx.BlockTime(),
	}

	queue := types.VoterQueue{}

	if len(genState.Voters) == 0 {
		panic("No relayer voters")
	}

	keySet := make(map[string]bool)
	for addr, v := range genState.Voters {
		if _, err := k.AddrCodec.StringToBytes(addr); err != nil {
			panic(err)
		}

		if err := v.Validate(); err != nil {
			panic(err)
		}

		if keySet[string(v.VoteKey)] {
			panic(fmt.Sprintf("duplicated key: %x", v.VoteKey))
		}
		keySet[string(v.VoteKey)] = true

		relayer.Voters = append(relayer.Voters, addr)
		if err := k.Voters.Set(ctx, addr, *v); err != nil {
			panic(err)
		}

		switch v.Status {
		case types.VOTER_STATUS_ON_BOARDING:
			queue.OnBoarding = append(queue.OnBoarding, addr)
		case types.VOTER_STATUS_OFF_BOARDING:
			queue.OffBoarding = append(queue.OnBoarding, addr)
		}
	}

	slices.Sort(relayer.Voters)
	relayer.Proposer = relayer.Voters[0]
	relayer.Voters = relayer.Voters[1:]
	if err := k.Relayer.Set(ctx, relayer); err != nil {
		panic(err)
	}

	if err := k.Sequence.Set(ctx, genState.Sequence); err != nil {
		panic(err)
	}

	if err := k.Randao.Set(ctx, bytes.Repeat([]byte{0}, 32)); err != nil {
		panic(err)
	}

	if err := k.Queue.Set(ctx, queue); err != nil {
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

	genesis.Epoch = relayer.Epoch

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
		genesis.Voters[kv.Key] = &kv.Value
	}

	genesis.Sequence, err = k.Sequence.Peek(ctx)
	if err != nil {
		panic(err)
	}

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
