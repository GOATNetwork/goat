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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
)

func VerifySystemAddressScript(pubkey *relayer.PublicKey, script []byte) bool {
	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		if len(script) != P2WPKHScriptSize {
			return false
		}
		if script[0] != txscript.OP_0 || script[1] != txscript.OP_DATA_20 {
			return false
		}
		return bytes.Equal(goatcrypto.Hash160Sum(v.Secp256K1), script[2:])
	case *relayer.PublicKey_Schnorr:
		if len(script) != P2TRScriptSize {
			return false
		}
		if script[0] != txscript.OP_1 || script[1] != txscript.OP_DATA_32 {
			return false
		}
		pubkey, err := schnorr.ParsePubKey(v.Schnorr)
		if err != nil {
			return false
		}
		witnessProg := schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(pubkey))
		return bytes.Equal(witnessProg, script[2:])
	}
	return false
}

func DepositAddressV0(pubkey *relayer.PublicKey, evmAddress []byte, netwk *chaincfg.Params) (btcutil.Address, error) {
	if len(evmAddress) != common.AddressLength {
		return nil, fmt.Errorf("invalid evm address")
	}

	if err := pubkey.Validate(); err != nil {
		return nil, err
	}

	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		script, err := txscript.NewScriptBuilder().AddData(evmAddress).AddOp(txscript.OP_DROP).
			AddData(v.Secp256K1).AddOp(txscript.OP_CHECKSIG).Script()
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

func DepositAddressV1(pubkey *relayer.PublicKey, magicPrefix, evmAddress []byte, netwk *chaincfg.Params) (btcutil.Address, []byte, error) {
	if len(evmAddress) != common.AddressLength {
		return nil, nil, errors.New("invalid evm address")
	}

	if len(magicPrefix) != DepositMagicLen {
		return nil, nil, errors.New("invalid deposit prefix")
	}

	if err := pubkey.Validate(); err != nil {
		return nil, nil, err
	}

	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		addr, err := btcutil.NewAddressWitnessPubKeyHash(goatcrypto.Hash160Sum(v.Secp256K1), netwk)
		if err != nil {
			return nil, nil, err
		}

		script, err := txscript.NewScriptBuilder().
			AddOp(txscript.OP_RETURN).AddFullData(slices.Concat(magicPrefix, evmAddress)).Script()
		if err != nil {
			return nil, nil, err
		}
		return addr, script, nil
	}
	return nil, nil, errors.New("unknown key type")
}

func VerifyDespositScriptV0(pubkey *relayer.PublicKey, evmAddress, txout []byte) error {
	if len(evmAddress) != common.AddressLength {
		return errors.New("invalid evm address")
	}

	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		if len(txout) != P2WSHScriptSize {
			return errors.New("invalid ouptut length")
		}

		if txout[0] != txscript.OP_0 || txout[1] != txscript.OP_DATA_32 {
			return errors.New("invalid p2wsh output")
		}
		script, err := txscript.NewScriptBuilder().AddData(evmAddress).AddOp(txscript.OP_DROP).
			AddData(v.Secp256K1).AddOp(txscript.OP_CHECKSIG).Script()
		if err != nil {
			return err
		}
		witnessProg := goatcrypto.SHA256Sum(script)
		if !bytes.Equal(witnessProg, txout[2:]) {
			return errors.New("p2wsh script mismatched")
		}
		return nil
	case *relayer.PublicKey_Schnorr:
		if len(txout) != P2TRScriptSize {
			return errors.New("invalid ouptut length")
		}

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
	if len(magicPrefix) != DepositMagicLen {
		return errors.New("invalid deposit prefix")
	}

	if len(evmAddress) != common.AddressLength {
		return errors.New("invalid evm address")
	}

	switch v := pubkey.GetKey().(type) {
	case *relayer.PublicKey_Secp256K1:
		if len(txout0) != P2WPKHScriptSize {
			return errors.New("invalid output script")
		}

		if txout0[0] != txscript.OP_0 || txout0[1] != txscript.OP_DATA_20 {
			return errors.New("invalid p2wpkh output")
		}

		if !bytes.Equal(goatcrypto.Hash160Sum(v.Secp256K1), txout0[2:]) {
			return errors.New("p2wpkh script mismatched")
		}

		if len(txout1) != DepositV1TxoutSize {
			return errors.New("invalid OP_RETURN script length")
		}

		if txout1[0] != txscript.OP_RETURN || txout1[1] != txscript.OP_DATA_24 {
			return errors.New("invalid OP_RETURN output")
		}

		if script := slices.Concat(magicPrefix, evmAddress); !bytes.Equal(txout1[2:], script) {
			return fmt.Errorf("OP_RETURN mismatched: expected 6a18%x got %x", script, txout1)
		}

		return nil
	}
	return errors.New("unknown key type")
}

func DecodeEthAddress(address string) ([]byte, error) {
	data, err := hexutil.Decode(address)
	if err != nil {
		return nil, err
	}
	if len(data) != common.AddressLength {
		return nil, errors.New("invalid address length")
	}
	return data, nil
}

// DecodeBtcAddress verifies if the address is valid and returns its payment script for later verification
func DecodeBtcAddress(address string, netwk *chaincfg.Params) ([]byte, error) {
	addr, err := btcutil.DecodeAddress(address, netwk)
	if err != nil {
		return nil, err
	}

	if !addr.IsForNet(netwk) {
		return nil, fmt.Errorf("not a %s network address", netwk.Name)
	}

	// the deprecated address, it takes more fee than others
	if _, ok := addr.(*btcutil.AddressPubKey); ok {
		return nil, errors.New("deprecated p2pk address")
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	return script, nil
}
