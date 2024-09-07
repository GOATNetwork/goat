package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
)

type RelayerKeeper interface {
	VerifyProposal(ctx context.Context, req relayer.IVoteMsg, verifyFn ...func(sigdoc []byte) error) (uint64, error)
	UpdateRandao(ctx context.Context, req relayer.IVoteMsg) error
	HasPubkey(ctx context.Context, raw []byte) (bool, error)
	AddNewKey(ctx context.Context, raw []byte) error
	SetProposalSeq(ctx context.Context, seq uint64) error
}

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
