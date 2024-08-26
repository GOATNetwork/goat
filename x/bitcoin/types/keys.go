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
	ParamsKey       = collections.NewPrefix(0)
	LatestPubkeyKey = collections.NewPrefix(1)
	LatestHeightKey = collections.NewPrefix(2)
	BlockHashsKey   = collections.NewPrefix(3)
	Depositedkey    = collections.NewPrefix(4)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
