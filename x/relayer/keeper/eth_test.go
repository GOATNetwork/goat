package keeper_test

import (
	"crypto/sha256"
	"slices"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/relayer/types"
)

func (suite *KeeperTestSuite) TestProcessRelayerRequest() {
	relayer := types.Relayer{
		Proposer:         suite.VoterKeys[0].Address,
		Voters:           []string{suite.VoterKeys[1].Address},
		LastElected:      time.Now().UTC(),
		ProposerAccepted: true,
	}

	err := suite.Keeper.Relayer.Set(suite.Context, relayer)
	suite.Require().NoError(err)

	err = suite.Keeper.ProcessRelayerRequest(suite.Context, goattypes.RelayerRequests{
		Adds: []*goattypes.AddVoterRequest{
			{Voter: common.Address(suite.Voters[0].Address), Pubkey: sha256.Sum256(suite.Voters[0].VoteKey)},
			{Voter: common.Address(suite.Voters[2].Address), Pubkey: sha256.Sum256(suite.Voters[2].VoteKey)},
		},
		Removes: []*goattypes.RemoveVoterRequest{
			{Voter: common.Address(suite.Voters[1].Address)},
			{Voter: common.Address(suite.Voters[3].Address)},
		},
	})
	suite.Require().NoError(err)

	iter, err := suite.Keeper.Voters.Iterate(suite.Context, nil)
	suite.Require().NoError(err)
	gotVoters, err := iter.Keys()
	suite.Require().NoError(err)
	suite.Require().NoError(iter.Close())

	gotRelayer, err := suite.Keeper.Relayer.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(relayer, gotRelayer, "relayer should have no changes")

	slices.Sort(gotVoters)
	expectVoters := []string{suite.VoterKeys[0].Address, suite.VoterKeys[1].Address, suite.VoterKeys[2].Address}
	slices.Sort(expectVoters)
	suite.Require().Equal(gotVoters, expectVoters)

	voter1, err := suite.Keeper.Voters.Get(suite.Context, suite.VoterKeys[1].Address)
	suite.Require().NoError(err)
	suite.Require().Equal(voter1.Address, suite.Voters[1].Address)
	suite.Require().Equal(voter1.VoteKey, suite.Voters[1].VoteKey)
	suite.Require().Equal(voter1.Height, suite.Voters[1].Height)
	suite.Require().Equal(voter1.Status, types.VOTER_STATUS_OFF_BOARDING)

	voter2, err := suite.Keeper.Voters.Get(suite.Context, suite.VoterKeys[2].Address)
	suite.Require().NoError(err)
	suite.Require().Equal(voter2.Address, suite.Voters[2].Address)
	suite.Require().Equal(voter2.Status, types.VOTER_STATUS_PENDING)
	suite.Require().EqualValues(voter2.Height, sdktypes.UnwrapSDKContext(suite.Context).BlockHeight())
	suite.Require().Equal([32]byte(voter2.VoteKey), sha256.Sum256(suite.Voters[2].VoteKey))
}
