package types

import (
	"encoding/binary"
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
)

type IVoteMsg interface {
	GetProposer() string
	GetVote() *Votes
	MethodName() string
	VoteSigDoc() []byte
}

func (v *Votes) Validate() error {
	if len(v.Signature) != goatcrypto.SignatureLength {
		return errors.New("invalid bls signature length")
	}
	return nil
}

func (v *Voter) Validate() error {
	if len(v.VoteKey) != goatcrypto.PubkeyLength {
		return errors.New("invalid bls pubkey length")
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

func (msg *MsgNewVoterRequest) Validate() error {
	if len(msg.VoterBlsKey) != goatcrypto.PubkeyLength {
		return errors.New("invalid vote pubkey length")
	}
	if len(msg.VoterBlsKeyProof) != goatcrypto.SignatureLength {
		return errors.New("invalid vote pubkey proof length")
	}

	if len(msg.VoterTxKey) != btcec.PubKeyBytesLenCompressed {
		return errors.New("invalid tx pubkey length")
	}

	if len(msg.VoterTxKeyProof) != goatcrypto.Secp256k1SigLength {
		return errors.New("invalid tx pubkey proof length")
	}

	return nil
}

func (msg *MsgNewVoterRequest) SignDoc(chainId string, height uint64, address, blsPKHash []byte) []byte {
	h := make([]byte, 8)
	binary.LittleEndian.PutUint64(h, height)
	return goatcrypto.SHA256Sum(
		[]byte(chainId),
		[]byte("NewVoter/Proof"),
		h,
		address,
		blsPKHash,
	)
}
