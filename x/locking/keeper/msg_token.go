package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) UpdateTokens(ctx context.Context, weights []*ethtypes.UpdateTokenWeight, thresholds []*ethtypes.UpdateTokenThreshold) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	for _, update := range weights {
		addr := hex.EncodeToString(update.Token.Bytes())
		token, err := k.Tokens.Get(sdkctx, addr)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				return err
			}
			token = types.Token{Weight: update.Weight}
		}

		if err := k.onWeightChanged(sdkctx, addr, token.Weight, update.Weight); err != nil {
			return err
		}

		token.Weight = update.Weight
		if err := k.Tokens.Set(sdkctx, addr, token); err != nil {
			return err
		}
	}

	if len(thresholds) == 0 {
		return nil
	}

	thrs, err := k.Threshold.Get(sdkctx)
	if err != nil {
		return err
	}

	for _, update := range thresholds {
		addr := hex.EncodeToString(update.Token.Bytes())
		token, err := k.Tokens.Get(sdkctx, addr)
		if err != nil {
			return err
		}

		sub := math.NewIntFromBigInt(update.Threshold).Sub(token.Threshold)
		thrs.List = thrs.List.Add(sdktypes.NewCoin(addr, sub))
		token.Threshold = math.NewIntFromBigInt(update.Threshold)
		if err := k.Tokens.Set(sdkctx, addr, token); err != nil {
			return err
		}
	}

	if err := k.Threshold.Set(sdkctx, thrs); err != nil {
		return err
	}
	return nil
}

func (k Keeper) onWeightChanged(ctx context.Context, token string, previous, current uint64) error {
	if previous == current {
		return nil
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	iter, err := k.Locking.Iterate(sdkctx,
		collections.NewPrefixedPairRange[string, sdktypes.ConsAddress](token))
	if err != nil {
		return err
	}
	defer iter.Close()

	isUp := current > previous
	for ; iter.Valid(); iter.Next() {
		kv, err := iter.KeyValue()
		if err != nil {
			return err
		}

		if t := kv.Key.K1(); t != token {
			return fmt.Errorf("invalid interator: inconsistent token: expected %s got %s", token, t)
		}

		valdtAddr, amount := kv.Key.K2(), kv.Value
		validator, err := k.Validators.Get(sdkctx, valdtAddr)
		if err != nil {
			return err
		}

		if err := k.PowerRanking.Remove(sdkctx,
			collections.Join(validator.Power, valdtAddr)); err != nil {
			return err
		}

		if isUp {
			diff := math.NewIntFromUint64(current - previous).Mul(amount).Quo(types.PowerReduction)
			if !diff.IsUint64() {
				return fmt.Errorf("power too large: %s", diff)
			}
			validator.Power += diff.Uint64()
		} else {
			diff := math.NewIntFromUint64(previous - current).Mul(amount).Quo(types.PowerReduction)
			if validator.Power < diff.Uint64() {
				return fmt.Errorf("previous power is less than updated power when weight is decreased (token %s validator %x)", token, valdtAddr.Bytes())
			}
			validator.Power -= diff.Uint64()
		}
		if err := k.Validators.Set(sdkctx, valdtAddr, validator); err != nil {
			return err
		}

		if validator.Power > 0 {
			if err := k.PowerRanking.Set(sdkctx,
				collections.Join(validator.Power, valdtAddr)); err != nil {
				return err
			}
		}
	}
	return nil
}
