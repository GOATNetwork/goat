package types

import "math"

var OsakaForkHeight = map[string]int64{
	"unitest": 10, // it's for unit test

	"goat-mainnet":  8821000, // estimate at 2025-12-19 16:00:00 UTC
	"goat-testnet3": 9695800, // estimate at 2025-12-15 15:00:00 UTC
}

var TzngForkHeight = map[string]int64{
	"unitest": 5, // it's for unit test

	"goat-mainnet":  math.MaxInt64, // TODO: set real height
	"goat-testnet3": math.MaxInt64, // TODO: set real height
}
