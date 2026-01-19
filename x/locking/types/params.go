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
		SlashFractionDoubleSign: math.LegacyNewDecWithPrec(5, 2),
		SlashFractionDowntime:   math.LegacyNewDecWithPrec(2, 2),
		HalvingInterval:         42048000,
		InitialBlockReward:      2378234400000000000,
	}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.MaxValidators > 100 || p.MaxValidators < 1 {
		return fmt.Errorf("invalid MaxValidators: %d", p.MaxValidators)
	}

	if p.MaxMissedPerWindow < 1 || p.SignedBlocksWindow < 1 {
		return fmt.Errorf("zero signed window values")
	}

	if p.MaxMissedPerWindow >= p.SignedBlocksWindow {
		return fmt.Errorf("MaxMissedPerWindow %d >= SignedBlocksWindow %d",
			p.MaxMissedPerWindow, p.SignedBlocksWindow)
	}

	if p.SlashFractionDoubleSign.GTE(math.LegacyNewDec(1)) {
		return fmt.Errorf("SlashFractionDoubleSign too high: %s", p.SlashFractionDoubleSign.String())
	}

	if p.SlashFractionDoubleSign.IsZero() || p.SlashFractionDoubleSign.IsNegative() {
		return fmt.Errorf("SlashFractionDoubleSign too low: %s", p.SlashFractionDoubleSign.String())
	}

	if p.SlashFractionDowntime.GTE(math.LegacyNewDec(1)) {
		return fmt.Errorf("SlashFractionDowntime too high: %s", p.SlashFractionDowntime.String())
	}

	if p.SlashFractionDowntime.IsZero() || p.SlashFractionDowntime.IsNegative() {
		return fmt.Errorf("SlashFractionDowntime too low: %s", p.SlashFractionDowntime.String())
	}

	if p.DowntimeJailDuration < time.Minute {
		return fmt.Errorf("DowntimeJailDuration too low: %s", p.DowntimeJailDuration.String())
	}

	if p.ExitingDuration < p.UnlockDuration {
		return fmt.Errorf("DowntimeJailDuration too low")
	}

	if p.InitialBlockReward < 1 {
		return fmt.Errorf("InitialBlockReward too low")
	}

	if p.HalvingInterval < 1 {
		return fmt.Errorf("HalvingInterval too low")
	}

	return nil
}
