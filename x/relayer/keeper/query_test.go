package keeper_test

import (
	"bytes"
	"slices"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

	suite.Require().Equal(relayer, resp.Relayer)
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

	for _, voter := range suite.Voters {
		{
			resp, err := query.Voter(suite.Context, &types.QueryVoterRequest{Address: hexutil.Encode(voter.Address)})
			suite.Require().NoError(err)
			suite.Require().Equal(resp.Voter, voter)
		}

		{
			address, err := suite.Keeper.AddrCodec.BytesToString(voter.Address)
			suite.Require().NoError(err)
			resp, err := query.Voter(suite.Context, &types.QueryVoterRequest{Address: address})
			suite.Require().NoError(err)
			suite.Require().Equal(resp.Voter, voter)
		}
	}

	{
		_, err := query.Voter(suite.Context, nil)
		suite.Require().ErrorContains(err, "invalid request")
	}

	{
		_, err := query.Voter(suite.Context, &types.QueryVoterRequest{Address: "0xxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxy"})
		suite.Require().ErrorContains(err, "invalid address(eth format)")
	}

	{
		_, err := query.Voter(suite.Context, &types.QueryVoterRequest{Address: "bc1qassh6388meyz0zqyjgwfffynfngjvup5dqhpsc"})
		suite.Require().ErrorContains(err, "invalid address(bech32 format)")
	}

	{
		_, err := query.Voter(suite.Context, &types.QueryVoterRequest{Address: "goat1gu2ttwaut55r4cguylq4l6xgjyjm5s5z4ugr9w"})
		suite.Require().ErrorContains(err, "not found")
	}
}
