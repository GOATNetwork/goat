package types

import (
	"crypto/sha256"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
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

	if l := len(req.NoWitnessTx); l < MinDepositTxSize || l > MaxAllowedBtcTxSize {
		return errors.New("invalid btc tx size")
	}

	if err := req.RelayerPubkey.Validate(); err != nil {
		return err
	}

	return nil
}
