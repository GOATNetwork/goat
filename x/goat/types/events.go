package types

import (
	"strconv"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	EventTypeNewEthBlock = "new_eth_block"
)

func NewEthBlockEvent(number uint64, hash []byte) sdktypes.Event {
	return sdktypes.NewEvent(
		EventTypeNewEthBlock,
		sdktypes.NewAttribute("number", strconv.FormatUint(number, 10)),
		sdktypes.NewAttribute("hash", hexutil.Encode(hash)),
	)
}
