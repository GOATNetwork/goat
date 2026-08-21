package keeper_test

import (
	"errors"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/beacon/engine"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"

	"github.com/goatnetwork/goat/x/goat/types"
	"go.uber.org/mock/gomock"
	protov2 "google.golang.org/protobuf/proto"
)

// The proposal path is the one place where a rotation is visible to code that
// predates it: CometBFT reports the address of the key currently in effect,
// while the account and the fee recipient stay keyed by the creation-time
// validator id. Every test below pins one half of that split.

// fakeTx is an sdk.Tx whose messages are whatever the test put in it. It exists
// so the tests can drive ProcessProposalHandler without a real tx encoder.
type fakeTx struct{ msgs []sdk.Msg }

func (t fakeTx) GetMsgs() []sdk.Msg { return t.msgs }

func (t fakeTx) GetMsgsV2() ([]protov2.Message, error) { return nil, nil }

// fakeVerifier hands back a canned tx per raw byte slice, keyed by index, so a
// test can say "the second transaction carries a MsgNewEthBlock" without
// encoding anything.
type fakeVerifier struct {
	txs []sdk.Tx
	err error
}

func (v *fakeVerifier) PrepareProposalVerifyTx(sdk.Tx) ([]byte, error) { return nil, v.err }

func (v *fakeVerifier) ProcessProposalVerifyTx(txBz []byte) (sdk.Tx, uint64, error) {
	if v.err != nil {
		return nil, 0, v.err
	}
	return v.txs[int(txBz[0])], 0, nil
}

func (v *fakeVerifier) TxDecode(txBz []byte) (sdk.Tx, error) { return v.txs[int(txBz[0])], nil }

func (v *fakeVerifier) TxEncode(sdk.Tx) ([]byte, error) { return nil, nil }

// rawTxs turns a count into the index-encoded byte slices fakeVerifier expects.
func rawTxs(n int) [][]byte {
	out := make([][]byte, n)
	for i := range out {
		out[i] = []byte{byte(i)}
	}
	return out
}

// rotatedProposer returns a consensus address unlike the validator id, which is
// what CometBFT reports once a validator has rotated.
func rotatedProposer() sdk.ConsAddress {
	return sdk.ConsAddress(secp256k1.GenPrivKey().PubKey().Address())
}

func (suite *KeeperTestSuite) processProposal(v *fakeVerifier, ctx sdk.Context, txs [][]byte) error {
	_, err := suite.Keeper.ProcessProposalHandler(v)(ctx, &abci.RequestProcessProposal{Txs: txs})
	return err
}

func (suite *KeeperTestSuite) TestProcessProposalRejectsEmptyBlock() {
	err := suite.processProposal(&fakeVerifier{}, suite.Context, nil)
	suite.Require().ErrorContains(err, "no transactions")
}

func (suite *KeeperTestSuite) TestProcessProposalRejectsOversizedBlock() {
	err := suite.processProposal(&fakeVerifier{}, suite.Context, rawTxs(17))
	suite.Require().ErrorContains(err, "too many transactions")
}

func (suite *KeeperTestSuite) TestProcessProposalRequiresEthBlockFirst() {
	v := &fakeVerifier{txs: []sdk.Tx{fakeTx{msgs: []sdk.Msg{&relayertypes.MsgAcceptProposerRequest{}}}}}
	err := suite.processProposal(v, suite.Context, rawTxs(1))
	suite.Require().ErrorContains(err, "the first tx should be MsgNewEthBlock")
}

