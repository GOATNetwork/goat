package types

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
	HasAccount(context.Context, sdk.AccAddress) bool
	SetAccount(context.Context, sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

type IVoteMsg interface {
	GetProposer() string
	GetVote() *Votes
	MethodName() string
	VoteSigDoc() []byte
}

type INonVoteMsg interface {
	GetProposer() string
}

type IRelayer interface {
	GetProposer() string
	GetEpoch() uint64
	GetLastElected() time.Time
	GetVoters() []string
}
