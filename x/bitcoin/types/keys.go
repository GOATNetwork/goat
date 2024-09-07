package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "bitcoin"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_bitcoin"
)

var (
	ParamsKey         = collections.NewPrefix(0)
	LatestPubkeyKey   = collections.NewPrefix(1)
	LatestHeightKey   = collections.NewPrefix(2)
	BlockHashsKey     = collections.NewPrefix(3)
	DepositedKey      = collections.NewPrefix(4)
	ExecuableQueueKey = collections.NewPrefix(5)
	EthTxNonceKey     = collections.NewPrefix(6)

	WithdrawalKey         = collections.NewPrefix(7)
	WithdrawalProposalKey = collections.NewPrefix(8)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
