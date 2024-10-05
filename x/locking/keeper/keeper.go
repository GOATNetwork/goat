package keeper

import (
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
		PowerRanking collections.KeySet[collections.Pair[uint64, sdktypes.ConsAddress]]
		ValidatorSet collections.Map[sdktypes.ConsAddress, uint64]
		Validators   collections.Map[sdktypes.ConsAddress, types.Validator]
		Tokens       collections.Map[string, types.Token]
		Threshold    collections.Item[types.Threshold]
		Slashed      collections.Map[string, math.Int]
		EthTxNonce   collections.Sequence
		RewardPool   collections.Item[types.RewardPool]
		EthTxQueue   collections.Item[types.EthTxQueue]
		UnlockQueue  collections.Map[time.Time, types.Unlocks]
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

		Params:       collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Locking:      collections.NewMap(sb, types.LockingKey, "locking", collections.PairKeyCodec(collections.StringKey, sdktypes.ConsAddressKey), sdktypes.IntValue),
		PowerRanking: collections.NewKeySet(sb, types.PowerRankingKey, "power_ranking", collections.PairKeyCodec(collections.Uint64Key, sdktypes.ConsAddressKey)),
		ValidatorSet: collections.NewMap(sb, types.ValidatorSetKey, "last_validator_set", sdktypes.ConsAddressKey, collections.Uint64Value),
		Validators:   collections.NewMap(sb, types.ValidatorsKey, "validator", sdktypes.ConsAddressKey, codec.CollValue[types.Validator](cdc)),
		Tokens:       collections.NewMap(sb, types.TokensKey, "token", collections.StringKey, codec.CollValue[types.Token](cdc)),
		Slashed:      collections.NewMap(sb, types.SlashedKey, "slashed", collections.StringKey, sdktypes.IntValue),
		EthTxNonce:   collections.NewSequence(sb, types.EthTxNonceKey, "eth_tx_nonce"),
		EthTxQueue:   collections.NewItem(sb, types.EthTxQueueKey, "eth_tx_queue", codec.CollValue[types.EthTxQueue](cdc)),
		RewardPool:   collections.NewItem(sb, types.RewardPoolKey, "reward_pool", codec.CollValue[types.RewardPool](cdc)),
		UnlockQueue:  collections.NewMap(sb, types.UnlockQueueKey, "unlock_queue", sdktypes.TimeKey, codec.CollValue[types.Unlocks](cdc)),
		Threshold:    collections.NewItem(sb, types.ThresholdKey, "threshold", codec.CollValue[types.Threshold](cdc)),
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
