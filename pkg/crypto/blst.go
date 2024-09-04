package crypto

import (
	"crypto/rand"
	"errors"

	blst "github.com/supranational/blst/bindings/go"
)

var blsMode = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

type PrivateKey = blst.SecretKey
type PublicKey = blst.P2Affine
type Signature = blst.P1Affine

type AggregatePublicKey = blst.P2Aggregate
type AggregateSignature = blst.P1Aggregate

const (
	PubkeyLength    = blst.BLST_P2_COMPRESS_BYTES
	SignatureLength = blst.BLST_P1_COMPRESS_BYTES
)

var (
	ErrorAggregation = errors.New("crypto: failed to aggregate bls signatures")
)

func GenPrivKey() *PrivateKey {
	var raw [32]byte
	_, _ = rand.Read(raw[:])
	return blst.KeyGenV3(raw[:])
}

func AggregateVerify(pks [][]byte, msg, sig []byte) bool {
	if len(pks) == 0 {
		return false
	}

	signature := new(Signature).Uncompress(sig)
	if signature == nil {
		return false
	}

	pubkeys := make([]*PublicKey, 0, len(pks))
	for _, v := range pks {
		pk := new(PublicKey).Uncompress(v)
		if pk == nil {
			return false
		}
		pubkeys = append(pubkeys, pk)
	}
	return signature.FastAggregateVerify(true, pubkeys, msg, blsMode)
}

func AggregateSignatures(sigs [][]byte) ([]byte, error) {
	signature := new(AggregateSignature)
	if !signature.AggregateCompressed(sigs, true) {
		return nil, ErrorAggregation
	}
	return signature.ToAffine().Compress(), nil
}

func Verify(pk *PublicKey, msg, sig []byte) bool {
	signature := new(Signature).Uncompress(sig)
	if signature == nil {
		return false
	}
	return signature.Verify(true, pk, false, msg, blsMode)
}

func Sign(sk *PrivateKey, msg []byte) []byte {
	return new(Signature).Sign(sk, msg, blsMode).Compress()
}
