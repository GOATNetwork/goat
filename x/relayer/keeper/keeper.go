package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/relayer/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		schema       collections.Schema

		AddrCodec address.Codec
		Params    collections.Item[types.Params]
		Relayer   collections.Item[types.Relayer]
		Sequence  collections.Sequence
		Voters    collections.Map[string, types.Voter]
		Queue     collections.Item[types.VoterQueue]
		Pubkeys   collections.KeySet[[]byte]
		Randao    collections.Item[[]byte]

		accountKeeper types.AccountKeeper
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
		AddrCodec:     addressCodec,
		storeService:  storeService,
		logger:        logger,
		accountKeeper: accountKeeper,

		Params:   collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Relayer:  collections.NewItem(sb, types.RelayerKey, "relayer", codec.CollValue[types.Relayer](cdc)),
		Sequence: collections.NewSequence(sb, types.SequenceKey, "sequence"),
		Voters:   collections.NewMap(sb, types.VotersKey, "voters", collections.StringKey, codec.CollValue[types.Voter](cdc)),
		Pubkeys:  collections.NewKeySet(sb, types.PubkeysKey, "pubkeys", collections.BytesKey),
		Queue:    collections.NewItem(sb, types.QueueKey, "queue", codec.CollValue[types.VoterQueue](cdc)),
		Randao:   collections.NewItem(sb, types.RandDAOKey, "randao", collections.BytesValue),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.schema = schema

	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) UpdateRandao(ctx context.Context, req types.IVoteMsg) error {
	randao, err := k.Randao.Get(ctx)
	if err != nil {
		return err
	}

	newRandao := goatcrypto.SHA256Sum(randao, req.GetVote().Signature)
	if err := k.Randao.Set(ctx, newRandao); err != nil {
		return err
	}
	k.Logger().Info("Randao updated", "current", hex.EncodeToString(newRandao))
	return nil
}

func (k Keeper) HasPubkey(ctx context.Context, raw []byte) (bool, error) {
	return k.Pubkeys.Has(ctx, raw)
}

func (k Keeper) AddNewKey(ctx context.Context, raw []byte) error {
	return k.Pubkeys.Set(ctx, raw)
}

func (k Keeper) SetProposalSeq(ctx context.Context, seq uint64) error {
	return k.Sequence.Set(ctx, seq)
}

func (k Keeper) GetCurrentProposer(ctx context.Context) (sdktypes.AccAddress, error) {
	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return nil, err
	}
	return k.AddrCodec.StringToBytes(relayer.Proposer)
}
