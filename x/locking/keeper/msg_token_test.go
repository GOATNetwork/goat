package keeper_test

import (
	"math/big"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (suite *KeeperTestSuite) TestUpdateToken() {
	for address, token := range suite.Token {
		err := suite.Keeper.Tokens.Set(suite.Context, address, token)
		suite.Require().NoError(err)
	}

	err := suite.Keeper.Threshold.Set(suite.Context, types.Threshold{List: suite.Threshold})
	suite.Require().NoError(err)

	err = suite.Keeper.UpdateTokens(suite.Context,
		[]*goattypes.UpdateTokenWeightRequest{{Token: TestToken, Weight: 10}},
		[]*goattypes.UpdateTokenThresholdRequest{{Token: TestToken, Threshold: big.NewInt(10)}})
	suite.Require().NoError(err)

	thres1, err := suite.Keeper.Threshold.Get(suite.Context)
	suite.Require().NoError(err)
	existed, amount := thres1.List.Find(TestTokenDenom)
	suite.Require().True(existed)
	suite.Require().Equal(amount, sdktypes.NewCoin(TestTokenDenom, math.NewInt(10)))

	token, err := suite.Keeper.Tokens.Get(suite.Context, TestTokenDenom)
	suite.Require().NoError(err)
	suite.Require().Equal(token, types.Token{Weight: 10, Threshold: math.NewInt(10)})

	err = suite.Keeper.UpdateTokens(suite.Context, nil,
		[]*goattypes.UpdateTokenThresholdRequest{{Token: TestToken, Threshold: big.NewInt(1)}})
	suite.Require().NoError(err)

	thres1, err = suite.Keeper.Threshold.Get(suite.Context)
	suite.Require().NoError(err)
	existed, amount = thres1.List.Find(TestTokenDenom)
	suite.Require().True(existed)
	suite.Require().Equal(amount, sdktypes.NewCoin(TestTokenDenom, math.NewInt(1)))

	token, err = suite.Keeper.Tokens.Get(suite.Context, TestTokenDenom)
	suite.Require().NoError(err)
	suite.Require().Equal(token, types.Token{Weight: 10, Threshold: math.NewInt(1)})

	err = suite.Keeper.UpdateTokens(suite.Context, nil,
		[]*goattypes.UpdateTokenThresholdRequest{{Token: TestToken, Threshold: big.NewInt(0)}})
	suite.Require().NoError(err)

	thres1, err = suite.Keeper.Threshold.Get(suite.Context)
	suite.Require().NoError(err)
	existed, _ = thres1.List.Find(TestTokenDenom)
	suite.Require().False(existed)

	token, err = suite.Keeper.Tokens.Get(suite.Context, TestTokenDenom)
	suite.Require().NoError(err)
	suite.Require().Equal(token, types.Token{Weight: 10, Threshold: math.NewInt(0)})

	newToken := common.HexToAddress("0xBC171dC497CC3410d3Ab989B7dc7Da37830c3a33")
	newTokenDenom := types.TokenDenom(newToken)
	err = suite.Keeper.UpdateTokens(suite.Context,
		[]*goattypes.UpdateTokenWeightRequest{{Token: newToken, Weight: 10}},
		[]*goattypes.UpdateTokenThresholdRequest{{Token: newToken, Threshold: big.NewInt(100)}})
	suite.Require().NoError(err)

	token, err = suite.Keeper.Tokens.Get(suite.Context, newTokenDenom)
	suite.Require().NoError(err)
	suite.Require().Equal(token, types.Token{Weight: 10, Threshold: math.NewInt(100)})

	thres1, err = suite.Keeper.Threshold.Get(suite.Context)
	suite.Require().NoError(err)
	existed, trsv := thres1.List.Find(newTokenDenom)
	suite.Require().True(existed)
	suite.Require().Equal(sdktypes.NewCoin(newTokenDenom, math.NewInt(100)), trsv)
}

func (suite *KeeperTestSuite) TestOnWeightChanged() {
	amount, ok := new(big.Int).SetString("100000000000000000000000", 10)
	suite.Require().True(ok)

	Addresses := []sdktypes.ConsAddress{
		sdktypes.ConsAddress(common.Hex2Bytes("f52a75aa5be8d8c9e3580ea6ba818e68de4fb76e")),
		sdktypes.ConsAddress(common.Hex2Bytes("108ca95b90e680f7e4374f911521941fe78b85ce")),
	}

	Validators := []types.Validator{
		{
			Pubkey:    common.Hex2Bytes("03df9b92d37f4e3dec8ea95de7ab9f54879b978cebafd77a8528e7d832594a2af5"),
			Reward:    math.ZeroInt(),
			Power:     100000,
			GasReward: math.ZeroInt(),
			Status:    types.Pending,
			Locking: sdktypes.NewCoins(
				sdktypes.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
			),
		},
		{
			Pubkey:    common.Hex2Bytes("03baf046326e0d1f48ad417b7336727e4454a286461ce1b2d01d50b3029468fd63"),
			Reward:    math.ZeroInt(),
			Power:     110000,
			GasReward: math.ZeroInt(),
			Status:    types.Active,
			Locking: sdktypes.NewCoins(
				sdktypes.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
				sdktypes.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
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

	for address, token := range suite.Token {
		err := suite.Keeper.Tokens.Set(suite.Context, address, token)
		suite.Require().NoError(err)
	}

	err := suite.Keeper.Threshold.Set(suite.Context, types.Threshold{List: suite.Threshold})
	suite.Require().NoError(err)

	{
		err = suite.Keeper.UpdateTokens(suite.Context,
			[]*goattypes.UpdateTokenWeightRequest{{Token: goattypes.GoatTokenContract}}, nil)
		suite.Require().NoError(err)

		iter, err := suite.Keeper.PowerRanking.Iterate(suite.Context, nil)
		suite.Require().NoError(err)

		powerSet := make(map[string]uint64)
		for ; iter.Valid(); iter.Next() {
			key, err := iter.Key()
			suite.Require().NoError(err)
			powerSet[string(key.K2())] = key.K1()
		}

		_, ok := powerSet[string(Addresses[0])]
		suite.Require().False(ok)

		power, ok := powerSet[string(Addresses[1])]
		suite.Require().True(ok)
		suite.Require().EqualValues(power, 10000)

		updates := []types.Validator{
			{
				Pubkey:    common.Hex2Bytes("03df9b92d37f4e3dec8ea95de7ab9f54879b978cebafd77a8528e7d832594a2af5"),
				Reward:    math.ZeroInt(),
				Power:     0,
				GasReward: math.ZeroInt(),
				Status:    types.Pending,
				Locking: sdktypes.NewCoins(
					sdktypes.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
				),
			},
			{
				Pubkey:    common.Hex2Bytes("03baf046326e0d1f48ad417b7336727e4454a286461ce1b2d01d50b3029468fd63"),
				Reward:    math.ZeroInt(),
				Power:     10000,
				GasReward: math.ZeroInt(),
				Status:    types.Active,
				Locking: sdktypes.NewCoins(
					sdktypes.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
					sdktypes.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
				),
			},
		}

		for idx, update := range updates {
			validator, err := suite.Keeper.Validators.Get(suite.Context, Addresses[idx])
			suite.Require().NoError(err)
			suite.Require().Equal(update, validator, "down test: idx", idx)

			for _, locking := range validator.Locking {
				existed, err := suite.Keeper.Locking.Has(suite.Context, collections.Join(locking.Denom, Addresses[idx]))
				suite.Require().NoError(err)
				suite.Require().True(existed)
			}
		}
	}

	{
		err = suite.Keeper.UpdateTokens(suite.Context,
			[]*goattypes.UpdateTokenWeightRequest{{Token: goattypes.GoatTokenContract, Weight: 1}}, nil)
		suite.Require().NoError(err)

		iter, err := suite.Keeper.PowerRanking.Iterate(suite.Context, nil)
		suite.Require().NoError(err)

		powerSet := make(map[string]uint64)
		for ; iter.Valid(); iter.Next() {
			key, err := iter.Key()
			suite.Require().NoError(err)
			powerSet[string(key.K2())] = key.K1()
		}

		power, ok := powerSet[string(Addresses[0])]
		suite.Require().True(ok)
		suite.Require().EqualValues(power, 100000)

		power, ok = powerSet[string(Addresses[1])]
		suite.Require().True(ok)
		suite.Require().EqualValues(power, 110000)

		for idx, update := range Validators {
			validator, err := suite.Keeper.Validators.Get(suite.Context, Addresses[idx])
			suite.Require().NoError(err)
			suite.Require().Equal(update, validator, "up test: idx", idx)

			for _, locking := range validator.Locking {
				existed, err := suite.Keeper.Locking.Has(suite.Context, collections.Join(locking.Denom, Addresses[idx]))
				suite.Require().NoError(err)
				suite.Require().True(existed)
			}
		}
	}
}
