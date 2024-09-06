package types

import (
	"crypto/sha256"
	"errors"
	"slices"

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
	if len(v.Voters) > sha256.Size { // we have max 256 voters
		return errors.New("voter bitmap too large")
	}

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

func VoteSignDoc(method, chainId, proposer string, sequence, epoch uint64, data []byte) []byte {
	sigdoc := goatcrypto.SHA256Sum(
		[]byte(chainId),
		goatcrypto.Uint64LE(sequence, epoch),
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

func NewOnBoardingVoterRequest(height uint64, txKeyHash, voteKeyHash []byte) *OnBoardingVoterRequest {
	return &OnBoardingVoterRequest{
		Height:      height,
		TxKeyHash:   txKeyHash,
		VoteKeyHash: voteKeyHash,
	}
}

func (msg *OnBoardingVoterRequest) SignDoc() []byte {
	return slices.Concat(goatcrypto.Uint64LE(msg.Height), msg.TxKeyHash, msg.VoteKeyHash)
}

func (msg *OnBoardingVoterRequest) MethodName() string {
	return "Relayer/NewVoter"
}
