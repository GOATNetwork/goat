package types

import (
	"bytes"
	"errors"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/goatnetwork/goat/pkg/crypto"
	blst "github.com/supranational/blst/bindings/go"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	ikm := bytes.Repeat([]byte{0xde, 0xad}, 16)

	txkey := &secp256k1.PrivKey{Key: ikm}
	voteKey := new(crypto.PublicKey).From(blst.KeyGenV3(ikm[:32]))

	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:    DefaultParams(),
		Threshold: 1,
		Voters: []*Voter{
			{
				TxKey:   txkey.PubKey().Bytes(),
				VoteKey: voteKey.Compress(),
				Status:  Activated,
			},
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	if err := gs.Params.Validate(); err != nil {
		return err
	}

	if len(gs.Voters) == 0 {
		return errors.New("No voters")
	}

	if gs.Threshold == 0 || gs.Threshold > uint64(len(gs.Voters)) {
		return errors.New("invalid proposal threshold")
	}

	return nil
}
