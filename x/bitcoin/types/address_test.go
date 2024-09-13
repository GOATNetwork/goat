package types

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
)

var BitcoinNetworks = map[string]*chaincfg.Params{
	chaincfg.MainNetParams.Name:       &chaincfg.MainNetParams,
	chaincfg.TestNet3Params.Name:      &chaincfg.TestNet3Params,
	chaincfg.SigNetParams.Name:        &chaincfg.SigNetParams,
	chaincfg.RegressionNetParams.Name: &chaincfg.RegressionNetParams,
}

var DepositMagicPreifxs = map[string]string{
	chaincfg.MainNetParams.Name:       "GTV2",
	chaincfg.TestNet3Params.Name:      "GTV1",
	chaincfg.SigNetParams.Name:        "GTV1",
	chaincfg.RegressionNetParams.Name: "GTT0",
}

func TestDecodeEthAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "invalid",
			args:    args{"invalid"},
			wantErr: true,
		},
		{
			name:    "invalid length",
			args:    args{"0x00"},
			wantErr: true,
		},
		{
			name:    "valid",
			args:    args{"0xBc12C40A3675a1289ab8F286a4B7FdAfBCf4F8e3"},
			want:    common.HexToAddress("0xBc12C40A3675a1289ab8F286a4B7FdAfBCf4F8e3").Bytes(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeEthAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeEthAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeEthAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepositAddressV1(t *testing.T) {
	type DepositTest struct {
		Network    string
		Type       string
		Script     hexutil.Bytes
		OpReturn   hexutil.Bytes
		Prvkey     hexutil.Bytes
		Pubkey     hexutil.Bytes
		EthAddress common.Address
		BtcAddress string
	}

	file, err := os.ReadFile("testdata/deposit-v1.json")
	if err != nil {
		t.Errorf("failed to read file: %s", err)
		return
	}

	var tests []DepositTest
	if err := json.Unmarshal(file, &tests); err != nil {
		t.Errorf("failed to decode test: %s", err)
		return
	}

	t.Parallel()

	for idx, item := range tests {
		idx := idx
		item := item
		t.Run(fmt.Sprintf("idx-%d", idx), func(t *testing.T) {
			network, ok := BitcoinNetworks[item.Network]
			if !ok {
				t.Errorf("network not found %s", item.Network)
				return
			}

			preifix, ok := DepositMagicPreifxs[item.Network]
			if !ok {
				t.Errorf("prefix not found %s", item.Network)
				return
			}

			pubkey := &relayer.PublicKey{Key: &relayer.PublicKey_Secp256K1{Secp256K1: item.Pubkey}}
			if err := VerifyDespositScriptV1(pubkey, []byte(preifix), item.EthAddress[:], item.Script, item.OpReturn[:]); err != nil {
				t.Errorf("VerifyDespositScriptV1() error = %v", err)
			}

			got, opReturns, err := DepositAddressV1(pubkey, []byte(preifix), item.EthAddress[:], network)
			if err != nil {
				t.Errorf("DepositAddressV1() error = %v", err)
				return
			}

			if addr := got.EncodeAddress(); addr != item.BtcAddress {
				t.Errorf("DepositAddressV1() want = %v got = %v", item.BtcAddress, addr)
			}

			if !bytes.Equal(opReturns[:], item.OpReturn) {
				t.Errorf("DepositAddressV1() opReturns: want = %x got = %x", item.OpReturn, opReturns)
			}
		})
	}
}

func TestDepositAddressV0(t *testing.T) {
	type DepositTest struct {
		Network    string
		Type       string
		Script     hexutil.Bytes
		Prvkey     hexutil.Bytes
		Pubkey     hexutil.Bytes
		EthAddress common.Address
		BtcAddress string
	}

	file, err := os.ReadFile("testdata/deposit-v0.json")
	if err != nil {
		t.Errorf("failed to read file: %s", err)
		return
	}

	var tests []DepositTest
	if err := json.Unmarshal(file, &tests); err != nil {
		t.Errorf("failed to decode test: %s", err)
		return
	}

	t.Parallel()

	for idx, item := range tests {
		idx := idx
		item := item
		t.Run(fmt.Sprintf("address-%d", idx), func(t *testing.T) {
			network, ok := BitcoinNetworks[item.Network]
			if !ok {
				t.Errorf("DepositAddressV0() network not found %s", item.Network)
				return
			}

			pubkey := &relayer.PublicKey{}
			if item.Type == "secp256k1" {
				pubkey.Key = &relayer.PublicKey_Secp256K1{Secp256K1: item.Pubkey}
			} else if item.Type == "schnorr" {
				pubkey.Key = &relayer.PublicKey_Schnorr{Schnorr: item.Pubkey}
			}

			if err := VerifyDespositScriptV0(pubkey, item.EthAddress[:], item.Script); err != nil {
				t.Errorf("VerifyDespositScriptV0() error = %v", err)
			}

			got, err := DepositAddressV0(pubkey, item.EthAddress[:], network)
			if err != nil {
				t.Errorf("DepositAddressV0() error = %v", err)
				return
			}

			if addr := got.EncodeAddress(); addr != item.BtcAddress {
				t.Errorf("DepositAddressV0() want = %v got = %v", item.BtcAddress, addr)
			}
		})
	}
}

func TestVerifySystemAddressScript(t *testing.T) {
	type DepositTest struct {
		Network    string
		Type       string
		Script     hexutil.Bytes
		Prvkey     hexutil.Bytes
		Pubkey     hexutil.Bytes
		EthAddress common.Address
		BtcAddress string
	}

	file, err := os.ReadFile("testdata/system-address.json")
	if err != nil {
		t.Errorf("failed to read file: %s", err)
		return
	}

	var tests []DepositTest
	if err := json.Unmarshal(file, &tests); err != nil {
		t.Errorf("failed to decode test: %s", err)
		return
	}

	t.Parallel()

	for idx, item := range tests {
		idx := idx
		item := item
		t.Run(fmt.Sprintf("address-%d", idx), func(t *testing.T) {
			pubkey := &relayer.PublicKey{}
			if item.Type == "secp256k1" {
				pubkey.Key = &relayer.PublicKey_Secp256K1{Secp256K1: item.Pubkey}
			} else if item.Type == "schnorr" {
				pubkey.Key = &relayer.PublicKey_Schnorr{Schnorr: item.Pubkey}
			}

			if !VerifySystemAddressScript(pubkey, item.Script) {
				t.Errorf("VerifySystemAddressScript() not true")
			}

			if VerifySystemAddressScript(nil, nil) {
				t.Errorf("VerifySystemAddressScript() not false")
			}
		})
	}
}
