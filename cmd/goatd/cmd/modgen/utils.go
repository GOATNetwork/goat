package modgen

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

func DecodeHexOrBase64String(str string) ([]byte, error) {
	pubkeyRaw, err := hex.DecodeString(str)
	if err != nil {
		pubkeyRaw, err = base64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil, fmt.Errorf("pubkey %s doesn't use base64 or hex encoding", str)
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
