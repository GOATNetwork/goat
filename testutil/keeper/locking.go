package keeper

import (
	"testing"

	"cosmossdk.io/log/v2"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/store/v2"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/locking/keeper"
	"github.com/goatnetwork/goat/x/locking/types"
	"github.com/stretchr/testify/require"
)

func LockingKeeper(tb testing.TB, accountKeeper types.AccountKeeper) (keeper.Keeper, sdk.Context, *ConsensusParamsStore) {
	tb.Helper()

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(tb, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	consensusParams := NewConsensusParamsStore(types.KeyTypeSecp256k1)

	k := keeper.NewKeeper(
		cdc,
		addressCodec,
		runtime.NewKVStoreService(storeKey),
		accountKeeper,
		consensusParams,
		log.NewNopLogger(),
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		tb.Fatalf("failed to set params: %v", err)
	}

	return k, ctx, consensusParams
}
