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

func (suite *KeeperTestSuite) TestUnlock() {
	now := time.Now().UTC()
	newctx := suite.Context.WithBlockTime(now)

	param, err := suite.Keeper.Params.Get(newctx)
	suite.Require().NoError(err)

	for address, token := range suite.Token {
		err := suite.Keeper.Tokens.Set(newctx, address, token)
		suite.Require().NoError(err)
	}

	err = suite.Keeper.Threshold.Set(newctx, types.Threshold{List: suite.Threshold})
	suite.Require().NoError(err)

	amount, ok := new(big.Int).SetString("100000000000000000000000", 10)
	suite.Require().True(ok)

	Validators := []types.Validator{
		// 0
		{
			Pubkey:    common.Hex2Bytes("0242d59fb617d2cb150966bf42c3c7e943f7194f82e41b68a12124d6c2d2f69ae2"),
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Power:     11000,
			Status:    types.Pending,
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18+1e17)),
			),
		},
		// 1
		{
			Pubkey:    common.Hex2Bytes("03a43c6d25cd52c8557e89a047f627d39315f1940f147ff8f191f2d058f33c1211"),
			Reward:    math.ZeroInt(),
			Power:     110000,
			GasReward: math.ZeroInt(),
			Status:    types.Active,
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
				sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
			),
		},
		// 2
		{
			Pubkey:    common.Hex2Bytes("02ae8db8dd5cc564b42d86ed353242ad0e0cfd1d95b65a172e445e3554c9bde9ae"),
			Reward:    math.ZeroInt(),
			Power:     200010,
			GasReward: math.ZeroInt(),
			Status:    types.Active,
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
				sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
				sdk.NewCoin(TestTokenDenom, math.NewIntFromUint64(1e19)),
			),
		},
		// 3
		{
			Pubkey:      common.Hex2Bytes("02269670f94151626e04ecf71792cf7ff89bd5bd40b07a408c2a6a11e10085ef76"),
			Reward:      math.ZeroInt(),
			GasReward:   math.ZeroInt(),
			Status:      types.Downgrade,
			JailedUntil: now.Add(time.Hour),
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e19)),
			),
		},
		// 4
		{
			Pubkey:    common.Hex2Bytes("037f812b3c1da7682afff4774ca7606405d81329a920dbc86c49762205489143d1"),
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Inactive,
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
			),
		},
		// 5
		{
			Pubkey:    common.Hex2Bytes("034d96d43803ff6b1579587b8a515066f7acebf5041a09706ee9f38a69f8f67ff6"),
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Tombstoned,
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
			),
		},
	}

	Addresses := []sdk.ConsAddress{
		sdk.ConsAddress(common.Hex2Bytes("c7bde79e58bdcfc8f058749b48feb2558c9c4adb")),
		sdk.ConsAddress(common.Hex2Bytes("c387bcb76c985d4873f04698eeec98fc5e3d4a7d")),
		sdk.ConsAddress(common.Hex2Bytes("044d1624630b729ab072fb3338b762dbd34e41ab")),
		sdk.ConsAddress(common.Hex2Bytes("d466208fb91604fc16170712ad186caa34c9b005")),
		sdk.ConsAddress(common.Hex2Bytes("d15780226764eadae9911acc7c2f6fdb3a342680")),
		sdk.ConsAddress(common.Hex2Bytes("359bc904aa5800469eb9b7a18f8456f47f62415d")),
	}

	Recipients := []common.Address{
		common.HexToAddress("0xbdccd1119e70d7e78a2bf03e0db494ad036cc32e"),
		common.HexToAddress("0x9fb2cd73e2a5aa5d880f2ff7695e37292c9bec84"),
		common.HexToAddress("0x77a1f9a9da582ae33fe73472380be99a802135b0"),
		common.HexToAddress("0xf8cd0c0a6f6568a7bf015134ffc5850cb9f53a57"),
		common.HexToAddress("0xb12ddf42601dfa555077fcd6eb8d3ac7e875f047"),
		common.HexToAddress("0x6893602d68a4f6ce6949ab51b9651e7cda0bf52d"),
	}

	for idx, validator := range Validators {
		err := suite.Keeper.Validators.Set(newctx, Addresses[idx], validator)
		suite.Require().NoError(err)

		if validator.Status == types.Active || validator.Status == types.Pending {
			for _, locking := range validator.Locking {
				err = suite.Keeper.Locking.Set(newctx,
					collections.Join(locking.Denom, Addresses[idx]), locking.Amount)
				suite.Require().NoError(err)
			}
			if validator.Power > 0 {
				err = suite.Keeper.PowerRanking.Set(newctx, collections.Join(validator.Power, Addresses[idx]))
				suite.Require().NoError(err)
			}
		}
	}

	// no unlock requests
	err = suite.Keeper.Unlock(newctx, nil)
	suite.Require().NoError(err)

	reqs := []*goattypes.UnlockRequest{
		// unlock
		{Id: 0, Validator: common.Address(Addresses[0]), Recipient: Recipients[0], Token: common.Address{}, Amount: big.NewInt(1e17)},
		// exit
		{Id: 1, Validator: common.Address(Addresses[0]), Recipient: Recipients[0], Token: common.Address{}, Amount: big.NewInt(1e4)},
		// exit and was slashed
		{Id: 2, Validator: common.Address(Addresses[1]), Recipient: Recipients[1], Token: common.Address{}, Amount: new(big.Int).SetUint64(1e19)},
		// unlock
		{Id: 3, Validator: common.Address(Addresses[2]), Recipient: Recipients[2], Token: TestToken, Amount: new(big.Int).SetUint64(1e19)},
		// inactive unlock
		{Id: 4, Validator: common.Address(Addresses[4]), Recipient: Recipients[4], Token: common.Address{}, Amount: new(big.Int).SetUint64(1e18)},
		// tombstoned unlock
		{Id: 5, Validator: common.Address(Addresses[5]), Recipient: Recipients[5], Token: common.Address{}, Amount: new(big.Int).SetUint64(1e18)},
		// downgrade unlock
		{Id: 6, Validator: common.Address(Addresses[3]), Recipient: Recipients[3], Token: common.Address{}, Amount: new(big.Int).SetUint64(1e4)},
	}

	err = suite.Keeper.Unlock(newctx, reqs)
	suite.Require().NoError(err)

	gotUnlocks := make(map[int64][]*types.Unlock)
	unlockIter, err := suite.Keeper.UnlockQueue.Iterate(newctx, nil)
	suite.Require().NoError(err)
	for ; unlockIter.Valid(); unlockIter.Next() {
		kv, err := unlockIter.KeyValue()
		suite.Require().NoError(err)
		gotUnlocks[kv.Key.Unix()] = kv.Value.Unlocks
	}

	exitTime, unlockTime := now.Add(param.ExitingDuration).Unix(), now.Add(param.UnlockDuration).Unix()
	expectedUnlocks := map[int64][]*types.Unlock{
		unlockTime: {
			{Id: reqs[0].Id, Token: reqs[0].Token[:], Recipient: reqs[0].Recipient[:], Amount: math.NewIntFromBigInt(reqs[0].Amount)},
			{Id: reqs[3].Id, Token: reqs[3].Token[:], Recipient: reqs[3].Recipient[:], Amount: math.NewIntFromBigInt(reqs[3].Amount)},
			{Id: reqs[6].Id, Token: reqs[6].Token[:], Recipient: reqs[6].Recipient[:], Amount: math.NewIntFromBigInt(reqs[6].Amount)},
		},
		exitTime: {
			{Id: reqs[1].Id, Token: reqs[1].Token[:], Recipient: reqs[1].Recipient[:], Amount: math.NewIntFromBigInt(reqs[1].Amount)},
			{Id: reqs[2].Id, Token: reqs[2].Token[:], Recipient: reqs[2].Recipient[:], Amount: math.NewInt(1e18)},
			{Id: reqs[4].Id, Token: reqs[4].Token[:], Recipient: reqs[4].Recipient[:], Amount: math.NewInt(1e18)},
			{Id: reqs[5].Id, Token: reqs[5].Token[:], Recipient: reqs[5].Recipient[:], Amount: math.NewInt(1e18)},
		},
	}

	suite.Require().Equal(len(expectedUnlocks), len(gotUnlocks))
	for key, values := range gotUnlocks {
		suite.Require().Equal(len(expectedUnlocks[key]), len(values))
		for idx, item := range values {
			suite.Require().Equal(expectedUnlocks[key][idx], item, "unlock %d: %d", key, idx)
		}
	}

	Updated := []types.Validator{
		// 0
		{
			Pubkey:    common.Hex2Bytes("0242d59fb617d2cb150966bf42c3c7e943f7194f82e41b68a12124d6c2d2f69ae2"),
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Power:     0,
			Status:    types.Inactive,
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18-1e4)),
			),
		},
		// 1
		{
			Pubkey:    common.Hex2Bytes("03a43c6d25cd52c8557e89a047f627d39315f1940f147ff8f191f2d058f33c1211"),
			Reward:    math.ZeroInt(),
			Power:     0,
			GasReward: math.ZeroInt(),
			Status:    types.Inactive,
			Locking: sdk.NewCoins(
				sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
			),
		},
		// 2
		{
			Pubkey:    common.Hex2Bytes("02ae8db8dd5cc564b42d86ed353242ad0e0cfd1d95b65a172e445e3554c9bde9ae"),
			Reward:    math.ZeroInt(),
			Power:     200000,
			GasReward: math.ZeroInt(),
			Status:    types.Active,
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
				sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
			),
		},
		// 3
		{
			Pubkey:      common.Hex2Bytes("02269670f94151626e04ecf71792cf7ff89bd5bd40b07a408c2a6a11e10085ef76"),
			Reward:      math.ZeroInt(),
			GasReward:   math.ZeroInt(),
			Status:      types.Downgrade,
			JailedUntil: now.Add(time.Hour),
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e19-1e4)),
			),
		},
		// 4
		{
			Pubkey:    common.Hex2Bytes("037f812b3c1da7682afff4774ca7606405d81329a920dbc86c49762205489143d1"),
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Inactive,
			Locking:   nil,
		},
		// 5
		{
			Pubkey:    common.Hex2Bytes("034d96d43803ff6b1579587b8a515066f7acebf5041a09706ee9f38a69f8f67ff6"),
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Tombstoned,
			Locking:   nil,
		},
	}

	for i, address := range Addresses {
		validator, err := suite.Keeper.Validators.Get(newctx, address)
		suite.Require().NoError(err)
		suite.Require().Equal(Updated[i], validator, i)
	}

	lockingIter, err := suite.Keeper.Locking.Iterate(newctx, nil)
	suite.Require().NoError(err)
	expectedLocking := map[string]sdk.Coins{}
	for ; lockingIter.Valid(); lockingIter.Next() {
		kv, err := lockingIter.KeyValue()
		suite.Require().NoError(err)
		address := string(kv.Key.K2())
		expectedLocking[address] = expectedLocking[address].Add(sdk.NewCoin(kv.Key.K1(), kv.Value))
	}

	suite.Require().Equal(len(expectedLocking), 2)
	suite.Require().Equal(expectedLocking, map[string]sdk.Coins{
		string(Addresses[2]): sdk.NewCoins(
			sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
			sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
		),
		string(Addresses[3]): sdk.NewCoins(sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e19-1e4))),
	})

	powerIter, err := suite.Keeper.PowerRanking.Iterate(newctx, nil)
	suite.Require().NoError(err)
	powers, err := powerIter.Keys()
	suite.Require().NoError(err)
	suite.Require().Len(powers, 1)
	suite.Require().EqualValues(powers[0].K1(), 200000)
	suite.Require().Equal(powers[0].K2(), Addresses[2])
}

