package types

import (
	"slices"

	"cosmossdk.io/math"

	tmcrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
)

var (
	PowerReduction = math.NewIntFromUint64(1e18)
)

func (v *Validator) CMPubkey() tmcrypto.PublicKey {
	return tmcrypto.PublicKey{Sum: &tmcrypto.PublicKey_Secp256K1{Secp256K1: slices.Clone(v.Pubkey)}}
}
