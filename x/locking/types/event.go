package types

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeDowngraded = "validator_downgraded"
	EventTypeTombstoned = "validator_tombstoned"
	EventTypeRotated    = "validator_rotated"
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

// ValidatorRotatedEvent reports an accepted consensus key rotation. It fires
// when the rotation is staged, one block before CometBFT sees the change.
func ValidatorRotatedEvent(validator, consAddr sdktypes.ConsAddress) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeRotated,
		sdktypes.NewAttribute("validator", ValidatorName(validator.Bytes())),
		sdktypes.NewAttribute("cons_address", ValidatorName(consAddr.Bytes())),
	)
}
