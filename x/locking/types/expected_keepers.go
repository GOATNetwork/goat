package types

import (
	"context"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
	HasAccount(context.Context, sdk.AccAddress) bool
	SetAccount(context.Context, sdk.AccountI)
	RemoveAccount(ctx context.Context, acc sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// ConsensusParamsKeeper is the slice of x/consensus this module needs to widen
// the validator pub_key_types whitelist at a fork height. There is no
// transaction that could do it: GOAT has no x/gov, and the ante handler does
// not admit cosmos.consensus.v1.MsgUpdateParams.
type ConsensusParamsKeeper interface {
	Get(ctx context.Context) (cmtproto.ConsensusParams, error)
	Set(ctx context.Context, params cmtproto.ConsensusParams) error
}
