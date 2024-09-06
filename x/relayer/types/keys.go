package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "relayer"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_relayer"
)

var (
	ParamsKey   = collections.NewPrefix(0)
	RelayerKey  = collections.NewPrefix(1)
	VotersKey   = collections.NewPrefix(2)
	PubkeysKey  = collections.NewPrefix(3)
	SequenceKey = collections.NewPrefix(4)
	QueueKey    = collections.NewPrefix(5)
	RandDAOKey  = collections.NewPrefix(6)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
