package types

import (
	"crypto/sha256"
	"errors"
	"slices"

	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
)

func (req *MsgProcessWithdrawal) MethodName() string {
	return ProcessWithdrawalMethodSigName
}

func (req *MsgProcessWithdrawal) VoteSigDoc() []byte {
	ids := goatcrypto.Uint64LE(req.Id...)
	fee := goatcrypto.Uint64LE(req.TxFee)
	tx := goatcrypto.SHA256Sum(req.NoWitnessTx)
	return slices.Concat(ids, tx, fee)
}

func (req *MsgProcessWithdrawal) Validate() error {
	if req == nil {
		return errors.New("empty MsgProcessWithdrawal")
	}

	if req.Vote == nil {
		return errors.New("empty Vote")
	}

	if txSize := len(req.NoWitnessTx); txSize < MinBtcTxSize || txSize > MaxAllowedBtcTxSize {
		return errors.New("invalid non-witness tx size")
	}

	if req.TxFee == 0 {
		return errors.New("invalid tx fee")
	}

	if len(req.Id) == 0 {
		return errors.New("no withdrawal ids to process")
	}

	if len(req.Id) > 32 {
		return errors.New("associate with too many withdrawal")
	}

	return nil
}

func (req *MsgReplaceWithdrawal) MethodName() string {
	return ReplaceWithdrawalMethodSigName
}

func (req *MsgReplaceWithdrawal) Validate() error {
	if req == nil {
		return errors.New("empty MsgReplaceWithdrawal")
	}

	if req.Vote == nil {
		return errors.New("empty Vote")
	}

	if txSize := len(req.NewNoWitnessTx); txSize < MinBtcTxSize || txSize > MaxAllowedBtcTxSize {
		return errors.New("invalid non-witness tx size")
	}

	if req.NewTxFee == 0 {
		return errors.New("invalid tx fee")
	}
	return nil
}

func (req *MsgReplaceWithdrawal) VoteSigDoc() []byte {
	ids := goatcrypto.Uint64LE(req.Pid, req.NewTxFee)
	tx := goatcrypto.SHA256Sum(req.NewNoWitnessTx)
	return slices.Concat(ids, tx)
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
