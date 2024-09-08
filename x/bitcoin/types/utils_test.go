package types

import (
	"encoding/hex"
	"testing"
)

func TestBtcTxid(t *testing.T) {
	hash := "e540f395191e788fb1bf9dc4772da301a6e099ee962ad95f1afc811487fbc6e1"

	raw, err := hex.DecodeString(hash)
	if err != nil {
		t.Fatal(err)
	}

	txid := "e1c6fb871481fc1a5fd92a96ee99e0a601a32d77c49dbfb18f781e1995f340e5"

	if got := BtcTxid(raw); got != txid {
		t.Errorf("BtcTxid() = %v, want %v", got, txid)
	}

	if got := hex.EncodeToString(raw); got != hash {
		t.Errorf("origin hash modified by BtcTxid() = %v, want %v", got, hash)
	}
}
