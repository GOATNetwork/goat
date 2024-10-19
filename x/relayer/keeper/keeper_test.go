package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	keepertest "github.com/goatnetwork/goat/testutil/keeper"
	"github.com/goatnetwork/goat/testutil/mock"
	"github.com/goatnetwork/goat/x/relayer/keeper"
	"github.com/goatnetwork/goat/x/relayer/types"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestKeeper(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite
	Ctrl        *gomock.Controller
	Account     *mock.MockAccountKeeper
	VoteMsgMock *mock.MockIVoteMsg
	Keeper      keeper.Keeper
	Context     sdk.Context
	Param       types.Params
	TestKey     types.PublicKey
	Voters      []types.Voter
	VoterKeys   []VoterWithKey
}

type VoterWithKey struct {
	Address string
	TxKey   []byte
	VoteKey []byte
}

func (suite *KeeperTestSuite) SetupTest() {
	ctl := gomock.NewController(suite.T())
	accountKeeper := mock.NewMockAccountKeeper(ctl)

	keeper, ctx, _ := keepertest.RelayerKeeper(suite.T(), accountKeeper)

	suite.Keeper = keeper
	suite.Context = ctx
	suite.Param = types.DefaultParams()
	suite.Account = accountKeeper
	suite.VoteMsgMock = mock.NewMockIVoteMsg(ctl)

	suite.TestKey = types.PublicKey{Key: &types.PublicKey_Secp256K1{
		Secp256K1: common.Hex2Bytes("0383560def84048edefe637d0119a4428dd12a42765a118b2bf77984057633c50e"),
	}}

	suite.Ctrl = ctl

	suite.Voters = []types.Voter{
		{
			// txPubkey 030ccc6707050ff1f30da3ae4dad458c0b3023372f49d07b6262de66efe1e39e1c
			Address: common.Hex2Bytes("6c76e7d2bf7e08fa16389dcd5d098f2ff18f9dca"),
			VoteKey: common.Hex2Bytes("931e41003cdbb46fa624f0636bceb743ff4f16240e3c175eeafa75fb29d2f3e06e9cd0515840af9d5899621ccb7e16a21251431e69361702478f584b4314e2178443b802302f22be9e38d9bb15abc692fa560ab3e2724f03d27c611db1227552"),
			Status:  types.VOTER_STATUS_ACTIVATED,
			Height:  100,
		},
		{
			// txPubkey 02f8c5fa7583cc77b5e75723d277538eb4c08dd80e499478032816515c9869e666
			Address: common.Hex2Bytes("1f7a29e5f1b99c6ea2bc6b4ec90a7654eaa33598"),
			VoteKey: common.Hex2Bytes("a05ad80960006177d2d5f424bcb8b68c0089cbb50de7482d1658ce3df4f01261db7df3c67cb4b459b441479e0b6aae8101a776335615749a61646f32fdd5812031f0318a63c7a13963b210dd64a68bc9abd530746190160aa781295e850ec1f5"),
			Status:  types.VOTER_STATUS_ACTIVATED,
			Height:  100,
		},
		{
			// txPubkey 029089e07f96cbb19e9cb95992b633f45029eaf5a3481f989d750628225ce997ed
			Address: common.Hex2Bytes("f3959b9137f61a69c6c33f91a4e6fb780c771164"),
			VoteKey: common.Hex2Bytes("90856b5bd5b3916c5163f2afc27ca7bed1dbd743f7c348bf345f603f2f23e428dadd11c675a75da600413abbb4334fb4061544fa0736f6f3554a1827bc58e820eb02433bb8c025f9c0f6814c38e4674fc86b024c5304657a714dfada3a2488d9"),
			Status:  types.VOTER_STATUS_PENDING,
		},
		{
			// txPubkey 02d9e0fa8fc514f85131e748fa8b58b88ca1e820cfc1caa5e3124753c724c3f515
			Address: common.Hex2Bytes("d45b120a2eb7c62f319d9c44e6caaad1752f5773"),
			VoteKey: common.Hex2Bytes("91fcc71f6c7b6922d7cf6e1d6f2a13e72e13c78cf72e8ef1e24d8f6bea76b4bc040dfb864a541981129c55162e77829f0c346091d45a2cc5a5fc620ca4ab63b27a1192a2715211ae4a59c21044e60d947a59cb79a03e7ffc5fb81f76431436a1"),
			Status:  types.VOTER_STATUS_PENDING,
		},
	}

	suite.VoterKeys = []VoterWithKey{
		{Address: "goat1d3mw054l0cy0593cnhx46zv09lccl8w2crw529", TxKey: common.Hex2Bytes("dcfce7f45069b8d81af211ea61de390b67dd86e78b1d9c4f3d8701c22a82fd29"), VoteKey: common.Hex2Bytes("37d02e811d7147a86b77c95dcb0ce8f249da5c31f28dbfd0e0a55f37341d0249")},
		{Address: "goat1raazne03hxwxag4udd8vjznk2n42xdvc4tfsnw", TxKey: common.Hex2Bytes("94a42c7d90d0d82b895d3b750877a2cd3d97c5fcb3a69fa789fc6ac7a98b1a96"), VoteKey: common.Hex2Bytes("27845014abb73960d2ed55c93feee780d8f4aea7bcb140f0e35bb57621b118e9")},
		{Address: "goat17w2ehyfh7cdxn3kr87g6fehm0qx8wytyevyllx", TxKey: common.Hex2Bytes("c2195372b1041fc8f9dea7c2b19b7014f4afe29953b08a14be038b10435cc255"), VoteKey: common.Hex2Bytes("490ce61c00399491e4e0d2ecb729e0d24fb7a3cd64d8a85387329598cfc9e3fe")},
		{Address: "goat163d3yz3wklrz7vvan3zwdj42696j74mnvj322d", TxKey: common.Hex2Bytes("39e7d9b0b9640ede2ab7b77c64d04e100357df5e6af9aeb1ec5d80c858b2b2f2"), VoteKey: common.Hex2Bytes("14a0d0cc11c712d476cc6c8184364dd65dfad3413a2d1b57aed01fb7e52051f9")},
	}

	for _, voter := range suite.Voters {
		if voter.Status != types.VOTER_STATUS_ACTIVATED {
			continue
		}
		address, err := suite.Keeper.AddrCodec.BytesToString(voter.Address)
		suite.Require().NoError(err)
		err = suite.Keeper.Voters.Set(suite.Context, address, voter)
		suite.Require().NoError(err)
	}
}

