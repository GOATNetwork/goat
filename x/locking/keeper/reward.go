package keeper

import (
	"context"
	"errors"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) UpdateRewardPool(ctx context.Context, gas []*goattypes.GasRequest, grants []*goattypes.GrantRequest) error {
	if l := len(gas); l != 1 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "expected gas revenue request length 1 but got %d", l)
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	pool, err := k.RewardPool.Get(sdkctx)
	if err != nil {
		return err
	}

	for _, revenue := range gas {
		if revenue.Amount.Sign() > 0 {
			pool.Gas = pool.Gas.Add(math.NewIntFromBigIntMut(revenue.Amount))
			k.Logger().Debug("Add gas fee", "amount", revenue.Amount.String())
		}
	}

	for _, grant := range grants {
		pool.Remain = pool.Remain.Add(math.NewIntFromBigIntMut(grant.Amount))
		k.Logger().Debug("Grant reward", "amount", grant.Amount.String())
	}

	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return err
	}

	reward := big.NewInt(param.InitialBlockReward)
	if halvings := sdkctx.BlockHeight() / param.HalvingInterval; halvings > 0 {
		count := big.NewInt(2)
		count.Exp(count, big.NewInt(halvings), nil)
		reward.Div(reward, count)
	}

	if remain := pool.Remain.BigInt(); reward.Cmp(remain) > 0 {
		reward = remain
	}

	if reward.Sign() != 0 {
		r := math.NewIntFromBigInt(reward)
		pool.Goat = pool.Goat.Add(r)
		pool.Remain = pool.Remain.Sub(r)
	}
	k.Logger().Debug("Add reward to pool", "amount", reward.String())

	if err := k.RewardPool.Set(sdkctx, pool); err != nil {
		return err
	}

	return nil
}

func (k Keeper) DistributeReward(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	// the cometbft consensus rule
	if sdkctx.BlockHeight() < 2 {
		return nil
	}

	pool, err := k.RewardPool.Get(sdkctx)
	if err != nil {
		return err
	}

	var totalPower int64 // previous block
	for _, voteInfo := range sdkctx.VoteInfos() {
		totalPower += voteInfo.Validator.Power
	}

	if totalPower == 0 { // should never happened
		return errors.New("invalid zero power")
	}

	remainGas, remainReward := pool.Gas.BigInt(), pool.Goat.BigInt()
	// Here we send reward to all validators even if it didn't vote for the block
	// it prevents proposer only including 2/3 votes to get more reward
	for _, voteInfo := range sdkctx.VoteInfos() {
		validator, err := k.Validators.Get(sdkctx, voteInfo.Validator.Address)
		if err != nil {
			return err
		}

		power := math.LegacyNewDec(voteInfo.Validator.Power).Quo(math.LegacyNewDec(totalPower))
		if !pool.Gas.IsZero() {
			share := math.LegacyNewDecFromBigInt(pool.Gas.BigInt()).MulTruncate(power).TruncateInt()
			if !share.IsZero() {
				remainGas.Sub(remainGas, share.BigIntMut())
				validator.GasReward = validator.GasReward.Add(share)
				k.Logger().Debug("Distribute gas reward",
					"address", types.ValidatorName(voteInfo.Validator.Address), "amount", share)
			}
		}

		if !pool.Goat.IsZero() {
			share := math.LegacyNewDecFromBigInt(pool.Goat.BigInt()).MulTruncate(power).TruncateInt()
			if !share.IsZero() {
				remainReward.Sub(remainReward, share.BigIntMut())
				validator.Reward = validator.Reward.Add(share)
				k.Logger().Debug("Distribute goat reward",
					"address", types.ValidatorName(voteInfo.Validator.Address), "amount", share)
			}
		}

		if err := k.Validators.Set(sdkctx, voteInfo.Validator.Address, validator); err != nil {
			return err
		}
	}

	// give back the dust to the pool again
	pool.Gas = math.NewIntFromBigIntMut(remainGas)
	pool.Goat = math.NewIntFromBigIntMut(remainReward)
	if err := k.RewardPool.Set(sdkctx, pool); err != nil {
		return err
	}
	return nil
}
