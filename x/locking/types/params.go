package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
)

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		UnlockDuration:          time.Hour * 24 * 7,
		ExitingDuration:         time.Hour * 24 * 21,
		DowntimeJailDuration:    time.Hour * 3,
		MaxValidators:           20,
		SignedBlocksWindow:      1200,
		MaxMissedPerWindow:      200,
		SlashFractionDoubleSign: math.LegacyNewDec(5).QuoInt64(100),
		SlashFractionDowntime:   math.LegacyNewDec(2).QuoInt64(100),
		HalvingInterval:         42048000,
		InitialBlockReward:      2378234400000000000,
	}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.MaxValidators > 100 {
		return fmt.Errorf("max validator should be less than 100")
	}
	if p.MaxMissedPerWindow == 0 || p.SignedBlocksWindow == 0 {
		return fmt.Errorf("zero signed window values")
	}

	if p.MaxMissedPerWindow > p.SignedBlocksWindow {
		return fmt.Errorf("MaxMissedPerWindow %d > SignedBlocksWindow %d",
			p.MaxMissedPerWindow, p.SignedBlocksWindow)
	}

	return nil
}
