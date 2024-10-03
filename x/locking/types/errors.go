package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/locking module sentinel errors
var (
	ErrInvalid = sdkerrors.Register(ModuleName, 1100, "invalid error")
)
