package keeper

import (
	"context"
	"errors"
	"math/big"

	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) UpdateRewardPool(ctx context.Context, gas []*ethtypes.GasRevenue, grants []*ethtypes.GoatGrant, hasTxs bool) error {
	if l := len(gas); l != 1 {
		return types.ErrInvalid.Wrapf("expected gas revenue request length 1 but got %d", l)
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	pool, err := k.RewardPool.Get(sdkctx)
	if err != nil {
		return err
	}

	for _, revenue := range gas {
		pool.Gas = pool.Gas.Add(math.NewIntFromBigInt(revenue.Amount))
	}

	for _, grant := range grants {
		pool.Remain = pool.Remain.Add(math.NewIntFromBigInt(grant.Amount))
	}

	if hasTxs {
		param, err := k.Params.Get(sdkctx)
		if err != nil {
			return err
		}

		reward := big.NewInt(param.InitialBlockReward)
		if halvings := pool.Index / param.HalvingInterval; halvings > 0 {
			count := big.NewInt(2)
			count.Exp(count, big.NewInt(halvings), nil)
			reward.Div(reward, count)
		}

		if remain := pool.Remain.BigInt(); reward.Cmp(remain) < 0 {
			reward = remain
		}

		if reward.Sign() != 0 {
			r := math.NewIntFromBigInt(reward)
			pool.Goat = pool.Goat.Add(r)
			pool.Remain = pool.Remain.Sub(r)
			pool.Index++
		}
	}

	if err := k.RewardPool.Set(sdkctx, pool); err != nil {
		return err
	}

	return nil
}

func (k Keeper) distributeReward(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	// the cometbft consensus rule
	if sdkctx.BlockHeight() < 2 {
		return nil
	}

	pool, err := k.RewardPool.Get(sdkctx)
	if err != nil {
		return err
	}

	if pool.Gas.IsZero() && pool.Goat.IsZero() { // todo: remove if reward distribution is changed
		return nil
	}

	var totalPower int64 // previous block
	for _, voteInfo := range sdkctx.VoteInfos() {
		totalPower += voteInfo.Validator.Power
	}

	if totalPower == 0 { // should never happend
		return errors.New("invalid zero power")
	}

	remainGas, remainReward := pool.Gas.BigInt(), pool.Goat.BigInt()
	// Here we send reward to all validators even if it didn't vote for the block
	// it prevents proposer only including 2/3 votes to get more reward
	for _, voteInfo := range sdkctx.VoteInfos() {
		power := math.LegacyNewDec(voteInfo.Validator.Power).QuoTruncate(math.LegacyNewDec(totalPower))

		validator, err := k.Validators.Get(sdkctx, voteInfo.Validator.Address)
		if err != nil {
			return err
		}

		if !pool.Gas.IsZero() {
			share := math.LegacyNewDecFromBigInt(pool.Gas.BigInt()).MulTruncate(power).TruncateInt()
			remainGas.Sub(remainGas, share.BigIntMut())
			validator.GasReward = validator.GasReward.Add(share)
		}

		if !pool.Goat.IsZero() {
			share := math.LegacyNewDecFromBigInt(pool.Goat.BigInt()).MulTruncate(power).TruncateInt()
			remainReward.Sub(remainReward, share.BigIntMut())
			validator.Reward = validator.Reward.Add(share)
		}
	}

	// give back the dust to pool again
	pool.Gas = math.NewIntFromBigIntMut(remainGas)
	pool.Goat = math.NewIntFromBigIntMut(remainReward)
	if err := k.RewardPool.Set(sdkctx, pool); err != nil {
		return err
	}
	return nil
}
