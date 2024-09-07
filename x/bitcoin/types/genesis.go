package types

import (
	"bytes"
	"errors"
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:           DefaultParams(),
		StartBlockNumber: 0,
		BlockHash:        [][]byte{bytes.Repeat([]byte{0}, 32)},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	if err := gs.Params.Validate(); err != nil {
		return err
	}

	if gs.Pubkey != nil {
		if err := gs.Pubkey.Validate(); err != nil {
			return err
		}
	}

	if len(gs.BlockHash) == 0 {
		return errors.New("No block hash provided in the genesis state")
	}

	for _, v := range gs.BlockHash {
		if len(v) != 32 {
			return fmt.Errorf("invalid block hash: %x", v)
		}
	}

	return nil
}
