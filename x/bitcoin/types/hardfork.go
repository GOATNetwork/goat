package types

import "time"

type Hardfork struct {
	WithdrawalV2 int64
}

func (h *Hardfork) IsWithdrawalV2Enable(time time.Time) bool {
	return time.Unix() >= h.WithdrawalV2
}

var Hardforks = map[string]*Hardfork{
	"goat-testnet3": {WithdrawalV2: 1735099200}, // 2024-12-25T04:00:00Z
	"goat-mainnet":  {WithdrawalV2: 1735135200}, // 2024-12-25T14:00:00Z
}
