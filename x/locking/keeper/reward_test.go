package keeper_test

import (
	"math/big"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (suite *KeeperTestSuite) TestUpdateRewardPool() {
	err := suite.Keeper.RewardPool.Set(suite.Context, types.RewardPool{
		Goat:   math.ZeroInt(),
		Gas:    math.ZeroInt(),
		Remain: math.NewInt(300),
		Index:  1,
	})
	suite.Require().NoError(err)

	err = suite.Keeper.Params.Set(suite.Context, types.Params{
		InitialBlockReward: 200,
		HalvingInterval:    2,
	})
	suite.Require().NoError(err)

	err = suite.Keeper.UpdateRewardPool(suite.Context,
		[]*goattypes.GasRequest{{Amount: big.NewInt(100)}},
		[]*goattypes.GrantRequest{{Amount: big.NewInt(100)}}, true)
	suite.Require().NoError(err)

	updated, err := suite.Keeper.RewardPool.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(updated, types.RewardPool{
		Gas:    math.NewInt(100),
		Goat:   math.NewInt(200),
		Remain: math.NewInt(200),
		Index:  2,
	})

	err = suite.Keeper.UpdateRewardPool(suite.Context,
		[]*goattypes.GasRequest{{Amount: new(big.Int)}}, nil, false)
	suite.Require().NoError(err)

	updated, err = suite.Keeper.RewardPool.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(updated, types.RewardPool{
		Gas:    math.NewInt(100),
		Goat:   math.NewInt(200),
		Remain: math.NewInt(200),
		Index:  2,
	})

	err = suite.Keeper.UpdateRewardPool(suite.Context,
		[]*goattypes.GasRequest{{Amount: big.NewInt(100)}}, nil, true)
	suite.Require().NoError(err)

	updated, err = suite.Keeper.RewardPool.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(updated, types.RewardPool{
		Gas:    math.NewInt(200),
		Goat:   math.NewInt(300),
		Remain: math.NewInt(100),
		Index:  3,
	})

	err = suite.Keeper.UpdateRewardPool(suite.Context,
		[]*goattypes.GasRequest{{Amount: new(big.Int)}}, nil, true)
	suite.Require().NoError(err)

	updated, err = suite.Keeper.RewardPool.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(updated, types.RewardPool{
		Gas:    math.NewInt(200),
		Goat:   math.NewInt(400),
		Remain: math.NewInt(0),
		Index:  4,
	})

	err = suite.Keeper.UpdateRewardPool(suite.Context,
		[]*goattypes.GasRequest{{Amount: big.NewInt(100)}}, nil, true)
	suite.Require().NoError(err)

	updated, err = suite.Keeper.RewardPool.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(updated, types.RewardPool{
		Gas:    math.NewInt(300),
		Goat:   math.NewInt(400),
		Remain: math.NewInt(0),
		Index:  4,
	})
}

func (suite *KeeperTestSuite) TestDistributeReward() {
	amount, ok := new(big.Int).SetString("100000000000000000000000", 10)
	suite.Require().True(ok)

	Addresses := []sdk.ConsAddress{
		sdk.ConsAddress(common.Hex2Bytes("f52a75aa5be8d8c9e3580ea6ba818e68de4fb76e")),
		sdk.ConsAddress(common.Hex2Bytes("108ca95b90e680f7e4374f911521941fe78b85ce")),
	}

	Validators := []types.Validator{
		{
			Pubkey:    common.Hex2Bytes("03df9b92d37f4e3dec8ea95de7ab9f54879b978cebafd77a8528e7d832594a2af5"),
			Power:     223049,
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Active,
			Locking: sdk.NewCoins(
				sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
			),
		},
		{
			Pubkey:    common.Hex2Bytes("03baf046326e0d1f48ad417b7336727e4454a286461ce1b2d01d50b3029468fd63"),
			Power:     39237,
			Reward:    math.NewInt(1e9),
			GasReward: math.NewInt(1e8),
			Status:    types.Active,
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
				sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
			),
		},
	}

	for idx, validator := range Validators {
		err := suite.Keeper.Validators.Set(suite.Context, Addresses[idx], validator)
		suite.Require().NoError(err)

		for _, locking := range validator.Locking {
			err = suite.Keeper.Locking.Set(suite.Context,
				collections.Join(locking.Denom, Addresses[idx]), locking.Amount)
			suite.Require().NoError(err)
		}
		if validator.Power > 0 {
			err = suite.Keeper.PowerRanking.Set(suite.Context,
				collections.Join(validator.Power, Addresses[idx]))
			suite.Require().NoError(err)
		}
	}

	err := suite.Keeper.RewardPool.Set(suite.Context, types.RewardPool{
		Goat:   math.NewInt(288285324),
		Gas:    math.NewInt(632676834),
		Remain: math.NewInt(300),
		Index:  1,
	})
	suite.Require().NoError(err)

	newctx := suite.Context.WithBlockHeight(10).WithVoteInfos([]abci.VoteInfo{
		{
			Validator:   abci.Validator{Address: Addresses[0], Power: int64(Validators[0].Power)},
			BlockIdFlag: tmtypes.BlockIDFlagCommit,
		},
		{
			Validator:   abci.Validator{Address: Addresses[1], Power: int64(Validators[1].Power)},
			BlockIdFlag: tmtypes.BlockIDFlagAbsent,
		}})

	err = suite.Keeper.DistributeReward(newctx)
	suite.Require().NoError(err)

	gasReward := []math.Int{math.NewInt(538030757), math.NewInt(194646076)}
	goatReward := []math.Int{math.NewInt(245158922), math.NewInt(1043126401)}

	for idx, address := range Addresses {
		validator, err := suite.Keeper.Validators.Get(newctx, address)
		suite.Require().NoError(err)
		suite.Require().Equal(validator.GasReward, gasReward[idx])
		suite.Require().Equal(validator.Reward, goatReward[idx])
	}

	pool, err := suite.Keeper.RewardPool.Get(newctx)
	suite.Require().NoError(err)
	suite.Require().Equal(pool.Gas, math.NewInt(1))
	suite.Require().Equal(pool.Goat, math.NewInt(1))
}
