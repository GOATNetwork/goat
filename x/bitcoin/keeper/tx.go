package keeper

import (
	"bytes"
	"context"
	"encoding/json"

	"cosmossdk.io/collections"
	"github.com/btcsuite/btcd/wire"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
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
		return nil, types.ErrInvalidRequest.Wrap(err.Error())
	}

	var headers map[uint64][]byte
	if err := json.Unmarshal(req.BlockHeaders, &headers); err != nil {
		return nil, types.ErrInvalidRequest.Wrap("invalid block header json")
	}

	if h := len(headers); h == 0 || h > len(req.Deposits) {
		return nil, types.ErrInvalidRequest.Wrap("invalid headers size")
	}

	if _, err := k.relayerKeeper.VerifyNonProposal(ctx, req); err != nil {
		return nil, err
	}

	events := make(sdktypes.Events, 0, len(req.Deposits))
	deposits := make([]*types.DepositExecReceipt, 0, len(req.Deposits))
	for _, v := range req.Deposits {
		if err := v.Validate(); err != nil {
			return nil, types.ErrInvalidRequest.Wrap(err.Error())
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

	queue, err := k.ExecuableQueue.Get(ctx)
	if err != nil {
		return nil, err
	}

	queue.Deposits = append(queue.Deposits, deposits...)
	if err := k.ExecuableQueue.Set(ctx, queue); err != nil {
		return nil, err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(events)
	return &types.MsgNewDepositsResponse{}, nil
}

func (k msgServer) NewBlockHashes(ctx context.Context, req *types.MsgNewBlockHashes) (*types.MsgNewBlockHashesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, types.ErrInvalidRequest.Wrap(err.Error())
	}

	parentHeight, err := k.BlockTip.Peek(ctx)
	if err != nil {
		return nil, err
	}
	if req.StartBlockNumber != parentHeight+1 {
		return nil, types.ErrInvalidRequest.Wrapf("block number is not the next of the block %d", parentHeight)
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

	sequence++
	if err := k.relayerKeeper.SetProposalSeq(ctx, sequence); err != nil {
		return nil, err
	}
	events = append(events, relayertypes.FinalizedProposalEvent(sequence))

	if err := k.relayerKeeper.UpdateRandao(ctx, req); err != nil {
		return nil, err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(events)
	return &types.MsgNewBlockHashesResponse{}, nil
}

func (k msgServer) NewPubkey(ctx context.Context, req *types.MsgNewPubkey) (*types.MsgNewPubkeyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, types.ErrInvalidRequest.Wrap(err.Error())
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
		return nil, relayertypes.ErrInvalidPubkey.Wrap("the key already existed")
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

func (k msgServer) InitializeWithdrawal(ctx context.Context, req *types.MsgInitializeWithdrawal) (*types.MsgInitializeWithdrawalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, types.ErrInvalidRequest.Wrap(err.Error())
	}

	tx, txrd := new(wire.MsgTx), bytes.NewReader(req.Proposal.NoWitnessTx)
	if err := tx.DeserializeNoWitness(txrd); err != nil || txrd.Len() > 0 {
		return nil, types.ErrInvalidRequest.Wrap("invalid non-witness tx")
	}

	txoutLen, withdrawalLen := len(tx.TxOut), len(req.Proposal.Id)
	if txoutLen != withdrawalLen && txoutLen != withdrawalLen+1 { // change output up to 1
		return nil, types.ErrInvalidRequest.Wrap("invalid tx output size for withdrawals")
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	sequence, err := k.relayerKeeper.VerifyProposal(sdkctx, req)
	if err != nil {
		return nil, err
	}

	txid := goatcrypto.DoubleSHA256Sum(req.Proposal.NoWitnessTx)

	/*
		Note:

		we don't check if relayer owns tx inputs since we don't manage utxo database on the chain

		we should allow relayer to spend the changes of the withdrawals which not reach the confirmation number

		to reduce complexity, we only have validations for the outputs.
	*/

	// Sat Per Byte
	txPrice := float64(req.Proposal.TxFee) / float64(len(req.Proposal.NoWitnessTx))

	for idx, wid := range req.Proposal.Id {
		withdrawal, err := k.Withdrawals.Get(sdkctx, wid)
		if err != nil {
			return nil, err
		}

		if withdrawal.Status != types.WITHDRAWAL_STATUS_PENDING && withdrawal.Status != types.WITHDRAWAL_STATUS_CANCELING {
			return nil, types.ErrInvalidRequest.Wrapf("witdhrawal %d is not pending or canceling", wid)
		}

		if txPrice > float64(withdrawal.MaxTxPrice) {
			return nil, types.ErrInvalidRequest.Wrapf("tx price is larger than user request for witdhrawal %d", wid)
		}

		txout := tx.TxOut[idx]
		if !bytes.Equal(withdrawal.OutputScript, txout.PkScript) {
			return nil, types.ErrInvalidRequest.Wrapf("witdhrawal %d script not matched", wid)
		}

		if withdrawal.RequestAmount < uint64(txout.Value) {
			return nil, types.ErrInvalidRequest.Wrapf("witdhrawal %d amount too large", wid)
		}

		// the withdrawal id can't be duplicated since we update the status here
		withdrawal.Status = types.WITHDRAWAL_STATUS_PROCESSING
		withdrawal.Receipt = &types.WithdrawalReceipt{Txid: txid, Txout: uint32(idx), Amount: uint64(txout.Value)}
		if err := k.Withdrawals.Set(sdkctx, wid, withdrawal); err != nil {
			return nil, err
		}
	}

	// check if the change output should be for current pubkey
	if txoutLen != withdrawalLen {
		change := tx.TxOut[withdrawalLen]
		pubkey, err := k.Pubkey.Get(ctx)
		if err != nil {
			return nil, err
		}
		if !types.VerifySystemAddressScript(&pubkey, change.PkScript) {
			return nil, types.ErrInvalidRequest.Wrap("give change to not a latest relayer pubkey")
		}
	}

	if err := k.Processing.Set(sdkctx, txid, types.WithdrawalIds{Id: req.Proposal.Id}); err != nil {
		return nil, err
	}

	sequence++
	if err := k.relayerKeeper.SetProposalSeq(sdkctx, sequence); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.UpdateRandao(sdkctx, req); err != nil {
		return nil, err
	}

	sdkctx.EventManager().EmitEvents(
		sdktypes.Events{types.InitializeWithdrawalEvent(txid), relayertypes.FinalizedProposalEvent(sequence)},
	)

	return &types.MsgInitializeWithdrawalResponse{}, nil
}

func (k msgServer) FinalizeWithdrawal(ctx context.Context, req *types.MsgFinalizeWithdrawal) (*types.MsgFinalizeWithdrawalResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	if _, err := k.relayerKeeper.VerifyNonProposal(sdkctx, req); err != nil {
		return nil, err
	}

	// check if the block is voted
	blockHash, err := k.BlockHashes.Get(sdkctx, req.BlockNumber)
	if err != nil {
		return nil, err
	}

	proposal, err := k.Processing.Get(sdkctx, req.Txid)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(blockHash, goatcrypto.DoubleSHA256Sum(req.BlockHeader)) {
		return nil, types.ErrInvalidRequest.Wrap("inconsistent block hash")
	}

	// check if the spv is valid
	if !types.VerifyMerkelProof(req.Txid, req.BlockHeader[36:68], req.IntermediateProof, req.TxIndex) {
		return nil, types.ErrInvalidRequest.Wrap("invalid spv")
	}

	queue, err := k.ExecuableQueue.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	for _, wid := range proposal.Id {
		withdrawal, err := k.Withdrawals.Get(sdkctx, wid)
		if err != nil {
			return nil, err
		}
		if withdrawal.Status != types.WITHDRAWAL_STATUS_PROCESSING {
			return nil, types.ErrInvalidRequest.Wrapf("witdhrawal %d is not processing", wid)
		}

		if withdrawal.Receipt == nil {
			return nil, types.ErrInvalidRequest.Wrapf("witdhrawal %d receipt is nil", wid)
		}

		withdrawal.Status = types.WITHDRAWAL_STATUS_PAID
		if err := k.Withdrawals.Set(sdkctx, wid, withdrawal); err != nil {
			return nil, err
		}
		queue.PaidWithdrawals = append(queue.PaidWithdrawals, &types.WithdrawalExecReceipt{
			Id:      wid,
			Receipt: withdrawal.Receipt,
		})
	}

	if err := k.ExecuableQueue.Set(sdkctx, queue); err != nil {
		return nil, err
	}

	// we don't use it anymore
	if err := k.Processing.Remove(sdkctx, req.Txid); err != nil {
		return nil, err
	}

	sdkctx.EventManager().EmitEvent(types.FinalizeWithdrawalEvent(req.Txid))
	return &types.MsgFinalizeWithdrawalResponse{}, nil
}

func (k msgServer) ApproveCancellation(ctx context.Context, req *types.MsgApproveCancellation) (*types.MsgApproveCancellationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, types.ErrInvalidRequest.Wrap(err.Error())
	}

	if _, err := k.relayerKeeper.VerifyNonProposal(ctx, req); err != nil {
		return nil, err
	}

	queue, err := k.ExecuableQueue.Get(ctx)
	if err != nil {
		return nil, err
	}

	var events sdktypes.Events
	for _, wid := range req.Id {
		withdrawal, err := k.Withdrawals.Get(ctx, wid)
		if err != nil {
			return nil, err
		}
		if withdrawal.Status != types.WITHDRAWAL_STATUS_CANCELING {
			return nil, types.ErrInvalidRequest.Wrapf("witdhrawal %d is not canceling", wid)
		}
		withdrawal.Status = types.WITHDRAWAL_STATUS_CANCELED
		if err := k.Withdrawals.Set(ctx, wid, withdrawal); err != nil {
			return nil, err
		}
		events = append(events, types.ApproveCancellationEvent(wid))
	}

	queue.RejectedWithdrawals = append(queue.RejectedWithdrawals, req.Id...)
	if err := k.ExecuableQueue.Set(ctx, queue); err != nil {
		return nil, err
	}

	sdktypes.UnwrapSDKContext(ctx).EventManager().EmitEvents(events)
	return &types.MsgApproveCancellationResponse{}, nil
}

func (k msgServer) NewConsolidation(ctx context.Context, req *types.MsgNewConsolidation) (*types.MsgNewConsolidationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, types.ErrInvalidRequest.Wrap(err.Error())
	}

	tx, txrd := new(wire.MsgTx), bytes.NewReader(req.NoWitnessTx)
	if err := tx.DeserializeNoWitness(txrd); err != nil || txrd.Len() > 0 {
		return nil, types.ErrInvalidRequest.Wrap("invalid non-witness tx")
	}

	if len(tx.TxOut) != 1 {
		return nil, types.ErrInvalidRequest.Wrap("consolidation should have only 1 output")
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	pubkey, err := k.Pubkey.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	if !types.VerifySystemAddressScript(&pubkey, tx.TxOut[0].PkScript) {
		return nil, types.ErrInvalidRequest.Wrap("not pay to the latest relayer pubkey")
	}

	sequence, err := k.relayerKeeper.VerifyProposal(sdkctx, req)
	if err != nil {
		return nil, err
	}

	txid := goatcrypto.DoubleSHA256Sum(req.NoWitnessTx)

	sequence++
	if err := k.relayerKeeper.SetProposalSeq(sdkctx, sequence); err != nil {
		return nil, err
	}

	if err := k.relayerKeeper.UpdateRandao(sdkctx, req); err != nil {
		return nil, err
	}

	sdkctx.EventManager().EmitEvents(
		sdktypes.Events{types.NewConsolidationEvent(txid), relayertypes.FinalizedProposalEvent(sequence)},
	)
	return nil, nil
}
