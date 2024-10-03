package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/goatnetwork/goat/x/locking/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		addressCodec  address.Codec
		storeService  store.KVStoreService
		logger        log.Logger
		accountKeeper types.AccountKeeper

		Schema collections.Schema
		Params collections.Item[types.Params]
		// (token,validator) => locking, it's used for updating power when the token weight is updated
		Locking collections.Map[collections.Pair[string, sdktypes.ConsAddress], math.Int]
		// (power,validator) => int64(power), it's used for getting validators of top-k power
		PowerRanking   collections.KeySet[collections.Pair[uint64, sdktypes.ConsAddress]]
		ValidatorSet   collections.KeySet[sdktypes.ConsAddress]
		Validators     collections.Map[sdktypes.ConsAddress, types.Validator]
		Tokens         collections.Map[string, types.Token]
		Slashed        collections.Map[string, math.Int]
		EthTxNonce     collections.Sequence
		RewardPool     collections.Item[types.RewardPool]
		ExecuableQueue collections.Item[types.ExecuableQueue]
		UnlockQueue    collections.Map[time.Time, types.UnlockQueue]
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	logger log.Logger,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:           cdc,
		addressCodec:  addressCodec,
		storeService:  storeService,
		logger:        logger,
		accountKeeper: accountKeeper,

		Params:         collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Locking:        collections.NewMap(sb, types.LockingKey, "locking", collections.PairKeyCodec(collections.StringKey, sdktypes.ConsAddressKey), sdktypes.IntValue),
		PowerRanking:   collections.NewKeySet(sb, types.PowerRankingKey, "power_ranking", collections.PairKeyCodec(collections.Uint64Key, sdktypes.ConsAddressKey)),
		ValidatorSet:   collections.NewKeySet(sb, types.ValidatorSetKey, "validator_set", sdktypes.ConsAddressKey),
		Validators:     collections.NewMap(sb, types.ValidatorsKey, "validator", sdktypes.ConsAddressKey, codec.CollValue[types.Validator](cdc)),
		Tokens:         collections.NewMap(sb, types.TokensKey, "token", collections.StringKey, codec.CollValue[types.Token](cdc)),
		Slashed:        collections.NewMap(sb, types.SlashedKey, "slashed", collections.StringKey, sdktypes.IntValue),
		EthTxNonce:     collections.NewSequence(sb, types.EthTxNonceKey, "eth_tx_nonce"),
		RewardPool:     collections.NewItem(sb, types.RewardPoolKey, "reward_pool", codec.CollValue[types.RewardPool](cdc)),
		ExecuableQueue: collections.NewItem(sb, types.ExecuableQueueKey, "execuable_queue", codec.CollValue[types.ExecuableQueue](cdc)),
		UnlockQueue:    collections.NewMap(sb, types.UnlockQueueKey, "unlock_queue", sdktypes.TimeKey, codec.CollValue[types.UnlockQueue](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) ThresholdList(ctx context.Context) (sdktypes.Coins, error) {
	iter, err := k.Tokens.Iterate(sdktypes.UnwrapSDKContext(ctx), nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	res := sdktypes.Coins{}
	for ; iter.Valid(); iter.Next() {
		kv, err := iter.KeyValue()
		if err != nil {
			return nil, err
		}
		if !kv.Value.Threshold.IsZero() {
			res = append(res, sdktypes.NewCoin(kv.Key, kv.Value.Threshold))
		}
	}

	return res.Sort(), nil
}
