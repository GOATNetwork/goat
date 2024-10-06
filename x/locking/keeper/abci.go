package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	abci "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) BeginBlocker(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	if err := k.distributeReward(sdkctx); err != nil {
		return err
	}
	if err := k.dequeueMatureUnlocks(sdkctx); err != nil {
		return err
	}
	if err := k.HandleVoteInfo(sdkctx); err != nil {
		return err
	}
	if err := k.HandleEvidences(sdkctx); err != nil {
		return err
	}
	return nil
}

func (k Keeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	lastSet := make(map[string]uint64)
	{
		iter, err := k.ValidatorSet.Iterate(sdkctx, nil)
		if err != nil {
			return nil, err
		}
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			kv, err := iter.KeyValue()
			if err != nil {
				return nil, err
			}
			lastSet[string(kv.Key)] = kv.Value
		}
	}

	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	pwIter, err := k.PowerRanking.Iterate(sdkctx,
		new(collections.PairRange[uint64, sdktypes.ConsAddress]).Descending())
	if err != nil {
		return nil, err
	}
	defer pwIter.Close()

	var newSet []abci.ValidatorUpdate
	var lastPower uint64
	for count := 0; pwIter.Valid() && count < int(param.MaxValidators); pwIter.Next() {
		key, err := pwIter.Key()
		if err != nil {
			return nil, err
		}
		curPower, valAddr := key.K1(), key.K2()
		if curPower > lastPower {
			return nil, errors.New("invalid iterator: validator power is bigger than before")
		}
		lastPower = curPower

		validator, err := k.Validators.Get(sdkctx, valAddr)
		if err != nil {
			return nil, err
		}
		valstr := string(valAddr)

		switch validator.Status {
		case types.ValidatorStatus_Active:
			if power := lastSet[valstr]; power != validator.Power { // the power is changed
				newSet = append(newSet, abci.ValidatorUpdate{
					Power: int64(validator.Power), PubKey: validator.CMPubkey()})
			}
			delete(lastSet, valstr)
		case types.ValidatorStatus_Pending:
			if _, ok := lastSet[valstr]; ok {
				return nil, fmt.Errorf("pending validator %x existed in the last validator set", valAddr.Bytes())
			}
			validator.Status = types.ValidatorStatus_Active
			validator.SigningInfo.Missed = 0
			validator.SigningInfo.Offset = 0
			if err := k.Validators.Set(sdkctx, valAddr, validator); err != nil {
				return nil, err
			}
			if err := k.ValidatorSet.Set(sdkctx, valAddr, validator.Power); err != nil {
				return nil, err
			}
			newSet = append(newSet, abci.ValidatorUpdate{
				Power: int64(validator.Power), PubKey: validator.CMPubkey()})
		default:
			return nil, fmt.Errorf("%s validator %x in power ranking", validator.Status, valAddr.Bytes())
		}
		count++
	}

	for val := range lastSet {
		valAddr := []byte(val)
		validator, err := k.Validators.Get(sdkctx, valAddr)
		if err != nil {
			return nil, err
		}
		if validator.Status == types.ValidatorStatus_Active {
			validator.Status = types.ValidatorStatus_Pending
			if err := k.Validators.Set(sdkctx, valAddr, validator); err != nil {
				return nil, err
			}
		}
		if err := k.ValidatorSet.Remove(sdkctx, valAddr); err != nil {
			return nil, err
		}
		newSet = append(newSet, abci.ValidatorUpdate{PubKey: validator.CMPubkey()})
	}
	return newSet, nil
}
