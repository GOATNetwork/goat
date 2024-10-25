package testutil

import (
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func EventEquals(tb testing.TB, ev1, ev2 []sdktypes.Event) {
	tb.Helper()
	require.Equal(tb, len(ev1), len(ev2))

	for idx := range len(ev1) {
		require.Equal(tb, ev1[idx].Type, ev2[idx].Type, "idx %d", idx)
		require.Equal(tb, len(ev1[idx].Attributes), len(ev2[idx].Attributes), "idx %d attr len", idx)

		for atIdx := range len(ev1[idx].Attributes) {
			require.Equal(tb, ev1[idx].Attributes[atIdx].Key, ev2[idx].Attributes[atIdx].Key, "idx %d attr key %d", idx, atIdx)
			require.Equal(tb, ev1[idx].Attributes[atIdx].Value, ev2[idx].Attributes[atIdx].Value, "idx %d attr val %d", idx, atIdx)
		}
	}
}
