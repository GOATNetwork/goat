package keeper

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
)

func (k Keeper) BeginBlock(ctx context.Context) error {
	if err := k.distributeReward(ctx); err != nil {
		return err
	}
	return nil
}

func (k Keeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	if err := k.dequeueMatureUnlocks(ctx); err != nil {
		return nil, err
	}
	return k.blockValidatorUpdates(ctx)
}
