package types

import (
	"encoding/binary"
	"errors"

	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	blst "github.com/supranational/blst/bindings/go"
)

type IVoteMsg interface {
	GetProposer() string
	GetVote() *Votes
	MethodName() string
	VoteSigDoc() []byte
}

func (v *Votes) Validate() error {
	if len(v.Signature) != blst.BLST_P1_COMPRESS_BYTES {
		return errors.New("invalid bls signature length")
	}
	return nil
}

func VoteSignDoc(method, chainId, proposer string, sequence uint64, data []byte) []byte {
	var seqRaw [8]byte
	binary.LittleEndian.PutUint64(seqRaw[:], sequence)

	sigdoc := goatcrypto.SHA256Sum(
		[]byte(chainId),
		seqRaw[:],
		[]byte(method),
		[]byte(proposer),
		data,
	)
	return sigdoc
}
