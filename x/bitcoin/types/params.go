package types

import (
	"errors"
	"fmt"
)

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{
		NetworkName:        "regtest",
		ConfirmationNumber: 1,
		MinDepositAmount:   1e4,
		DepositMagicPrefix: []byte("GTT0"),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params.
func (p Params) Validate() error {
	network := BitcoinNetworks[p.NetworkName]
	if network == nil {
		return fmt.Errorf("network %s not found", p.NetworkName)
	}

	if p.MinDepositAmount < DustTxoutAmount {
		return errors.New("minimal deposit amount can't be less than dust value")
	}

	if len(p.DepositMagicPrefix) != DepositMagicLen {
		return errors.New("invalid DepositMagicPrefix length")
	}

	if p.ConfirmationNumber == 0 {
		return errors.New("confirmation number can't set to zero(mempool txs are not reliable)")
	}

	if p.DepositTaxRate > 0 {
		if p.MaxDepositTax == 0 || p.DepositTaxRate > 1e4 {
			return fmt.Errorf("invalid deposit tax: DepositTaxRate(%d) MaxDepositTax(%d)",
				p.DepositTaxRate, p.MaxDepositTax)
		}
		if p.MaxDepositTax > 1e8 {
			return fmt.Errorf("MaxDepositTax is too large: %d", p.MaxDepositTax)
		}
	} else if p.MaxDepositTax != 0 {
		return fmt.Errorf("invalid deposit tax: DepositTaxRate(%d) MaxDepositTax(%d)",
			p.DepositTaxRate, p.MaxDepositTax)
	}

	return nil
}
