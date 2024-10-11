package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
//
//go:generate mockgen -source=expected_keepers.go -destination=../../../testutil/mock/account_keeper.go -package=mock
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
	HasAccount(context.Context, sdk.AccAddress) bool
	SetAccount(context.Context, sdk.AccountI)
	RemoveAccount(ctx context.Context, acc sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}
