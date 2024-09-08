package types

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
)

const (
	EventTypeNewKey              = "new_key"
	EventTypeNewDeposit          = "new_deposit"
	EventTypeNewBlockHash        = "new_block_hash"
	EventTypeNewWithdrawal       = "new_withdrawal"
	EventTypeApproveCancellation = "approve_cancellation_withdrawal"
	EventTypeFinalizeWithdrawal  = "finalize_withdrawal"
)

func NewKeyEvent(key *relayertypes.PublicKey) sdktypes.Event {
	var typ, raw string
	switch v := key.Key.(type) {
	case *relayertypes.PublicKey_Secp256K1:
		typ = "secp256k1"
		raw = base64.StdEncoding.EncodeToString(v.Secp256K1)
	case *relayertypes.PublicKey_Schnorr:
		typ = "schnorr"
		raw = base64.StdEncoding.EncodeToString(v.Schnorr)
	default:
		panic("unknown relayer public key type")
	}

	return sdktypes.NewEvent(
		EventTypeNewKey,
		sdktypes.NewAttribute("type", typ),
		sdktypes.NewAttribute("key", raw),
	)
}

func NewDepositEvent(deposit *DepositExecReceipt) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeNewDeposit,
		sdktypes.NewAttribute("txid", BtcTxid(deposit.Txid)), // we must use big endian
		sdktypes.NewAttribute("txout", strconv.FormatUint(uint64(deposit.Txout), 10)),
		sdktypes.NewAttribute("address", hex.EncodeToString(deposit.Address)),
		sdktypes.NewAttribute("amount", strconv.FormatUint(deposit.Amount, 10)),
	)
}

func NewBlockHashEvent(height uint64, hash []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeNewBlockHash,
		sdktypes.NewAttribute("height", strconv.FormatUint(height, 10)),
		sdktypes.NewAttribute("hash", BtcTxid(hash)), // we must use big endian
	)
}

func NewWithdrawalEvent(hash []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeNewWithdrawal,
		sdktypes.NewAttribute("txid", BtcTxid(hash)), // we must use big endian
	)
}

func FinalizeWithdrawalEvent(hash []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeFinalizeWithdrawal,
		sdktypes.NewAttribute("txid", BtcTxid(hash)), // we must use big endian
	)
}

func ApproveCancellationEvent(id uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeApproveCancellation,
		sdktypes.NewAttribute("id", strconv.FormatUint(id, 10)),
	)
}
