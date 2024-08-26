package types

import (
	"bytes"
	"errors"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	txkey := &secp256k1.PrivKey{Key: bytes.Repeat([]byte{0xde, 0xad}, 16)}

	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:           DefaultParams(),
		StartBlockNumber: 0,
		BlockHash:        [][]byte{bytes.Repeat([]byte{0}, 32)},
		Pubkey: &relayer.PublicKey{Key: &relayer.PublicKey_Secp256K1{
			Secp256K1: txkey.PubKey().Bytes(),
		}},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	if err := gs.Params.Validate(); err != nil {
		return err
	}

	if err := gs.Pubkey.Validate(); err != nil {
		return err
	}

	if len(gs.BlockHash) == 0 {
		return errors.New("No block hash provided in the genesis state")
	}

	return nil
}
