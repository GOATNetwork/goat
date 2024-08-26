package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/goatnetwork/goat/x/relayer/types"
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

func (q queryServer) Relayer(ctx context.Context, req *types.QueryRelayerRequest) (*types.QueryRelayerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	epoch, err := q.k.Epoch.Peek(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	relayer, err := q.k.Relayer.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryRelayerResponse{Relayer: relayer, Epoch: epoch}, nil
}

func (q queryServer) Voters(ctx context.Context, req *types.QueryVotersRequest) (*types.QueryVotersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// todo: Pagination

	iter, err := q.k.Voters.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var res types.QueryVotersResponse
	for ; iter.Valid(); iter.Next() {
		kv, err := iter.Value()
		if err != nil {
			return nil, err
		}
		res.Voters = append(res.Voters, &kv)
	}
	return &res, nil
}

func (q queryServer) Pubkeys(ctx context.Context, req *types.QueryPubkeysRequest) (*types.QueryPubkeysResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	iter, err := q.k.Pubkeys.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var pubkeys []*types.PublicKey
	for ; iter.Valid(); iter.Next() {
		value, err := iter.Key()
		if err != nil {
			return nil, err
		}

		res := types.DecodePublicKey(value)
		if res == nil {
			return nil, status.Error(codes.Internal, "invalid public key to decode")
		}

		pubkeys = append(pubkeys, res)
	}

	return &types.QueryPubkeysResponse{PublicKeys: pubkeys}, nil
}
