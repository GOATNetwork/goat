package types

import (
	"encoding/binary"
	"errors"

	"github.com/goatnetwork/goat/pkg/crypto"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
)

type IVoteMsg interface {
	GetProposer() string
	GetVote() *Votes
	MethodName() string
	VoteSigDoc() []byte
}

func (v *Votes) Validate() error {
	if len(v.Signature) != crypto.SignatureLength {
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
