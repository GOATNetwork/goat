package types

import (
	"crypto/sha256"
	"errors"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
)

const (
	NewPubkeyMethodSigName            = "Bitcoin/NewPubkey"
	NewBlocksMethodSigName            = "Bitcoin/NewBlocks"
	InitializeWithdrawalMethodSigName = "Bitcoin/InitializeWithdrawal"
)

const (
	DepositMagicLen    = 4
	DustTxoutAmount    = 546
	RawBtcHeaderSize   = 80
	P2TRScriptSize     = 34
	P2WSHScriptSize    = 34
	P2WPKHScriptSize   = 22
	MinBtcTxSize       = 4 + 1 + 32 + 4 + 1 + 0 + 4 + 1 + 8 + 1 + 22 + 4
	DepositV1TxoutSize = 26
	// 4 version
	// 1 input length
	// 41 = 32 prevTxid + 4 prevTxOut + 1 sigScriptLength + 0 sigScript(witness) + 4 sequence
	// 1 output length
	// 8 value + 1 pkScriptLength + 34 p2wsh/p2tr
	// || 8 value + 1 pkScriptLength + 22 p2wph +  8 value + 1 pkScriptLength + 26 data OP_RETURN
	// 4 lockTime
	MinDepositTxSize    = 4 + 1 + 32 + 4 + 1 + 0 + 4 + 1 + 8 + 1 + 34 + 4
	MaxAllowedBtcTxSize = 32 * 1024
)

func (req *MsgNewPubkey) Validate() error {
	if req == nil {
		return errors.New("empty MsgNewKey")
	}
	if err := req.Pubkey.Validate(); err != nil {
		return err
	}
	if err := req.Vote.Validate(); err != nil {
		return err
	}
	return nil
}

func (req *MsgNewPubkey) MethodName() string {
	return NewPubkeyMethodSigName
}

func (req *MsgNewPubkey) VoteSigDoc() []byte {
	return relayertypes.EncodePublicKey(req.Pubkey)
}

func (req *MsgNewBlockHashes) Validate() error {
	if req == nil {
		return errors.New("empty MsgNewBlockHashes")
	}

	if req.StartBlockNumber == 0 {
		return errors.New("block number is 0")
	}

	if len(req.BlockHash) > 16 {
		return errors.New("block hash list too large")
	}

	for _, v := range req.BlockHash {
		if len(v) != sha256.Size {
			return errors.New("block hash should be 32 bytes")
		}
	}

	if err := req.Vote.Validate(); err != nil {
		return err
	}
	return nil
}

func (req *MsgNewBlockHashes) MethodName() string {
	return NewBlocksMethodSigName
}

func (req *MsgNewBlockHashes) VoteSigDoc() []byte {
	data := make([]byte, 8, 8+len(req.BlockHash)*32)
	data = append(data, goatcrypto.Uint64LE(req.StartBlockNumber)...)
	for _, v := range req.BlockHash {
		data = append(data, v...)
	}
	return data
}

func (req *MsgNewDeposits) Validate() error {
	if req == nil {
		return errors.New("empty MsgNewDeposits")
	}

	depositLen := len(req.Deposits)
	if depositLen == 0 || depositLen > 16 {
		return errors.New("invalid deposit list length")
	}

	if h := len(req.BlockHeaders); h == 0 || h > depositLen {
		return errors.New("invalid headers size")
	}

	for _, v := range req.BlockHeaders {
		if len(v) != RawBtcHeaderSize {
			return errors.New("invalid raw header bytes size")
		}
	}

	for _, v := range req.Deposits {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

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
		return errors.New("empty MsgNewWithdrawal")
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

func (req *MsgApproveCancellation) Validate() error {
	if req == nil {
		return errors.New("empty MsgNewWithdrawal")
	}

	if len(req.Id) == 0 {
		return errors.New("no withdrawal ids to process")
	}

	if len(req.Id) > 32 {
		return errors.New("associate with too many withdrawal")
	}

	return nil
}

func (req *Deposit) Validate() error {
	if len(req.EvmAddress) != common.AddressLength {
		return errors.New("invalid evm address")
	}

	if l := len(req.NoWitnessTx); l < MinDepositTxSize || l > MaxAllowedBtcTxSize {
		return errors.New("invalid btc tx size")
	}

	if err := req.RelayerPubkey.Validate(); err != nil {
		return err
	}

	return nil
}

func (req *MsgFinalizeWithdrawal) Validate() error {
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
