package types

import (
	"encoding/hex"

	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
)

const (
	GoatTokenDenom = "goat"
	BitcoinDenom   = "btc"
)

var PowerReduction = math.NewIntFromUint64(1e18)

func TokenDenom(token common.Address) string {
	switch token {
	case common.Address{}:
		return BitcoinDenom
	case goattypes.GoatTokenContract:
		return GoatTokenDenom
	}
	return "tkn:" + hex.EncodeToString(token.Bytes())
}

func ValidatorName(val sdktypes.ConsAddress) string {
	return hex.EncodeToString(val)
}
