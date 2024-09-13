package relayer

import (
	"fmt"

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

	if genState.Relayer == nil {
		panic("no relayer")
	}

	if _, ok := genState.Voters[genState.Relayer.Proposer]; !ok {
		panic(fmt.Sprintf("missing proposer %s in the voter state", genState.Relayer.Proposer))
	}

	if _, err := k.AddrCodec.StringToBytes(genState.Relayer.Proposer); err != nil {
		panic(err)
	}

	if err := k.Relayer.Set(ctx, *genState.Relayer); err != nil {
		panic(err)
	}

	keySet := make(map[string]bool)
	for _, voter := range genState.Relayer.Voters {
		if keySet[voter] {
			panic("duplicated voter: " + voter)
		}
		keySet[voter] = true

		if voter == genState.Relayer.Proposer {
			panic("voter should not be a proposer")
		}

		if _, err := k.AddrCodec.StringToBytes(genState.Relayer.Proposer); err != nil {
			panic(err)
		}

		if _, ok := genState.Voters[voter]; !ok {
			panic(fmt.Sprintf("missing proposer %s in the voter state", genState.Relayer.Proposer))
		}
	}

	if err := k.Sequence.Set(ctx, genState.Sequence); err != nil {
		panic(err)
	}

	for _, pubkey := range genState.Pubkeys {
		if err := pubkey.Validate(); err != nil {
			panic(err)
		}

		if err := k.Pubkeys.Set(ctx, types.EncodePublicKey(pubkey)); err != nil {
			panic(err)
		}
	}

	queue := types.VoterQueue{}
	if len(genState.Voters) == 0 {
		panic("No relayer voters")
	}

	clear(keySet)
	for addr, v := range genState.Voters {
		if _, err := k.AddrCodec.StringToBytes(addr); err != nil {
			panic(err)
		}

		if err := v.Validate(); err != nil {
			panic(err)
		}

		if keySet[string(v.VoteKey)] {
			panic(fmt.Sprintf("duplicated vote key: %x", v.VoteKey))
		}

		keySet[string(v.VoteKey)] = true
		switch v.Status {
		case types.VOTER_STATUS_ON_BOARDING:
			queue.OnBoarding = append(queue.OnBoarding, addr)
		case types.VOTER_STATUS_OFF_BOARDING:
			queue.OffBoarding = append(queue.OnBoarding, addr)
		}
	}

	if err := k.Queue.Set(ctx, queue); err != nil {
		panic(err)
	}

	if err := k.Randao.Set(ctx, genState.Randao); err != nil {
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

	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		panic(err)
	}
	genesis.Relayer = &relayer

	genesis.Sequence, err = k.Sequence.Peek(ctx)
	if err != nil {
		panic(err)
	}

	// voters
	{
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
	}

	// public keys
	{
		iter, err := k.Pubkeys.Iterate(ctx, nil)
		if err != nil {
			panic(err)
		}
		defer iter.Close()

		for ; iter.Valid(); iter.Next() {
			value, err := iter.Key()
			if err != nil {
				panic(err)
			}
			pubkey := types.DecodePublicKey(value)
			if pubkey == nil {
				panic(fmt.Sprintf("invalid public key %x to decode", value))
			}
			genesis.Pubkeys = append(genesis.Pubkeys, pubkey)
		}
	}

	genesis.Randao, err = k.Randao.Get(ctx)
	if err != nil {
		panic(err)
	}

	return genesis
}
