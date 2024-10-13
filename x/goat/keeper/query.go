package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/goatnetwork/goat/x/goat/types"
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

func (q queryServer) EthBlockTip(ctx context.Context, req *types.QueryEthBlockTipRequest) (*types.QueryEthBlockTipResponse, error) {
	block, err := q.k.Block.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &types.QueryEthBlockTipResponse{Block: block}, nil
}
