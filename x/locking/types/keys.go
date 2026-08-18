package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "locking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_locking"
)

var (
	ParamsKey        = collections.NewPrefix(0)
	LockingKey       = collections.NewPrefix(1)
	PowerRankingKey  = collections.NewPrefix(2)
	ValidatorSetKey  = collections.NewPrefix(3)
	ValidatorsKey    = collections.NewPrefix(4)
	TokensKey        = collections.NewPrefix(5)
	SlashedKey       = collections.NewPrefix(6)
	EthTxNonceKey    = collections.NewPrefix(7)
	EthTxQueueKey    = collections.NewPrefix(8)
	RewardPoolKey    = collections.NewPrefix(9)
	UnlockQueueKey   = collections.NewPrefix(10)
	ThresholdKey     = collections.NewPrefix(11)
	FinalizedTimeKey = collections.NewPrefix(12)

	// ConsAddrIndexKey maps a consensus address to the validator id.
	// It only holds entries for validators whose consensus pubkey has been
	// rotated; for every other validator the id is the consensus address
	// itself, see Keeper.ResolveValidatorID.
	ConsAddrIndexKey = collections.NewPrefix(13)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
