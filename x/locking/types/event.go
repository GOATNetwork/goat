package types

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeDowngraded = "validator_downgraded"
	EventTypeTombstoned = "validator_tombstoned"
)

func ValidatorDowngradedEvent(validator sdktypes.ConsAddress) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeDowngraded,
		sdktypes.NewAttribute("validator", ValidatorName(validator.Bytes())),
	)
}

func ValidatorTombstonedEvent(validator sdktypes.ConsAddress) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeTombstoned,
		sdktypes.NewAttribute("validator", ValidatorName(validator.Bytes())),
	)
}
