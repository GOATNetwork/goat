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
	EventVoterActivated    = "voter_activated"
	EventVoterDischarged   = "voter_discharged"
)

func PendingVoterEvent(voter string) sdktypes.Event {
	return sdktypes.NewEvent(
		EventVoterPending,
		sdktypes.NewAttribute("voter", voter),
	)
}

func RemovingVoterEvent(voter string) sdktypes.Event {
	return sdktypes.NewEvent(
		EventVoterOffBoarding,
		sdktypes.NewAttribute("voter", voter),
	)
}

func NewEpochEvent(epoch uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeNewEpoch,
		sdktypes.NewAttribute("epoch", strconv.FormatUint(epoch, 10)),
	)
}

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

func VoterChangedEvent(epoch uint64, added, removed []string) sdktypes.Events {
	events := make(sdktypes.Events, 0, len(added)+len(removed)+2) // +2 for EventElectedProposer and EventTypeNewEpoch
	epochStr := strconv.FormatUint(epoch, 10)
	for _, v := range added {
		events = append(events, sdktypes.NewEvent(
			EventVoterActivated,
			sdktypes.NewAttribute("epoch", epochStr),
			sdktypes.NewAttribute("voter", v),
		))
	}

	for _, v := range removed {
		events = append(events, sdktypes.NewEvent(
			EventVoterDischarged,
			sdktypes.NewAttribute("epoch", epochStr),
			sdktypes.NewAttribute("voter", v),
		))
	}

	return events
}

func AcceptedProposerEvent(proposer string, epoch uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventAceeptedProposer,
		sdktypes.NewAttribute("epoch", strconv.FormatUint(epoch, 10)),
		sdktypes.NewAttribute("proposer", proposer),
	)
}
