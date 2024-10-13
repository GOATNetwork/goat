package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
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
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "consensus proposer mismatched")
	}

	block, err := k.Block.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	payload := req.Payload
	if payload == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "empty payload")
	}

	if !bytes.Equal(block.BlockHash, payload.ParentHash) || block.BlockNumber+1 != payload.BlockNumber {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "incorrect parent block")
	}

	if payload.BlobGasUsed > 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "blob tx is not allowed")
	}

	beaconRoot, err := k.BeaconRoot.Get(sdkctx)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(beaconRoot, payload.BeaconRoot) {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "invalid beacon root")
	}

	if err := k.VerifyDequeue(sdkctx, payload.ExtraData, payload.Transactions); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "dequeue mismatched")
	}

	bridgeReq, relayerReq, lockingReq, err := goattypes.DecodeRequests(payload.Requests)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "invalid execution requests")
	}

	if err := k.lockingKeeper.ProcessLockingRequest(sdkctx, lockingReq); err != nil {
		return nil, err
	}

	if err := k.bitcoinKeeper.ProcessBridgeRequest(sdkctx, bridgeReq); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.ProcessRelayerRequest(sdkctx, relayerReq); err != nil {
		return nil, err
	}

	if err := k.Block.Set(sdkctx, *payload); err != nil {
		return nil, err
	}

	// Update beacon root
	if err := k.BeaconRoot.Set(sdkctx, sdkctx.HeaderHash()); err != nil {
		return nil, err
	}

	sdkctx.EventManager().EmitEvent(types.NewEthBlockEvent(req.Payload.BlockNumber, req.Payload.BlockHash))
	return &types.MsgNewEthBlockResponse{}, nil
}
