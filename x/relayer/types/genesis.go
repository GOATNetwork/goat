package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Voters: []*Voter{},
		Randao: make([]byte, 32),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	for _, voter := range gs.Voters {
		if err := voter.Validate(); err != nil {
			return fmt.Errorf("invalid voter of %s: %w", voter.Address, err)
		}
	}

	return nil
}
