package consensusfork

// The default value for unknown networks is 0, meaning the fork is active from genesis.
// key: network name, value: fork height

// unitest is the chain id the module tests run under.
const unitest = "unitest"

// OsakaForkHeight defines the fork height for Osaka upgrade on different networks
var OsakaForkHeight = map[string]int64{
	unitest: 10, // it's for unit test

	"goat-mainnet":  8821000, // estimate at 2025-12-19 16:00:00 UTC
	"goat-testnet3": 9695800, // estimate at 2025-12-15 15:00:00 UTC
}

// ReeseForkHeight is for fixing unlock queue and beacon root processing bug
var ReeseForkHeight = map[string]int64{
	unitest: 5, // it's for unit test

	"goat-mainnet":  11130000, // estimate at 2026-03-26 02:18:00 UTC
	"goat-testnet3": 11044000, // estimate at 2026-02-05 13:00:00 UTC
}

// RotationForkHeight is when the consensus layer starts accepting consensus
// key rotations.
//
// Nothing can produce a rotation before Locking.sol has the rotate function,
// and Locking.sol is a genesis predeploy, so on a live network its code can
// only change at a coordinated fork. This gate exists so that the consensus
// layer's behavior is pinned to a height of its own rather than depending on
// exactly when that happens.
var RotationForkHeight = map[string]int64{
	unitest: 2, // it's for unit test

	// deliberately unset for the public networks: the height can only be
	// chosen once the execution layer can carry the new Locking code
	// "goat-mainnet":  0,
	// "goat-testnet3": 0,
}

// PubKeyTypesForkHeight is when ml_dsa_65 is added to the consensus
// pub_key_types whitelist.
//
// GOAT registers neither x/gov nor x/upgrade and the ante handler admits no
// message that could carry cosmos.consensus.v1.MsgUpdateParams, so a consensus
// parameter can only be changed by code at a hardcoded height.
//
// This has to land at or before RotationForkHeight. A validator update naming
// a key type that is not on the whitelist is rejected by CometBFT, and a
// rejected update fails the block.
var PubKeyTypesForkHeight = map[string]int64{
	unitest: 2, // it's for unit test

	// deliberately unset for the public networks, see RotationForkHeight
	// "goat-mainnet":  0,
	// "goat-testnet3": 0,
}
