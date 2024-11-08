package keeper

import (
	"context"
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/bitcoin/types"
)

func (k Keeper) DequeueBitcoinModuleTx(ctx context.Context) (txs []*ethtypes.Transaction, err error) {
	const (
		MaxDeposit    = 8
		MaxWithdrawal = 8
	)

	queue, err := k.EthTxQueue.Get(ctx)
	if err != nil {
		return nil, err
	}

	txNonce, err := k.EthTxNonce.Peek(ctx)
	if err != nil {
		return nil, err
	}

	// pop block hash up to 1
	{
		tip, err := k.BlockTip.Peek(ctx)
		if err != nil {
			return nil, err
		}

		if queue.BlockNumber < tip {
			queue.BlockNumber++

			blockHash, err := k.BlockHashes.Get(ctx, queue.BlockNumber)
			if err != nil {
				return nil, err
			}
			txs = append(txs, types.NewBitcoinHashEthTx(txNonce, blockHash))
			txNonce++
		}
	}

	// pop deposit up to 8
	{
		var n int
		for ; n < len(queue.Deposits) && n < MaxDeposit; n++ {
			deposit := queue.Deposits[n]
			txs = append(txs, deposit.EthTx(txNonce))

			txNonce++
		}
		queue.Deposits = queue.Deposits[n:]
	}

	// pop paid and reject withdrwal up to 8
	{
		var n int
		for ; n < len(queue.PaidWithdrawals) && n < MaxWithdrawal; n++ {
			paid := queue.PaidWithdrawals[n]
			txs = append(txs, paid.EthTx(txNonce))

			txNonce++
		}
		queue.PaidWithdrawals = queue.PaidWithdrawals[n:]

		var i int
		for ; i < len(queue.RejectedWithdrawals) && n < MaxWithdrawal; i++ {
			txs = append(txs, types.NewRejectEthTx(queue.RejectedWithdrawals[i], txNonce))

			n++
			txNonce++
		}
		queue.RejectedWithdrawals = queue.RejectedWithdrawals[i:]
	}

	if len(txs) > 0 {
		if err := k.EthTxQueue.Set(ctx, queue); err != nil {
			return nil, err
		}
		if err := k.EthTxNonce.Set(ctx, txNonce); err != nil {
			return nil, err
		}
	}
	return txs, nil
}

func (k Keeper) ProcessBridgeRequest(ctx context.Context, reqs goattypes.BridgeRequests) error {
	reqLens := len(reqs.Withdraws) + len(reqs.ReplaceByFees) + len(reqs.Cancel1s) +
		len(reqs.DepositTax) + len(reqs.Confirmation) + len(reqs.MinDeposit)
	if reqLens == 0 {
		return nil
	}

	events := make(sdktypes.Events, 0, reqLens)

	param, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	netwk := types.BitcoinNetworks[param.NetworkName]
	if netwk == nil {
		return fmt.Errorf("%s network is not defined", param.NetworkName)
	}

	var rejecting []uint64
	for _, v := range reqs.Withdraws {
		status := types.WITHDRAWAL_STATUS_PENDING
		// reject if we have an invalid address
		if _, err := types.DecodeBtcAddress(v.Address, netwk); err != nil {
			k.Logger().Info("invalid withdrawal address", "id", v.Id, "address", v.Address, "err", err.Error())
			rejecting = append(rejecting, v.Id)
			status = types.WITHDRAWAL_STATUS_CANCELED
		}

		if err := k.Withdrawals.Set(ctx, v.Id, types.Withdrawal{
			Address:       v.Address,
			RequestAmount: v.Amount,
			MaxTxPrice:    v.TxPrice,
			Status:        status,
		}); err != nil {
			return err
		}
		if status == types.WITHDRAWAL_STATUS_PENDING {
			events = append(events, types.NewWithdrawalInitEvent(v.Address, v.Id, v.TxPrice, v.Amount))
		} else {
			events = append(events, types.NewWithdrawalRelayerCancelEvent(v.Id))
		}
	}

	if len(rejecting) > 0 {
		queue, err := k.EthTxQueue.Get(ctx)
		if err != nil {
			return err
		}
		queue.RejectedWithdrawals = append(queue.RejectedWithdrawals, rejecting...)
		if err := k.EthTxQueue.Set(ctx, queue); err != nil {
			return err
		}
	}

	for _, v := range reqs.ReplaceByFees {
		withdrawal, err := k.Withdrawals.Get(ctx, v.Id)
		if err != nil {
			return err
		}
		// event if it's processing, we should allow the change
		// relayer can use the latest tx price to do the rbf
		if withdrawal.Status != types.WITHDRAWAL_STATUS_PENDING && withdrawal.Status != types.WITHDRAWAL_STATUS_PROCESSING {
			k.Logger().Info("disregard rbf request", "id", v.Id)
			continue
		}
		withdrawal.MaxTxPrice = v.TxPrice
		if err := k.Withdrawals.Set(ctx, v.Id, withdrawal); err != nil {
			return err
		}
		events = append(events, types.NewWithdrawalUserReplaceEvent(v.Id, v.TxPrice))
	}

	for _, v := range reqs.Cancel1s {
		withdrawal, err := k.Withdrawals.Get(ctx, v.Id)
		if err != nil {
			return err
		}

		if withdrawal.Status != types.WITHDRAWAL_STATUS_PENDING {
			k.Logger().Info("disregard cancellation request due to it's processing", "id", v.Id)
			continue
		}

		withdrawal.Status = types.WITHDRAWAL_STATUS_CANCELING
		if err := k.Withdrawals.Set(ctx, v.Id, withdrawal); err != nil {
			return err
		}
		events = append(events, types.NewWithdrawalUserCancelEvent(v.Id))
	}

	for _, v := range reqs.DepositTax {
		param.DepositTaxRate = v.Rate
		param.MaxDepositTax = v.Max
	}

	for _, v := range reqs.Confirmation {
		param.ConfirmationNumber = v.Number
		events = append(events, types.UpdateConfirmationNumberEvent(v.Number))
	}

	for _, v := range reqs.MinDeposit {
		param.MinDepositAmount = v.Satoshi
		events = append(events, types.UpdateMinDepositEvent(v.Satoshi))
	}

	if err := k.Params.Set(ctx, param); err != nil {
		return err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(events)
	return nil
}
