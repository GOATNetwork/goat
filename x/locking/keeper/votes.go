package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/comet"
	"cosmossdk.io/math"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) HandleVoteInfos(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return err
	}

	for _, voteInfo := range sdkctx.VoteInfos() {
		signed := comet.BlockIDFlag(voteInfo.BlockIdFlag)
		address := sdktypes.ConsAddress(voteInfo.Validator.Address)
		if err := k.handleVoteInfo(sdkctx, address, signed, &param); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) handleVoteInfo(sdkctx sdktypes.Context, address sdktypes.ConsAddress, signed comet.BlockIDFlag, param *types.Params) error {
	validator, err := k.Validators.Get(sdkctx, address)
	if err != nil {
		return err
	}

	// It was slashed
	if validator.Status != types.Active {
		return nil
	}

	if signed == comet.BlockIDFlagAbsent {
		validator.SigningInfo.Missed++
	}
	isDown := validator.SigningInfo.Missed >= param.MaxMissedPerWindow

	validator.SigningInfo.Offset++
	if validator.SigningInfo.Offset >= param.SignedBlocksWindow {
		validator.SigningInfo.Missed = 0
		validator.SigningInfo.Offset = 0
	}

	if isDown {
		k.Logger().Info(
			"Offline confirmed", "validator", types.ValidatorName(address), "height", sdkctx.BlockHeight(),
		)

		// remove it from power ranking
		if err := k.PowerRanking.Remove(sdkctx,
			collections.Join(validator.Power, address)); err != nil {
			return err
		}

		var updated = sdktypes.Coins{}
		for _, locking := range validator.Locking {
			if err := k.Locking.Remove(sdkctx, collections.Join(locking.Denom, address)); err != nil {
				return err
			}

			amount := math.LegacyNewDecFromInt(locking.Amount).Mul(param.SlashFractionDowntime).TruncateInt()
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
		validator.Status = types.Downgrade
		validator.Power = 0
		validator.JailedUntil = sdkctx.BlockTime().Add(param.DowntimeJailDuration)
		sdkctx.EventManager().EmitEvent(types.ValidatorDowngradedEvent(address))
	}

	if err := k.Validators.Set(sdkctx, address, validator); err != nil {
		return err
	}

	return nil
}
