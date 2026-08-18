package keeper_test

import (
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/goatnetwork/goat/x/locking/types"
)

// the id a validator was created with, it never changes
var rotationID = sdk.ConsAddress(common.Hex2Bytes("f0933654a540830e283b87bba9ff2eb16b5acd1d"))

// the consensus address of the pubkey the validator rotated to
var rotatedConsAddr = sdk.ConsAddress(common.Hex2Bytes("2b1c5cd45e0c0eb1a9b9e0f8a2b3c4d5e6f70819"))

// setupRotatedValidator writes an Active validator keyed by rotationID and
// points rotatedConsAddr at it, i.e. the state a validator is in right after
// its consensus pubkey has been rotated.
func (suite *KeeperTestSuite) setupRotatedValidator(param types.Params) types.Validator {
	validator := types.Validator{
		Pubkey:    common.Hex2Bytes("03ac22905ded6095255f498cd5cb217b6ebf0d82c7df2c89bce6e9089dd51e6f50"),
		Power:     10000,
		Reward:    math.ZeroInt(),
		GasReward: math.ZeroInt(),
		Status:    types.Active,
		Locking: sdk.NewCoins(
			sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
			sdk.NewCoin(GoatToekenDenom, math.NewIntFromUint64(1e9)),
		),
	}

	err := suite.Keeper.Validators.Set(suite.Context, rotationID, validator)
	suite.Require().NoError(err)
	for _, locking := range validator.Locking {
		err = suite.Keeper.Locking.Set(suite.Context,
			collections.Join(locking.Denom, rotationID), locking.Amount)
		suite.Require().NoError(err)
	}
	err = suite.Keeper.PowerRanking.Set(suite.Context,
		collections.Join(validator.Power, rotationID))
	suite.Require().NoError(err)

	for denom, token := range suite.Token {
		err := suite.Keeper.Tokens.Set(suite.Context, denom, token)
		suite.Require().NoError(err)
	}
	err = suite.Keeper.Threshold.Set(suite.Context, types.Threshold{List: suite.Threshold})
	suite.Require().NoError(err)
	err = suite.Keeper.Params.Set(suite.Context, param)
	suite.Require().NoError(err)

	// the rotation itself: CometBFT will start reporting rotatedConsAddr while
	// every table stays keyed by rotationID
	err = suite.Keeper.ConsAddrIndex.Set(suite.Context, rotatedConsAddr, rotationID)
	suite.Require().NoError(err)

	return validator
}

// A validator that never rotated has no ConsAddrIndex entry and must resolve
// to itself. This is what keeps the index proportional to the number of
// rotated validators rather than to the validator set size.
func (suite *KeeperTestSuite) TestResolveValidatorIDFallsBackToItself() {
	for _, address := range suite.Address {
		id, err := suite.Keeper.ResolveValidatorID(suite.Context, address)
		suite.Require().NoError(err)
		suite.Require().Equal(address, id)
	}
}

func (suite *KeeperTestSuite) TestResolveValidatorIDUsesIndex() {
	err := suite.Keeper.ConsAddrIndex.Set(suite.Context, rotatedConsAddr, rotationID)
	suite.Require().NoError(err)

	id, err := suite.Keeper.ResolveValidatorID(suite.Context, rotatedConsAddr)
	suite.Require().NoError(err)
	suite.Require().Equal(rotationID, id)

	// the id itself keeps resolving to itself, so both the old and the new
	// address are usable during the rotation window
	id, err = suite.Keeper.ResolveValidatorID(suite.Context, rotationID)
	suite.Require().NoError(err)
	suite.Require().Equal(rotationID, id)
}

// BeginBlocker must not fail when CometBFT reports the rotated consensus
// address while the state is still keyed by the validator id. Before the
// resolver was introduced this was a Validators.Get miss, i.e. a halted chain.
func (suite *KeeperTestSuite) TestHandleVoteInfosAfterRotation() {
	param := types.Params{
		SignedBlocksWindow:    3,
		MaxMissedPerWindow:    1,
		DowntimeJailDuration:  time.Hour,
		SlashFractionDowntime: math.LegacyNewDec(2).QuoInt64(100),
	}
	validator := suite.setupRotatedValidator(param)

	// CometBFT reports the new address
	newctx := suite.Context.WithBlockHeight(1).WithVoteInfos([]abci.VoteInfo{{
		Validator:   abci.Validator{Address: rotatedConsAddr, Power: int64(validator.Power)},
		BlockIdFlag: cmtproto.BlockIDFlagCommit,
	}})
	suite.Require().NoError(suite.Keeper.HandleVoteInfos(newctx))

	// and the signing info lands on the id, not on the reported address
	updated, err := suite.Keeper.Validators.Get(newctx, rotationID)
	suite.Require().NoError(err)
	suite.Require().Equal(types.SigningInfo{Offset: 1}, updated.SigningInfo)

	has, err := suite.Keeper.Validators.Has(newctx, rotatedConsAddr)
	suite.Require().NoError(err)
	suite.Require().False(has, "no record may be created under the reported address")

	// the old address keeps working too, it is still reported in LastCommit
	// for two more blocks after the update takes effect
	newctx = suite.Context.WithBlockHeight(2).WithVoteInfos([]abci.VoteInfo{{
		Validator:   abci.Validator{Address: rotationID, Power: int64(validator.Power)},
		BlockIdFlag: cmtproto.BlockIDFlagCommit,
	}})
	suite.Require().NoError(suite.Keeper.HandleVoteInfos(newctx))

	updated, err = suite.Keeper.Validators.Get(newctx, rotationID)
	suite.Require().NoError(err)
	suite.Require().Equal(types.SigningInfo{Offset: 2}, updated.SigningInfo)
}

