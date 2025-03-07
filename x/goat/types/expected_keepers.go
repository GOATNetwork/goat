package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
)

type BitcoinKeeper interface {
	DequeueBitcoinModuleTx(ctx context.Context) ([]*ethtypes.Transaction, error)
	ProcessBridgeRequest(ctx context.Context, req goattypes.BridgeRequests) error
}

type LockingKeeper interface {
	DequeueLockingModuleTx(ctx context.Context) ([]*ethtypes.Transaction, error)
	ProcessLockingRequest(ctx context.Context, req goattypes.LockingRequests) error
}

type RelayerKeeper interface {
	GetCurrentProposer(ctx context.Context) (sdk.AccAddress, error)
	ProcessRelayerRequest(ctx context.Context, req goattypes.RelayerRequests) error
	VerifyProposal(ctx context.Context, req relayer.IVoteMsg, verifyFn ...func(sigdoc []byte) error) (uint64, error)
	VerifyNonProposal(ctx context.Context, req relayer.INonVoteMsg) (relayer.IRelayer, error)
	UpdateRandao(ctx context.Context, req relayer.IVoteMsg) error
	HasPubkey(ctx context.Context, raw []byte) (bool, error)
	AddNewKey(ctx context.Context, raw []byte) error
	SetProposalSeq(ctx context.Context, seq uint64) error
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

	RemoveAccount(ctx context.Context, acc sdk.AccountI)
}
