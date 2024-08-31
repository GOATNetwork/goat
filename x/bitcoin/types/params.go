package types

import (
	"errors"
	"fmt"
)

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{
		SafeConfirmationBlock: 3,
		HardConfirmationBlock: 6,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.HardConfirmationBlock == 0 || p.SafeConfirmationBlock == 0 {
		return errors.New("mempool txs are not reliable (confirmation number can't set to zero)")
	}
	if p.HardConfirmationBlock < p.SafeConfirmationBlock {
		return fmt.Errorf("hard block(%d) < safe block(%d)", p.HardConfirmationBlock, p.SafeConfirmationBlock)
	}
	return nil
}
