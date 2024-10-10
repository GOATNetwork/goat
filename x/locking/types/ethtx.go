package types

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	goattypes "github.com/ethereum/go-ethereum/core/types/goattypes"
)

func (r *Reward) EthTx(seq uint64) *ethtypes.Transaction {
	return ethtypes.NewTx(
		ethtypes.NewGoatTx(
			goattypes.LockingModule,
			goattypes.LockingDistributeRewardAction,
			seq,
			&goattypes.DistributeRewardTx{
				Id:        r.Id,
				Recipient: common.BytesToAddress(r.Recipient),
				Goat:      r.Goat.BigInt(),
				GasReward: r.Gas.BigInt(),
			},
		),
	)
}

func (l *Unlock) EthTx(seq uint64) *ethtypes.Transaction {
	return ethtypes.NewTx(
		ethtypes.NewGoatTx(
			goattypes.LockingModule,
			goattypes.LockingCompleteUnlockAction,
			seq,
			&goattypes.CompleteUnlockTx{
				Id:        l.Id,
				Recipient: common.BytesToAddress(l.Recipient),
				Token:     common.BytesToAddress(l.Token),
				Amount:    l.Amount.BigInt(),
			},
		),
	)
}
