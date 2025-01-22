package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/goatnetwork/goat/x/bitcoin/types"
)

func (k msgServer) ProcessWithdrawal(ctx context.Context, req *types.MsgProcessWithdrawal) (*types.MsgProcessWithdrawalResponse, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	// disable by default
	if hk := types.Hardforks[sdkctx.ChainID()]; hk == nil || hk.IsWithdrawalV2Enable(sdkctx.BlockTime()) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "ProcessWithdrawal is disabled")
	}
	if err := k.processWithdrawal(sdkctx, req); err != nil {
		return nil, err
	}
	return &types.MsgProcessWithdrawalResponse{}, nil
}

func (k msgServer) ProcessWithdrawalV2(ctx context.Context, req *types.MsgProcessWithdrawalV2) (*types.MsgProcessWithdrawalV2Response, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	if hk := types.Hardforks[sdkctx.ChainID()]; hk != nil && !hk.IsWithdrawalV2Enable(sdkctx.BlockTime()) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "tx v2 hardfork is not activated")
	}
	if err := k.processWithdrawal(ctx, req); err != nil {
		return nil, err
	}
	return &types.MsgProcessWithdrawalV2Response{}, nil
}

func (k msgServer) ReplaceWithdrawal(ctx context.Context, req *types.MsgReplaceWithdrawal) (*types.MsgReplaceWithdrawalResponse, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	// disable by default
	if hk := types.Hardforks[sdkctx.ChainID()]; hk == nil || hk.IsWithdrawalV2Enable(sdkctx.BlockTime()) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "ReplaceWithdrawal is disabled")
	}
	if err := k.replaceWithdrawal(sdkctx, req); err != nil {
		return nil, err
	}
	return &types.MsgReplaceWithdrawalResponse{}, nil
}

func (k msgServer) ReplaceWithdrawalV2(ctx context.Context, req *types.MsgReplaceWithdrawalV2) (*types.MsgReplaceWithdrawalV2Response, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	if hk := types.Hardforks[sdkctx.ChainID()]; hk != nil && !hk.IsWithdrawalV2Enable(sdkctx.BlockTime()) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "tx v2 hardfork is not activated")
	}
	if err := k.replaceWithdrawal(sdkctx, req); err != nil {
		return nil, err
	}
	return &types.MsgReplaceWithdrawalV2Response{}, nil
}
