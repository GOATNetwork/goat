package types

import (
	"encoding/hex"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeDowngraded = "validator_downgraded"
	EventTypeTombstoned = "validator_tombstoned"
)

func ValidatorDowngradedEvent(validator sdktypes.ConsAddress) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeDowngraded,
		sdktypes.NewAttribute("validator", hex.EncodeToString(validator.Bytes())),
	)
}

func ValidatorTombstonedEvent(validator sdktypes.ConsAddress) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeTombstoned,
		sdktypes.NewAttribute("validator", hex.EncodeToString(validator.Bytes())),
	)
}
