package keeper

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/params"
)

func (k Keeper) Dequeue(ctx context.Context) ([][]byte, error) {
	btcTxs, err := k.bitcoinKeeper.DequeueBitcoinModuleTx(ctx)
	if err != nil {
		return nil, err
	}

	lockingTxs, err := k.lockingKeeper.DequeueLockingModuleTx(ctx)
	if err != nil {
		return nil, err
	}

	res := make([][]byte, 0, len(btcTxs)+len(lockingTxs))
	for _, tx := range btcTxs {
		raw, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		res = append(res, raw)
	}

	for _, tx := range lockingTxs {
		raw, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		res = append(res, raw)
	}

	return res, nil
}

// VerifyDequeue verifies if the goat transactions of the new eth block are consistent with expected here
func (k Keeper) VerifyDequeue(ctx context.Context, txRoot []byte, txs [][]byte) error {
	// goat-geth will check the tx root hash
	if len(txRoot) != params.GoatHeaderExtraLengthV0 {
		return errors.New("invalid goat tx root")
	}

	goatTxLen := int(txRoot[0])
	if len(txs) < goatTxLen {
		return errors.New("tx length is less than expected")
	}

	btcTxs, err := k.bitcoinKeeper.DequeueBitcoinModuleTx(ctx)
	if err != nil {
		return err
	}

	if len(txs) < len(btcTxs) {
		return fmt.Errorf("tx mismatched: len(txs)=%d len(btc)=%d", len(txs), len(btcTxs))
	}

	for idx, tx := range btcTxs {
		if !tx.IsGoatTx() {
			return fmt.Errorf("not a goat tx %d", idx)
		}

		raw, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		if !bytes.Equal(raw, txs[idx]) {
			return fmt.Errorf("bridge tx %d bytes mismatched", idx)
		}
		goatTxLen--
	}

	lockingTxs, err := k.lockingKeeper.DequeueLockingModuleTx(ctx)
	if err != nil {
		return err
	}

	txs = txs[len(btcTxs):]
	if len(txs) < len(lockingTxs) {
		return fmt.Errorf("tx mismatched: len(txs)=%d len(btc)=%d len(locking)=%d", len(txs), len(btcTxs), len(lockingTxs))
	}

	for idx, tx := range lockingTxs {
		if !tx.IsGoatTx() {
			return fmt.Errorf("not a goat tx %d", idx+len(btcTxs))
		}

		raw, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		if !bytes.Equal(raw, txs[idx]) {
			return fmt.Errorf("locking tx %d bytes mismatched", idx)
		}
		goatTxLen--
	}

	if goatTxLen != 0 {
		return errors.New("goat txs length mismatched")
	}
	return nil
}
