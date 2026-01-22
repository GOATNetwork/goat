package keeper_test

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/goatnetwork/goat/x/locking/keeper"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (suite *KeeperTestSuite) TestParamsQuery() {
	qs := keeper.NewQueryServerImpl(suite.Keeper)
	params := types.DefaultParams()
	suite.Require().NoError(suite.Keeper.Params.Set(suite.Context, params))

	response, err := qs.Params(suite.Context, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryParamsResponse{Params: params}, response)
}

func (suite *KeeperTestSuite) TestValidatorQuery() {
	qs := keeper.NewQueryServerImpl(suite.Keeper)

	const height = 100
	suite.Context = suite.Context.WithBlockHeight(height)

	for idx, validator := range suite.Validator {
		err := suite.Keeper.Validators.Set(suite.Context, suite.Address[idx], validator)
		suite.Require().NoError(err)

		if validator.Status == types.Active || validator.Status == types.Pending {
			for _, locking := range validator.Locking {
				err = suite.Keeper.Locking.Set(suite.Context,
					collections.Join(locking.Denom, suite.Address[idx]), locking.Amount)
				suite.Require().NoError(err)
			}
			if validator.Power > 0 {
				err = suite.Keeper.PowerRanking.Set(suite.Context, collections.Join(validator.Power, suite.Address[idx]))
				suite.Require().NoError(err)
			}
		}
	}

	for idx, validator := range suite.Validator {
		response, err := qs.Validator(suite.Context, &types.QueryValidatorRequest{
			Address: common.BytesToAddress(suite.Address[idx]).String(),
		})
		suite.Require().NoError(err)
		suite.Require().Equal(&types.QueryValidatorResponse{Validator: validator, Height: height}, response)
	}

	_, err := qs.Validator(suite.Context, &types.QueryValidatorRequest{
		Address: (common.Address{}).String(),
	})
	suite.Require().ErrorContains(err, "not found")

	_, err = qs.Validator(suite.Context, &types.QueryValidatorRequest{
		Address: "invalid",
	})
	suite.Require().ErrorContains(err, "invalid address")
}

func (suite *KeeperTestSuite) TestActiveValidatorsQuery() {
	qs := keeper.NewQueryServerImpl(suite.Keeper)

	const height int64 = 100
	suite.Context = suite.Context.WithBlockHeight(height)

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
	validatorAddr := "0xf0933654a540830e283b87bba9ff2eb16b5acd1d"
	Address := sdk.ConsAddress(hexutil.MustDecode(validatorAddr))

	err := suite.Keeper.Validators.Set(suite.Context, Address, Validator)
	suite.Require().NoError(err)
	err = suite.Keeper.ValidatorSet.Set(suite.Context, Address, Validator.Power)
	suite.Require().NoError(err)

	resp, err := qs.ActiveValidators(suite.Context, &types.QueryActiveValidatorsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(height, resp.Height)
	suite.Require().Equal(resp.TotalPower, int64(10000))
	suite.Require().Len(resp.Validators, 1)
	suite.Require().Equal(validatorAddr, resp.Validators[0].ValidatorAddress)
	suite.Require().Equal(Validator, resp.Validators[0].Validator)
	suite.Require().Equal(math.LegacyNewDecWithPrec(100, 0), resp.Validators[0].PowerPercentage)
}
