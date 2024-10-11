package keeper_test

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (suite *KeeperTestSuite) TestEndBlocker() {
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

			if validator.Status == types.Active {
				err = suite.Keeper.ValidatorSet.Set(suite.Context, suite.Address[idx], validator.Power)
				suite.Require().NoError(err)
			}
		}
	}

	for address, token := range suite.Token {
		err := suite.Keeper.Tokens.Set(suite.Context, address, token)
		suite.Require().NoError(err)
	}

	{
		err := suite.Keeper.Threshold.Set(suite.Context, types.Threshold{List: suite.Threshold})
		suite.Require().NoError(err)

		err = suite.Keeper.Params.Set(suite.Context, types.Params{MaxValidators: 2})
		suite.Require().NoError(err)

		vs, err := suite.Keeper.EndBlocker(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Len(vs, 0)
	}

	// Add power to validator 0
	{
		var power uint64 = 10000
		err := suite.Keeper.Validators.Set(suite.Context, suite.Address[0], types.Validator{
			Pubkey:    suite.Validator[0].Pubkey,
			Power:     power,
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Pending,
			Locking:   sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewInt(1e18))),
		})
		suite.Require().NoError(err)

		err = suite.Keeper.PowerRanking.Set(suite.Context, collections.Join(power, suite.Address[0]))
		suite.Require().NoError(err)

		vs, err := suite.Keeper.EndBlocker(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Len(vs, 1)
		suite.Require().Equal(vs[0].PubKey.GetSecp256K1(), suite.Validator[0].Pubkey)
		suite.Require().EqualValues(vs[0].Power, power)

		updated, err := suite.Keeper.Validators.Get(suite.Context, suite.Address[0])
		suite.Require().NoError(err)
		suite.Require().Equal(types.Validator{
			Pubkey:    suite.Validator[0].Pubkey,
			Power:     power,
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Active,
			Locking:   sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewInt(1e18))),
		}, updated)

		newPower, err := suite.Keeper.ValidatorSet.Get(suite.Context, suite.Address[0])
		suite.Require().NoError(err)
		suite.Require().EqualValues(newPower, power)
	}

	// Add power to validator 0 again
	{
		var power uint64 = 20000
		err := suite.Keeper.Validators.Set(suite.Context, suite.Address[0], types.Validator{
			Pubkey:    suite.Validator[0].Pubkey,
			Power:     power,
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Active,
			Locking:   sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewInt(1e18).Mul(math.NewInt(2)))),
		})
		suite.Require().NoError(err)

		err = suite.Keeper.PowerRanking.Remove(suite.Context, collections.Join(uint64(10000), suite.Address[0]))
		suite.Require().NoError(err)

		err = suite.Keeper.PowerRanking.Set(suite.Context, collections.Join(uint64(power), suite.Address[0]))
		suite.Require().NoError(err)

		vs, err := suite.Keeper.EndBlocker(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Len(vs, 1)
		suite.Require().Equal(vs[0].PubKey.GetSecp256K1(), suite.Validator[0].Pubkey)
		suite.Require().EqualValues(vs[0].Power, power)

		updated, err := suite.Keeper.Validators.Get(suite.Context, suite.Address[0])
		suite.Require().NoError(err)
		suite.Require().Equal(types.Validator{
			Pubkey:    suite.Validator[0].Pubkey,
			Power:     power,
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Active,
			Locking:   sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewInt(1e18).Mul(math.NewInt(2)))),
		}, updated)

		newPower, err := suite.Keeper.ValidatorSet.Get(suite.Context, suite.Address[0])
		suite.Require().NoError(err)
		suite.Require().EqualValues(newPower, power)
	}

	// Add validator 2
	{
		var power uint64 = 10000
		err := suite.Keeper.Validators.Set(suite.Context, suite.Address[2], types.Validator{
			Pubkey:    suite.Validator[2].Pubkey,
			Power:     power,
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Pending,
			Locking:   sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewInt(1e18))),
		})
		suite.Require().NoError(err)

		err = suite.Keeper.PowerRanking.Set(suite.Context, collections.Join(power, suite.Address[2]))
		suite.Require().NoError(err)

		vs, err := suite.Keeper.EndBlocker(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Len(vs, 0)

		updated, err := suite.Keeper.Validators.Get(suite.Context, suite.Address[2])
		suite.Require().NoError(err)
		suite.Require().Equal(types.Validator{
			Pubkey:    suite.Validator[2].Pubkey,
			Power:     power,
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Pending,
			Locking:   sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewInt(1e18))),
		}, updated)

		exists, err := suite.Keeper.ValidatorSet.Has(suite.Context, suite.Address[2])
		suite.Require().NoError(err)
		suite.Require().False(exists)
	}

	// Unlock
	{
		var power uint64 = 9000
		var prevPower uint64 = 20000
		err := suite.Keeper.Validators.Set(suite.Context, suite.Address[0], types.Validator{
			Pubkey:    suite.Validator[0].Pubkey,
			Power:     power,
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Active,
			Locking:   sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewInt(1e18))),
		})
		suite.Require().NoError(err)

		err = suite.Keeper.PowerRanking.Remove(suite.Context, collections.Join(prevPower, suite.Address[0]))
		suite.Require().NoError(err)

		err = suite.Keeper.PowerRanking.Set(suite.Context, collections.Join(power, suite.Address[0]))
		suite.Require().NoError(err)

		vs, err := suite.Keeper.EndBlocker(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Len(vs, 2)

		{
			var v2Power uint64 = 10000
			suite.Require().Equal(vs[0].PubKey.GetSecp256K1(), suite.Validator[2].Pubkey)
			suite.Require().EqualValues(vs[0].Power, v2Power)

			updated, err := suite.Keeper.Validators.Get(suite.Context, suite.Address[2])
			suite.Require().NoError(err)
			suite.Require().Equal(types.Validator{
				Pubkey:    suite.Validator[2].Pubkey,
				Power:     v2Power,
				Reward:    math.ZeroInt(),
				GasReward: math.ZeroInt(),
				Status:    types.Active,
				Locking:   sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewInt(1e18))),
			}, updated)

			power, err := suite.Keeper.ValidatorSet.Get(suite.Context, suite.Address[2])
			suite.Require().NoError(err)
			suite.Require().Equal(power, v2Power)
		}

		{
			suite.Require().Equal(vs[1].PubKey.GetSecp256K1(), suite.Validator[0].Pubkey)
			suite.Require().EqualValues(vs[1].Power, 0)

			updated, err := suite.Keeper.Validators.Get(suite.Context, suite.Address[0])
			suite.Require().NoError(err)
			suite.Require().Equal(types.Validator{
				Pubkey:    suite.Validator[0].Pubkey,
				Power:     power,
				Reward:    math.ZeroInt(),
				GasReward: math.ZeroInt(),
				Status:    types.Pending,
				Locking:   sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewInt(1e18))),
			}, updated)

			exists, err := suite.Keeper.ValidatorSet.Has(suite.Context, suite.Address[0])
			suite.Require().NoError(err)
			suite.Require().False(exists)
		}
	}
}
