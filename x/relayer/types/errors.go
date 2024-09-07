package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/relayer module sentinel errors
var (
	ErrInvalidSigner            = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalid                  = sdkerrors.Register(ModuleName, 1101, "invalid error")
	ErrNotProposer              = sdkerrors.Register(ModuleName, 1102, "not current proposer")
	ErrInvalidProposalSignature = sdkerrors.Register(ModuleName, 1103, "invalid propsal signature")
	ErrInvalidPubkey            = sdkerrors.Register(ModuleName, 1104, "invalid pubkey")
)
