package locking

import (
	"cosmossdk.io/collections"
	abci "github.com/cometbft/cometbft/abci/types"
	tmcrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/locking/keeper"
	"github.com/goatnetwork/goat/x/locking/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) (vs []abci.ValidatorUpdate) {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		panic(err)
	}

	// validator
	for _, validator := range genState.Validators {
		pubkey := &secp256k1.PubKey{Key: validator.Pubkey}
		address := sdk.ConsAddress(pubkey.Address())

		err := k.Validators.Set(ctx, address, validator)
		if err != nil {
			panic(err)
		}

		if validator.Status != types.Active && validator.Status != types.Pending {
			continue
		}

		for _, locking := range validator.Locking {
			err = k.Locking.Set(ctx,
				collections.Join(locking.Denom, address), locking.Amount)
			if err != nil {
				panic(err)
			}
		}

		if validator.Status == types.Active {
			if err := k.ValidatorSet.Set(ctx, address, validator.Power); err != nil {
				panic(err)
			}
			vs = append(vs, abci.ValidatorUpdate{
				Power: int64(validator.Power),
				PubKey: tmcrypto.PublicKey{
					Sum: &tmcrypto.PublicKey_Secp256K1{Secp256K1: validator.Pubkey},
				},
			})
		}

		err = k.PowerRanking.Set(ctx, collections.Join(validator.Power, address))
		if err != nil {
			panic(err)
		}
	}

	// tokens
	{
		var threshold sdk.Coins
		for _, token := range genState.Tokens {
			err := k.Tokens.Set(ctx, token.Denom, token.Token)
			if err != nil {
				panic(err)
			}
			if !token.Token.Threshold.IsZero() {
				threshold = threshold.Add(sdk.NewCoin(token.Denom, token.Token.Threshold))
			}
		}
		if err := k.Threshold.Set(ctx, types.Threshold{List: threshold}); err != nil {
			panic(err)
		}
	}

	// slashed
	{
		for _, token := range genState.Slashed {
			err := k.Slashed.Set(ctx, token.Denom, token.Amount)
			if err != nil {
				panic(err)
			}
		}
	}

	// eth tx queue
	{
		if err := k.EthTxNonce.Set(ctx, genState.EthTxNonce); err != nil {
			panic(err)
		}

		if err := k.EthTxQueue.Set(ctx, genState.EthTxQueue); err != nil {
			panic(err)
		}
	}

	// reward pool
	{
		err := k.RewardPool.Set(ctx, genState.RewardPool)
		if err != nil {
			panic(err)
		}
	}

	// unlocks
	{
		for _, unlock := range genState.UnlockQueue {
			err := k.UnlockQueue.Set(ctx, unlock.Timestamp, types.Unlocks{Unlocks: unlock.Unlocks})
			if err != nil {
				panic(err)
			}
		}
	}

	return vs
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var err error

	genesis := new(types.GenesisState)

	{
		genesis.Params, err = k.Params.Get(ctx)
		if err != nil {
			panic(err)
		}
	}

	// validator
	{
		iter, err := k.Validators.Iterate(ctx, nil)
		if err != nil {
			panic(err)
		}
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			kv, err := iter.Value()
			if err != nil {
				panic(err)
			}
			genesis.Validators = append(genesis.Validators, kv)
		}
	}

	// tokens
	{
		iter, err := k.Tokens.Iterate(ctx, nil)
		if err != nil {
			panic(err)
		}
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			kv, err := iter.KeyValue()
			if err != nil {
				panic(err)
			}
			genesis.Tokens = append(genesis.Tokens, &types.TokenGenesis{
				Denom: kv.Key,
				Token: kv.Value,
			})
		}
	}

	// slashed
	{
		iter, err := k.Slashed.Iterate(ctx, nil)
		if err != nil {
			panic(err)
		}
		defer iter.Close()

		coins := sdk.Coins{}
		for ; iter.Valid(); iter.Next() {
			kv, err := iter.KeyValue()
			if err != nil {
				panic(err)
			}
			coins = coins.Add(sdk.NewCoin(kv.Key, kv.Value))
		}
		genesis.Slashed = coins
	}

	// eth tx queue
	{
		genesis.EthTxNonce, err = k.EthTxNonce.Peek(ctx)
		if err != nil {
			panic(err)
		}
		genesis.EthTxQueue, err = k.EthTxQueue.Get(ctx)
		if err != nil {
			panic(err)
		}
	}

	// reward pool
	{
		genesis.RewardPool, err = k.RewardPool.Get(ctx)
		if err != nil {
			panic(err)
		}
	}

	// unlocking
	{
		iter, err := k.UnlockQueue.Iterate(ctx, nil)
		if err != nil {
			panic(err)
		}
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			kv, err := iter.KeyValue()
			if err != nil {
				panic(err)
			}
			genesis.UnlockQueue = append(genesis.UnlockQueue, &types.UnlockQueueGenesis{
				Timestamp: kv.Key,
				Unlocks:   kv.Value.Unlocks,
			})
		}
	}

	return genesis
}
