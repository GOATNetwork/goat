package keeper_test

import (
	"math/big"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (suite *KeeperTestSuite) TestLocks() {
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

	for address, token := range suite.Token {
		err := suite.Keeper.Tokens.Set(suite.Context, address, token)
		suite.Require().NoError(err)
	}

	err := suite.Keeper.Threshold.Set(suite.Context, types.Threshold{List: suite.Threshold})
	suite.Require().NoError(err)

	amount, ok := new(big.Int).SetString("100000000000000000000000", 10)
	suite.Require().True(ok)

	err = suite.Keeper.Lock(suite.Context, nil)
	suite.Require().NoError(err)

	{
		reqs := []*goattypes.LockRequest{
			{Validator: common.BytesToAddress(suite.Address[0]), Token: goattypes.GoatTokenContract, Amount: amount},
			{Validator: common.BytesToAddress(suite.Address[1]), Token: goattypes.GoatTokenContract, Amount: amount},
			{Validator: common.BytesToAddress(suite.Address[1]), Token: goattypes.GoatTokenContract, Amount: amount},
			{Validator: common.BytesToAddress(suite.Address[0]), Token: NativeToken, Amount: big.NewInt(1e18)},
			{Validator: common.BytesToAddress(suite.Address[1]), Token: NativeToken, Amount: big.NewInt(1e18)},
			{Validator: common.BytesToAddress(suite.Address[2]), Token: NativeToken, Amount: big.NewInt(1e18)},
			{Validator: common.BytesToAddress(suite.Address[3]), Token: NativeToken, Amount: big.NewInt(1e18)},
		}
		err = suite.Keeper.Lock(suite.Context, reqs)
		suite.Require().NoError(err)

		updated := []types.Validator{
			{
				Pubkey:    common.Hex2Bytes("03df9b92d37f4e3dec8ea95de7ab9f54879b978cebafd77a8528e7d832594a2af5"),
				Reward:    math.ZeroInt(),
				GasReward: math.ZeroInt(),
				Power:     110000,
				Status:    types.Pending,
				Locking: sdk.NewCoins(
					sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
					sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
				),
			},
			{
				Pubkey:    common.Hex2Bytes("03baf046326e0d1f48ad417b7336727e4454a286461ce1b2d01d50b3029468fd63"),
				Reward:    math.ZeroInt(),
				Power:     110000 + 110000 + 100000,
				GasReward: math.ZeroInt(),
				Status:    types.Active,
				Locking: sdk.NewCoins(
					sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18).Mul(math.NewInt(2))),
					sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount).Mul(math.NewInt(3))),
				),
			},
			{
				Pubkey:      common.Hex2Bytes("0236ea9615e7d0931ff24701a154d9e53ea4722b754710f22c8136058a7f251d73"),
				Power:       0,
				Reward:      math.ZeroInt(),
				GasReward:   math.ZeroInt(),
				Status:      types.Downgrade,
				JailedUntil: suite.Validator[2].JailedUntil,
				Locking:     sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18))),
			},
			{
				Pubkey:    common.Hex2Bytes("0293192d2e29b9713a5ad3a66cd7690a9232da4555c90d291d38612986421e8fa1"),
				Reward:    math.ZeroInt(),
				GasReward: math.ZeroInt(),
				Status:    types.Inactive,
				Locking:   sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18))),
			},
		}

		for idx, address := range suite.Address {
			validator, err := suite.Keeper.Validators.Get(suite.Context, address)
			suite.Require().NoError(err)
			suite.Require().Equal(validator, updated[idx], "idx %d", idx)
			if validator.Status == types.Active || validator.Status == types.Pending {
				for _, locking := range updated[idx].Locking {
					val, err := suite.Keeper.Locking.Get(suite.Context, collections.Join(locking.Denom, suite.Address[idx]))
					suite.Require().NoError(err)
					suite.Require().Equal(val, locking.Amount, "idx locking %d", idx)
				}
				has, err := suite.Keeper.PowerRanking.Has(suite.Context, collections.Join(validator.Power, suite.Address[idx]))
				suite.Require().NoError(err)
				suite.Require().True(has, "idx power %d", idx)
			} else {
				for _, locking := range updated[idx].Locking {
					has, err := suite.Keeper.Locking.Has(suite.Context, collections.Join(locking.Denom, suite.Address[idx]))
					suite.Require().NoError(err)
					suite.Require().False(has, "idx locking %d", idx)
				}
			}
		}
	}

	// unjail
	{
		newctx := suite.Context.WithBlockTime(time.Now().UTC().Add(time.Hour))
		reqs := []*goattypes.LockRequest{
			{Validator: common.BytesToAddress(suite.Address[2]), Token: goattypes.GoatTokenContract, Amount: amount},
		}
		err = suite.Keeper.Lock(newctx, reqs)
		suite.Require().NoError(err)

		validator, err := suite.Keeper.Validators.Get(suite.Context, suite.Address[2])
		suite.Require().NoError(err)
		suite.Require().Equal(validator, types.Validator{
			Pubkey:      common.Hex2Bytes("0236ea9615e7d0931ff24701a154d9e53ea4722b754710f22c8136058a7f251d73"),
			Power:       110000,
			Reward:      math.ZeroInt(),
			GasReward:   math.ZeroInt(),
			Status:      types.Pending,
			JailedUntil: suite.Validator[2].JailedUntil,
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
				sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
			),
		})

		for _, locking := range validator.Locking {
			val, err := suite.Keeper.Locking.Get(suite.Context, collections.Join(locking.Denom, suite.Address[2]))
			suite.Require().NoError(err)
			suite.Require().Equal(val, locking.Amount, "unjail locking")
		}
		has, err := suite.Keeper.PowerRanking.Has(suite.Context, collections.Join(validator.Power, suite.Address[2]))
		suite.Require().NoError(err)
		suite.Require().True(has)
	}
}
