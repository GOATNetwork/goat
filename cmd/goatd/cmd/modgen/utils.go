package modgen

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

var BitcoinNetworks = map[string]*chaincfg.Params{
	chaincfg.MainNetParams.Name:       &chaincfg.MainNetParams,
	chaincfg.TestNet3Params.Name:      &chaincfg.TestNet3Params,
	chaincfg.SigNetParams.Name:        &chaincfg.SigNetParams,
	chaincfg.RegressionNetParams.Name: &chaincfg.RegressionNetParams,
}

func DecodeHexOrBase64String(str string) ([]byte, error) {
	pubkeyRaw, err := hex.DecodeString(str)
	if err != nil {
		pubkeyRaw, err = base64.StdEncoding.DecodeString(str)
		if err != nil {
			pubkeyRaw, err = hex.DecodeString(strings.TrimPrefix(str, "0x"))
			if err != nil {
				return nil, fmt.Errorf("pubkey %s doesn't use base64 or hex encoding", str)
			}
		}
	}
	return pubkeyRaw, nil
}

func IsValidSecp256Pubkey(key []byte) error {
	if len(key) != secp256k1.PubKeySize {
		return errors.New("invalid secp256k1 pubkey length")
	}
	if key[0] != 2 && key[0] != 3 {
		return errors.New("invalid secp256k1 pubkey prefix")
	}
	return nil
}
