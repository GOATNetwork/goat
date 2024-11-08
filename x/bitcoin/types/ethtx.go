package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
)

var satoshi = big.NewInt(1e10)

func (d *DepositExecReceipt) EthTx(seq uint64) *ethtypes.Transaction {
	amount := new(big.Int).SetUint64(d.Amount)
	amount.Mul(amount, satoshi)

	tax := new(big.Int).SetUint64(d.Tax)
	tax.Mul(tax, satoshi)

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
				Tax:    tax,
			},
		),
	)
}

func (d *WithdrawalExecReceipt) EthTx(seq uint64) *ethtypes.Transaction {
	amount := new(big.Int).SetUint64(d.Receipt.Amount)
	amount.Mul(amount, satoshi)

	return ethtypes.NewTx(
		ethtypes.NewGoatTx(
			goattypes.BirdgeModule,
			goattypes.BridgePaidAction,
			seq,
			&goattypes.PaidTx{
				Id:     new(big.Int).SetUint64(d.Id),
				Txid:   common.BytesToHash(d.Receipt.Txid),
				TxOut:  d.Receipt.Txout,
				Amount: amount,
			},
		),
	)
}

func NewRejectEthTx(wid, seq uint64) *ethtypes.Transaction {
	return ethtypes.NewTx(
		ethtypes.NewGoatTx(
			goattypes.BirdgeModule,
			goattypes.BridgeCancel2Action,
			seq,
			&goattypes.Cancel2Tx{
				Id: new(big.Int).SetUint64(wid),
			},
		),
	)
}

func NewBitcoinHashEthTx(seq uint64, hash []byte) *ethtypes.Transaction {
	return ethtypes.NewTx(
		ethtypes.NewGoatTx(
			goattypes.BirdgeModule,
			goattypes.BitcoinNewBlockAction,
			seq,
			&goattypes.NewBtcBlockTx{
				Hash: common.BytesToHash(hash),
			},
		),
	)
}
