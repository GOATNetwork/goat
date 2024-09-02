package types

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"

	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
)

const (
	NewPubkeyMethodSigName = "Bitcoin/NewPubkey"
	NewBlocksMethodSigName = "Bitcoin/NewBlocks"
)

const (
	EVMAddressLen   = 20
	DepositMagicLen = 4
	DustTxoutAmount = 546
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

	if len(req.Deposits) > 16 {
		return errors.New("deposit list too large")
	}

	for _, v := range req.Deposits {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	// todo: more...

	return nil
}

func (req *Deposit) Validate() error {
	if len(req.EvmAddress) != 20 {
		return errors.New("invalid evm address")
	}

	if len(req.BlockHeader) != 80 {
		return errors.New("invalid block header")
	}

	if err := req.RelayerPubkey.Validate(); err != nil {
		return err
	}

	return nil
}
