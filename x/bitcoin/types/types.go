package types

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
)

const (
	NewPubkeyMethodSigName = "Bitcoin/NewPubkey"
	NewBlocksMethodSigName = "Bitcoin/NewBlocks"
)

const (
	DepositMagicLen    = 4
	DustTxoutAmount    = 546
	RawBtcHeaderSize   = 80
	DepositV0TxoutSize = 34
	P2whScriptSize     = 22
	DepositV1TxoutSize = 26
	// 4 version
	// 1 input length
	// 32 prevTxid + 4 prevTxOut + 1 sigScriptLength + 0 sigScript(witness) + 4 sequence
	// 1 output length
	// 8 value + 1 pkScriptLength + 34 p2wsh/p2tr
	// || 8 value + 1 pkScriptLength + 22 p2wph +  8 value + 1 pkScriptLength + 26 data OP_RETURN
	// 4 lockTime
	MinDepositTxSize = 4 + 1 + 32 + 4 + 1 + 0 + 4 + 1 + 8 + 1 + 34 + 4
	MaxBtcTxSize     = 1024 * 1024 // the consensus hard limit
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
	binary.LittleEndian.AppendUint64(data[:8], req.StartBlockNumber)
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
		return errors.New("deposit list too large")
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

func (req *Deposit) Validate() error {
	if len(req.EvmAddress) != common.AddressLength {
		return errors.New("invalid evm address")
	}

	if l := len(req.NoWitnessTx); l < MinDepositTxSize || l > MaxBtcTxSize {
		return errors.New("invalid btc tx size")
	}

	if err := req.RelayerPubkey.Validate(); err != nil {
		return err
	}

	return nil
}
