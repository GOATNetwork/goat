package types

import (
	"math"

	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
)

func VoteSignDoc(method, chainId, proposer string, sequence, epoch uint64, data []byte) []byte {
	return goatcrypto.SHA256Sum(
		[]byte(chainId),
		goatcrypto.Uint64LE(sequence, epoch),
		[]byte(method),
		[]byte(proposer),
		data,
	)
}

func (relayer *Relayer) Threshold() int {
	return int(math.Ceil(float64(1+len(relayer.Voters)) * 2 / 3))
}
