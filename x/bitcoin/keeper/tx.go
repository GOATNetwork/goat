package keeper

import (
	"bytes"
	"context"
	"encoding/json"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"github.com/btcsuite/btcd/wire"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) NewDeposits(ctx context.Context, req *types.MsgNewDeposits) (*types.MsgNewDepositsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	var headers map[uint64][]byte
	if err := json.Unmarshal(req.BlockHeaders, &headers); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "invalid block header json")
	}

	if h := len(headers); h == 0 || h > len(req.Deposits) {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "invalid headers size")
	}

	if _, err := k.relayerKeeper.VerifyNonProposal(ctx, req); err != nil {
		return nil, err
	}

	events := make(sdktypes.Events, 0, len(req.Deposits))
	deposits := make([]*types.DepositExecReceipt, 0, len(req.Deposits))
	for _, v := range req.Deposits {
		if err := v.Validate(); err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
		}

		deposit, err := k.VerifyDeposit(ctx, headers, v)
		if err != nil {
			return nil, err
		}
		if err := k.Deposited.Set(ctx,
			collections.Join(deposit.Txid, deposit.Txout), deposit.Amount); err != nil {
			return nil, err
		}
		events = append(events, types.NewDepositEvent(deposit))
		deposits = append(deposits, deposit)
	}

	queue, err := k.EthTxQueue.Get(ctx)
	if err != nil {
		return nil, err
	}

	queue.Deposits = append(queue.Deposits, deposits...)
	if err := k.EthTxQueue.Set(ctx, queue); err != nil {
		return nil, err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(events)
	return &types.MsgNewDepositsResponse{}, nil
}

func (k msgServer) NewBlockHashes(ctx context.Context, req *types.MsgNewBlockHashes) (*types.MsgNewBlockHashesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	parentHeight, err := k.BlockTip.Peek(ctx)
	if err != nil {
		return nil, err
	}
	if req.StartBlockNumber != parentHeight+1 {
		return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "block number is not the next of the block %d", parentHeight)
	}

	sequence, err := k.relayerKeeper.VerifyProposal(ctx, req)
	if err != nil {
		return nil, err
	}

	events := make(sdktypes.Events, 0, len(req.BlockHash)+1)
	for _, v := range req.BlockHash {
		parentHeight++
		if err := k.BlockHashes.Set(ctx, parentHeight, v); err != nil {
			return nil, err
		}
		events = append(events, types.NewBlockHashEvent(parentHeight, v))
	}

	if err := k.BlockTip.Set(ctx, parentHeight); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.SetProposalSeq(ctx, sequence+1); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.UpdateRandao(ctx, req); err != nil {
		return nil, err
	}

	events = append(events, relayertypes.FinalizedProposalEvent(sequence))
	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(events)
	return &types.MsgNewBlockHashesResponse{}, nil
}

func (k msgServer) NewPubkey(ctx context.Context, req *types.MsgNewPubkey) (*types.MsgNewPubkeyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	sequence, err := k.relayerKeeper.VerifyProposal(ctx, req)
	if err != nil {
		return nil, err
	}

	rawKey := relayertypes.EncodePublicKey(req.Pubkey)
	exists, err := k.relayerKeeper.HasPubkey(ctx, rawKey)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "the key already existed")
	}

	if err := k.relayerKeeper.AddNewKey(ctx, rawKey); err != nil {
		return nil, err
	}
	if err := k.relayerKeeper.SetProposalSeq(ctx, sequence+1); err != nil {
		return nil, err
	}
	if err := k.Pubkey.Set(ctx, *req.Pubkey); err != nil {
		return nil, err
	}
	if err := k.relayerKeeper.UpdateRandao(ctx, req); err != nil {
		return nil, err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(
		sdktypes.Events{types.NewKeyEvent(req.Pubkey), relayertypes.FinalizedProposalEvent(sequence)},
	)

	return &types.MsgNewPubkeyResponse{}, nil
}

