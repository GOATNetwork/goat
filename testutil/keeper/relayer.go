package keeper

import (
	"testing"
	"time"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/goatnetwork/goat/x/relayer/keeper"
	"github.com/goatnetwork/goat/x/relayer/types"
)

func RelayerKeeper(t testing.TB, accountKeeper types.AccountKeeper) (keeper.Keeper, sdk.Context, address.Codec) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	k := keeper.NewKeeper(
		cdc,
		addressCodec,
		runtime.NewKVStoreService(storeKey),
		accountKeeper,
		log.NewNopLogger(),
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{ChainID: "goat-unit-test", Time: time.Now().UTC()}, false, log.NewNopLogger())

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	if err := k.Queue.Set(ctx, types.VoterQueue{}); err != nil {
		t.Fatalf("failed to set VoterQueue: %v", err)
	}

	if err := k.Randao.Set(ctx, make([]byte, 32)); err != nil {
		t.Fatalf("failed to set Randao: %v", err)
	}

	if err := k.Sequence.Set(ctx, 0); err != nil {
		t.Fatalf("failed to set Sequence: %v", err)
	}

	return k, ctx, addressCodec
}
