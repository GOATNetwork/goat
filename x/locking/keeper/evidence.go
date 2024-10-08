package keeper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/comet"
	"cosmossdk.io/math"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) HandleEvidences(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	evidences := sdkctx.CometInfo().GetEvidence()
	if evidences.Len() == 0 {
		return nil
	}

	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return err
	}

	for i := 0; i < evidences.Len(); i++ {
		evidence := evidences.Get(i)
		switch evidence.Type() {
		case comet.LightClientAttack, comet.DuplicateVote:
			if err := k.handleEvidence(sdkctx, evidence, &param); err != nil {
				return err
			}
		default:
			k.Logger().Error(fmt.Sprintf("ignored unknown evidence type: %x", evidence.Type()))
		}
	}
	return nil
}

func (k Keeper) handleEvidence(ctx context.Context, evidence comet.Evidence, param *types.Params) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	address := sdktypes.ConsAddress(evidence.Validator().Address())

	cp := sdkctx.ConsensusParams()
	if cp.Evidence != nil {
		ageDuration := sdkctx.BlockTime().Sub(evidence.Time())
		ageBlocks := sdkctx.BlockHeight() - evidence.Height()
		if ageDuration > cp.Evidence.MaxAgeDuration && ageBlocks > cp.Evidence.MaxAgeNumBlocks {
			return nil
		}
	}

	validator, err := k.Validators.Get(sdkctx, address)
	if err != nil {
		return err
	}

	switch validator.Status {
	case types.Tombstoned:
		return nil
	}

	k.Logger().Info(
		"confirmed equivocation", "validator", hex.EncodeToString(address),
		"infraction_height", evidence.Height(), "infraction_time", evidence.Time(),
	)

	if err := k.PowerRanking.Remove(sdkctx,
		collections.Join(validator.Power, address)); err != nil {
		return err
	}

	var updated = sdktypes.Coins{}
	for _, locking := range validator.Locking {
		if err := k.Locking.Remove(sdkctx, collections.Join(locking.Denom, address)); err != nil {
			return err
		}

		amount := math.LegacyNewDecFromInt(locking.Amount).Mul(param.SlashFractionDoubleSign).TruncateInt()
		if amount.IsZero() {
			amount = math.NewIntFromBigIntMut(locking.Amount.BigInt())
		} else {
			updated = updated.Add(sdktypes.NewCoin(locking.Denom, locking.Amount.Sub(amount)))
		}

		slashed, err := k.Slashed.Get(sdkctx, locking.Denom)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				return err
			}
			slashed = math.ZeroInt()
		}

		if err := k.Slashed.Set(sdkctx, locking.Denom, slashed.Add(amount)); err != nil {
			return err
		}
	}

	validator.Locking = updated
	validator.Status = types.Tombstoned
	validator.Power = 0
	if err := k.Validators.Set(sdkctx, address, validator); err != nil {
		return err
	}

	sdkctx.EventManager().EmitEvent(types.ValidatorTombstonedEvent(address))
	return nil
}
