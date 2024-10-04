package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	abci "github.com/cometbft/cometbft/abci/types"
	tmcrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) activeValdiators(ctx context.Context) (map[string]bool, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	vsIter, err := k.ValidatorSet.Iterate(sdkctx, nil)
	if err != nil {
		return nil, err
	}
	defer vsIter.Close()

	lastVs, err := vsIter.Keys()
	if err != nil {
		return nil, err
	}

	res := make(map[string]bool)
	for i := 0; i < len(lastVs); i++ {
		res[string(lastVs[i])] = true
	}
	return res, nil
}

func (k Keeper) blockValidatorUpdates(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	lastValidators, err := k.activeValdiators(sdkctx)
	if err != nil {
		return nil, err
	}

	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	powerIter, err := k.PowerRanking.Iterate(sdkctx,
		new(collections.PairRange[uint64, sdktypes.ConsAddress]).Descending())
	if err != nil {
		return nil, err
	}
	defer powerIter.Close()

	var updated []abci.ValidatorUpdate
	for count := 0; powerIter.Valid() && count < int(param.MaxValidators); powerIter.Next() {
		key, err := powerIter.Key()
		if err != nil {
			return nil, err
		}
		_, valAddr := key.K1(), key.K2()
		validator, err := k.Validators.Get(sdkctx, valAddr)
		if err != nil {
			return nil, err
		}
		switch validator.Status {
		case types.ValidatorStatus_Active:
			delete(lastValidators, string(valAddr))
		case types.ValidatorStatus_Pending:
			validator.Status = types.ValidatorStatus_Active
			if err := k.Validators.Set(sdkctx, valAddr, validator); err != nil {
				return nil, err
			}
			if !lastValidators[string(valAddr)] { // unlocking
				if err := k.ValidatorSet.Set(sdkctx, valAddr); err != nil {
					return nil, err
				}
			} else {
				delete(lastValidators, string(valAddr))
			}
		default:
			return nil, fmt.Errorf("validator %x should in the power ranking(%s)", valAddr, validator.Status)
		}
		updated = append(updated, abci.ValidatorUpdate{
			Power:  int64(validator.Power),
			PubKey: tmcrypto.PublicKey{Sum: &tmcrypto.PublicKey_Secp256K1{Secp256K1: validator.Pubkey}},
		})
	}

	for val := range lastValidators {
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
		updated = append(updated, abci.ValidatorUpdate{
			Power:  0,
			PubKey: tmcrypto.PublicKey{Sum: &tmcrypto.PublicKey_Secp256K1{Secp256K1: validator.Pubkey}},
		})
	}
	return updated, nil
}
