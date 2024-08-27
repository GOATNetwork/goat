package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BitcoinKeeper interface {
	// TODO Add methods imported from bitcoin should be defined here
}

type LockingKeeper interface {
	// TODO Add methods imported from locking should be defined here
}

type RelayerKeeper interface {
	// TODO Add methods imported from relayer should be defined here
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