func (suite *KeeperTestSuite) TestProcessProposalRejectsSecondEthBlock() {
	proposer := rotatedProposer()
	id := sdk.ConsAddress(secp256k1.GenPrivKey().PubKey().Address())
	suite.Locking.EXPECT().ResolveValidatorID(gomock.Any(), proposer).Return(id, nil).AnyTimes()
	suite.Engine.EXPECT().NewPayloadV4(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&engine.PayloadStatusV1{Status: engine.VALID}, nil).AnyTimes()

	// The first tx has to get past verifyEthBlockProposal for the loop to reach
	// the second one, so give it an empty payload and stop caring about the
	// verdict; what this pins is that a second MsgNewEthBlock is refused.
	v := &fakeVerifier{txs: []sdk.Tx{
		fakeTx{msgs: []sdk.Msg{&types.MsgNewEthBlock{}}},
		fakeTx{msgs: []sdk.Msg{&types.MsgNewEthBlock{}}},
	}}
	err := suite.processProposal(v, suite.ctxWithProposer(proposer), rawTxs(2))
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) ctxWithProposer(addr sdk.ConsAddress) sdk.Context {
	return suite.Context.WithCometInfo(baseapp.NewBlockInfo(nil, nil, addr, abci.CommitInfo{}))
}

// ethBlockTx builds the first transaction of a proposal.
func ethBlockTx(proposer string, feeRecipient []byte) sdk.Tx {
	return fakeTx{msgs: []sdk.Msg{&types.MsgNewEthBlock{
		Proposer: proposer,
		Payload:  &types.ExecutionPayload{FeeRecipient: feeRecipient},
	}}}
}

func (suite *KeeperTestSuite) TestProcessProposalRejectsEmptyPayload() {
	proposer := rotatedProposer()
	suite.Engine.EXPECT().NewPayloadV4(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&engine.PayloadStatusV1{Status: engine.VALID}, nil).AnyTimes()

	v := &fakeVerifier{txs: []sdk.Tx{fakeTx{msgs: []sdk.Msg{&types.MsgNewEthBlock{}}}}}
	err := suite.processProposal(v, suite.ctxWithProposer(proposer), rawTxs(1))
	suite.Require().ErrorContains(err, "empty payload")
}

// A rotated validator proposes under a consensus address that is not its id.
// The identity checks have to resolve it rather than compare it directly.
func (suite *KeeperTestSuite) TestProcessProposalResolvesRotatedProposer() {
	comet := rotatedProposer()
	id := sdk.ConsAddress(secp256k1.GenPrivKey().PubKey().Address())
	idStr, err := suite.Codec.BytesToString(id)
	suite.Require().NoError(err)

	suite.Locking.EXPECT().ResolveValidatorID(gomock.Any(), comet).Return(id, nil)
	suite.Engine.EXPECT().NewPayloadV4(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&engine.PayloadStatusV1{Status: engine.VALID}, nil).AnyTimes()

	v := &fakeVerifier{txs: []sdk.Tx{ethBlockTx(idStr, id)}}
	err = suite.processProposal(v, suite.ctxWithProposer(comet), rawTxs(1))

	// It still fails further down, on state this test does not set up. What
	// matters is that neither of the two identity checks is what rejected it:
	// before the rotation work those compared against the comet address and
	// would have refused a rotated proposer outright.
	suite.Require().Error(err)
	suite.Require().NotContains(err.Error(), "invalid MsgNewEthBlock proposer")
	suite.Require().NotContains(err.Error(), "fee recipient mismatched")
}

func (suite *KeeperTestSuite) TestProcessProposalRejectsForeignProposer() {
	comet := rotatedProposer()
	id := sdk.ConsAddress(secp256k1.GenPrivKey().PubKey().Address())
	other := sdk.ConsAddress(secp256k1.GenPrivKey().PubKey().Address())
	otherStr, err := suite.Codec.BytesToString(other)
	suite.Require().NoError(err)

	suite.Locking.EXPECT().ResolveValidatorID(gomock.Any(), comet).Return(id, nil)
	suite.Engine.EXPECT().NewPayloadV4(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&engine.PayloadStatusV1{Status: engine.VALID}, nil).AnyTimes()

	v := &fakeVerifier{txs: []sdk.Tx{ethBlockTx(otherStr, other)}}
	err = suite.processProposal(v, suite.ctxWithProposer(comet), rawTxs(1))
	suite.Require().ErrorContains(err, "invalid MsgNewEthBlock proposer")
}