func (k msgServer) ProcessWithdrawal(ctx context.Context, req *types.MsgProcessWithdrawal) (*types.MsgProcessWithdrawalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	tx, txrd := new(wire.MsgTx), bytes.NewReader(req.NoWitnessTx)
	if err := tx.DeserializeNoWitness(txrd); err != nil || txrd.Len() > 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid non-witness tx")
	}

	txoutLen, withdrawalLen := len(tx.TxOut), len(req.Id)
	if txoutLen != withdrawalLen && txoutLen != withdrawalLen+1 { // change output up to 1
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid tx output size for withdrawals")
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	sequence, err := k.relayerKeeper.VerifyProposal(sdkctx, req)
	if err != nil {
		return nil, err
	}

	txid := goatcrypto.DoubleSHA256Sum(req.NoWitnessTx)

	/*
		Note:

		we don't check if relayer owns tx inputs since we don't manage utxo database on the chain

		we should allow relayer to spend the changes of the withdrawals which not reach the confirmation number

		to reduce complexity, we only have validations for the outputs.
	*/

	// Sat Per Byte
	txPrice := float64(req.TxFee) / float64(len(req.NoWitnessTx))

	// get the network config
	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return nil, err
	}
	netwk := types.BitcoinNetworks[param.NetworkName]
	if netwk == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrAppConfig, "%s network is not defined", param.NetworkName)
	}

	txOutput := types.TxOuptut{Values: make([]uint64, withdrawalLen)}
	for idx, wid := range req.Id {
		withdrawal, err := k.Withdrawals.Get(sdkctx, wid)
		if err != nil {
			return nil, err
		}

		if withdrawal.Status != types.WITHDRAWAL_STATUS_PENDING && withdrawal.Status != types.WITHDRAWAL_STATUS_CANCELING {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "witdhrawal %d is not pending or canceling", wid)
		}

		if txPrice > float64(withdrawal.MaxTxPrice) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "tx price is larger than user request for witdhrawal %d", wid)
		}

		txout := tx.TxOut[idx]
		outputScript, err := types.DecodeBtcAddress(withdrawal.Address, netwk)
		if err != nil { // It should not happen
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address to process %d", wid)
		}

		if !bytes.Equal(outputScript, txout.PkScript) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "witdhrawal %d script not matched", wid)
		}

		if withdrawal.RequestAmount < uint64(txout.Value) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "witdhrawal %d amount too large", wid)
		}

		// the withdrawal id can't be duplicated since we update the status here
		withdrawal.Status = types.WITHDRAWAL_STATUS_PROCESSING
		withdrawal.Receipt = &types.WithdrawalReceipt{Txid: txid, Txout: uint32(idx), Amount: uint64(txout.Value)}
		if err := k.Withdrawals.Set(sdkctx, wid, withdrawal); err != nil {
			return nil, err
		}
		txOutput.Values[idx] = uint64(txout.Value)
	}

	// check if the change output should be for current pubkey
	if txoutLen != withdrawalLen {
		change := tx.TxOut[withdrawalLen]
		pubkey, err := k.Pubkey.Get(ctx)
		if err != nil {
			return nil, err
		}
		if !types.VerifySystemAddressScript(&pubkey, change.PkScript) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "give change to not a latest relayer pubkey")
		}
	}

	// Add processing staus
	pid, err := k.ProcessID.Peek(sdkctx)
	if err != nil {
		return nil, err
	}
	if err := k.Processing.Set(sdkctx, pid, types.Processing{
		Txid: [][]byte{txid}, Output: []types.TxOuptut{txOutput},
		Withdrawals: req.Id, Fee: req.TxFee,
	}); err != nil {
		return nil, err
	}
	if err := k.ProcessID.Set(sdkctx, pid+1); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.SetProposalSeq(sdkctx, sequence+1); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.UpdateRandao(sdkctx, req); err != nil {
		return nil, err
	}

	sdkctx.EventManager().EmitEvents(sdktypes.Events{
		types.NewWithdrawalProcessingEvent(pid, txid),
		relayertypes.FinalizedProposalEvent(sequence),
	})

	return &types.MsgProcessWithdrawalResponse{}, nil
}