// Evidence references the consensus address that was active at the height of
// the misbehaviour, which may be a pubkey the validator has rotated away from.
// Keeping the index entry for a whole evidence window is what makes that
// evidence slashable.
func (suite *KeeperTestSuite) TestHandleEvidencesAfterRotation() {
	param := types.Params{
		SignedBlocksWindow:      3,
		MaxMissedPerWindow:      1,
		DowntimeJailDuration:    time.Hour,
		SlashFractionDoubleSign: math.LegacyNewDec(2).QuoInt64(100),
	}
	validator := suite.setupRotatedValidator(param)

	now := time.Now().UTC()
	newctx := suite.Context.
		WithConsensusParams(cmtproto.ConsensusParams{Evidence: &cmtproto.EvidenceParams{
			MaxAgeNumBlocks: 100,
			MaxAgeDuration:  time.Hour,
		}}).
		WithBlockHeight(11).
		WithBlockTime(now).
		WithCometInfo(baseapp.NewBlockInfo([]abci.Misbehavior{
			{
				Type: abci.MisbehaviorType_DUPLICATE_VOTE,
				// the address of the pubkey that was active when it double signed
				Validator: abci.Validator{Address: rotatedConsAddr, Power: int64(validator.Power)},
				Time:      now.Add(-time.Minute),
				Height:    1,
			},
		}, nil, rotationID, abci.CommitInfo{}))

	suite.Require().NoError(suite.Keeper.HandleEvidences(newctx))

	updated, err := suite.Keeper.Validators.Get(newctx, rotationID)
	suite.Require().NoError(err)
	suite.Require().Equal(types.Tombstoned, updated.Status)
	suite.Require().EqualValues(0, updated.Power)
	suite.Require().Equal(sdk.NewCoins(
		sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(980000000000000000)),
		sdk.NewCoin(GoatToekenDenom, math.NewIntFromUint64(980000000)),
	), updated.Locking)

	// and the power ranking entry keyed by the id is gone
	rankIter, err := suite.Keeper.PowerRanking.Iterate(newctx, nil)
	suite.Require().NoError(err)
	ranking, err := rankIter.Keys()
	suite.Require().NoError(err)
	suite.Require().Len(ranking, 0)
}

// DistributeReward runs in BeginBlocker and also keys Validators by the
// address CometBFT reports. It is the same halt shape as the vote and
// evidence handlers, and the reward has to land on the id.
func (suite *KeeperTestSuite) TestDistributeRewardAfterRotation() {
	param := types.Params{
		SignedBlocksWindow:   3,
		MaxMissedPerWindow:   1,
		DowntimeJailDuration: time.Hour,
	}
	validator := suite.setupRotatedValidator(param)

	err := suite.Keeper.RewardPool.Set(suite.Context, types.RewardPool{
		Goat:   math.NewInt(1000),
		Gas:    math.NewInt(2000),
		Remain: math.NewInt(0),
	})
	suite.Require().NoError(err)

	// CometBFT reports the rotated address, it is the only validator so it
	// takes the whole pool
	newctx := suite.Context.WithBlockHeight(10).WithVoteInfos([]abci.VoteInfo{{
		Validator:   abci.Validator{Address: rotatedConsAddr, Power: int64(validator.Power)},
		BlockIdFlag: cmtproto.BlockIDFlagCommit,
	}})
	suite.Require().NoError(suite.Keeper.DistributeReward(newctx))

	updated, err := suite.Keeper.Validators.Get(newctx, rotationID)
	suite.Require().NoError(err)
	suite.Require().Equal(math.NewInt(1000), updated.Reward)
	suite.Require().Equal(math.NewInt(2000), updated.GasReward)

	has, err := suite.Keeper.Validators.Has(newctx, rotatedConsAddr)
	suite.Require().NoError(err)
	suite.Require().False(has, "reward must not create a record under the reported address")
}
