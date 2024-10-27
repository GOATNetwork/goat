package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"

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
	sequence, err := q.k.Sequence.Peek(ctx)
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

	return &types.QueryRelayerResponse{Relayer: relayer, Sequence: sequence}, nil
}

func (q queryServer) Voter(ctx context.Context, req *types.QueryVoterRequest) (*types.QueryVoterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if len(req.Address) == 42 && strings.HasPrefix(req.Address, "0x") {
		addrBytes, err := hex.DecodeString(req.Address[2:])
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid address(eth format)")
		}
		address, err := q.k.AddrCodec.BytesToString(addrBytes)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to encode address")
		}
		req.Address = address
	} else {
		if _, err := q.k.AddrCodec.StringToBytes(req.Address); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid address(bech32 format)")
		}
	}

	voter, err := q.k.Voters.Get(ctx, req.Address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryVoterResponse{Voter: voter}, nil
}

func (q queryServer) Pubkeys(ctx context.Context, req *types.QueryPubkeysRequest) (*types.QueryPubkeysResponse, error) {
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
