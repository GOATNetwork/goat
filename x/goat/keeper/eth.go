package keeper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
)

func (k Keeper) Dequeue(ctx context.Context) ([][]byte, error) {
	btcTxs, err := k.bitcoinKeeper.DequeueBitcoinModuleTx(ctx)
	if err != nil {
		return nil, err
	}
	k.Logger().Debug("dequeue bitcoin module txs", "len", len(btcTxs))

	lockingTxs, err := k.lockingKeeper.DequeueLockingModuleTx(ctx)
	if err != nil {
		return nil, err
	}
	k.Logger().Debug("dequeue locking module txs", "len", len(lockingTxs))

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
	if len(txRoot) != 33 {
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
	k.Logger().Debug("verifying dequeue bitcoin module txs", "len", len(btcTxs))

	if len(txs) < len(btcTxs) {
		return fmt.Errorf("bitcoin module txs length mismatched: len(txs)=%d len(mod)=%d", len(txs), len(btcTxs))
	}

	for idx, tx := range btcTxs {
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

	k.Logger().Debug("verifing dequeue locking module txs", "len", len(lockingTxs))
	txs = txs[len(btcTxs):]
	if len(txs) < len(lockingTxs) {
		return fmt.Errorf("locking module txs length mismatched: len(txs)=%d len(mod)=%d", len(txs), len(lockingTxs))
	}

	for idx, tx := range lockingTxs {
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
