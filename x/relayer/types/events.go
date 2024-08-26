package types

import (
	"strconv"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeNewEpoch = "new_epoch"
	EventProposalDone = "proposal_done"
	EventNewProposer  = "new_proposer"
)

func ProposalDoneEvent(sequence uint64) sdktypes.Event {
	return sdktypes.NewEvent(
		EventProposalDone,
		sdktypes.NewAttribute("sequence", strconv.FormatUint(sequence, 10)),
	)
}

func NewProposer(proposer string) sdktypes.Event {
	return sdktypes.NewEvent(
		EventNewProposer,
		sdktypes.NewAttribute("proposer", proposer),
	)
}
