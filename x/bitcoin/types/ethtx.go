package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
)

var satoshi = big.NewInt(1e10)

func (d *DepositReceipt) EthTx(seq uint64) *ethtypes.Transaction {
	amount := new(big.Int).SetUint64(d.Amount)
	amount.Mul(amount, satoshi)

	return ethtypes.NewTx(
		ethtypes.NewGoatTx(
			goattypes.BirdgeModule,
			goattypes.BridgeDepoitAction,
			seq,
			&goattypes.DepositTx{
				Txid:   common.BytesToHash(d.Txid),
				TxOut:  d.Txout,
				Target: common.BytesToAddress(d.Address),
				Amount: amount,
			},
		),
	)
}

func (d *WithdrawalReceipt) EthTx(seq uint64) *ethtypes.Transaction {
	amount := new(big.Int).SetUint64(d.Amount)
	amount.Mul(amount, satoshi)

	return ethtypes.NewTx(
		ethtypes.NewGoatTx(
			goattypes.BirdgeModule,
			goattypes.BridgePaidAction,
			seq,
			&goattypes.PaidTx{
				Id:     new(big.Int).SetUint64(d.Id),
				Txid:   common.BytesToHash(d.Txid),
				TxOut:  d.Txout,
				Amount: amount,
			},
		),
	)
}

func NewBitcoinHashEthTx(seq uint64, hash []byte) *ethtypes.Transaction {
	return ethtypes.NewTx(
		ethtypes.NewGoatTx(
			goattypes.BirdgeModule,
			goattypes.BitcoinNewHashAction,
			seq,
			&goattypes.AppendBitcoinHash{
				Hash: common.BytesToHash(hash),
			},
		),
	)
}
