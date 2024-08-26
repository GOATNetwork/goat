package types

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeNewKey     = "new_key"
	EventTypeNewDeposit = "new_deposit"
)

func NewKeyEvent(key []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeNewKey,
		sdktypes.NewAttribute("key", base64.RawStdEncoding.EncodeToString(key)),
	)
}

func NewDepositEvent(deposit *ExecuableDeposit) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeNewDeposit,
		sdktypes.NewAttribute("txid", hex.EncodeToString(deposit.Txid)),
		sdktypes.NewAttribute("txout", strconv.FormatUint(uint64(deposit.Txout), 10)),
		sdktypes.NewAttribute("address", hex.EncodeToString(deposit.Address)),
		sdktypes.NewAttribute("amount", deposit.Amount.String()),
	)
}
