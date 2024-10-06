package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/relayer/types"
)

func (k Keeper) ProcessRelayerRequest(ctx context.Context, req types.ExecRequests) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	events := make(sdktypes.Events, 0, len(req.AddVoters)+len(req.RemoveVoters))

	height := uint64(sdkctx.BlockHeight())
	for _, add := range req.AddVoters {
		addr, err := k.AddrCodec.BytesToString(add.Voter[:])
		if err != nil {
			return err
		}
		exists, err := k.Voters.Has(sdkctx, addr)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		if err := k.Voters.Set(sdkctx, addr, types.Voter{
			Address: add.Voter.Bytes(),
			VoteKey: add.Pubkey.Bytes(),
			Height:  height,
			Status:  types.VOTER_STATUS_PENDING,
		}); err != nil {
			return err
		}
		k.Logger().Info("New on-boarding voter", "voter", addr)
		events = append(events, types.PendingVoterEvent(addr))
	}

	if len(req.RemoveVoters) == 0 {
		sdkctx.EventManager().EmitEvents(events)
		return nil
	}

	queue, err := k.Queue.Get(sdkctx)
	if err != nil {
		return err
	}

	// get current voter count in the active set
	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return err
	}

	active := len(relayer.Voters) + 1 - len(queue.OffBoarding)
	for _, rm := range req.RemoveVoters {
		addr, err := k.AddrCodec.BytesToString(rm.Voter.Bytes())
		if err != nil {
			return err
		}

		voter, err := k.Voters.Get(sdkctx, addr)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				continue
			}
			return err
		}

		if voter.Status != types.VOTER_STATUS_ACTIVATED {
			continue
		}

		active--
		if active < 1 {
			k.Logger().Warn("requires 1 voter at least in active set, disregard the removal", "voter", addr)
			break
		}

		k.Logger().Info("New off-boarding voter", "voter", addr)
		voter.Status = types.VOTER_STATUS_OFF_BOARDING
		if err := k.Voters.Set(sdkctx, addr, voter); err != nil {
			return err
		}

		events = append(events, types.RemovingVoterEvent(addr))
		queue.OffBoarding = append(queue.OffBoarding, addr)
	}

	if err := k.Queue.Set(sdkctx, queue); err != nil {
		return err
	}

	sdkctx.EventManager().EmitEvents(events)
	return nil
}
