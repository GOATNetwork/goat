package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/goatnetwork/goat/testutil/keeper"
	"github.com/goatnetwork/goat/x/locking/keeper"
	"github.com/goatnetwork/goat/x/locking/types"
)

func TestParamsQuery(t *testing.T) {
	k, ctx, _ := keepertest.LockingKeeper(t)

	qs := keeper.NewQueryServerImpl(k)
	params := types.DefaultParams()
	require.NoError(t, k.Params.Set(ctx, params))

	response, err := qs.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
