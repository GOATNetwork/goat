package types

import (
	"bytes"
	"crypto/sha256"

	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
)

func VerifyMerkelProof(txid, root, proof []byte, index uint32) bool {
	if len(txid) != sha256.Size || len(root) != sha256.Size || len(proof)%sha256.Size != 0 {
		return false
	}

	var buf []byte
	if len(proof) > 0 {
		buf = make([]byte, 64)
	}

	current := txid
	for i := 0; i < len(proof)/sha256.Size; i++ {
		start := i * sha256.Size
		end := start + sha256.Size
		next := proof[start:end]
		if index&1 == 0 {
			copy(buf[:32], current)
			copy(buf[32:], next)
			current = goatcrypto.DoubleSHA256Sum(buf)
		} else {
			copy(buf[:32], next)
			copy(buf[32:], current)
			current = goatcrypto.DoubleSHA256Sum(buf)
		}
		index >>= 1
	}

	return bytes.Equal(current, root)
}
