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

	events := make(sdktypes.Events, 0, len(req.Deposits))
	deposits := make([]*types.DepositReceipt, 0, len(req.Deposits))
	for _, v := range req.Deposits {
		depoist, err := k.VerifyDeposit(ctx, v)
		if err != nil {
			return nil, err
		}
		events = append(events, types.NewDepositEvent(depoist))
		deposits = append(deposits, depoist)
	}

	queue, err := k.ExecuableQueue.Get(ctx)
	if err != nil {
		return nil, err
	}

	queue.Deposits = append(queue.Deposits, deposits...)
	if err := k.ExecuableQueue.Set(ctx, queue); err != nil {
		return nil, err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(events)
	return &types.MsgNewDepositsResponse{}, nil
}

func (k msgServer) NewBlockHashes(ctx context.Context, req *types.MsgNewBlockHashes) (*types.MsgNewBlockHashesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, types.ErrInvalidRequest.Wrap(err.Error())
	}

	parentHeight, err := k.BlockTip.Peek(ctx)
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

	events := make(sdktypes.Events, 0, len(req.BlockHash))
	for _, v := range req.BlockHash {
		parentHeight++
		if err := k.BlockHashes.Set(ctx, parentHeight, v); err != nil {
			return nil, err
		}
		events = append(events, types.NewBlockHashEvent(parentHeight, v))
	}

	if err := k.BlockTip.Set(ctx, parentHeight); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.SetProposalSeq(ctx, sequence+1); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.UpdateRandao(ctx, req); err != nil {
		return nil, err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(events)
	return &types.MsgNewBlockHashesResponse{}, nil
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
	if err := k.Pubkey.Set(ctx, *req.Pubkey); err != nil {
		return nil, err
	}
	if err := k.relayerKeeper.UpdateRandao(ctx, req); err != nil {
		return nil, err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(
		sdktypes.Events{types.NewKeyEvent(req.Pubkey), relayertypes.ProposalDoneEvent(sequence)},
	)

	k.Logger().Debug("NewKey added", "type", rawKey[0], "key", hex.EncodeToString(rawKey[1:]))
	return &types.MsgNewPubkeyResponse{}, nil
}
