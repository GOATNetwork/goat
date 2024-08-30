package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "goat"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_goat"
)

var (
	ParamsKey   = collections.NewPrefix(0)
	BlockKey    = collections.NewPrefix(1)
	ConsHashKey = collections.NewPrefix(2)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
