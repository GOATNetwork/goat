package consensusfork

import "math"

// The default value for unknown networks is 0, meaning the fork is active from genesis.
// key: network name, value: fork height

// OsakaForkHeight defines the fork height for Osaka upgrade on different networks
var OsakaForkHeight = map[string]int64{
	"unitest": 10, // it's for unit test

	"goat-mainnet":  8821000, // estimate at 2025-12-19 16:00:00 UTC
	"goat-testnet3": 9695800, // estimate at 2025-12-15 15:00:00 UTC
}

// ReeseForkHeight is for fixing unlock queue and beacon root processing bug
var ReeseForkHeight = map[string]int64{
	"unitest": 5, // it's for unit test

	"goat-mainnet":  math.MaxInt64, // TODO: set real height
	"goat-testnet3": 11044000,      // estimate at 2026-02-05 13:00:00 UTC
}
