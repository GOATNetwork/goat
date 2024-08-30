package types

import (
	"errors"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:    DefaultParams(),
		Threshold: 0,
		Voters:    make(map[string]*Voter),
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
