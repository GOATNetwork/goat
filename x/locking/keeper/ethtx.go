package keeper

import (
	"context"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func (k Keeper) DequeueLockingModuleTx(ctx context.Context) ([]*ethtypes.Transaction, error) {
	return nil, nil
}
