package keeper_test

import (
	"cosmossdk.io/collections"
	"github.com/ethereum/go-ethereum/common"
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
