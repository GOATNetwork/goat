package types

import (
	"errors"
	"fmt"
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

	if gs.Threshold != 0 && gs.Threshold > uint64(len(gs.Voters)) {
		return errors.New("invalid proposal threshold")
	}

	if gs.Threshold == 0 && len(gs.Voters) != 0 {
		return errors.New("threshold shuould not be 0 if voter length is not 0")
	}

	for addr, vote := range gs.Voters {
		if err := vote.Validate(); err != nil {
			return fmt.Errorf("invalid vote key of %s", addr)
		}
	}

	return nil
}
