package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) Lock(ctx context.Context, reqs []*goattypes.LockRequest) error {
	if len(reqs) == 0 {
		return nil
	}

	// aggregate
	updates := make(map[common.Address]sdktypes.Coins)
	for _, req := range reqs {
		if _, ok := updates[req.Validator]; !ok {
			updates[req.Validator] = sdktypes.Coins{}
		}
		coin := sdktypes.NewCoin(types.TokenDenom(req.Token), math.NewIntFromBigInt(req.Amount))
		updates[req.Validator] = updates[req.Validator].Add(coin)
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	for validator, coins := range updates {
		if err := k.lock(sdkctx, validator, coins); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) lock(ctx context.Context, target common.Address, coins sdktypes.Coins) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	valdtAddr := sdktypes.ConsAddress(target.Bytes())

	validator, err := k.Validators.Get(sdkctx, valdtAddr)
	if err != nil {
		return err
	}
	validator.Locking = validator.Locking.Add(coins...)

	switch validator.Status {
	case types.Pending, types.Active:
		// remove it from power ranking
		if err := k.PowerRanking.Remove(ctx,
			collections.Join(validator.Power, valdtAddr)); err != nil {
			return err
		}

		// caculate the new power and update locking state
		for _, coin := range coins {
			token, err := k.Tokens.Get(sdkctx, coin.Denom)
			if err != nil {
				return err
			}

			if token.Weight > 0 {
				power := math.NewIntFromUint64(token.Weight).Mul(coin.Amount).Quo(types.PowerReduction)
				if !power.IsUint64() {
					return errorsmod.Wrapf(sdkerrors.ErrLogic, "power too large: %s", power)
				}
				validator.Power += power.Uint64()
			}

			pair := collections.Join(coin.Denom, valdtAddr)
			val, err := k.Locking.Get(sdkctx, pair)
			if err != nil {
				if !errors.Is(err, collections.ErrNotFound) {
					return err
				}
				val = math.ZeroInt()
			}

			if err := k.Locking.Set(sdkctx, pair, val.Add(coin.Amount)); err != nil {
				return err
			}
		}

		// update power ranking
		if err := k.PowerRanking.Set(ctx,
			collections.Join(validator.Power, valdtAddr)); err != nil {
			return err
		}
		k.Logger().Info("Lock", "validator", types.ValidatorName(valdtAddr), "power", validator.Power)
	case types.Downgrade:
		threshold, err := k.Threshold.Get(sdkctx)
		if err != nil {
			return err
		}
		// check if it's unjailed and locking is enough
		if sdkctx.BlockTime().After(validator.JailedUntil) && validator.Locking.IsAllGTE(threshold.List) {
			validator.Status = types.Pending

			for _, coin := range validator.Locking {
				if err := k.Locking.Set(sdkctx,
					collections.Join(coin.Denom, valdtAddr), coin.Amount); err != nil {
					return err
				}

				token, err := k.Tokens.Get(sdkctx, coin.Denom)
				if err != nil {
					return err
				}

				if token.Weight > 0 {
					power := math.NewIntFromUint64(token.Weight).Mul(coin.Amount).Quo(types.PowerReduction)
					if !power.IsUint64() {
						return errorsmod.Wrapf(sdkerrors.ErrLogic, "power too large: %s", power)
					}
					validator.Power += power.Uint64()
				}
			}

			if err := k.PowerRanking.Set(ctx,
				collections.Join(validator.Power, valdtAddr)); err != nil {
				return err
			}
			k.Logger().Info("Unjail", "validator", types.ValidatorName(valdtAddr), "power", validator.Power)
		}
	case types.Tombstoned, types.Inactive:
		// don't do anything
		k.Logger().Info("No effect lock", "validator", types.ValidatorName(valdtAddr), "power", validator.Power)
	}

	// update validator state
	if err := k.Validators.Set(sdkctx, valdtAddr, validator); err != nil {
		return err
	}

	return nil
}
