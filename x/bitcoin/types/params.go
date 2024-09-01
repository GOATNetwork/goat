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
		MinDepositAmount:      546,
		DepositMagicPrefix:    []byte("GTT0"),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.MinDepositAmount < 546 {
		return errors.New("minimal deposit amount can't be less than dust value")
	}

	if len(p.DepositMagicPrefix) != 4 {
		return errors.New("invalid DepositMagicPrefix length")
	}

	if p.HardConfirmationBlock == 0 || p.SafeConfirmationBlock == 0 {
		return errors.New("mempool txs are not reliable (confirmation number can't set to zero)")
	}
	if p.HardConfirmationBlock < p.SafeConfirmationBlock {
		return fmt.Errorf("hard block(%d) < safe block(%d)", p.HardConfirmationBlock, p.SafeConfirmationBlock)
	}
	return nil
}
