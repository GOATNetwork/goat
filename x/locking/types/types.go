package types

import (
	"slices"

	"cosmossdk.io/math"

	tmcrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	PowerReduction = math.NewIntFromUint64(1e18)
)

func (v *Validator) CMPubkey() tmcrypto.PublicKey {
	return tmcrypto.PublicKey{Sum: &tmcrypto.PublicKey_Secp256K1{Secp256K1: slices.Clone(v.Pubkey)}}
}

type ExecRequests struct {
	GasRevenues      []*ethtypes.GasRevenue
	Creates          []*ethtypes.CreateValidator
	Locks            []*ethtypes.GoatLock
	Unlocks          []*ethtypes.GoatUnlock
	Claims           []*ethtypes.GoatClaimReward
	Grants           []*ethtypes.GoatGrant
	UpdateWeights    []*ethtypes.UpdateTokenWeight
	UpdateThresholds []*ethtypes.UpdateTokenThreshold
}
