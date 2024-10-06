package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) Unlock(ctx context.Context, reqs []*ethtypes.GoatUnlock) error {
	if len(reqs) == 0 {
		return nil
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return err
	}
	for i := 0; i < len(reqs); i++ {
		if err := k.unlock(sdkctx, reqs[i], &param); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) unlock(ctx context.Context, req *ethtypes.GoatUnlock, param *types.Params) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	tokenAddr := hex.EncodeToString(req.Token.Bytes())
	valdtAddr := sdktypes.ConsAddress(req.Validator.Bytes())

	validator, err := k.Validators.Get(sdkctx, valdtAddr)
	if err != nil {
		return err
	}

	if err := k.PowerRanking.Remove(sdkctx,
		collections.Join(validator.Power, valdtAddr)); err != nil {
		return err
	}

	token, err := k.Tokens.Get(sdkctx, tokenAddr)
	if err != nil {
		return err
	}

	amount := math.NewIntFromBigInt(req.Amount) // max amount to send back
	lockingAmount := validator.Locking.AmountOf(tokenAddr)
	if lockingAmount.LT(amount) { // the validator was slashed
		amount = math.NewIntFromBigIntMut(lockingAmount.BigInt())
	}
	validator.Locking = validator.Locking.Sub(sdktypes.NewCoin(tokenAddr, amount))
	lockingAmount = lockingAmount.Sub(amount)

	var powerU64 uint64
	if !amount.IsZero() && token.Weight > 0 && validator.Status != types.ValidatorStatus_Tombstoned {
		p := math.NewIntFromUint64(token.Weight).Mul(amount).Quo(types.PowerReduction)
		if !p.IsUint64() {
			return types.ErrInvalid.Wrapf("power too large: %s", p)
		}
		powerU64 = p.Uint64()
	}

	if validator.Power > powerU64 {
		validator.Power -= powerU64
	} else { // the latest weight is bigger than before
		validator.Power = 0
	}

	exiting := validator.Status == types.ValidatorStatus_Inactive ||
		validator.Status == types.ValidatorStatus_Tombstoned || lockingAmount.LT(token.Threshold)
	var unlockTime time.Time
	if exiting {
		unlockTime = sdkctx.BlockTime().Add(param.ExitingDuration)

		switch validator.Status {
		case types.ValidatorStatus_Active, types.ValidatorStatus_Pending:
			validator.Status = types.ValidatorStatus_Inactive
		}

		// remove all from locking state
		for _, coin := range validator.Locking {
			if err := k.Locking.Remove(sdkctx, collections.Join(coin.Denom, valdtAddr)); err != nil {
				return err
			}
		}
	} else {
		unlockTime = sdkctx.BlockTime().Add(param.UnlockDuration)

		// remove if there is no locking
		if lockingAmount.IsZero() {
			if err := k.Locking.Remove(sdkctx, collections.Join(tokenAddr, valdtAddr)); err != nil {
				return err
			}
		} else {
			if err := k.Locking.Set(sdkctx,
				collections.Join(tokenAddr, valdtAddr), lockingAmount); err != nil {
				return err
			}
		}

		if err := k.PowerRanking.Set(sdkctx, collections.Join(validator.Power, valdtAddr)); err != nil {
			return err
		}
	}

	// get the unlock queue by the time
	queue, err := k.UnlockQueue.Get(sdkctx, unlockTime)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return err
	}

	// append it
	queue.Unlocks = append(queue.Unlocks, &types.Unlock{
		Id:        req.Id,
		Token:     req.Token.Bytes(),
		Recipient: req.Recipient.Bytes(),
		Amount:    amount,
	})

	if err := k.UnlockQueue.Set(sdkctx, unlockTime, queue); err != nil {
		return err
	}

	if err := k.Validators.Set(sdkctx, valdtAddr, validator); err != nil {
		return err
	}

	return nil
}

func (k Keeper) dequeueMatureUnlocks(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	var keys []time.Time
	var values []*types.Unlock

	rng := (&collections.Range[time.Time]{}).EndInclusive(sdkctx.BlockTime())
	if err := k.UnlockQueue.Walk(ctx, rng, func(key time.Time, value types.Unlocks) (bool, error) {
		keys = append(keys, key)
		values = append(values, value.Unlocks...)
		return false, nil
	}); err != nil {
		return err
	}

	for _, key := range keys {
		err := k.UnlockQueue.Remove(ctx, key)
		if err != nil {
			return err
		}
	}

	if len(values) == 0 {
		return nil
	}

	execQueue, err := k.EthTxQueue.Get(sdkctx)
	if err != nil {
		return err
	}

	execQueue.Unlocks = append(execQueue.Unlocks, values...)
	if err := k.EthTxQueue.Set(sdkctx, execQueue); err != nil {
		return err
	}

	return nil
}
