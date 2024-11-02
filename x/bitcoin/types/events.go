package types

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
)

const (
	EventTypeNewKey           = "new_key"
	EventTypeNewDeposit       = "new_deposit"
	EventTypeNewBlockHash     = "new_block_hash"
	EventTypeNewConsolidation = "new_consolidation"

	EventTypeWithdrawalInit           = "withdrawal_init"
	EventTypeWithdrawalProcessing     = "withdrawal_processing"
	EventTypeWithdrawalFinalized      = "withdrawal_finalized"
	EventTypeWithdrawalUserReplace    = "withdrawal_user_rbf"
	EventTypeWithdrawalUserCancel     = "withdrawal_user_cancel"
	EventTypeWithdrawalRelayerReplace = "withdrawal_relayer_rbf"
	EventTypeWithdrawalRelayerCancel  = "withdrawal_relayer_cancel"
)

const (
	Secp256K1Name = "secp256k1"
	SchnorrName   = "schnorr"
)

func NewKeyEvent(key *relayertypes.PublicKey) sdktypes.Event {
	var typ, raw string
	switch v := key.Key.(type) {
	case *relayertypes.PublicKey_Secp256K1:
		typ = Secp256K1Name
		raw = base64.StdEncoding.EncodeToString(v.Secp256K1)
	case *relayertypes.PublicKey_Schnorr:
		typ = SchnorrName
		raw = base64.StdEncoding.EncodeToString(v.Schnorr)
	default:
		typ = "unknown"
	}

	return sdktypes.NewEvent(
		EventTypeNewKey,
		sdktypes.NewAttribute("type", typ),
		sdktypes.NewAttribute("key", raw),
	)
}

func NewBlockHashEvent(height uint64, hash []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeNewBlockHash,
		sdktypes.NewAttribute("height", strconv.FormatUint(height, 10)),
		sdktypes.NewAttribute("hash", BtcTxid(hash)), // we must use big endian
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

func NewWithdrawalInitEvent(address string, id, txPrice, amount uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeWithdrawalInit,
		sdktypes.NewAttribute("id", strconv.FormatUint(id, 10)),
		sdktypes.NewAttribute("address", address),
		sdktypes.NewAttribute("tx_price", strconv.FormatUint(txPrice, 10)),
		sdktypes.NewAttribute("amount", strconv.FormatUint(amount, 10)),
	)
}

func NewWithdrawalUserReplaceEvent(id, txPrice uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeWithdrawalUserReplace,
		sdktypes.NewAttribute("id", strconv.FormatUint(id, 10)),
		sdktypes.NewAttribute("tx_price", strconv.FormatUint(txPrice, 10)),
	)
}

func NewWithdrawalUserCancelEvent(id uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeWithdrawalUserCancel,
		sdktypes.NewAttribute("id", strconv.FormatUint(id, 10)),
	)
}

func NewWithdrawalProcessingEvent(id uint64, hash []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeWithdrawalProcessing,
		sdktypes.NewAttribute("pid", strconv.FormatUint(id, 10)), // the process id, used for replacing and Fianlizing
		sdktypes.NewAttribute("txid", BtcTxid(hash)),             // we must use big endian
	)
}

func NewWithdrawalRelayerReplaceEvent(id uint64, hash []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeWithdrawalRelayerReplace,
		sdktypes.NewAttribute("pid", strconv.FormatUint(id, 10)),
		sdktypes.NewAttribute("txid", BtcTxid(hash)), // we must use big endian
	)
}

func NewWithdrawalFinalizedEvent(id uint64, hash []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeWithdrawalFinalized,
		sdktypes.NewAttribute("pid", strconv.FormatUint(id, 10)),
		sdktypes.NewAttribute("txid", BtcTxid(hash)), // we must use big endian
	)
}

func NewWithdrawalRelayerCancelEvent(id uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeWithdrawalRelayerCancel,
		sdktypes.NewAttribute("id", strconv.FormatUint(id, 10)),
	)
}

func NewConsolidationEvent(hash []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeNewConsolidation,
		sdktypes.NewAttribute("txid", BtcTxid(hash)), // we must use big endian
	)
}