func (suite *KeeperTestSuite) TearDownSuite() {
	suite.Ctrl.Finish()
}

func (suite *KeeperTestSuite) TestUpdateRandao() {
	err := suite.Keeper.Randao.Set(suite.Context,
		common.Hex2Bytes("631ce70cc1e6818ab1b0dd4c7d8c9af4b7a893ff9aed518a886f1c3c9823a970"))
	suite.Require().NoError(err)

	suite.VoteMsgMock.EXPECT().GetVote().Return(&types.Votes{
		Signature: common.Hex2Bytes("1d3ae4583cce4fe713e39782131ef6a61c70e58a9952d10fccc1f1f4f74b0b6e7ef59be37e08ca6ecb5666ad0de94132"),
	})

	err = suite.Keeper.UpdateRandao(suite.Context, suite.VoteMsgMock)
	suite.Require().NoError(err)

	updated, err := suite.Keeper.Randao.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(updated, common.Hex2Bytes("a9c74af54809f927cec82a77d5e23bc24c213480c6cf8dc656e70fffec7308ea"))
}

func (suite *KeeperTestSuite) TestNewPubkey() {
	encoded := types.EncodePublicKey(&suite.TestKey)
	exists, err := suite.Keeper.HasPubkey(suite.Context, encoded)
	suite.Require().NoError(err)
	suite.Require().False(exists)

	err = suite.Keeper.AddNewKey(suite.Context, encoded)
	suite.Require().NoError(err)
	exists, err = suite.Keeper.HasPubkey(suite.Context, encoded)
	suite.Require().NoError(err)
	suite.Require().True(exists)
}

func (suite *KeeperTestSuite) TestRelayerSeqAndProposer() {
	relayer := types.Relayer{
		Proposer:         "goat1vpa50kdf63rfuurelvzsmpp99ffmvuc20h68qh",
		Voters:           []string{},
		LastElected:      time.Now().UTC(),
		ProposerAccepted: true,
	}

	err := suite.Keeper.Relayer.Set(suite.Context, relayer)
	suite.Require().NoError(err)

	address, err := suite.Keeper.GetCurrentProposer(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(address.Bytes(), common.Hex2Bytes("607b47d9a9d4469e7079fb050d84252a53b6730a"))

	err = suite.Keeper.SetProposalSeq(suite.Context, 100)
	suite.Require().NoError(err)

	seq, err := suite.Keeper.Sequence.Peek(suite.Context)
	suite.Require().NoError(err)
	suite.Require().EqualValues(seq, 100)
}