func (suite *KeeperTestSuite) TestProcessProposalRejectsForeignFeeRecipient() {
	comet := rotatedProposer()
	id := sdk.ConsAddress(secp256k1.GenPrivKey().PubKey().Address())
	idStr, err := suite.Codec.BytesToString(id)
	suite.Require().NoError(err)

	suite.Locking.EXPECT().ResolveValidatorID(gomock.Any(), comet).Return(id, nil)
	suite.Engine.EXPECT().NewPayloadV4(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&engine.PayloadStatusV1{Status: engine.VALID}, nil).AnyTimes()

	// The fee recipient is the comet address, which is exactly what a node that
	// had not been taught about rotation would have put there.
	v := &fakeVerifier{txs: []sdk.Tx{ethBlockTx(idStr, comet)}}
	err = suite.processProposal(v, suite.ctxWithProposer(comet), rawTxs(1))
	suite.Require().ErrorContains(err, "fee recipient mismatched")
}

// --- PrepareProposal ---------------------------------------------------
//
// Everything below stops inside createEthBlockProposal before the engine is
// reached, so no payload has to be faked.

// PrepareProposal reads the proposer off the request, not off CometInfo the way
// ProcessProposal does. Both end up at ResolveValidatorID, by different routes.
func (suite *KeeperTestSuite) prepareProposal(key cryptotypes.PrivKey, proposer sdk.ConsAddress) error {
	handler := suite.Keeper.PrepareProposalHandler(mempool.NoOpMempool{}, &fakeVerifier{}, key, nil)
	_, err := handler(suite.Context, &abci.RequestPrepareProposal{Height: 1, ProposerAddress: proposer})
	return err
}

func (suite *KeeperTestSuite) TestPrepareProposalPropagatesResolveFailure() {
	comet := rotatedProposer()
	suite.Locking.EXPECT().ResolveValidatorID(gomock.Any(), comet).
		Return(nil, errors.New("unknown consensus address"))

	err := suite.prepareProposal(secp256k1.GenPrivKey(), comet)
	suite.Require().ErrorContains(err, "unknown consensus address")
}

// The account is looked up by the resolved id, never by the address CometBFT
// reported. A rotated validator has no account at its new consensus address, so
// getting this wrong makes it unable to propose.
func (suite *KeeperTestSuite) TestPrepareProposalLooksUpAccountByResolvedID() {
	comet := rotatedProposer()
	id := sdk.ConsAddress(secp256k1.GenPrivKey().PubKey().Address())

	suite.Locking.EXPECT().ResolveValidatorID(gomock.Any(), comet).Return(id, nil)
	suite.Account.EXPECT().GetAccount(gomock.Any(), sdk.AccAddress(id)).Return(nil)

	err := suite.prepareProposal(secp256k1.GenPrivKey(), comet)
	suite.Require().ErrorContains(err, "nil validator account")
}

// The signing key the node holds has to be the one the account was registered
// with. This is the check a validator trips when it swaps priv_validator_key
// without splitting the transaction signing key out first.
func (suite *KeeperTestSuite) TestPrepareProposalRejectsForeignSigningKey() {
	comet := rotatedProposer()
	id := sdk.ConsAddress(secp256k1.GenPrivKey().PubKey().Address())
	accountKey := secp256k1.GenPrivKey()
	nodeKey := secp256k1.GenPrivKey()

	acc := authtypes.NewBaseAccount(sdk.AccAddress(id), accountKey.PubKey(), 0, 0)
	suite.Locking.EXPECT().ResolveValidatorID(gomock.Any(), comet).Return(id, nil)
	suite.Account.EXPECT().GetAccount(gomock.Any(), sdk.AccAddress(id)).Return(acc)

	err := suite.prepareProposal(nodeKey, comet)
	suite.Require().ErrorContains(err, "validator pubkey mismatched")
}
