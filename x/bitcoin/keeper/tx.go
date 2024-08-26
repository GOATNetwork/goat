package keeper

import (
	"context"
	"encoding/hex"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
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

func (k msgServer) NewDeposits(ctx context.Context, req *types.MsgNewDeposits) (*types.MsgNewDepositsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, types.ErrInvalidRequest.Wrap(err.Error())
	}

	panic("todo")
}

func (k msgServer) NewBlocks(ctx context.Context, req *types.MsgNewBlocks) (*types.MsgNewBlocksResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, types.ErrInvalidRequest.Wrap(err.Error())
	}

	parentHeight, err := k.LastestHeight.Peek(ctx)
	if err != nil {
		return nil, err
	}
	if req.StartBlockNumber != parentHeight+1 {
		return nil, types.ErrInvalidRequest.Wrapf("block number is not the next of the current %d", parentHeight)
	}

	sequence, err := k.relayerKeeper.VerifyProposal(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, v := range req.BlockHash {
		parentHeight++
		if err := k.BlockHashs.Set(ctx, parentHeight, v); err != nil {
			return nil, err
		}
	}

	if err := k.LastestHeight.Set(ctx, parentHeight); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.SetProposalSeq(ctx, sequence+1); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.UpdateRandao(ctx, req); err != nil {
		return nil, err
	}

	return &types.MsgNewBlocksResponse{}, nil
}

func (k msgServer) NewPubkey(ctx context.Context, req *types.MsgNewPubkey) (*types.MsgNewPubkeyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, types.ErrInvalidRequest.Wrap(err.Error())
	}

	sequence, err := k.relayerKeeper.VerifyProposal(ctx, req)
	if err != nil {
		return nil, err
	}

	rawKey := relayertypes.EncodePublicKey(req.Pubkey)
	exists, err := k.relayerKeeper.HasPubkey(ctx, rawKey)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, relayertypes.ErrInvalidPubkey.Wrap("the key already existed")
	}

	if err := k.relayerKeeper.AddNewKey(ctx, rawKey); err != nil {
		return nil, err
	}
	if err := k.relayerKeeper.SetProposalSeq(ctx, sequence+1); err != nil {
		return nil, err
	}
	if err := k.LatestPubkey.Set(ctx, *req.Pubkey); err != nil {
		return nil, err
	}
	if err := k.relayerKeeper.UpdateRandao(ctx, req); err != nil {
		return nil, err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(
		sdktypes.Events{types.NewKeyEvent(rawKey), relayertypes.ProposalDoneEvent(sequence)},
	)

	k.Logger().Debug("NewKey added", "type", rawKey[0], "key", hex.EncodeToString(rawKey[1:]))
	return &types.MsgNewPubkeyResponse{}, nil
}
