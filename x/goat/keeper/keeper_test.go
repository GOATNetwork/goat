package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/goatnetwork/goat/testutil/keeper"
	"github.com/goatnetwork/goat/testutil/mock"
	"github.com/goatnetwork/goat/x/goat/keeper"
	"github.com/goatnetwork/goat/x/goat/types"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type KeeperTestSuite struct {
	suite.Suite
	Ctrl    *gomock.Controller
	Account *mock.MockAccountKeeper
	Bitcoin *mock.MockBitcoinKeeper
	Locking *mock.MockLockingKeeper
	Relayer *mock.MockRelayerKeeper
	Keeper  keeper.Keeper
	Context sdk.Context
	Param   types.Params
}

func TestKeeper(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	ctl := gomock.NewController(suite.T())
	accountKeeper := mock.NewMockAccountKeeper(ctl)
	bitcoinKeeper := mock.NewMockBitcoinKeeper(ctl)
	lockingKeeper := mock.NewMockLockingKeeper(ctl)
	relayerKeeper := mock.NewMockRelayerKeeper(ctl)
	ethClient := mock.NewMockEngineClient(ctl)

	keeper, ctx, _ := keepertest.GoatKeeper(suite.T(),
		bitcoinKeeper, lockingKeeper, relayerKeeper, accountKeeper, ethClient)

	suite.Keeper = keeper
	suite.Account = accountKeeper
	suite.Bitcoin = bitcoinKeeper
	suite.Locking = lockingKeeper
	suite.Relayer = relayerKeeper

	suite.Context = ctx
	suite.Param = types.DefaultParams()
	suite.Ctrl = ctl
}

func (suite *KeeperTestSuite) TearDownSuite() {
	suite.Ctrl.Finish()
}
