package keeper_test

import (
	"time"

	"cosmossdk.io/collections"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/relayer/types"
)

func (suite *KeeperTestSuite) TestEndBlocker() {
	suite.Run("no-action", func() {
		err := suite.Keeper.Params.Set(suite.Context, types.Params{
			ElectingPeriod:        time.Minute,
			AcceptProposerTimeout: time.Minute,
		})
		suite.Require().NoError(err)

		relayer := types.Relayer{
			Proposer:         suite.VoterKeys[0].Address,
			Voters:           []string{suite.VoterKeys[1].Address},
			LastElected:      time.Now().UTC(),
			ProposerAccepted: true,
		}
		err = suite.Keeper.Relayer.Set(suite.Context, relayer)
		suite.Require().NoError(err)

		for i := 2; i < 4; i++ {
			v := types.Voter{
				Address: suite.Voters[i].Address,
				VoteKey: suite.Voters[i].VoteKey,
				Status:  types.VOTER_STATUS_PENDING,
				Height:  111,
			}
			err = suite.Keeper.Voters.Set(suite.Context, suite.VoterKeys[i].Address, v)
			suite.Require().NoError(err)
		}

		err = suite.Keeper.Queue.Set(suite.Context, types.VoterQueue{
			OnBoarding:  []string{suite.VoterKeys[2].Address, suite.VoterKeys[3].Address},
			OffBoarding: []string{suite.VoterKeys[0].Address},
		})
		suite.Require().NoError(err)

		for i := 0; i < 2; i++ {
			voter, err := suite.Keeper.Voters.Get(suite.Context, suite.VoterKeys[i].Address)
			suite.Require().NoError(err)
			suite.Require().Equal(voter, types.Voter{
				Address: suite.Voters[i].Address,
				VoteKey: suite.Voters[i].VoteKey,
				Status:  types.VOTER_STATUS_ACTIVATED,
				Height:  100,
			})
		}

		for i := 2; i < 4; i++ {
			voter, err := suite.Keeper.Voters.Get(suite.Context, suite.VoterKeys[i].Address)
			suite.Require().NoError(err)
			suite.Require().Equal(voter, types.Voter{
				Address: suite.Voters[i].Address,
				VoteKey: suite.Voters[i].VoteKey,
				Status:  types.VOTER_STATUS_PENDING,
				Height:  111,
			})
		}
	})

	suite.Run("remove-proposer", func() {
		err := suite.Keeper.Params.Set(suite.Context, types.Params{
			ElectingPeriod:        time.Minute * 10,
			AcceptProposerTimeout: time.Minute,
		})
		suite.Require().NoError(err)

		relayer := types.Relayer{
			Proposer:         suite.VoterKeys[0].Address,
			Voters:           []string{suite.VoterKeys[1].Address},
			LastElected:      time.Now().UTC().Add(0 - time.Minute*10),
			ProposerAccepted: false,
		}
		err = suite.Keeper.Relayer.Set(suite.Context, relayer)
		suite.Require().NoError(err)

		for i := 2; i < 4; i++ {
			err = suite.Keeper.Voters.Set(suite.Context, suite.VoterKeys[i].Address, types.Voter{
				Address: suite.Voters[i].Address,
				VoteKey: suite.Voters[i].VoteKey,
				Status:  types.VOTER_STATUS_PENDING,
			})
			suite.Require().NoError(err)
		}

		err = suite.Keeper.Queue.Set(suite.Context, types.VoterQueue{
			OnBoarding:  []string{suite.VoterKeys[2].Address, suite.VoterKeys[3].Address},
			OffBoarding: []string{suite.VoterKeys[0].Address},
		})
		suite.Require().NoError(err)

		err = suite.Keeper.EndBlocker(suite.Context)
		suite.Require().NoError(err)

		queue, err := suite.Keeper.Queue.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(queue, types.VoterQueue{})

		sdkctx := sdktypes.UnwrapSDKContext(suite.Context)
		newRelayer := types.Relayer{
			Epoch:            1,
			Proposer:         suite.VoterKeys[1].Address,
			Voters:           []string{suite.VoterKeys[2].Address, suite.VoterKeys[3].Address},
			LastElected:      sdkctx.BlockTime(),
			ProposerAccepted: false,
		}
		gotRelayer, err := suite.Keeper.Relayer.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(newRelayer, gotRelayer)

		_, err = suite.Keeper.Voters.Get(suite.Context, suite.VoterKeys[0].Address)
		suite.Require().ErrorIs(err, collections.ErrNotFound)

		for i := 2; i < 4; i++ {
			voter, err := suite.Keeper.Voters.Get(suite.Context, suite.VoterKeys[i].Address)
			suite.Require().NoError(err)
			suite.Require().Equal(voter, types.Voter{
				Address: suite.Voters[i].Address,
				VoteKey: suite.Voters[i].VoteKey,
				Status:  types.VOTER_STATUS_ACTIVATED,
			})
		}
	})

	suite.Run("remove-voter-1", func() {
		err := suite.Keeper.Params.Set(suite.Context, types.Params{
			ElectingPeriod: time.Second,
		})
		suite.Require().NoError(err)

		relayer := types.Relayer{
			Proposer:         suite.VoterKeys[0].Address,
			Voters:           []string{suite.VoterKeys[1].Address},
			LastElected:      time.Now().UTC().Add(-time.Minute),
			ProposerAccepted: true,
		}
		err = suite.Keeper.Relayer.Set(suite.Context, relayer)
		suite.Require().NoError(err)

		err = suite.Keeper.Queue.Set(suite.Context, types.VoterQueue{
			OffBoarding: []string{suite.VoterKeys[1].Address},
		})
		suite.Require().NoError(err)

		err = suite.Keeper.EndBlocker(suite.Context)
		suite.Require().NoError(err)

		queue, err := suite.Keeper.Queue.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(queue, types.VoterQueue{})

		sdkctx := sdktypes.UnwrapSDKContext(suite.Context)
		newRelayer := types.Relayer{
			Epoch:            1,
			Proposer:         suite.VoterKeys[0].Address,
			Voters:           nil,
			LastElected:      sdkctx.BlockTime(),
			ProposerAccepted: true,
		}
		gotRelayer, err := suite.Keeper.Relayer.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(newRelayer, gotRelayer)

		_, err = suite.Keeper.Voters.Get(suite.Context, suite.VoterKeys[1].Address)
		suite.Require().ErrorIs(err, collections.ErrNotFound)
	})

	suite.Run("remove-voter-2", func() {
		err := suite.Keeper.Params.Set(suite.Context, types.Params{
			ElectingPeriod: time.Second,
		})
		suite.Require().NoError(err)

		relayer := types.Relayer{
			Proposer:         suite.VoterKeys[0].Address,
			Voters:           []string{suite.VoterKeys[1].Address},
			LastElected:      time.Now().UTC().Add(-time.Minute),
			ProposerAccepted: true,
		}
		err = suite.Keeper.Relayer.Set(suite.Context, relayer)
		suite.Require().NoError(err)

		for i := 2; i < 4; i++ {
			err = suite.Keeper.Voters.Set(suite.Context, suite.VoterKeys[i].Address, types.Voter{
				Address: suite.Voters[i].Address,
				VoteKey: suite.Voters[i].VoteKey,
				Status:  types.VOTER_STATUS_PENDING,
			})
			suite.Require().NoError(err)
		}

		err = suite.Keeper.Queue.Set(suite.Context, types.VoterQueue{
			OnBoarding:  []string{suite.VoterKeys[2].Address},
			OffBoarding: []string{suite.VoterKeys[1].Address},
		})
		suite.Require().NoError(err)

		err = suite.Keeper.EndBlocker(suite.Context)
		suite.Require().NoError(err)

		queue, err := suite.Keeper.Queue.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(queue, types.VoterQueue{})

		sdkctx := sdktypes.UnwrapSDKContext(suite.Context)
		newRelayer := types.Relayer{
			Epoch:            1,
			Proposer:         suite.VoterKeys[2].Address,
			Voters:           []string{suite.VoterKeys[0].Address},
			LastElected:      sdkctx.BlockTime(),
			ProposerAccepted: false,
		}
		gotRelayer, err := suite.Keeper.Relayer.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(newRelayer, gotRelayer)

		_, err = suite.Keeper.Voters.Get(suite.Context, suite.VoterKeys[1].Address)
		suite.Require().ErrorIs(err, collections.ErrNotFound)

		voter, err := suite.Keeper.Voters.Get(suite.Context, suite.VoterKeys[2].Address)
		suite.Require().NoError(err)
		suite.Require().Equal(voter, types.Voter{
			Address: suite.Voters[2].Address,
			VoteKey: suite.Voters[2].VoteKey,
			Status:  types.VOTER_STATUS_ACTIVATED,
		})
	})

	suite.Run("remove-voter-3", func() {
		err := suite.Keeper.Params.Set(suite.Context, types.Params{
			ElectingPeriod: time.Second,
		})
		suite.Require().NoError(err)

		relayer := types.Relayer{
			Proposer:         suite.VoterKeys[0].Address,
			Voters:           []string{suite.VoterKeys[1].Address},
			LastElected:      time.Now().UTC().Add(-time.Minute),
			ProposerAccepted: true,
		}
		err = suite.Keeper.Relayer.Set(suite.Context, relayer)
		suite.Require().NoError(err)

		for i := 2; i < 4; i++ {
			err = suite.Keeper.Voters.Set(suite.Context, suite.VoterKeys[i].Address, types.Voter{
				Address: suite.Voters[i].Address,
				VoteKey: suite.Voters[i].VoteKey,
				Status:  types.VOTER_STATUS_PENDING,
			})
			suite.Require().NoError(err)
		}

		err = suite.Keeper.Queue.Set(suite.Context, types.VoterQueue{
			OnBoarding:  []string{suite.VoterKeys[2].Address, suite.VoterKeys[3].Address},
			OffBoarding: []string{suite.VoterKeys[1].Address},
		})
		suite.Require().NoError(err)

		err = suite.Keeper.EndBlocker(suite.Context)
		suite.Require().NoError(err)

		queue, err := suite.Keeper.Queue.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(queue, types.VoterQueue{})

		sdkctx := sdktypes.UnwrapSDKContext(suite.Context)
		newRelayer := types.Relayer{
			Epoch:            1,
			Proposer:         suite.VoterKeys[3].Address,
			Voters:           []string{suite.VoterKeys[2].Address, suite.VoterKeys[0].Address},
			LastElected:      sdkctx.BlockTime(),
			ProposerAccepted: false,
		}
		gotRelayer, err := suite.Keeper.Relayer.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(newRelayer, gotRelayer)
	})
}
