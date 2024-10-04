package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) Claim(ctx context.Context, reqs []*ethtypes.GoatClaimReward) error {
	if len(reqs) == 0 {
		return nil
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	queue, err := k.EthTxQueue.Get(sdkctx)
	if err != nil {
		return err
	}

	for _, req := range reqs {
		valdtAddr := sdktypes.ConsAddress(req.Validator.Bytes())

		validator, err := k.Validators.Get(sdkctx, valdtAddr)
		if err != nil {
			return err
		}

		queue.Rewards = append(queue.Rewards, &types.Reward{
			Id:        req.Id,
			Recipient: req.Recipient.Bytes(),
			Goat:      validator.GoatReward,
			Gas:       validator.GasReward,
		})

		validator.GoatReward = math.ZeroInt()
		validator.GasReward = math.ZeroInt()
		if err := k.Validators.Set(sdkctx, valdtAddr, validator); err != nil {
			return err
		}
	}

	if err := k.EthTxQueue.Set(sdkctx, queue); err != nil {
		return err
	}

	return nil
}