func (suite *KeeperTestSuite) TestDequeueMatureUnlocks() {
	now := time.Now().UTC()

	unlocks := []*types.Unlock{}
	for i := int64(1); i < 10; i++ {
		d := now.Add(time.Minute * time.Duration(i-10))
		err := suite.Keeper.UnlockQueue.Set(suite.Context, d, types.Unlocks{
			Unlocks: []*types.Unlock{{Id: uint64(i), Amount: math.NewInt(i)}},
		})
		suite.Require().NoError(err)
		unlocks = append(unlocks, &types.Unlock{Id: uint64(i), Amount: math.NewInt(i)})
	}

	err := suite.Keeper.EthTxQueue.Set(suite.Context, types.EthTxQueue{})
	suite.Require().NoError(err)

	newctx := suite.Context.WithBlockTime(now.Add(time.Minute * -20))
	err = suite.Keeper.DequeueMatureUnlocks(newctx)
	suite.Require().NoError(err)
	queue, err := suite.Keeper.EthTxQueue.Get(newctx)
	suite.Require().NoError(err)
	suite.Require().Len(queue.Rewards, 0)
	suite.Require().Len(queue.Unlocks, 0)

	newctx = suite.Context.WithBlockTime(now.Add(time.Minute * -5))
	err = suite.Keeper.DequeueMatureUnlocks(newctx)
	suite.Require().NoError(err)

	queue, err = suite.Keeper.EthTxQueue.Get(newctx)
	suite.Require().NoError(err)
	suite.Require().Len(queue.Rewards, 0)
	suite.Require().Len(queue.Unlocks, 5)
	suite.Require().Equal(queue.Unlocks, unlocks[:5])

	newctx = suite.Context.WithBlockTime(now)
	err = suite.Keeper.DequeueMatureUnlocks(newctx)
	suite.Require().NoError(err)

	queue, err = suite.Keeper.EthTxQueue.Get(newctx)
	suite.Require().NoError(err)
	suite.Require().Len(queue.Rewards, 0)
	suite.Require().Len(queue.Unlocks, 9)
	suite.Require().Equal(queue.Unlocks, unlocks)
}
