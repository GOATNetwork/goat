package types

import (
	"bytes"
	"errors"
	"fmt"
	"slices"

	cmtcrypto "github.com/cometbft/cometbft/crypto"
	cmmldsa65 "github.com/cometbft/cometbft/crypto/mldsa65"
	cmsecp256k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	tmcrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mldsa65"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// Key types accepted for a validator consensus key. Validator.pubkey is a bare
// bytes field with no type tag of its own, so the tag lives in
// Validator.key_type and the empty string has to keep meaning secp256k1: every
// record written before rotation existed carries no tag at all.
const (
	KeyTypeSecp256k1 = cmsecp256k1.KeyType
	KeyTypeMlDsa65   = cmmldsa65.KeyType
)

// KeyTypeName normalises a key type, mapping the empty tag of a pre-rotation
// record onto secp256k1.
func KeyTypeName(keyType string) string {
	if keyType == "" {
		return KeyTypeSecp256k1
	}
	return keyType
}

// ParsePubkey turns a stored consensus public key back into a typed key.
func ParsePubkey(keyType string, pubkey []byte) (cryptotypes.PubKey, error) {
	switch KeyTypeName(keyType) {
	case KeyTypeSecp256k1:
		if len(pubkey) != secp256k1.PubKeySize {
			return nil, fmt.Errorf("invalid secp256k1 pubkey length %d", len(pubkey))
		}
		if first := pubkey[0]; first != 2 && first != 3 {
			return nil, fmt.Errorf("invalid secp256k1 pubkey prefix %d", first)
		}
		return &secp256k1.PubKey{Key: slices.Clone(pubkey)}, nil
	case KeyTypeMlDsa65:
		if len(pubkey) != cmmldsa65.PubKeySize {
			return nil, fmt.Errorf("invalid ml_dsa_65 pubkey length %d", len(pubkey))
		}
		return &mldsa65.PubKey{Key: slices.Clone(pubkey)}, nil
	default:
		return nil, fmt.Errorf("unsupported consensus key type %q", keyType)
	}
}

// ConsAddress derives the consensus address a key would be known by.
func ConsAddress(keyType string, pubkey []byte) (sdktypes.ConsAddress, error) {
	key, err := ParsePubkey(keyType, pubkey)
	if err != nil {
		return nil, err
	}
	return sdktypes.ConsAddress(key.Address()), nil
}

// CMPubkey renders the validator's current consensus key in the form CometBFT
// expects in a ValidatorUpdate.
func (v *Validator) CMPubkey() tmcrypto.PublicKey {
	pubkey, err := cmPubkey(v.KeyType, v.Pubkey)
	if err != nil {
		// unreachable: nothing reaches the store without going through
		// ParsePubkey first
		panic(err)
	}
	return pubkey
}

// PrevCMPubkey renders the key being rotated away from. It is only valid while
// a rotation is in flight.
func (v *Validator) PrevCMPubkey() (tmcrypto.PublicKey, error) {
	if len(v.PrevPubkey) == 0 {
		return tmcrypto.PublicKey{}, errors.New("validator is not rotating")
	}
	return cmPubkey(v.PrevKeyType, v.PrevPubkey)
}

func cmPubkey(keyType string, pubkey []byte) (tmcrypto.PublicKey, error) {
	switch KeyTypeName(keyType) {
	case KeyTypeSecp256k1:
		return tmcrypto.PublicKey{
			Sum: &tmcrypto.PublicKey_Secp256K1{Secp256K1: slices.Clone(pubkey)},
		}, nil
	case KeyTypeMlDsa65:
		return tmcrypto.PublicKey{
			Sum: &tmcrypto.PublicKey_Mldsa65{Mldsa65: slices.Clone(pubkey)},
		}, nil
	default:
		return tmcrypto.PublicKey{}, fmt.Errorf("unsupported consensus key type %q", keyType)
	}
}

// IsRotating reports whether a rotation has been accepted but not yet applied.
func (v *Validator) IsRotating() bool {
	return v.RotationApplyHeight != 0
}

// CMTPubkey renders a consensus key as a CometBFT crypto.PubKey, which is what
// genesis export wants; CMPubkey renders the protobuf form used in a
// ValidatorUpdate.
func CMTPubkey(keyType string, pubkey []byte) (cmtcrypto.PubKey, error) {
	switch KeyTypeName(keyType) {
	case KeyTypeSecp256k1:
		if len(pubkey) != cmsecp256k1.PubKeySize {
			return nil, fmt.Errorf("invalid secp256k1 pubkey length %d", len(pubkey))
		}
		return cmsecp256k1.PubKey(slices.Clone(pubkey)), nil
	case KeyTypeMlDsa65:
		return cmmldsa65.NewPubKeyFromBytes(pubkey)
	default:
		return nil, fmt.Errorf("unsupported consensus key type %q", keyType)
	}
}

// Key type tags as they travel on the wire in a RotateRequest. The execution
// layer has one byte to spend and no notion of these names.
const (
	WireKeyTypeSecp256k1 uint8 = 0
	WireKeyTypeMlDsa65   uint8 = 1
)

// KeyTypeFromRequest maps the wire tag onto the stored name.
func KeyTypeFromRequest(wire uint8) (string, error) {
	switch wire {
	case WireKeyTypeSecp256k1:
		return KeyTypeSecp256k1, nil
	case WireKeyTypeMlDsa65:
		return KeyTypeMlDsa65, nil
	default:
		return "", fmt.Errorf("unknown consensus key type %d", wire)
	}
}

// SameConsensusKey reports whether the key given is the one the validator
// already uses.
func SameConsensusKey(v *Validator, keyType string, pubkey []byte) bool {
	return KeyTypeName(v.KeyType) == KeyTypeName(keyType) && bytes.Equal(v.Pubkey, pubkey)
}

// RotationProofMessage is what a validator signs with the key it wants to move
// to. It names the chain, so a proof cannot be replayed onto another network,
// and the validator id, so it cannot be replayed onto another validator.
func RotationProofMessage(chainID string, validator sdktypes.ConsAddress, keyType string, pubkey []byte) []byte {
	msg := make([]byte, 0, len(rotationProofDomain)+len(chainID)+len(validator)+len(keyType)+len(pubkey)+5)
	msg = append(msg, rotationProofDomain...)
	msg = append(msg, byte(len(chainID)))
	msg = append(msg, chainID...)
	msg = append(msg, byte(len(validator)))
	msg = append(msg, validator...)
	msg = append(msg, byte(len(keyType)))
	msg = append(msg, keyType...)
	msg = append(msg, pubkey...)
	return msg
}

const rotationProofDomain = "goat-rotate-v1"

// VerifyRotationProof checks that whoever asked for the rotation holds the
// private half of the key being rotated to.
//
// create does the equivalent on the execution layer with ECDSA.recover. That
// stops working the moment the target key is not secp256k1, because the EVM has
// no precompile for anything else, which is why this check lives here.
func VerifyRotationProof(chainID string, validator sdktypes.ConsAddress, keyType string, pubkey, proof []byte) error {
	key, err := ParsePubkey(keyType, pubkey)
	if err != nil {
		return err
	}
	if !key.VerifySignature(RotationProofMessage(chainID, validator, keyType, pubkey), proof) {
		return errors.New("signature does not verify under the new key")
	}
	return nil
}
