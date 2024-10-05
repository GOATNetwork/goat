package types

import (
	"slices"

	"cosmossdk.io/math"

	tmcrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
)

var (
	PowerReduction = math.NewIntFromUint64(1e18)
)

const (
	HalvingInterval    = 42048000
	InitialBlockReward = 2378234400000000000
)

func (v *Validator) CMPubkey() tmcrypto.PublicKey {
	return tmcrypto.PublicKey{Sum: &tmcrypto.PublicKey_Secp256K1{Secp256K1: slices.Clone(v.Pubkey)}}
}
