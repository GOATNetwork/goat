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

func (suite *KeeperTestSuite) TestHandleEvidences() {
	Validator := types.Validator{
		Pubkey:    common.Hex2Bytes("03ac22905ded6095255f498cd5cb217b6ebf0d82c7df2c89bce6e9089dd51e6f50"),
		Power:     10000,
		Reward:    math.ZeroInt(),
		GasReward: math.ZeroInt(),
		Status:    types.Active,
		Locking: sdk.NewCoins(
			sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
			sdk.NewCoin(GoatToekenDenom, math.NewIntFromUint64(1e9)),
			sdk.NewCoin(TestTokenDenom, math.NewIntFromUint64(10)),
		),
	}
	Address := sdk.ConsAddress(common.Hex2Bytes("f0933654a540830e283b87bba9ff2eb16b5acd1d"))

	err := suite.Keeper.Validators.Set(suite.Context, Address, Validator)
	suite.Require().NoError(err)
	for _, locking := range Validator.Locking {
		err = suite.Keeper.Locking.Set(suite.Context,
			collections.Join(locking.Denom, Address), locking.Amount)
		suite.Require().NoError(err)
	}

	err = suite.Keeper.PowerRanking.Set(suite.Context,
		collections.Join(Validator.Power, Address))
	suite.Require().NoError(err)

	for address, token := range suite.Token {
		err := suite.Keeper.Tokens.Set(suite.Context, address, token)
		suite.Require().NoError(err)
	}

	err = suite.Keeper.Threshold.Set(suite.Context, types.Threshold{List: suite.Threshold})
	suite.Require().NoError(err)

	param := types.Params{
		SignedBlocksWindow:      3,
		MaxMissedPerWindow:      1,
		DowntimeJailDuration:    time.Hour,
		SlashFractionDoubleSign: math.LegacyNewDec(2).QuoInt64(100),
	}
	err = suite.Keeper.Params.Set(suite.Context, param)
	suite.Require().NoError(err)

	newctx := suite.Context.WithBlockHeight(10).WithCometInfo(baseapp.NewBlockInfo(nil, nil, Address, abci.CommitInfo{}))
	err = suite.Keeper.HandleEvidences(newctx)
	suite.Require().NoError(err)

	now := time.Now().UTC()

	newctx = suite.Context.
		WithConsensusParams(cmtproto.ConsensusParams{Evidence: &cmtproto.EvidenceParams{
			MaxAgeNumBlocks: 1,
			MaxAgeDuration:  time.Second * 10,
		}}).
		WithBlockHeight(11).
		WithBlockTime(now).
		WithCometInfo(baseapp.NewBlockInfo([]abci.Misbehavior{
			{
				Type:      abci.MisbehaviorType_DUPLICATE_VOTE,
				Validator: abci.Validator{Address: Address, Power: int64(Validator.Power)},
				Time:      now.Add(-time.Minute),
				Height:    1,
			},
		}, nil, Address, abci.CommitInfo{}))
	err = suite.Keeper.HandleEvidences(newctx)
	suite.Require().NoError(err)

	updated, err := suite.Keeper.Validators.Get(newctx, Address)
	suite.Require().NoError(err)
	suite.Require().Equal(updated, Validator)

	newctx = newctx.WithConsensusParams(cmtproto.ConsensusParams{Evidence: &cmtproto.EvidenceParams{
		MaxAgeNumBlocks: 100,
		MaxAgeDuration:  time.Hour,
	}})

	err = suite.Keeper.HandleEvidences(newctx)
	suite.Require().NoError(err)

	updated, err = suite.Keeper.Validators.Get(newctx, Address)
	suite.Require().NoError(err)

	expected := types.Validator{
		Pubkey:      Validator.Pubkey,
		Power:       0,
		Reward:      math.ZeroInt(),
		GasReward:   math.ZeroInt(),
		Status:      types.Tombstoned,
		SigningInfo: Validator.SigningInfo,
		JailedUntil: Validator.JailedUntil,
		Locking: sdk.NewCoins(
			sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(980000000000000000)),
			sdk.NewCoin(GoatToekenDenom, math.NewIntFromUint64(980000000)),
		),
	}
	suite.Require().Equal(updated, expected)

	iter, err := suite.Keeper.Slashed.Iterate(newctx, nil)
	suite.Require().NoError(err)

	var slashed = sdk.NewCoins()
	for ; iter.Valid(); iter.Next() {
		kv, err := iter.KeyValue()
		suite.Require().NoError(err)
		slashed = slashed.Add(sdk.NewCoin(kv.Key, kv.Value))
	}

	suite.Require().Equal(slashed, sdk.NewCoins(
		sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(20000000000000000)),
		sdk.NewCoin(GoatToekenDenom, math.NewIntFromUint64(20000000)),
		sdk.NewCoin(TestTokenDenom, math.NewIntFromUint64(10)),
	))

	rankIter, err := suite.Keeper.PowerRanking.Iterate(newctx, nil)
	suite.Require().NoError(err)
	ranking, err := rankIter.Keys()
	suite.Require().NoError(err)
	suite.Require().Len(ranking, 0)

	err = suite.Keeper.HandleEvidences(newctx)
	suite.Require().NoError(err)

	updated, err = suite.Keeper.Validators.Get(newctx, Address)
	suite.Require().NoError(err)
	suite.Require().Equal(updated, expected)
}
