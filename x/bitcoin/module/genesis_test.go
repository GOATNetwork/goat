package bitcoin_test

import (
	"testing"

	keepertest "github.com/goatnetwork/goat/testutil/keeper"
	"github.com/goatnetwork/goat/testutil/nullify"
	bitcoin "github.com/goatnetwork/goat/x/bitcoin/module"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, _ := keepertest.BitcoinKeeper(t, nil)
	bitcoin.InitGenesis(ctx, k, genesisState)
	got := bitcoin.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
