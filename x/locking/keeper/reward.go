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

func (k Keeper) UpdateRewardPool(ctx context.Context, gas []*ethtypes.GasRevenue, grants []*ethtypes.GoatGrant) error {
	if len(grants) == 0 {
		return nil
	}

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

	// block gas is not zero, which means the block has transactions
	// we only give out the reward if there are transactions in the block
	if !pool.Gas.IsZero() { // todo: how about system txs like deposits?
		reward := big.NewInt(types.InitialBlockReward)
		if halvings := pool.Index / types.HalvingInterval; halvings > 0 {
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

		sdkctx.EventManager().EmitEvent(types.AddRewardEvent(sdkctx.BlockHeight(), reward))
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

	if pool.Gas.IsZero() && pool.Goat.IsZero() {
		return nil
	}

	var totalPower int64 // previous block
	for _, voteInfo := range sdkctx.VoteInfos() {
		totalPower += voteInfo.Validator.Power
	}

	if totalPower == 0 { // should never happend
		return errors.New("invalid zero power")
	}

	// Here we send reward to all validators even if it didn't vote for the block
	// it prevents proposer only including 2/3 votes to get more reward
	// if a validator is down for long, it will be slashed
	for _, voteInfo := range sdkctx.VoteInfos() {
		power := math.LegacyNewDec(voteInfo.Validator.Power).QuoTruncate(math.LegacyNewDec(totalPower))
		if power.IsZero() {
			continue
		}

		validator, err := k.Validators.Get(sdkctx, voteInfo.Validator.Address)
		if err != nil {
			return err
		}

		if !pool.Gas.IsZero() {
			share := math.LegacyNewDecFromBigInt(pool.Gas.BigInt()).MulTruncate(power).RoundInt()
			validator.GasReward = validator.GasReward.Add(share)
		}

		if !pool.Goat.IsZero() {
			share := math.LegacyNewDecFromBigInt(pool.Goat.BigInt()).MulTruncate(power).RoundInt()
			validator.GoatReward = validator.GoatReward.Add(share)
		}
	}

	pool.Gas = math.ZeroInt()
	pool.Goat = math.ZeroInt()
	if err := k.RewardPool.Set(sdkctx, pool); err != nil {
		return err
	}
	return nil
}
