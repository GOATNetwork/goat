package app

import (
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"

	bitcoinkeeper "github.com/goatnetwork/goat/x/bitcoin/keeper"
	goattypes "github.com/goatnetwork/goat/x/goat/types"
	lockingkeeper "github.com/goatnetwork/goat/x/locking/keeper"
	lockingtypes "github.com/goatnetwork/goat/x/locking/types"
	relayerkeeper "github.com/goatnetwork/goat/x/relayer/keeper"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
)

// depinject wires the modules' expected keeper interfaces by implicit
// interface binding, so when an upstream signature drifts the container
// simply fails to find an implementer and the binary panics on startup.
// Neither the compiler nor the module tests catch that, because the tests
// run against mocks generated from these very interfaces.
//
// These assertions move that failure to build time.
var (
	_ goattypes.AccountKeeper    = authkeeper.AccountKeeper{}
	_ lockingtypes.AccountKeeper = authkeeper.AccountKeeper{}
	_ relayertypes.AccountKeeper = authkeeper.AccountKeeper{}

	_ goattypes.LockingKeeper = lockingkeeper.Keeper{}
	_ goattypes.RelayerKeeper = relayerkeeper.Keeper{}
	_ goattypes.BitcoinKeeper = bitcoinkeeper.Keeper{}

	// x/consensus provides the keeper, and it is its params store that
	// implements the interface, not the keeper itself
	_ lockingtypes.ConsensusParamsKeeper = consensuskeeper.Keeper{}.ParamsStore
)
