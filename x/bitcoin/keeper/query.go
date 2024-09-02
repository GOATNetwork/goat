package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	params, err := q.k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) Pubkey(ctx context.Context, req *types.QueryPubkeyRequest) (*types.QueryPubkeyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	key, err := q.k.Pubkey.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryPubkeyResponse{PublicKey: key}, nil
}

func (q queryServer) DepositAddress(ctx context.Context, req *types.QueryDepositAddress) (*types.QueryDepositAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	pubkey, err := q.k.Pubkey.Get(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	param, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	chainConfig, ok := types.BitcoinNetworks[param.NetworkName]
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal error: undefined network %s", param.NetworkName)
	}

	switch req.Version {
	case 0:
		address, err := types.DepositAddressV0(&pubkey, req.EvmAddress, chainConfig)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid request: %s", err.Error())
		}
		return &types.QueryDepositAddressResponse{
			Address:     address.EncodeAddress(),
			PublicKey:   &pubkey,
			NetworkName: param.NetworkName,
		}, nil
	case 1:
		address, script, err := types.DepositAddressV1(&pubkey, param.DepositMagicPrefix, req.EvmAddress, chainConfig)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid request: %s", err.Error())
		}
		return &types.QueryDepositAddressResponse{
			NetworkName:    param.NetworkName,
			PublicKey:      &pubkey,
			Address:        address.EncodeAddress(),
			OpReturnScript: script,
		}, nil
	}
	return nil, status.Error(codes.InvalidArgument, "unknown deposit version")
}

func (q queryServer) WithdrawalAddress(ctx context.Context, req *types.QueryWithdrawalAddress) (*types.QueryWithdrawalAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	param, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	chainConfig, ok := types.BitcoinNetworks[param.NetworkName]
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal error: undefined network %s", param.NetworkName)
	}

	address, err := types.WithdrawalAddress(req.Address, chainConfig)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %s", err.Error())
	}

	return &types.QueryWithdrawalAddressResponse{Address: address}, nil
}
