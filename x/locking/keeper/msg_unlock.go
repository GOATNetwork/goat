package keeper

import (
	"context"
	"errors"
	"time"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) Unlock(ctx context.Context, reqs []*goattypes.UnlockRequest) error {
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

func (k Keeper) unlock(sdkctx sdktypes.Context, req *goattypes.UnlockRequest, param *types.Params) error {
	tokenAddr := types.TokenDenom(req.Token)
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
	updatedLocking := validator.Locking.Sub(sdktypes.NewCoin(tokenAddr, amount))
	lockingAmount = lockingAmount.Sub(amount)

	isExiting := validator.Status == types.Inactive ||
		validator.Status == types.Tombstoned || lockingAmount.LT(token.Threshold)

	if !amount.IsZero() && token.Weight > 0 && !isExiting {
		if validator.Status == types.Active || validator.Status == types.Pending {
			p := math.NewIntFromUint64(token.Weight).Mul(amount).Quo(types.PowerReduction)
			if !p.IsUint64() {
				return errorsmod.Wrapf(sdkerrors.ErrLogic, "power too large: %s", p)
			}

			if powerU64 := p.Uint64(); validator.Power > powerU64 {
				validator.Power -= powerU64
			} else {
				validator.Power = 0
			}
		}
	}

	var unlockTime time.Time
	if isExiting {
		unlockTime = sdkctx.BlockTime().Add(param.ExitingDuration)

		validator.Power = 0
		switch validator.Status {
		case types.Active, types.Pending:
			validator.Status = types.Inactive
		}

		// remove all from locking state
		for _, coin := range validator.Locking {
			if err := k.Locking.Remove(sdkctx, collections.Join(coin.Denom, valdtAddr)); err != nil {
				return err
			}
		}
		validator.Locking = updatedLocking
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

		// insert the power to the ranking again
		if validator.Power > 0 {
			if err := k.PowerRanking.Set(sdkctx,
				collections.Join(validator.Power, valdtAddr)); err != nil {
				return err
			}
		}
		validator.Locking = updatedLocking
	}

	// update the validator state
	if err := k.Validators.Set(sdkctx, valdtAddr, validator); err != nil {
		return err
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

	k.Logger().Info("Unlock", "id", req.Id, "validator", types.ValidatorName(valdtAddr),
		"token", tokenAddr, "amount", amount, "unlockTime", unlockTime)
	return nil
}

func (k Keeper) DequeueMatureUnlocks(ctx context.Context) error {
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

	if len(keys) == 0 {
		return nil
	}

	for _, key := range keys {
		err := k.UnlockQueue.Remove(ctx, key)
		if err != nil {
			return err
		}
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
