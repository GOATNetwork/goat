package keeper

import (
	"bytes"
	"context"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/goat/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) NewEthBlock(ctx context.Context, req *types.MsgNewEthBlock) (*types.MsgNewEthBlockResponse, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	proposer, err := k.addressCodec.StringToBytes(req.Proposer)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(proposer, sdkctx.CometInfo().GetProposerAddress()) || !bytes.Equal(proposer, req.Payload.FeeRecipient) {
		return nil, types.ErrInvalidRequest.Wrap("invalid proposer")
	}

	block, err := k.Block.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	payload := req.Payload
	if !bytes.Equal(block.ParentHash, payload.ParentHash) || block.BlockNumber+1 != payload.BlockNumber {
		return nil, types.ErrInvalidRequest.Wrap("refer to incorrect parent block")
	}

	beaconRoot, err := k.BeaconRoot.Get(sdkctx)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(beaconRoot, payload.BeaconRoot) {
		return nil, types.ErrInvalidRequest.Wrap("refer to incorrect beacon root")
	}

	if err := k.VerifyDequeue(ctx, req.Payload.Transactions); err != nil {
		return nil, types.ErrInvalidRequest.Wrapf("dequeue mismatched: %s", err.Error())
	}

	// todo: handle request from execution node

	if err := k.Block.Set(ctx, req.Payload); err != nil {
		return nil, err
	}

	return &types.MsgNewEthBlockResponse{}, nil
}
