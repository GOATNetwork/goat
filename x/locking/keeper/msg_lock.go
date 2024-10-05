package keeper

import (
	"context"
	"encoding/hex"
	"errors"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) Lock(ctx context.Context, reqs []*ethtypes.GoatLock) error {
	if len(reqs) == 0 {
		return nil
	}

	updates := make(map[common.Address]sdktypes.Coins)
	for _, req := range reqs {
		if _, ok := updates[req.Validator]; !ok {
			updates[req.Validator] = sdktypes.Coins{}
		}
		coin := sdktypes.NewCoin(hex.EncodeToString(req.Token.Bytes()), math.NewIntFromBigInt(req.Amount))
		updates[req.Validator] = updates[req.Validator].Add(coin)
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return err
	}

	for validator, coins := range updates {
		if err := k.lock(sdkctx, validator, coins, &param); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) lock(ctx context.Context, target common.Address, coins sdktypes.Coins, param *types.Params) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	valdtAddr := sdktypes.ConsAddress(target.Bytes())

	validator, err := k.Validators.Get(sdkctx, valdtAddr)
	if err != nil {
		return err
	}
	validator.Locking = validator.Locking.Add(coins...)

	switch validator.Status {
	case types.ValidatorStatus_Pending, types.ValidatorStatus_Active:
		if err := k.PowerRanking.Remove(ctx,
			collections.Join(validator.Power, valdtAddr)); err != nil {
			return err
		}

		for _, coin := range coins {
			token, err := k.Tokens.Get(sdkctx, coin.Denom)
			if err != nil {
				return err
			}

			if token.Weight > 0 {
				power := math.NewIntFromUint64(token.Weight).Mul(coin.Amount).Quo(types.PowerReduction)
				if !power.IsUint64() {
					return types.ErrInvalid.Wrapf("power too large: %s", power)
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

		if err := k.PowerRanking.Set(ctx,
			collections.Join(validator.Power, valdtAddr)); err != nil {
			return err
		}

		if err := k.Validators.Set(sdkctx, valdtAddr, validator); err != nil {
			return err
		}
	case types.ValidatorStatus_Downgrade:
		threshold, err := k.Threshold.Get(sdkctx)
		if err != nil {
			return err
		}

		if sdkctx.BlockTime().After(validator.SigningInfo.JailedUntil) && validator.Locking.IsAllGTE(threshold.List) {
			validator.Status = types.ValidatorStatus_Pending

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
						return types.ErrInvalid.Wrapf("power too large: %s", power)
					}
					validator.Power += power.Uint64()
				}
			}

			if err := k.PowerRanking.Set(ctx,
				collections.Join(validator.Power, valdtAddr)); err != nil {
				return err
			}
		}
		if err := k.Validators.Set(sdkctx, valdtAddr, validator); err != nil {
			return err
		}
	case types.ValidatorStatus_Tombstoned, types.ValidatorStatus_Inactive:
		if err := k.Validators.Set(sdkctx, valdtAddr, validator); err != nil {
			return err
		}
		return nil
	}

	return nil
}
