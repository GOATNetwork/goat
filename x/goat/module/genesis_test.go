package goat_test

import (
	"testing"

	keepertest "github.com/goatnetwork/goat/testutil/keeper"
	"github.com/goatnetwork/goat/testutil/nullify"
	goat "github.com/goatnetwork/goat/x/goat/module"
	"github.com/goatnetwork/goat/x/goat/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, _ := keepertest.GoatKeeper(t)
	goat.InitGenesis(ctx, k, genesisState)
	got := goat.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
