package keeper_test

import (
	"encoding/hex"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goatnetwork/goat/x/bitcoin/keeper"
	"github.com/goatnetwork/goat/x/bitcoin/types"
)

func (suite *KeeperTestSuite) TestParamsQuery() {
	qs := keeper.NewQueryServerImpl(suite.Keeper)
	params := types.DefaultParams()
	suite.Require().NoError(suite.Keeper.Params.Set(suite.Context, params))

	response, err := qs.Params(suite.Context, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryParamsResponse{Params: params}, response)
}

func (suite *KeeperTestSuite) TestPubkeyQuery() {
	qs := keeper.NewQueryServerImpl(suite.Keeper)

	_, err := qs.Pubkey(suite.Context, &types.QueryPubkeyRequest{})
	suite.Require().EqualError(err, status.Error(codes.NotFound, "not found").Error())

	suite.Require().NoError(suite.Keeper.Pubkey.Set(suite.Context, suite.TestKey))
	got, err := qs.Pubkey(suite.Context, &types.QueryPubkeyRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(got.PublicKey, suite.TestKey)
}

func (suite *KeeperTestSuite) TestDepositAddress() {
	suite.Require().NoError(suite.Keeper.Pubkey.Set(suite.Context, suite.TestKey))

	qs := keeper.NewQueryServerImpl(suite.Keeper)

	const invalid = "invalid"
	const address = "0xBC1cb6A680cF76505F3C992Fa6ab4F91913511e8"

	_, err := qs.DepositAddress(suite.Context, &types.QueryDepositAddress{EvmAddress: invalid})
	suite.Require().EqualError(err, status.Error(codes.InvalidArgument, "invalid eth address").Error())

	_, err = qs.DepositAddress(suite.Context, &types.QueryDepositAddress{EvmAddress: address, Version: 0xff})
	suite.Require().EqualError(err, status.Error(codes.InvalidArgument, "unknown deposit version").Error())

	params := types.DefaultParams()
	chaincfg := types.BitcoinNetworks[params.NetworkName]

	{
		resp, err := qs.DepositAddress(suite.Context, &types.QueryDepositAddress{EvmAddress: address, Version: 0})
		suite.Require().NoError(err)

		suite.Require().Equal(resp.NetworkName, chaincfg.Name)
		suite.Require().Equal(*resp.PublicKey, suite.TestKey)
		suite.Require().Equal(resp.Address, "bcrt1q6dxx7mfels0u4f59c0mjvltukvgnnur7v377nrusdm0r3gm0ycjsxjx0uj")

		suite.Require().NoError(err)
		suite.Require().Equal(hex.EncodeToString(resp.OpReturnScript), "")
	}

	{
		resp, err := qs.DepositAddress(suite.Context, &types.QueryDepositAddress{EvmAddress: address, Version: 1})
		suite.Require().NoError(err)

		suite.Require().Equal(resp.NetworkName, chaincfg.Name)
		suite.Require().Equal(*resp.PublicKey, suite.TestKey)
		suite.Require().Equal(resp.Address, "bcrt1qjav7664wdt0y8tnx9z558guewnxjr3wllz2s9u")

		suite.Require().NoError(err)
		suite.Require().Equal(hex.EncodeToString(resp.OpReturnScript), "6a1847545430bc1cb6a680cf76505f3c992fa6ab4f91913511e8")
	}
}
