package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	keepertest "github.com/goatnetwork/goat/testutil/keeper"
	"github.com/goatnetwork/goat/testutil/mock"
	"github.com/goatnetwork/goat/x/locking/keeper"
	"github.com/goatnetwork/goat/x/locking/types"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type KeeperTestSuite struct {
	suite.Suite
	Ctrl      *gomock.Controller
	Account   *mock.MockAccountKeeper
	Keeper    keeper.Keeper
	Context   sdk.Context
	Param     types.Params
	Validator []types.Validator
	Address   []sdk.ConsAddress
	Token     map[string]types.Token
	Threshold sdk.Coins
}

func TestKeeper(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

var (
	TestToken   = common.HexToAddress("cafebeefcafebeefcafebeefcafebeefcafebeef")
	NativeToken = common.Address{}
)

var (
	GoatToekenDenom  = types.TokenDenom(goattypes.GoatTokenContract)
	NativeTokenDenom = types.TokenDenom(common.Address{})
	TestTokenDenom   = types.TokenDenom(TestToken)
)

func (suite *KeeperTestSuite) SetupTest() {
	ctl := gomock.NewController(suite.T())
	accountKeeper := mock.NewMockAccountKeeper(ctl)

	suite.Keeper, suite.Context = keepertest.LockingKeeper(suite.T(), accountKeeper)

	suite.Account = accountKeeper
	suite.Param = types.DefaultParams()
	suite.Ctrl = ctl

	amount, ok := new(big.Int).SetString("100000000000000000000000", 10)
	suite.Require().True(ok)
	suite.Token = map[string]types.Token{
		NativeTokenDenom: {Weight: 1e4, Threshold: math.NewIntFromUint64(1e18)},
		GoatToekenDenom:  {Weight: 1, Threshold: math.NewIntFromBigInt(amount)},
		TestTokenDenom:   {Weight: 1, Threshold: math.ZeroInt()},
	}

	suite.Threshold = sdk.NewCoins(
		sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
		sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
	)

	suite.Validator = []types.Validator{
		{
			Pubkey:    common.Hex2Bytes("03df9b92d37f4e3dec8ea95de7ab9f54879b978cebafd77a8528e7d832594a2af5"),
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Pending,
		},
		{
			Pubkey:    common.Hex2Bytes("03baf046326e0d1f48ad417b7336727e4454a286461ce1b2d01d50b3029468fd63"),
			Reward:    math.ZeroInt(),
			Power:     110000,
			GasReward: math.ZeroInt(),
			Status:    types.Active,
			Locking: sdk.NewCoins(
				sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
				sdk.NewCoin(GoatToekenDenom, math.NewIntFromBigInt(amount)),
			),
		},
		{
			Pubkey:      common.Hex2Bytes("0236ea9615e7d0931ff24701a154d9e53ea4722b754710f22c8136058a7f251d73"),
			Reward:      math.ZeroInt(),
			GasReward:   math.ZeroInt(),
			Status:      types.Downgrade,
			JailedUntil: time.Now().UTC(),
		},
		{
			Pubkey:    common.Hex2Bytes("0293192d2e29b9713a5ad3a66cd7690a9232da4555c90d291d38612986421e8fa1"),
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    types.Inactive,
		},
	}

	suite.Address = []sdk.ConsAddress{
		sdk.ConsAddress(common.Hex2Bytes("f52a75aa5be8d8c9e3580ea6ba818e68de4fb76e")),
		sdk.ConsAddress(common.Hex2Bytes("108ca95b90e680f7e4374f911521941fe78b85ce")),
		sdk.ConsAddress(common.Hex2Bytes("5f2354785046bc2d3b65681af66ee4ccc34b95f7")),
		sdk.ConsAddress(common.Hex2Bytes("34fc7515049f5bd9600b9116f90dd668762ab7ac")),
	}
}

func (suite *KeeperTestSuite) TearDownSuite() {
	suite.Ctrl.Finish()
}
