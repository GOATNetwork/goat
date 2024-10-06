package keeper

import (
	"context"
	"math/big"
	"slices"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/relayer/types"
)

// EndBlocker elects new proposer and increase epoch number at regular intervals
func (k Keeper) EndBlocker(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return err
	}

	param, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	if duration := sdkctx.BlockTime().Sub(relayer.LastElected); duration < param.ElectingPeriod &&
		(relayer.ProposerAccepted || param.AcceptProposerTimeout == 0 || duration < param.AcceptProposerTimeout) {
		return nil
	}

	queue, err := k.Queue.Get(ctx)
	if err != nil {
		return err
	}

	onBoarding, offBoarding := len(queue.OnBoarding) > 0, len(queue.OffBoarding) > 0
	if onBoarding {
		for _, v := range queue.OnBoarding {
			voter, err := k.Voters.Get(ctx, v)
			if err != nil {
				return err
			}
			voter.Status = types.VOTER_STATUS_ACTIVATED
			if err := k.Voters.Set(ctx, v, voter); err != nil {
				return err
			}
		}
		relayer.Voters = append(relayer.Voters, queue.OnBoarding...)
	}

	var isProposerRemvoed bool
	if offBoarding {
		set := make(map[string]bool, len(queue.OffBoarding))
		for _, v := range queue.OffBoarding {
			if err := k.Voters.Remove(ctx, v); err != nil {
				return err
			}
			set[v] = true
		}

		isProposerRemvoed = set[relayer.Proposer]
		newVoters := slices.DeleteFunc(relayer.Voters, func(e string) bool {
			return set[e]
		})

		if isProposerRemvoed {
			if len(newVoters) == 0 { // it should never happen
				k.Logger().Error("delete too many voters in ElectProposer")
				return nil
			}
			// use the first voter as the new proposer
			relayer.Proposer = newVoters[0]
			relayer.Voters = newVoters[1:]
		} else {
			relayer.Voters = newVoters[:]
		}
	}

	// epoch number is constantly increasing even if we don't have a new election
	relayer.Epoch++
	relayer.LastElected = sdkctx.BlockTime()

	var events = sdktypes.Events{types.NewEpochEvent(relayer.Epoch)}
	if offBoarding || onBoarding {
		events = append(types.VoterChangedEvent(relayer.Epoch, queue.OnBoarding, queue.OffBoarding), events...)

		queue.OnBoarding = queue.OnBoarding[:0]
		queue.OffBoarding = queue.OffBoarding[:0]
		if err := k.Queue.Set(ctx, queue); err != nil {
			return err
		}

		// if the proposer is removed, we don't make a election, just use the next voter as the new proposer
		if isProposerRemvoed {
			relayer.ProposerAccepted = false
			if err := k.Relayer.Set(ctx, relayer); err != nil {
				return err
			}

			k.Logger().Info("New proposer", "height", sdkctx.BlockHeight(), "epoch", relayer.Epoch, "proposer", relayer.Proposer)
			sdkctx.EventManager().EmitEvents(
				append(events, types.ElectedProposerEvent(relayer.Proposer, relayer.Epoch)),
			)
			return nil
		}
	}

	voterLen := len(relayer.Voters)
	// no voter no election
	if voterLen == 0 {
		relayer.ProposerAccepted = true
		if err := k.Relayer.Set(ctx, relayer); err != nil {
			return err
		}
		sdkctx.EventManager().EmitEvents(events)
		return nil
	}

	// only get hash when we have 2 voters at least
	if voterLen > 1 {
		randao, err := k.Randao.Get(ctx)
		if err != nil {
			return err
		}
		// hash with the current epoch to ensure always have a new randao value
		rand := new(big.Int).SetBytes(goatcrypto.SHA256Sum(randao, goatcrypto.Uint64LE(relayer.Epoch)))
		proposerIndex := rand.Mod(rand, big.NewInt(int64(voterLen))).Int64()
		relayer.Proposer, relayer.Voters[proposerIndex] = relayer.Voters[proposerIndex], relayer.Proposer
	} else {
		relayer.Proposer, relayer.Voters[0] = relayer.Voters[0], relayer.Proposer
	}

	relayer.ProposerAccepted = false
	if err := k.Relayer.Set(ctx, relayer); err != nil {
		return err
	}

	k.Logger().Info("New proposer", "height", sdkctx.BlockHeight(), "epoch", relayer.Epoch, "proposer", relayer.Proposer)
	sdkctx.EventManager().EmitEvents(
		append(events, types.ElectedProposerEvent(relayer.Proposer, relayer.Epoch)),
	)
	return nil
}
