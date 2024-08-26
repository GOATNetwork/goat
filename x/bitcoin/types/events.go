package types

import (
	"encoding/base64"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeNewKey = "new_key"
)

func NewKeyEvent(key []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeNewKey,
		sdktypes.NewAttribute("key", base64.RawStdEncoding.EncodeToString(key)),
	)
}
