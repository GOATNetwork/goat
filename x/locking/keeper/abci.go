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
	if err := k.UpdateForkParams(sdkctx); err != nil {
		return err
	}
	if err := k.DistributeReward(sdkctx); err != nil {
		return err
	}
	if height, ok := types.TzngForkHeight[sdkctx.ChainID()]; ok && sdkctx.BlockHeight() < height {
		if err := k.DequeueMatureUnlocks(sdkctx); err != nil {
			return err
		}
	}
	if err := k.HandleVoteInfos(sdkctx); err != nil {
		return err
	}
	if err := k.HandleEvidences(sdkctx); err != nil {
		return err
	}
	return nil
}

func (k Keeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	if height := types.TzngForkHeight[sdkctx.ChainID()]; sdkctx.BlockHeight() >= height {
		if err := k.DequeueMatureUnlocks(sdkctx); err != nil {
			return nil, err
		}
	}

	// set finalized time after osaka fork, enable by default(use >= symbol)
	if osakaHeight := types.OsakaForkHeight[sdkctx.ChainID()]; sdkctx.BlockHeight() >= osakaHeight {
		if err := k.FinalizedTime.Set(sdkctx, sdkctx.BlockTime()); err != nil {
			return nil, err
		}
	}

	return k.updateValidatorSet(sdkctx)
}

func (k Keeper) updateValidatorSet(sdkctx sdktypes.Context) ([]abci.ValidatorUpdate, error) {
	lastSet := make(map[string]uint64)
	{
		iter, err := k.ValidatorSet.Iterate(sdkctx, nil)
		if err != nil {
			return nil, err
		}
		for ; iter.Valid(); iter.Next() {
			kv, err := iter.KeyValue()
			if err != nil {
				return nil, err
			}
			lastSet[string(kv.Key)] = kv.Value
		}
		if err := iter.Close(); err != nil {
			return nil, err
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

	var newSet []abci.ValidatorUpdate
	var lastPower uint64
	for count := int64(0); pwIter.Valid() && count < param.MaxValidators; pwIter.Next() {
		key, err := pwIter.Key()
		if err != nil {
			return nil, err
		}

		curPower, valAddr := key.K1(), key.K2()
		if curPower > lastPower && count != 0 {
			return nil, errors.New("invalid iterator: validator power is bigger than before")
		}
		lastPower = curPower

		validator, err := k.Validators.Get(sdkctx, valAddr)
		if err != nil {
			return nil, err
		}
		valstr := string(valAddr)

		switch validator.Status {
		case types.Active:
			if power := lastSet[valstr]; power != validator.Power { // the power is changed
				if err := k.ValidatorSet.Set(sdkctx, valAddr, validator.Power); err != nil {
					return nil, err
				}
				newSet = append(newSet, abci.ValidatorUpdate{
					Power: int64(validator.Power), PubKey: validator.CMPubkey(),
				})
				k.Logger().Info("Update validator set", "address", types.ValidatorName(valAddr), "power", validator.Power)
			}
			delete(lastSet, valstr)
		case types.Pending:
			if _, ok := lastSet[valstr]; ok {
				return nil, fmt.Errorf("pending validator %x existed in the last validator set", valAddr.Bytes())
			}
			validator.Status = types.Active
			validator.SigningInfo = types.SigningInfo{}
			if err := k.Validators.Set(sdkctx, valAddr, validator); err != nil {
				return nil, err
			}
			if err := k.ValidatorSet.Set(sdkctx, valAddr, validator.Power); err != nil {
				return nil, err
			}
			newSet = append(newSet, abci.ValidatorUpdate{
				Power: int64(validator.Power), PubKey: validator.CMPubkey(),
			})
			k.Logger().Info("Add to validator set", "address", types.ValidatorName(valAddr), "power", validator.Power)
		default:
			return nil, fmt.Errorf("%s validator %x in power ranking", validator.Status, valAddr.Bytes())
		}
		count++
	}

	if err := pwIter.Close(); err != nil {
		return nil, err
	}

	// remove
	for val := range lastSet {
		valAddr := []byte(val)
		validator, err := k.Validators.Get(sdkctx, valAddr)
		if err != nil {
			return nil, err
		}
		if validator.Status == types.Active {
			validator.Status = types.Pending
			if err := k.Validators.Set(sdkctx, valAddr, validator); err != nil {
				return nil, err
			}
		}
		if err := k.ValidatorSet.Remove(sdkctx, valAddr); err != nil {
			return nil, err
		}
		newSet = append(newSet, abci.ValidatorUpdate{PubKey: validator.CMPubkey()})
		k.Logger().Info("Remove from validator set", "address", types.ValidatorName(valAddr), "power", validator.Power)
	}
	return newSet, nil
}
