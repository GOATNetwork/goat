package types

import (
	"bytes"
	"errors"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
)

func Version0Address(pubkey *relayer.PublicKey, evmAddress []byte, netwk *chaincfg.Params) (btcutil.Address, error) {
	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		script, err := txscript.NewScriptBuilder().
			AddData(evmAddress[:]).
			AddOp(txscript.OP_DROP).
			AddData(v.Secp256K1).
			AddOp(txscript.OP_CHECKSIGVERIFY).Script()
		if err != nil {
			return nil, err
		}
		witnessProg := goatcrypto.SHA256Sum(script)
		return btcutil.NewAddressWitnessScriptHash(witnessProg, netwk)
	case *relayer.PublicKey_Schnorr:
		pubkey, err := schnorr.ParsePubKey(v.Schnorr)
		if err != nil {
			return nil, err
		}

		script, err := txscript.NewScriptBuilder().AddFullData(evmAddress).Script()
		if err != nil {
			return nil, err
		}
		// tweek the pubkey with the fake script
		// so we don't need to use script tree to spend it
		witnessProg := schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(pubkey, script))
		return btcutil.NewAddressTaproot(witnessProg, netwk)
	}
	return nil, errors.New("unknown key type")
}

func ValidateDespositTxOut(pubkey *relayer.PublicKey, evmAddress, txout []byte) (bool, error) {
	if len(txout) != 34 {
		return false, nil
	}

	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		if txout[0] != txscript.OP_0 || txout[1] != txscript.OP_DATA_32 {
			return false, nil
		}
		script, err := txscript.NewScriptBuilder().
			AddData(evmAddress[:]).
			AddOp(txscript.OP_DROP).
			AddData(v.Secp256K1).
			AddOp(txscript.OP_CHECKSIGVERIFY).Script()
		if err != nil {
			return false, err
		}
		witnessProg := goatcrypto.SHA256Sum(script)
		return bytes.Equal(witnessProg, txout[2:]), nil
	case *relayer.PublicKey_Schnorr:
		if txout[0] != txscript.OP_1 || txout[1] != txscript.OP_DATA_32 {
			return false, nil
		}

		pubkey, err := schnorr.ParsePubKey(v.Schnorr)
		if err != nil {
			return false, err
		}

		script, err := txscript.NewScriptBuilder().AddFullData(evmAddress).Script()
		if err != nil {
			return false, err
		}
		witnessProg := schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(pubkey, script))
		return bytes.Equal(witnessProg, txout[2:]), nil
	}
	return false, errors.New("unknown key type")
}
