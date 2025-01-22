package types

import (
	"errors"
	"slices"

	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
)

type ProcessWithdrawalMsger interface {
	MethodName() string
	VoteSigDoc() []byte
	Validate() error
	CalTxPrice() float64
	GetProposer() string
	GetVote() *relayertypes.Votes
	GetId() []uint64
	GetNoWitnessTx() []byte
	GetTxFee() uint64
}

var _ ProcessWithdrawalMsger = &MsgProcessWithdrawalV2{}

func (req *MsgProcessWithdrawalV2) MethodName() string {
	return ProcessWithdrawalV2MethodSigName
}

func (req *MsgProcessWithdrawalV2) VoteSigDoc() []byte {
	ids := goatcrypto.Uint64LE(req.Id...)
	fee := goatcrypto.Uint64LE(req.TxFee, req.WitnessSize)
	tx := goatcrypto.SHA256Sum(req.NoWitnessTx)
	return slices.Concat(ids, tx, fee)
}

func (req *MsgProcessWithdrawalV2) Validate() error {
	if req == nil {
		return errors.New("empty MsgProcessWithdrawalV2")
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

	if req.WitnessSize == 0 {
		return errors.New("invalid witness size")
	}

	return nil
}

func (req *MsgProcessWithdrawalV2) CalTxPrice() float64 {
	// tx price = feeInSat / vbytes
	// total_size = stripped_size + witness_size
	// vbytes = (stripped_size * 3 + total_size) / 4 = stripped_size + witness_size / 4
	return float64(req.TxFee) / (float64(len(req.NoWitnessTx)) + float64(req.WitnessSize)/4)
}

type ReplaceWithdrawalMsger interface {
	MethodName() string
	VoteSigDoc() []byte
	Validate() error
	CalTxPrice() float64
	GetProposer() string
	GetVote() *relayertypes.Votes
	GetPid() uint64
	GetNewNoWitnessTx() []byte
	GetNewTxFee() uint64
}

var _ ReplaceWithdrawalMsger = &MsgReplaceWithdrawalV2{}

func (req *MsgReplaceWithdrawalV2) MethodName() string {
	return ReplaceWithdrawalV2MethodSigName
}

func (req *MsgReplaceWithdrawalV2) Validate() error {
	if req == nil {
		return errors.New("empty MsgReplaceWithdrawalV2")
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

	if req.WitnessSize == 0 {
		return errors.New("invalid witness size")
	}
	return nil
}

func (req *MsgReplaceWithdrawalV2) VoteSigDoc() []byte {
	ids := goatcrypto.Uint64LE(req.Pid, req.NewTxFee, req.WitnessSize)
	tx := goatcrypto.SHA256Sum(req.NewNoWitnessTx)
	return slices.Concat(ids, tx)
}

func (req *MsgReplaceWithdrawalV2) CalTxPrice() float64 {
	return float64(req.NewTxFee) / (float64(len(req.NewNoWitnessTx)) + float64(req.WitnessSize)/4)
}