func (k msgServer) ReplaceWithdrawal(ctx context.Context, req *types.MsgReplaceWithdrawal) (*types.MsgReplaceWithdrawalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	tx, txrd := new(wire.MsgTx), bytes.NewReader(req.NewNoWitnessTx)
	if err := tx.DeserializeNoWitness(txrd); err != nil || txrd.Len() > 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid non-witness tx")
	}
	txid := goatcrypto.DoubleSHA256Sum(req.NewNoWitnessTx)

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	processing, err := k.Processing.Get(sdkctx, req.Pid)
	if err != nil {
		return nil, err
	}

	if processing.Fee >= req.NewTxFee {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "new tx fee is less than before")
	}
	processing.Fee = req.NewTxFee

	for _, item := range processing.Txid {
		if bytes.Equal(item, txid) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "the tx doesn't have any change")
		}
	}

	txoutLen, withdrawalLen := len(tx.TxOut), len(processing.Withdrawals)
	if txoutLen != withdrawalLen && txoutLen != withdrawalLen+1 { // change output up to 1
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid tx output size for withdrawals")
	}

	txPrice := float64(req.NewTxFee) / float64(len(req.NewNoWitnessTx))

	// verify proposal vote
	sequence, err := k.relayerKeeper.VerifyProposal(sdkctx, req)
	if err != nil {
		return nil, err
	}

	// get the network config
	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return nil, err
	}
	netwk := types.BitcoinNetworks[param.NetworkName]
	if netwk == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrAppConfig, "%s network is not defined", param.NetworkName)
	}

	txOutput := types.TxOuptut{Values: make([]uint64, withdrawalLen)}
	for idx, wid := range processing.Withdrawals {
		withdrawal, err := k.Withdrawals.Get(sdkctx, wid)
		if err != nil {
			return nil, err
		}

		if withdrawal.Status != types.WITHDRAWAL_STATUS_PROCESSING || withdrawal.Receipt == nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "witdhrawal %d is not processing", wid)
		}

		if txPrice > float64(withdrawal.MaxTxPrice) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "tx price is larger than user request for witdhrawal %d", wid)
		}

		txout := tx.TxOut[idx]
		outputScript, err := types.DecodeBtcAddress(withdrawal.Address, netwk)
		if err != nil { // It should not happen
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address to process %d", wid)
		}

		if !bytes.Equal(outputScript, txout.PkScript) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "witdhrawal %d script not matched", wid)
		}

		if withdrawal.RequestAmount < uint64(txout.Value) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "witdhrawal %d amount too large", wid)
		}

		withdrawal.Receipt.Txid = txid
		withdrawal.Receipt.Amount = uint64(txout.Value)
		if err := k.Withdrawals.Set(sdkctx, wid, withdrawal); err != nil {
			return nil, err
		}
		txOutput.Values[idx] = uint64(txout.Value)
	}

	// check if the change output should be for current pubkey
	if txoutLen != withdrawalLen {
		change := tx.TxOut[withdrawalLen]
		pubkey, err := k.Pubkey.Get(ctx)
		if err != nil {
			return nil, err
		}
		if !types.VerifySystemAddressScript(&pubkey, change.PkScript) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "give change to not a latest relayer pubkey")
		}
	}

	processing.Txid = append(processing.Txid, txid)
	processing.Output = append(processing.Output, txOutput)
	if err := k.Processing.Set(sdkctx, req.Pid, processing); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.SetProposalSeq(sdkctx, sequence+1); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.UpdateRandao(sdkctx, req); err != nil {
		return nil, err
	}

	sdkctx.EventManager().EmitEvents(sdktypes.Events{
		types.NewWithdrawalRelayerReplaceEvent(req.Pid, txid),
		relayertypes.FinalizedProposalEvent(sequence),
	})

	return &types.MsgReplaceWithdrawalResponse{}, nil
}

