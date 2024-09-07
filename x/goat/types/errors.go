package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/goat module sentinel errors
var (
	ErrInvalidSigner  = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidRequest = sdkerrors.Register(ModuleName, 1101, "invalid request")
)
