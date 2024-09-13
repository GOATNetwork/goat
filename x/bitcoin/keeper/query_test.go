package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/goatnetwork/goat/testutil/keeper"
	"github.com/goatnetwork/goat/x/bitcoin/keeper"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
)

var TestingPubkey = relayer.PublicKey{Key: &relayer.PublicKey_Secp256K1{Secp256K1: common.Hex2Bytes("0383560def84048edefe637d0119a4428dd12a42765a118b2bf77984057633c50e")}}

func TestParamsQuery(t *testing.T) {
	k, ctx, _ := keepertest.BitcoinKeeper(t, nil)

	qs := keeper.NewQueryServerImpl(k)
	params := types.DefaultParams()
	require.NoError(t, k.Params.Set(ctx, params))

	response, err := qs.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}

func TestPubkeyQuery(t *testing.T) {
	k, ctx, _ := keepertest.BitcoinKeeper(t, nil)

	qs := keeper.NewQueryServerImpl(k)

	_, err := qs.Pubkey(ctx, &types.QueryPubkeyRequest{})
	require.EqualError(t, err, status.Error(codes.NotFound, "not found").Error())

	require.NoError(t, k.Pubkey.Set(ctx, TestingPubkey))

	got, err := qs.Pubkey(ctx, &types.QueryPubkeyRequest{})
	require.NoError(t, err)
	require.Equal(t, got.PublicKey, TestingPubkey)
}

func TestDepositAddress(t *testing.T) {
	k, ctx, _ := keepertest.BitcoinKeeper(t, nil)
	require.NoError(t, k.Pubkey.Set(ctx, TestingPubkey))

	qs := keeper.NewQueryServerImpl(k)

	const invalid = "invalid"
	const address = "0xBC1cb6A680cF76505F3C992Fa6ab4F91913511e8"

	_, err := qs.DepositAddress(ctx, &types.QueryDepositAddress{EvmAddress: invalid})
	require.EqualError(t, err, status.Error(codes.InvalidArgument, "invalid eth address").Error())

	_, err = qs.DepositAddress(ctx, &types.QueryDepositAddress{EvmAddress: address, Version: 0xff})
	require.EqualError(t, err, status.Error(codes.InvalidArgument, "unknown deposit version").Error())

	params := types.DefaultParams()
	chaincfg := params.ChainConfig.ToBtcdParam()

	{
		resp, err := qs.DepositAddress(ctx, &types.QueryDepositAddress{EvmAddress: address, Version: 1})
		require.NoError(t, err)

		require.Equal(t, resp.NetworkName, chaincfg.Name)
		require.Equal(t, *resp.PublicKey, TestingPubkey)
		require.Equal(t, resp.Address, "bcrt1qjav7664wdt0y8tnx9z558guewnxjr3wllz2s9u")

		require.NoError(t, err)
		require.Equal(t, resp.OpReturnScript, common.Hex2Bytes("1847545430bc1cb6a680cf76505f3c992fa6ab4f91913511e8"))
	}
}
