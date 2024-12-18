package types

import (
	"errors"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	// regtest geneis hash is always the same
	// https://github.com/bitcoin/bitcoin/blob/v27.0/src/kernel/chainparams.cpp#L404
	geneis, err := chainhash.NewHashFromStr("0f9188f13cb7b2c71f2a335e3a4fc328bf5beb436012afca590b1a11466e2206")
	if err != nil {
		panic(err)
	}

	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:      DefaultParams(),
		BlockTip:    0,
		BlockHashes: [][]byte{geneis[:]},
		EthTxNonce:  0,
		EthTxQueue:  EthTxQueue{BlockNumber: 0},
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

	if len(gs.BlockHashes) == 0 {
		return errors.New("no block hash provided in the genesis state")
	}

	return nil
}
