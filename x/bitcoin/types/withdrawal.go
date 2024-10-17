package types

import (
	"crypto/sha256"
	"errors"
	"slices"

	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
)

func (req *MsgInitializeWithdrawal) MethodName() string {
	return InitializeWithdrawalMethodSigName
}

func (req *MsgInitializeWithdrawal) VoteSigDoc() []byte {
	ids := goatcrypto.Uint64LE(req.Proposal.Id...)
	fee := goatcrypto.Uint64LE(req.Proposal.TxFee)
	tx := goatcrypto.SHA256Sum(req.Proposal.NoWitnessTx)
	return slices.Concat(ids, tx, fee)
}

func (req *MsgInitializeWithdrawal) Validate() error {
	if req == nil {
		return errors.New("empty MsgInitializeWithdrawal")
	}

	if req.Proposal == nil {
		return errors.New("empty proposal")
	}

	if txSize := len(req.Proposal.NoWitnessTx); txSize < MinBtcTxSize || txSize > MaxAllowedBtcTxSize {
		return errors.New("invalid non-witness tx size")
	}

	if req.Proposal.TxFee == 0 {
		return errors.New("invalid tx fee")
	}

	if len(req.Proposal.Id) == 0 {
		return errors.New("no withdrawal ids to process")
	}

	if len(req.Proposal.Id) > 32 {
		return errors.New("associate with too many withdrawal")
	}

	return nil
}

func (req *MsgFinalizeWithdrawal) Validate() error {
	if req == nil {
		return errors.New("empty MsgFinalizeWithdrawal")
	}

	if len(req.Txid) != sha256.Size {
		return errors.New("invalid txid")
	}
	if req.TxIndex == 0 || len(req.IntermediateProof) == 0 {
		return errors.New("withdrawal can't be a coinbase tx")
	}
	if len(req.BlockHeader) != RawBtcHeaderSize {
		return errors.New("invalid block header size")
	}
	return nil
}

func (req *MsgApproveCancellation) Validate() error {
	if req == nil {
		return errors.New("empty MsgApproveCancellation")
	}

	if len(req.Id) == 0 {
		return errors.New("no withdrawal ids to process")
	}

	if len(req.Id) > 32 {
		return errors.New("associate with too many withdrawal")
	}

	return nil
}
