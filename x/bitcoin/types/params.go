package types

import (
	"errors"
	"fmt"
)

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{
		ChainConfig: &ChainConfig{
			NetworkName:          "regtest",
			PubkeyHashAddrPrefix: 0x6f,
			ScriptHashAddrPrefix: 0xc4,
			Bech32Hrp:            "bcrt",
		},
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
	if p.ChainConfig.NetworkName == "" {
		return errors.New("emtpy network name")
	}

	if p.ChainConfig.Bech32Hrp == "" {
		return errors.New("emtpy bech32 hrp")
	}

	// uint8
	if p.ChainConfig.PubkeyHashAddrPrefix > 255 {
		return fmt.Errorf("overflow for PubkeyHashAddrPrefix: %d", p.ChainConfig.PubkeyHashAddrPrefix)
	}

	if p.ChainConfig.ScriptHashAddrPrefix > 255 {
		return fmt.Errorf("overflow for ScriptHashAddrPrefix: %d", p.ChainConfig.ScriptHashAddrPrefix)
	}

	if p.MinDepositAmount < DustTxoutAmount {
		return errors.New("minimal deposit amount can't be less than dust value")
	}

	if len(p.DepositMagicPrefix) != DepositMagicLen {
		return errors.New("invalid DepositMagicPrefix length")
	}

	if p.ConfirmationNumber == 0 {
		return errors.New("confirmation number can't set to zero(mempool txs are not reliable )")
	}
	return nil
}
