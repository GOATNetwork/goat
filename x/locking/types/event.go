package types

import (
	"math/big"
	"strconv"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeAddReward = "goat_add_reward"
)

func AddRewardEvent(blockNumber int64, amount *big.Int) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeAddReward,
		sdktypes.NewAttribute("block", strconv.FormatInt(blockNumber, 10)),
		sdktypes.NewAttribute("amount", amount.String()),
	)
}
