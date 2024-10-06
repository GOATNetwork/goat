package keeper

import (
	"context"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/relayer/types"
	"github.com/kelindar/bitmap"
)

func (k Keeper) VerifyProposal(ctx context.Context, req types.IVoteMsg, verifyFn ...func(sigdoc []byte) error) (uint64, error) {
	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return 0, err
	}

	if relayer.Proposer != req.GetProposer() {
		return 0, types.ErrNotProposer.Wrapf("not proposer")
	}

	voters := relayer.GetVoters()
	sequence, err := k.Sequence.Peek(ctx)
	if err != nil {
		return 0, err
	}

	if req.GetVote().GetSequence() != sequence {
		return 0, types.ErrInvalidProposalSignature.Wrap("incorrect seqeuence")
	}

	if req.GetVote().GetEpoch() != relayer.Epoch {
		return 0, types.ErrInvalidProposalSignature.Wrap("incorrect epoch")
	}

	bmp := bitmap.FromBytes(req.GetVote().GetVoters())

	bmpLen := bmp.Count()
	if bmpLen+1 < relayer.Threshold() || bmpLen > len(voters) {
		return 0, types.ErrInvalidProposalSignature.Wrapf("malformed signature length")
	}

	pubkeys := make([][]byte, 0, bmpLen+1)
	proposer, err := k.Voters.Get(ctx, relayer.Proposer)
	if err != nil {
		return 0, err
	}
	pubkeys = append(pubkeys, proposer.VoteKey)

	for i := 0; i < len(voters); i++ {
		if !bmp.Contains(uint32(i)) {
			continue
		}

		voter, err := k.Voters.Get(ctx, voters[i])
		if err != nil {
			return 0, err
		}
		pubkeys = append(pubkeys, voter.VoteKey)
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	sigdoc := types.VoteSignDoc(req.MethodName(), sdkctx.ChainID(), relayer.Proposer, sequence, relayer.Epoch, req.VoteSigDoc())
	if !goatcrypto.AggregateVerify(pubkeys, sigdoc, req.GetVote().GetSignature()) {
		return 0, types.ErrInvalidProposalSignature.Wrapf("invalid signature")
	}

	for _, fn := range verifyFn {
		if err := fn(sigdoc); err != nil {
			return 0, err
		}
	}

	// As long as the proposer sends a valid tx, it should be considered that the proposer is accepted.
	if !relayer.ProposerAccepted {
		relayer.ProposerAccepted = true
		if err := k.Relayer.Set(ctx, relayer); err != nil {
			return 0, err
		}
		k.Logger().Info("new proposer is accepted implicitly", "epoch", relayer.Epoch, "proposer", relayer.Proposer)
		sdkctx.EventManager().EmitEvent(types.AcceptedProposerEvent(relayer.Proposer, relayer.Epoch))
	}

	return sequence, nil
}

func (k Keeper) VerifyNonProposal(ctx context.Context, req types.INonVoteMsg) (types.IRelayer, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	relayer, err := k.Relayer.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	if relayer.Proposer != req.GetProposer() {
		return nil, types.ErrNotProposer.Wrapf("not proposer")
	}

	// As long as the proposer sends a valid tx, it should be considered that the proposer is accepted.
	if !relayer.ProposerAccepted {
		relayer.ProposerAccepted = true
		if err := k.Relayer.Set(sdkctx, relayer); err != nil {
			return nil, err
		}
		k.Logger().Info("new proposer is accepted implicitly", "epoch", relayer.Epoch, "proposer", relayer.Proposer)
		sdkctx.EventManager().EmitEvent(types.AcceptedProposerEvent(relayer.Proposer, relayer.Epoch))
	}

	return &relayer, nil
}
