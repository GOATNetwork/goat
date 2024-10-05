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
	ParamsKey       = collections.NewPrefix(0)
	LockingKey      = collections.NewPrefix(1)
	PowerRankingKey = collections.NewPrefix(2)
	ValidatorSetKey = collections.NewPrefix(3)
	ValidatorsKey   = collections.NewPrefix(4)
	TokensKey       = collections.NewPrefix(5)
	SlashedKey      = collections.NewPrefix(6)
	EthTxNonceKey   = collections.NewPrefix(7)
	EthTxQueueKey   = collections.NewPrefix(8)
	RewardPoolKey   = collections.NewPrefix(9)
	UnlockQueueKey  = collections.NewPrefix(10)
	ThresholdKey    = collections.NewPrefix(11)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