func (k msgServer) FinalizeWithdrawal(ctx context.Context, req *types.MsgFinalizeWithdrawal) (*types.MsgFinalizeWithdrawalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	if _, err := k.relayerKeeper.VerifyNonProposal(sdkctx, req); err != nil {
		return nil, err
	}

	processing, err := k.Processing.Get(sdkctx, req.Pid)
	if err != nil {
		return nil, err
	}

	if len(processing.Txid) != len(processing.Output) {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "internal error: txid length is not the same with outputs")
	}

	idx := -1
	for i, txid := range processing.Txid {
		if bytes.Equal(txid, req.Txid) {
			idx = i
			break
		}
	}

	if idx == -1 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "txid not found")
	}

	txOutput := processing.Output[idx]
	if len(txOutput.Values) != len(processing.Withdrawals) {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "internal error: output length is not the same with withdrawals")
	}

	// check if the block is voted
	blockHash, err := k.BlockHashes.Get(sdkctx, req.BlockNumber)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(blockHash, goatcrypto.DoubleSHA256Sum(req.BlockHeader)) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "inconsistent block hash")
	}

	// check if the spv is valid
	if !types.VerifyMerkelProof(req.Txid, req.BlockHeader[36:68], req.IntermediateProof, req.TxIndex) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid spv")
	}

	queue, err := k.EthTxQueue.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	for idx, wid := range processing.Withdrawals {
		withdrawal, err := k.Withdrawals.Get(sdkctx, wid)
		if err != nil {
			return nil, err
		}
		if withdrawal.Status != types.WITHDRAWAL_STATUS_PROCESSING {
			return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "witdhrawal %d is not processing", wid)
		}

		if withdrawal.Receipt == nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "witdhrawal %d receipt is nil", wid)
		}

		// the last RBFed txes might not be confirmed, so we should update the txid and amount
		withdrawal.Receipt.Txid = req.Txid
		withdrawal.Receipt.Amount = txOutput.Values[idx]
		withdrawal.Status = types.WITHDRAWAL_STATUS_PAID
		if err := k.Withdrawals.Set(sdkctx, wid, withdrawal); err != nil {
			return nil, err
		}
		queue.PaidWithdrawals = append(queue.PaidWithdrawals, &types.WithdrawalExecReceipt{
			Id:      wid,
			Receipt: withdrawal.Receipt,
		})
	}

	if err := k.EthTxQueue.Set(sdkctx, queue); err != nil {
		return nil, err
	}

	// we don't use it anymore
	if err := k.Processing.Remove(sdkctx, req.Pid); err != nil {
		return nil, err
	}

	sdkctx.EventManager().EmitEvent(types.NewWithdrawalFinalizedEvent(req.Pid))
	return &types.MsgFinalizeWithdrawalResponse{}, nil
}

func (k msgServer) ApproveCancellation(ctx context.Context, req *types.MsgApproveCancellation) (*types.MsgApproveCancellationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if _, err := k.relayerKeeper.VerifyNonProposal(ctx, req); err != nil {
		return nil, err
	}

	queue, err := k.EthTxQueue.Get(ctx)
	if err != nil {
		return nil, err
	}

	events := make(sdktypes.Events, 0, len(req.Id))
	for _, wid := range req.Id {
		withdrawal, err := k.Withdrawals.Get(ctx, wid)
		if err != nil {
			return nil, err
		}
		if withdrawal.Status != types.WITHDRAWAL_STATUS_CANCELING {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "witdhrawal %d is not canceling", wid)
		}
		withdrawal.Status = types.WITHDRAWAL_STATUS_CANCELED
		if err := k.Withdrawals.Set(ctx, wid, withdrawal); err != nil {
			return nil, err
		}
		events = append(events, types.NewWithdrawalRelayerCancelEvent(wid))
	}

	queue.RejectedWithdrawals = append(queue.RejectedWithdrawals, req.Id...)
	if err := k.EthTxQueue.Set(ctx, queue); err != nil {
		return nil, err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(events)
	return &types.MsgApproveCancellationResponse{}, nil
}

func (k msgServer) NewConsolidation(ctx context.Context, req *types.MsgNewConsolidation) (*types.MsgNewConsolidationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	tx, txrd := new(wire.MsgTx), bytes.NewReader(req.NoWitnessTx)
	if err := tx.DeserializeNoWitness(txrd); err != nil || txrd.Len() > 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid non-witness tx")
	}

	if len(tx.TxOut) != 1 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "consolidation should have only 1 output")
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	pubkey, err := k.Pubkey.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	if !types.VerifySystemAddressScript(&pubkey, tx.TxOut[0].PkScript) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "not pay to the latest relayer pubkey")
	}

	sequence, err := k.relayerKeeper.VerifyProposal(sdkctx, req)
	if err != nil {
		return nil, err
	}

	txid := goatcrypto.DoubleSHA256Sum(req.NoWitnessTx)
	if err := k.relayerKeeper.SetProposalSeq(sdkctx, sequence+1); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.UpdateRandao(sdkctx, req); err != nil {
		return nil, err
	}

	sdkctx.EventManager().EmitEvents(sdktypes.Events{
		types.NewConsolidationEvent(txid),
		relayertypes.FinalizedProposalEvent(sequence),
	})
	return &types.MsgNewConsolidationResponse{}, nil
}
