package keeper_test

import (
	"slices"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/relayer/keeper"
	"github.com/goatnetwork/goat/x/relayer/types"
)

func (suite *KeeperTestSuite) TestMsgServerNewVoter() {
	relayer := types.Relayer{
		Proposer:         suite.VoterKeys[0].Address,
		Voters:           []string{suite.VoterKeys[1].Address},
		LastElected:      time.Now().UTC(),
		ProposerAccepted: true,
	}
	err := suite.Keeper.Relayer.Set(suite.Context, relayer)
	suite.Require().NoError(err)

	txKeyHash := goatcrypto.SHA256Sum(suite.Voters[2].VoteKey)
	err = suite.Keeper.ProcessRelayerRequest(suite.Context, goattypes.RelayerRequests{
		Adds: []*goattypes.AddVoterRequest{
			{Voter: common.Address(suite.Voters[2].Address), Pubkey: common.Hash(txKeyHash)},
		},
	})
	suite.Require().NoError(err)

	prvkey := &secp256k1.PrivKey{Key: suite.VoterKeys[2].TxKey}
	account, err := authtypes.NewBaseAccountWithPubKey(prvkey.PubKey())
	suite.Require().NoError(err)

	sdkctx := sdktypes.UnwrapSDKContext(suite.Context)

	sigdoc := slices.Concat(
		[]byte(sdkctx.ChainID()),
		goatcrypto.Uint64LE(0, 0),
		[]byte("Relayer/NewVoter"),
		[]byte(suite.VoterKeys[0].Address),
		types.NewOnBoardingVoterRequest(uint64(sdkctx.BlockHeight()), account.GetAddress(), txKeyHash).SignDoc(),
	)
	txKeyProof, err := prvkey.Sign(sigdoc)
	suite.Require().NoError(err)

	voteKeyProof := goatcrypto.Sign(new(goatcrypto.PrivateKey).Deserialize(suite.VoterKeys[2].VoteKey), goatcrypto.SHA256Sum(sigdoc))

	server := keeper.NewMsgServerImpl(suite.Keeper)

	suite.Account.EXPECT().HasAccount(suite.Context, account.GetAddress()).Return(false)
	suite.Account.EXPECT().NewAccountWithAddress(suite.Context, account.GetAddress()).Return(account)
	suite.Account.EXPECT().SetAccount(suite.Context, account)

	_, err = server.NewVoter(suite.Context, &types.MsgNewVoterRequest{
		Proposer:         suite.VoterKeys[0].Address,
		VoterBlsKey:      suite.Voters[2].VoteKey,
		VoterTxKey:       prvkey.PubKey().Bytes(),
		VoterTxKeyProof:  txKeyProof,
		VoterBlsKeyProof: voteKeyProof,
	})
	suite.Require().NoError(err)

	voter, err := suite.Keeper.Voters.Get(suite.Context, suite.VoterKeys[2].Address)
	suite.Require().NoError(err)
	suite.Require().Equal(voter, types.Voter{
		Address: suite.Voters[2].Address,
		VoteKey: suite.Voters[2].VoteKey,
		Status:  types.VOTER_STATUS_ON_BOARDING,
		Height:  uint64(sdkctx.BlockHeight()),
	})

	queue, err := suite.Keeper.Queue.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(queue, types.VoterQueue{
		OnBoarding: []string{suite.VoterKeys[2].Address},
	})
}

func (suite *KeeperTestSuite) TestMsgServerAcceptProposer() {
	relayer := types.Relayer{
		Proposer:         suite.VoterKeys[0].Address,
		Voters:           []string{suite.VoterKeys[1].Address},
		LastElected:      time.Now().UTC(),
		ProposerAccepted: false,
	}
	err := suite.Keeper.Relayer.Set(suite.Context, relayer)
	suite.Require().NoError(err)

	server := keeper.NewMsgServerImpl(suite.Keeper)

	_, err = server.AcceptProposer(suite.Context, &types.MsgAcceptProposerRequest{
		Proposer: relayer.Proposer,
		Epoch:    relayer.Epoch,
	})
	suite.Require().NoError(err)

	newRelayer, err := suite.Keeper.Relayer.Get(suite.Context)
	suite.Require().NoError(err)
	relayer.ProposerAccepted = true

	suite.Require().Equal(newRelayer, relayer)
}
