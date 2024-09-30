package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type BitcoinKeeper interface {
	DequeueBitcoinModuleTx(ctx context.Context) ([]*ethtypes.Transaction, error)
	ProcessBridgeRequest(ctx context.Context, withdrawals []*WithdrawalReq, rbf []*ReplaceByFeeReq, cancel1 []*Cancel1Req) error
}

type LockingKeeper interface {
	DequeueLockingModuleTx(ctx context.Context) ([]*ethtypes.Transaction, error)
}

type RelayerKeeper interface {
	GetCurrentProposer(ctx context.Context) (sdk.AccAddress, error)
	ProcessRelayerRequest(ctx context.Context, adds []*AddVoterReq, rms []*RemoveVoterReq) error
}

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	// Return a new account with the next account number and the specified address. Does not save the new account to the store.
	NewAccountWithAddress(context.Context, sdk.AccAddress) sdk.AccountI

	// Return a new account with the next account number. Does not save the new account to the store.
	NewAccount(context.Context, sdk.AccountI) sdk.AccountI

	// Check if an account exists in the store.
	HasAccount(context.Context, sdk.AccAddress) bool

	// Retrieve an account from the store.
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI

	// Set an account in the store.
	SetAccount(context.Context, sdk.AccountI)

	// Fetch the sequence of an account at a specified address.
	GetSequence(context.Context, sdk.AccAddress) (uint64, error)

	// Fetch the next account number, and increment the internal counter.
	NextAccountNumber(context.Context) uint64

	// AddressCodec returns the account address codec.
	AddressCodec() address.Codec
	// Methods imported from account should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
