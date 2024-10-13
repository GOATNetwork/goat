package keeper_test

import (
	"bytes"
	"slices"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/goatnetwork/goat/x/relayer/keeper"
	"github.com/goatnetwork/goat/x/relayer/types"
)

func (suite *KeeperTestSuite) TestQueryPubkeys() {
	expected := []*types.PublicKey{}
	for i := 0; i < 5; i++ {
		prvkey := secp256k1.GenPrivKey()
		pubkey := prvkey.PubKey().Bytes()

		typed := &types.PublicKey{
			Key: &types.PublicKey_Secp256K1{Secp256K1: pubkey},
		}
		encoded := types.EncodePublicKey(typed)
		err := suite.Keeper.Pubkeys.Set(suite.Context, encoded)
		suite.Require().NoError(err)
		expected = append(expected, typed)
	}

	slices.SortFunc(expected, func(i, j *types.PublicKey) int {
		return bytes.Compare(i.GetSecp256K1(), j.GetSecp256K1())
	})

	query := keeper.NewQueryServerImpl(suite.Keeper)

	resp, err := query.Pubkeys(suite.Context, nil)
	suite.Require().NoError(err)

	slices.SortFunc(resp.PublicKeys, func(i, j *types.PublicKey) int {
		return bytes.Compare(i.GetSecp256K1(), j.GetSecp256K1())
	})
	suite.Require().Equal(resp.PublicKeys, expected)
}

func (suite *KeeperTestSuite) TestQueryParams() {
	query := keeper.NewQueryServerImpl(suite.Keeper)

	resp, err := query.Params(suite.Context, nil)
	suite.Require().NoError(err)
	suite.Require().Equal(resp.Params, suite.Param)
}

func (suite *KeeperTestSuite) TestQueryRelayer() {
	relayer := types.Relayer{
		Proposer:         "goat1vpa50kdf63rfuurelvzsmpp99ffmvuc20h68qh",
		Voters:           []string{"goat1kqvxakwshk52q8vs69rcd0p9jttufuk0drvzgp"},
		LastElected:      time.Now().UTC(),
		ProposerAccepted: true,
	}
	err := suite.Keeper.Relayer.Set(suite.Context, relayer)
	suite.Require().NoError(err)

	seq := uint64(100)
	err = suite.Keeper.Sequence.Set(suite.Context, seq)
	suite.Require().NoError(err)

	query := keeper.NewQueryServerImpl(suite.Keeper)

	resp, err := query.Relayer(suite.Context, nil)
	suite.Require().NoError(err)

	suite.Require().Equal(relayer, *resp.Relayer)
	suite.Require().Equal(seq, resp.Sequence)
}

func (suite *KeeperTestSuite) TestQueryVoters() {
	for i, voter := range suite.Voters {
		err := suite.Keeper.Voters.Set(suite.Context, suite.VoterKeys[i].Address, voter)
		suite.Require().NoError(err)
	}

	slices.SortFunc(suite.Voters, func(i, j types.Voter) int {
		return bytes.Compare(i.Address, j.Address)
	})

	query := keeper.NewQueryServerImpl(suite.Keeper)

	resp, err := query.Voters(suite.Context, nil)
	suite.Require().NoError(err)

	slices.SortFunc(resp.Voters, func(i, j types.Voter) int {
		return bytes.Compare(i.Address, j.Address)
	})
	suite.Require().Equal(resp.Voters, suite.Voters)
}
