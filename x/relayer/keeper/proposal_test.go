package keeper_test

import (
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/relayer/types"
	"github.com/kelindar/bitmap"
)

func (suite *KeeperTestSuite) TestVerifyProposal() {
	relayer := types.Relayer{
		Proposer:         "goat1d3mw054l0cy0593cnhx46zv09lccl8w2crw529",
		Voters:           []string{"goat1raazne03hxwxag4udd8vjznk2n42xdvc4tfsnw"},
		LastElected:      time.Now().UTC(),
		ProposerAccepted: false,
	}

	err := suite.Keeper.Relayer.Set(suite.Context, relayer)
	suite.Require().NoError(err)

	sdkctx := sdktypes.UnwrapSDKContext(suite.Context)
	reqMethod, reqSigDoc := "testing", []byte("testing sig doc")

	sigdoc := types.VoteSignDoc(reqMethod, sdkctx.ChainID(), relayer.Proposer, 0, relayer.Epoch, reqSigDoc)
	var sigs [][]byte
	for i := range 2 {
		sk := new(goatcrypto.PrivateKey).Deserialize(suite.VoterKeys[i].VoteKey)
		suite.Require().NotNil(sk)
		sigs = append(sigs, goatcrypto.Sign(sk, sigdoc))
	}

	aggsig, err := goatcrypto.AggregateSignatures(sigs)
	suite.Require().NoError(err)

	reqMsg := &types.Votes{
		Sequence:  0,
		Epoch:     relayer.Epoch,
		Voters:    common.Hex2Bytes("0100000000000000"),
		Signature: aggsig,
	}

	suite.VoteMsgMock.EXPECT().GetProposer().AnyTimes().Return(relayer.Proposer)
	suite.VoteMsgMock.EXPECT().GetVote().AnyTimes().Return(reqMsg)
	suite.VoteMsgMock.EXPECT().MethodName().AnyTimes().Return(reqMethod)
	suite.VoteMsgMock.EXPECT().VoteSigDoc().AnyTimes().Return(reqSigDoc)

	seq, err := suite.Keeper.VerifyProposal(suite.Context, suite.VoteMsgMock, func([]byte) error { return nil })
	suite.Require().NoError(err)
	suite.Require().EqualValues(seq, 0)

	relayer.ProposerAccepted = true
	newRelayer, err := suite.Keeper.Relayer.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(relayer, newRelayer)
}

func (suite *KeeperTestSuite) TestVerifyProposal_invalidVoterBitmap() {
	voters := make([]string, len(suite.VoterKeys)-1)
	for i := range voters {
		voters[i] = suite.VoterKeys[i+1].Address
	}
	relayer := types.Relayer{
		Proposer:         suite.VoterKeys[0].Address,
		Voters:           voters,
		LastElected:      time.Now().UTC(),
		ProposerAccepted: false,
	}
	suite.Require().NoError(suite.Keeper.Relayer.Set(suite.Context, relayer))
	th := relayer.Threshold()
	suite.Require().Equal(3, th)

	sdkctx := sdktypes.UnwrapSDKContext(suite.Context)
	reqMethod, reqSigDoc := "testing", []byte("testing sig doc")
	sigdoc := types.VoteSignDoc(reqMethod, sdkctx.ChainID(), relayer.Proposer, 0, relayer.Epoch, reqSigDoc)

	proposerSK := new(goatcrypto.PrivateKey).Deserialize(suite.VoterKeys[0].VoteKey)
	suite.Require().NotNil(proposerSK)
	sig := goatcrypto.Sign(proposerSK, sigdoc)

	var b bitmap.Bitmap
	for j := 0; j < th-1; j++ {
		b.Set(uint32(len(voters) + j)) // OUT-OF-RANGE bits only
	}
	voteBytes := b.ToBytes()

	reqMsg := &types.Votes{Sequence: 0, Epoch: relayer.Epoch, Voters: voteBytes, Signature: sig}
	suite.Require().NoError(reqMsg.Validate())
	suite.Require().Equal(th-1, bitmap.FromBytes(voteBytes).Count())

	suite.VoteMsgMock.EXPECT().GetProposer().AnyTimes().Return(relayer.Proposer)
	suite.VoteMsgMock.EXPECT().GetVote().AnyTimes().Return(reqMsg)
	suite.VoteMsgMock.EXPECT().MethodName().AnyTimes().Return(reqMethod)
	suite.VoteMsgMock.EXPECT().VoteSigDoc().AnyTimes().Return(reqSigDoc)

	seq, err := suite.Keeper.VerifyProposal(suite.Context, suite.VoteMsgMock)
	suite.Require().Error(err)
	suite.Require().EqualValues(0, seq)
}
