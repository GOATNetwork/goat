package keeper

import (
	"context"
	"errors"
	"math/big"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	"github.com/goatnetwork/goat/x/locking/types"
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

func (q queryServer) Validator(ctx context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	address, err := bitcointypes.DecodeEthAddress(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	validator, err := q.k.Validators.Get(sdkctx, sdktypes.ConsAddress(address))
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryValidatorResponse{Validator: validator, Height: sdkctx.BlockHeight()}, nil
}

func (q queryServer) ActiveValidators(ctx context.Context, req *types.QueryActiveValidatorsRequest) (*types.QueryActiveValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	resp := &types.QueryActiveValidatorsResponse{Height: sdkctx.BlockHeight()}
	iter, err := q.k.ValidatorSet.Iterate(sdkctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error: "+err.Error())
	}
	for ; iter.Valid(); iter.Next() {
		kv, err := iter.KeyValue()
		if err != nil {
			return nil, status.Error(codes.Internal, "internal error: "+err.Error())
		}
		validator, err := q.k.Validators.Get(sdkctx, kv.Key)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal error: "+err.Error())
		}
		pubkey := &secp256k1.PubKey{Key: validator.Pubkey}
		resp.Validators = append(resp.Validators, &types.ValidatorInfo{
			ValidatorAddress: hexutil.Encode(pubkey.Address()),
			Validator:        validator,
		})
		resp.TotalPower += int64(validator.Power)
	}
	if err := iter.Close(); err != nil {
		return nil, status.Error(codes.Internal, "internal error: "+err.Error())
	}
	total := math.NewIntFromBigIntMut(big.NewInt(resp.TotalPower))
	for _, v := range resp.Validators {
		v.PowerPercentage = math.LegacyNewDecWithPrec(int64(v.Validator.Power), 0).QuoInt(total).MulInt(math.NewInt(100))
	}
	return resp, nil
}
