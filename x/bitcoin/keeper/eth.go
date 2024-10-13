package keeper

import (
	"context"

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
	reqLens := len(reqs.Withdraws) + len(reqs.ReplaceByFees) + len(reqs.Cancel1s)
	if reqLens == 0 {
		return nil
	}

	param, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	events := make(sdktypes.Events, 0, reqLens)

	var rejecting []uint64
	for _, v := range reqs.Withdraws {
		// reject if we have an invalid address
		script, err := types.DecodeBtcAddress(v.Address, types.BitcoinNetworks[param.NetworkName])
		if err != nil {
			k.Logger().Info("invalid withdrawal address", "id", v.Id, "address", v.Address, "err", err.Error())
			rejecting = append(rejecting, v.Id)
			continue
		}

		if err := k.Withdrawals.Set(ctx, v.Id, types.Withdrawal{
			Address:       v.Address,
			RequestAmount: v.Amount,
			MaxTxPrice:    v.TxPrice,
			OutputScript:  script,
			Status:        types.WITHDRAWAL_STATUS_PENDING,
		}); err != nil {
			return err
		}
		events = append(events, types.NewWithdrawalRequestEvent(v.Address, v.Id, v.TxPrice, v.Amount))
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
		if withdrawal.Status != types.WITHDRAWAL_STATUS_PENDING {
			k.Logger().Info("disregard rbf request due to it's processing", "id", v.Id)
			continue
		}
		withdrawal.MaxTxPrice = v.TxPrice
		if err := k.Withdrawals.Set(ctx, v.Id, withdrawal); err != nil {
			return err
		}
		events = append(events, types.NewWithdrawalReplaceEvent(v.Id, v.TxPrice))
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
		events = append(events, types.NewWithdrawalCancellationEvent(v.Id))
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(events)
	return nil
}
