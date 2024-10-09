package keeper_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	keepertest "github.com/goatnetwork/goat/testutil/keeper"
	"github.com/goatnetwork/goat/x/bitcoin/keeper"
	"github.com/goatnetwork/goat/x/bitcoin/mock"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestKeeper(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite
	Ctrl          *gomock.Controller
	Keeper        keeper.Keeper
	Context       sdk.Context
	RelayerKeeper *mock.MockRelayerKeeper
	Param         types.Params
	TestKey       relayer.PublicKey
}

func (suite *KeeperTestSuite) SetupTest() {
	ctl := gomock.NewController(suite.T())
	relayerKeeper := mock.NewMockRelayerKeeper(ctl)

	keeper, ctx, _ := keepertest.BitcoinKeeper(suite.T(), relayerKeeper)

	suite.Keeper = keeper
	suite.Context = ctx
	suite.Param = types.DefaultParams()
	suite.RelayerKeeper = relayerKeeper

	suite.TestKey = relayer.PublicKey{Key: &relayer.PublicKey_Secp256K1{
		Secp256K1: common.Hex2Bytes("0383560def84048edefe637d0119a4428dd12a42765a118b2bf77984057633c50e"),
	}}

	suite.Ctrl = ctl
}

func (suite *KeeperTestSuite) TearDownSuite() {
	suite.Ctrl.Finish()
}

func (suite *KeeperTestSuite) TestNewPubkey() {
	suite.RelayerKeeper.EXPECT().AddNewKey(suite.Context, relayer.EncodePublicKey(&suite.TestKey)).Return(nil)

	suite.Require().NoError(suite.Keeper.NewPubkey(suite.Context, &suite.TestKey))

	pubkey, err := suite.Keeper.Pubkey.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().True(pubkey.Equal(suite.TestKey))
}

func (suite *KeeperTestSuite) TestVerifyDeposit() {
	suite.RelayerKeeper.EXPECT().HasPubkey(suite.Context, relayer.EncodePublicKey(&suite.TestKey)).Return(true, nil)

	evmAddress := common.HexToAddress("0xbC122aEc3EdD80433dfE3c708b2E549B5A7Ab96E")
	blockHash, _ := chainhash.NewHashFromStr("38fb77a25662f9eda5abef8a407ba45e8c3374b5a0724cfa9762f1f9cbf627e2")
	const height = 102
	suite.Require().NoError(suite.Keeper.BlockHashes.Set(suite.Context, 102, blockHash[:]))

	header := common.Hex2Bytes("00000020451119ce15cd42ceb7a00c2ef9843aa613a69f19f7b4fc483f0f28b099c54d1bc8df397f2235b299f7ca89e10f789e598f53dc89789b8a047bc78238ef4bd4daf9f8e466ffff7f2000000000")

	const txIndex = 1
	const txOutput = 1
	const amount = 1e8
	txid, _ := chainhash.NewHashFromStr("9a31c75d3676059c7483d29f12082b4df9e396df5c22612e50fa97b94bbf532c")
	tx := common.Hex2Bytes("0200000001e15e44fc827b0e1a3178b6e07f67e8339faae54e4241e5fa5c1ed61786a84bda0000000000fdffffff020dc74c0001000000225120098ad136e9ed8106af7c1b6b4934011f320b30f6e18871917e0d6fb1bdcb5d1400e1f50500000000220020f7608234b4bc67678cc5498dfe7db5dfda221d3ff669f1d9ee89fbcf14d104f366000000")

	proof := common.Hex2Bytes("4930ac654c3c2e487fcc2106a51ecaaf4188093686dfffcfe880798044aadc02")

	res, err := suite.Keeper.VerifyDeposit(suite.Context,
		map[uint64][]byte{height: header},
		&types.Deposit{
			Version:           0,
			BlockNumber:       height,
			TxIndex:           txIndex,
			NoWitnessTx:       tx,
			OutputIndex:       txOutput,
			IntermediateProof: proof,
			EvmAddress:        evmAddress.Bytes(),
			RelayerPubkey:     &suite.TestKey,
		})
	suite.Require().NoError(err)

	suite.Require().Equal(res, &types.DepositExecReceipt{
		Address: evmAddress.Bytes(),
		Txid:    txid[:],
		Txout:   txOutput,
		Amount:  amount,
	})
}
