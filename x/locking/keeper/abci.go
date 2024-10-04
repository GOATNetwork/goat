package keeper

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	if err := k.distributeReward(sdkctx); err != nil {
		return err
	}
	if err := k.dequeueMatureUnlocks(ctx); err != nil {
		return err
	}
	if err := k.HandleVoteInfo(sdkctx); err != nil {
		return err
	}
	if err := k.HandleEvidences(sdkctx); err != nil {
		return err
	}
	return nil
}

func (k Keeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	return k.blockValidatorUpdates(ctx)
}
