package relayer_test

import (
	"testing"

	keepertest "github.com/goatnetwork/goat/testutil/keeper"
	"github.com/goatnetwork/goat/testutil/nullify"
	relayer "github.com/goatnetwork/goat/x/relayer/module"
	"github.com/goatnetwork/goat/x/relayer/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, _ := keepertest.RelayerKeeper(t)
	relayer.InitGenesis(ctx, k, genesisState)
	got := relayer.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
