package keeper_test

import (
	"math/big"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (suite *KeeperTestSuite) TestClaim() {
	amount, ok := new(big.Int).SetString("100000000000000000000000", 10)
	suite.Require().True(ok)

	Addresses := []sdk.ConsAddress{
		sdk.ConsAddress(common.Hex2Bytes("f52a75aa5be8d8c9e3580ea6ba818e68de4fb76e")),
		sdk.ConsAddress(common.Hex2Bytes("108ca95b90e680f7e4374f911521941fe78b85ce")),
	}

	Recipients := []common.Address{
		common.HexToAddress("0xaf51c9cf596fe672785982ae061c9f4abf4f215e"),
		common.HexToAddress("0x672695d0a4e77222a86a4fd9e1288eca7d1f57be"),
	}

	Validators := []types.Validator{
		{
			Pubkey:    common.Hex2Bytes("03df9b92d37f4e3dec8ea95de7ab9f54879b978cebafd77a8528e7d832594a2af5"),
			Reward:    math.ZeroInt(),
			Power:     100000,
			GasReward: math.ZeroInt(),
			Status:    types.Pending,
			Locking: sdk.NewCoins(
				sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
			),
		},
		{
			Pubkey:    common.Hex2Bytes("03baf046326e0d1f48ad417b7336727e4454a286461ce1b2d01d50b3029468fd63"),
			Reward:    math.NewInt(1e9),
			Power:     110000,
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
	err := suite.Keeper.EthTxQueue.Set(suite.Context, types.EthTxQueue{})
	suite.Require().NoError(err)

	err = suite.Keeper.Claim(suite.Context, nil)
	suite.Require().NoError(err)

	err = suite.Keeper.Claim(suite.Context, []*goattypes.ClaimRequest{
		{Id: 0, Validator: common.Address(Addresses[0]), Recipient: Recipients[0]},
		{Id: 1, Validator: common.Address(Addresses[1]), Recipient: Recipients[1]},
	})
	suite.Require().NoError(err)
	for _, item := range Addresses {
		validator, err := suite.Keeper.Validators.Get(suite.Context, item)
		suite.Require().NoError(err)
		suite.Require().Equal(validator.GasReward, math.ZeroInt())
		suite.Require().Equal(validator.Reward, math.ZeroInt())
	}

	queue, err := suite.Keeper.EthTxQueue.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Len(queue.Unlocks, 0)
	suite.Require().Len(queue.Rewards, 2)

	suite.Require().Equal(queue.Rewards, []*types.Reward{
		{Id: 0, Recipient: Recipients[0].Bytes(), Goat: Validators[0].Reward, Gas: Validators[0].GasReward},
		{Id: 1, Recipient: Recipients[1].Bytes(), Goat: Validators[1].Reward, Gas: Validators[1].GasReward},
	})
}
