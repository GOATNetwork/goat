package crypto

const (
	Secp256k1SigLength = 64
)

func CompressP256k1Pubkey(pubkey [64]byte) []byte {
	if pubkey[63]&1 == 1 {
		return append([]byte{0x03}, pubkey[:32]...)
	}
	return append([]byte{0x02}, pubkey[:32]...)
}
