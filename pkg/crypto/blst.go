package crypto

import blst "github.com/supranational/blst/bindings/go"

var BLSMode = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

type PrivateKey = blst.SecretKey
type PublicKey = blst.P2Affine
type Signature = blst.P1Affine

type AggregatePublicKey = blst.P2Aggregate
type AggregateSignature = blst.P1Aggregate

func AggregateVerify(pks [][]byte, msg, sig []byte) bool {
	dsig := new(Signature).Uncompress(sig)
	if dsig == nil {
		return false
	}

	dpks := new(AggregatePublicKey)
	for _, v := range pks {
		pk := new(PublicKey).Uncompress(v)
		if pk == nil {
			return false
		}
		if !dpks.Add(pk, true) {
			return false
		}
	}
	return dsig.Verify(true, dpks.ToAffine(), true, msg, BLSMode)
}

func Verify(pk, msg, sig []byte) bool {
	signature := new(Signature).Uncompress(sig)
	if signature == nil {
		return false
	}
	pubkey := new(PublicKey).Uncompress(pk)
	if pubkey == nil {
		return false
	}
	return signature.Verify(true, pubkey, false, msg, BLSMode)
}

func Sign(sk, msg []byte) []byte {
	prvkey := new(PrivateKey).FromBEndian(sk)
	if prvkey == nil {
		return nil
	}
	sig2 := new(Signature).Sign(prvkey, msg, BLSMode)
	return sig2.Compress()
}
