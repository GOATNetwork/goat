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

	cometProposer := sdkctx.CometInfo().GetProposerAddress()
	if !bytes.Equal(proposer, cometProposer) || !bytes.Equal(proposer, req.Payload.FeeRecipient) {
		return nil, types.ErrInvalidRequest.Wrap("invalid proposer")
	}

	block, err := k.Block.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	payload := req.Payload
	if payload == nil {
		return nil, types.ErrInvalidRequest.Wrap("empty payload")
	}

	if payload.GasRevenue == nil {
		return nil, types.ErrInvalidRequest.Wrap("no gas revenue request")
	}

	if !bytes.Equal(block.BlockHash, payload.ParentHash) || block.BlockNumber+1 != payload.BlockNumber {
		return nil, types.ErrInvalidRequest.Wrap("refer to incorrect parent block")
	}

	if payload.BlobGasUsed > 0 {
		return nil, types.ErrInvalidRequest.Wrap("blob tx is not allowed")
	}

	if cometTime := uint64(sdkctx.BlockTime().UTC().Unix()); payload.Timestamp < cometTime {
		return nil, types.ErrInvalidRequest.Wrap("invalid timestamp")
	}

	beaconRoot, err := k.BeaconRoot.Get(sdkctx)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(beaconRoot, payload.BeaconRoot) {
		return nil, types.ErrInvalidRequest.Wrap("refer to inconsistent beacon root")
	}

	if err := k.VerifyDequeue(sdkctx, payload.ExtraData, payload.Transactions); err != nil {
		return nil, types.ErrInvalidRequest.Wrap("dequeue mismatched")
	}

	if err := k.Block.Set(sdkctx, *payload); err != nil {
		return nil, err
	}

	if err := k.bitcoinKeeper.ProcessBridgeRequest(sdkctx, payload.WithdrawalReq, payload.RbfReq, payload.Cancel1Req); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.ProcessRelayerRequest(sdkctx, payload.AddVoterReq, payload.RmVoterReq); err != nil {
		return nil, err
	}

	// Update beacon root
	if err := k.BeaconRoot.Set(sdkctx, sdkctx.HeaderHash()); err != nil {
		return nil, err
	}

	sdkctx.EventManager().EmitEvent(types.NewEthBlockEvent(req.Payload.BlockNumber, req.Payload.BlockHash))
	return &types.MsgNewEthBlockResponse{}, nil
}
