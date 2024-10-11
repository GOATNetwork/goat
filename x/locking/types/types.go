package types

import (
	"encoding/hex"
	"slices"

	"cosmossdk.io/math"

	tmcrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
)

var (
	PowerReduction = math.NewIntFromUint64(1e18)
)

func (v *Validator) CMPubkey() tmcrypto.PublicKey {
	return tmcrypto.PublicKey{Sum: &tmcrypto.PublicKey_Secp256K1{Secp256K1: slices.Clone(v.Pubkey)}}
}

func TokenDenom(token common.Address) string {
	switch token {
	case common.Address{}:
		return "btc"
	case goattypes.GoatTokenContract:
		return "goat"
	}
	return "tkn:" + hex.EncodeToString(token.Bytes())
}

func ValidatorName(val sdktypes.ConsAddress) string {
	return hex.EncodeToString(val)
}
