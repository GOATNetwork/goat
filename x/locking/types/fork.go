package types

import "math"

var OsakaForkHeight = map[string]int64{
	"unitest": 10, // it's for unit test

	"mainnet":  math.MaxInt64, // TODO: set the mainnet fork height
	"testnet3": math.MaxInt64, // TODO: set the testnet3 fork height
}
