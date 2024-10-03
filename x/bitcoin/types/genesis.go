package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	// regtest geneis hash is always the same
	// https://github.com/bitcoin/bitcoin/blob/v27.0/src/kernel/chainparams.cpp#L404
	geneis, err := hex.DecodeString("0f9188f13cb7b2c71f2a335e3a4fc328bf5beb436012afca590b1a11466e2206")
	if err != nil {
		panic(err)
	}
	// conver it to little endian format
	slices.Reverse(geneis)

	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:      DefaultParams(),
		BlockTip:    0,
		BlockHashes: [][]byte{geneis},
		EthTxNonce:  0,
		Queue: &ExecuableQueue{
			BlockNumber: 0,
		},
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
		return errors.New("No block hash provided in the genesis state")
	}

	if gs.Queue == nil {
		return fmt.Errorf("execution queue is empty")
	}

	return nil
}
