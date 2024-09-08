package types

import (
	"encoding/hex"
	"slices"
)

func BtcTxid(hash []byte) string {
	txid := slices.Clone(hash)
	slices.Reverse(txid)
	return hex.EncodeToString(txid)
}
