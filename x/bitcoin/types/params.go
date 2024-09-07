package types

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
)

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{
		ChainConfig: &ChainConfig{
			NetworkName:          chaincfg.RegressionNetParams.Name,
			PubkeyHashAddrPrefix: uint32(chaincfg.RegressionNetParams.PubKeyHashAddrID),
			ScriptHashAddrPrefix: uint32(chaincfg.RegressionNetParams.ScriptHashAddrID),
			Bech32Hrp:            chaincfg.RegressionNetParams.Bech32HRPSegwit,
		},
		ConfirmationNumber: 1,
		MinDepositAmount:   DustTxoutAmount,
		DepositMagicPrefix: []byte("GTT0"),
	}

	/*
		DepositMagicPrefix
			- regtest/devnet: GTT0
			- testnet: GTV1
			- mainnet: GTV2
	*/
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params.
func (p Params) Validate() error {
	switch p.ChainConfig.NetworkName {
	case chaincfg.MainNetParams.Name, chaincfg.TestNet3Params.Name, chaincfg.SigNetParams.Name, chaincfg.RegressionNetParams.Name:
	default:
		return fmt.Errorf("unknown bitcoin network: %s", p.ChainConfig.NetworkName)
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
