package types

import (
	"strconv"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeNewEpoch      = "new_epoch"
	EventFinalizedProposal = "finalized_proposal"
	EventElectedProposer   = "elected_proposer"
	EventAceeptedProposer  = "accepted_proposer"
	EventVoterPending      = "voter_pending"
	EventVoterOnBoarding   = "voter_on_boarding"
	EventVoterBoarded      = "voter_boarded"
	EventVoterOffBoarding  = "voter_off_boarding"
	EventVoterDischarged   = "voter_discharged"
)

func FinalizedProposalEvent(sequence uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventFinalizedProposal,
		sdktypes.NewAttribute("sequence", strconv.FormatUint(sequence, 10)),
	)
}

func VoterBoardedEvent(proposer, voter string) sdktypes.Event {
	return sdktypes.NewEvent(
		EventVoterBoarded,
		sdktypes.NewAttribute("proposer", proposer),
		sdktypes.NewAttribute("voter", voter),
	)
}

func ElectedProposerEvent(proposer string, epoch uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventElectedProposer,
		sdktypes.NewAttribute("epoch", strconv.FormatUint(epoch, 10)),
		sdktypes.NewAttribute("proposer", proposer),
	)
}

func AcceptedProposerEvent(proposer string, epoch uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventAceeptedProposer,
		sdktypes.NewAttribute("epoch", strconv.FormatUint(epoch, 10)),
		sdktypes.NewAttribute("proposer", proposer),
	)
}
