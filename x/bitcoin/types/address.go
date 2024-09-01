package types

import (
	"bytes"
	"errors"
	"fmt"
	"slices"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
)

func DepositAddressV0(pubkey *relayer.PublicKey, evmAddress []byte, netwk *chaincfg.Params) (btcutil.Address, error) {
	if len(evmAddress) != 20 {
		return nil, fmt.Errorf("invalid evm address")
	}

	if err := pubkey.Validate(); err != nil {
		return nil, err
	}

	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		script, err := txscript.NewScriptBuilder().AddData(evmAddress).AddOp(txscript.OP_DROP).
			AddData(v.Secp256K1).AddOp(txscript.OP_CHECKSIGVERIFY).Script()
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
		// tweek the pubkey with the evm address
		// so we don't need to build a script tree to spend it
		witnessProg := schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(pubkey, evmAddress))
		return btcutil.NewAddressTaproot(witnessProg, netwk)
	}
	return nil, errors.New("unknown key type")
}

func DepositAddressV1(pubkey *relayer.PublicKey, magic, evmAddress []byte, netwk *chaincfg.Params) (btcutil.Address, []byte, error) {
	if len(evmAddress) != 20 {
		return nil, nil, errors.New("invalid evm address")
	}

	if len(magic) != 4 {
		return nil, nil, errors.New("invalid deposit prefix")
	}

	if err := pubkey.Validate(); err != nil {
		return nil, nil, err
	}

	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(v.Secp256K1), netwk)
		if err != nil {
			return nil, nil, err
		}

		script, err := txscript.NewScriptBuilder().
			AddFullData(slices.Concat(magic, evmAddress)).Script()
		if err != nil {
			return nil, nil, err
		}
		return addr, script, nil
	}
	return nil, nil, errors.New("unknown key type")
}

func WithdrawalAddress(address string, netwk *chaincfg.Params) ([]byte, error) {
	addr, err := btcutil.DecodeAddress(address, netwk)
	if err != nil {
		return nil, err
	}

	switch v := addr.(type) {
	case *btcutil.AddressPubKeyHash:
		return append([]byte{0}, v.ScriptAddress()...), nil
	case *btcutil.AddressScriptHash:
		return append([]byte{1}, v.ScriptAddress()...), nil
	case *btcutil.AddressWitnessPubKeyHash:
		return append([]byte{2}, v.ScriptAddress()...), nil
	case *btcutil.AddressWitnessScriptHash:
		return append([]byte{3}, v.ScriptAddress()...), nil
	case *btcutil.AddressTaproot:
		return append([]byte{4}, v.ScriptAddress()...), nil
	}
	return nil, errors.New("unknown address type")
}

func VerifyDespositScriptV0(pubkey *relayer.PublicKey, evmAddress, txout []byte) error {
	if len(txout) != 34 {
		return errors.New("invalid output script")
	}

	if len(evmAddress) != 20 {
		return errors.New("invalid evm address")
	}

	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		if txout[0] != txscript.OP_0 || txout[1] != txscript.OP_DATA_32 {
			return errors.New("invalid p2wsh output")
		}
		script, err := txscript.NewScriptBuilder().AddData(evmAddress).AddOp(txscript.OP_DROP).
			AddData(v.Secp256K1).AddOp(txscript.OP_CHECKSIGVERIFY).Script()
		if err != nil {
			return err
		}
		witnessProg := goatcrypto.SHA256Sum(script)
		if !bytes.Equal(witnessProg, txout[2:]) {
			return errors.New("p2sh script mismatched")
		}
		return nil
	case *relayer.PublicKey_Schnorr:
		if txout[0] != txscript.OP_1 || txout[1] != txscript.OP_DATA_32 {
			return errors.New("invalid p2tr output")
		}

		pubkey, err := schnorr.ParsePubKey(v.Schnorr)
		if err != nil {
			return err
		}
		witnessProg := schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(pubkey, evmAddress))
		if !bytes.Equal(witnessProg, txout[2:]) {
			return errors.New("p2tr script mismatched")
		}
		return nil
	}
	return errors.New("unknown key type")
}

func VerifyDespositScriptV1(pubkey *relayer.PublicKey, magicPrefix, evmAddress, txout0, txout1 []byte) error {
	if len(magicPrefix) != 4 {
		return errors.New("invalid deposit prefix")
	}

	if len(evmAddress) != 20 {
		return errors.New("invalid evm address")
	}

	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		if len(txout0) != 22 {
			return errors.New("invalid output script")
		}

		if txout0[0] != txscript.OP_0 || txout0[1] != txscript.OP_DATA_20 {
			return errors.New("invalid p2wpkh output")
		}

		if !bytes.Equal(btcutil.Hash160(v.Secp256K1), txout0[2:]) {
			return errors.New("p2wpkh script mismatched")
		}

		if len(txout1) != 26 {
			return errors.New("invalid OP_RETURNS script length")
		}

		if txout0[0] != txscript.OP_RETURN || txout0[1] != txscript.OP_DATA_24 {
			return errors.New("invalid OP_RETURNS output")
		}

		if script := slices.Concat(magicPrefix, evmAddress); !bytes.Equal(txout1[2:], script) {
			return fmt.Errorf("OP_RETRURNS mismatched: expected 6a18%x got %x", script, txout1)
		}

		return nil
	}
	return errors.New("unknown key type")
}

func VerifyMerkelProof(txid, root, proof []byte, index uint32) bool {
	if len(txid) != 32 || len(root) != 32 || len(proof)%32 != 0 {
		return false
	}

	current := txid
	for i := 0; i < len(proof)/32; i++ {
		start := i * 32
		end := start + 32
		next := proof[start:end]
		if index&1 == 0 {
			current = goatcrypto.DoubleSHA256Sum(slices.Concat(current, next))
		} else {
			current = goatcrypto.DoubleSHA256Sum(slices.Concat(next, current))
		}
		index >>= 1
	}

	return bytes.Equal(current, root)
}
