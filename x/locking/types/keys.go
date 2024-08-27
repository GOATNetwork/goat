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
	ParamsKey = collections.NewPrefix("p_locking")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
