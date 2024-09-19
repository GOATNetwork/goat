package types

import (
	"errors"

	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
)

func (req *MsgNewConsolidation) Validate() error {
	if req == nil {
		return errors.New("empty MsgNewBlockHashes")
	}

	if txSize := len(req.NoWitnessTx); txSize < MinBtcTxSize || txSize > MaxAllowedBtcTxSize {
		return errors.New("invalid non-witness tx size")
	}

	if err := req.Vote.Validate(); err != nil {
		return err
	}
	return nil
}

func (req *MsgNewConsolidation) MethodName() string {
	return NewConsolidationMethodSigName
}

func (req *MsgNewConsolidation) VoteSigDoc() []byte {
	return goatcrypto.SHA256Sum(req.NoWitnessTx)
}
