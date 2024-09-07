package types

import (
	"fmt"
	"time"
)

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{ElectingPeriod: time.Minute * 10}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.ElectingPeriod == 0 {
		return fmt.Errorf("invalid electing period")
	}
	return nil
}
