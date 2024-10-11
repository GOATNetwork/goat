package types

import (
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
)

func TestTokenDenom(t *testing.T) {
	tests := []struct {
		name string
		args common.Address
		want string
	}{
		{"1", common.Address{}, "btc"},
		{"2", goattypes.GoatTokenContract, "goat"},
		{"3", goattypes.BridgeContract, "tkn:bc10000000000000000000000000000000000003"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TokenDenom(tt.args)
			if got != tt.want {
				t.Errorf("TokenDenom() = %v, want %v", got, tt.want)
			}
			if err := sdktypes.ValidateDenom(got); err != nil {
				t.Errorf("not valid address: %s", err)
			}
		})
	}
}
