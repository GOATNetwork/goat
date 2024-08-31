package types

import (
	"errors"
	"slices"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

const (
	secp256k1Type byte = iota
	schnoorType
)

func (p *PublicKey) Validate() error {
	if p == nil {
		return errors.New("empty public key")
	}

	switch t := p.Key.(type) {
	case *PublicKey_Secp256K1:
		if len(t.Secp256K1) != btcec.PubKeyBytesLenCompressed {
			return errors.New("invalid secp256k1 key")
		}
		if prefix := t.Secp256K1[0]; prefix != 2 && prefix != 3 {
			return errors.New("invalid compressed secp256k1 prefix")
		}
	case *PublicKey_Schnorr:
		if len(t.Schnorr) != schnorr.PubKeyBytesLen {
			return errors.New("invalid schnoor key")
		}
	default:
		return errors.New("unknown pubkey type")
	}
	return nil
}

func (p *PublicKey) VerifySign(msg, sig []byte) bool {
	// note: msg is 32 bytes, sig is 64 bytes
	if len(msg) != 32 || len(sig) != schnorr.SignatureSize {
		return false
	}

	switch t := p.Key.(type) {
	case *PublicKey_Secp256K1:
		if len(t.Secp256K1) != btcec.PubKeyBytesLenCompressed {
			return false
		}
		return ethcrypto.VerifySignature(t.Secp256K1, msg, sig)
	case *PublicKey_Schnorr:
		pub, err := schnorr.ParsePubKey(t.Schnorr)
		if err != nil {
			return false
		}
		signature, err := schnorr.ParseSignature(sig)
		if err != nil {
			return false
		}
		return signature.Verify(msg, pub)
	}
	return false
}

func EncodePublicKey(v *PublicKey) []byte {
	switch t := v.Key.(type) {
	case *PublicKey_Secp256K1:
		res := make([]byte, 0, 1+btcec.PubKeyBytesLenCompressed)
		res = append(res, byte(secp256k1Type))
		res = append(res, t.Secp256K1...)
		return res
	case *PublicKey_Schnorr:
		res := make([]byte, 0, 1+schnorr.PubKeyBytesLen)
		res = append(res, byte(schnoorType))
		res = append(res, t.Schnorr...)
		return res
	}
	return nil
}

func DecodePublicKey(raw []byte) *PublicKey {
	switch len(raw) - 1 {
	case btcec.PubKeyBytesLenCompressed:
		if raw[0] == secp256k1Type {
			return &PublicKey{&PublicKey_Secp256K1{slices.Clone(raw[1:])}}
		}
	case schnorr.PubKeyBytesLen:
		if raw[0] == schnoorType {
			return &PublicKey{&PublicKey_Schnorr{slices.Clone(raw[1:])}}
		}
	}
	return nil
}
