package types

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func ExecutableDataToPayload(data *engine.ExecutableData, beaconRoot []byte) *ExecutionPayload {
	var BlobGasUsed uint64
	if data.BlobGasUsed != nil {
		BlobGasUsed = *data.BlobGasUsed
	}

	var ExcessBlobGas uint64
	if data.ExcessBlobGas != nil {
		ExcessBlobGas = *data.ExcessBlobGas
	}

	res := &ExecutionPayload{
		ParentHash:    data.ParentHash.Bytes(),
		FeeRecipient:  data.FeeRecipient.Bytes(),
		StateRoot:     data.StateRoot.Bytes(),
		ReceiptsRoot:  data.ReceiptsRoot.Bytes(),
		LogsBloom:     data.LogsBloom,
		PrevRandao:    data.Random.Bytes(),
		BlockNumber:   data.Number,
		GasLimit:      data.GasLimit,
		GasUsed:       data.GasUsed,
		Timestamp:     data.Timestamp,
		ExtraData:     data.ExtraData,
		BaseFeePerGas: data.BaseFeePerGas.Bytes(),
		BlockHash:     data.BlockHash.Bytes(),
		Transactions:  data.Transactions,
		BeaconRoot:    beaconRoot,
		BlobGasUsed:   BlobGasUsed,
		ExcessBlobGas: ExcessBlobGas,
	}

	if len(data.GasRevenues) == 1 {
		v := data.GasRevenues[0]
		res.GasRevenue = &GasRevenueReq{Amount: math.NewIntFromBigInt(v.Amount)}
	}

	for _, v := range data.AddVoters {
		res.AddVoterReq = append(res.AddVoterReq, &AddVoterReq{Voter: v.Voter.Bytes(), PubkeyHash: v.Pubkey.Bytes()})
	}

	for _, v := range data.RemoveVoters {
		res.RmVoterReq = append(res.RmVoterReq, &RemoveVoterReq{Voter: v.Voter.Bytes()})
	}

	for _, v := range data.GoatWithdrawals {
		res.WithdrawalReq = append(res.WithdrawalReq, &WithdrawalReq{
			Id:         v.Id,
			Amount:     v.Amount,
			MaxTxPrice: v.MaxTxPrice,
			Address:    v.Address,
		})
	}

	for _, v := range data.ReplaceByFees {
		res.RbfReq = append(res.RbfReq, &ReplaceByFeeReq{Id: v.Id, MaxTxPrice: v.MaxTxPrice})
	}

	for _, v := range data.Cancel1s {
		res.Cancel1Req = append(res.Cancel1Req, &Cancel1Req{Id: v.Id})
	}
	return res
}

func PayloadToExecutableData(data *ExecutionPayload) *engine.ExecutableData {
	if data.Transactions == nil {
		data.Transactions = [][]byte{}
	}

	res := &engine.ExecutableData{
		ParentHash:    common.BytesToHash(data.ParentHash),
		FeeRecipient:  common.BytesToAddress(data.FeeRecipient),
		StateRoot:     common.BytesToHash(data.StateRoot),
		ReceiptsRoot:  common.BytesToHash(data.ReceiptsRoot),
		LogsBloom:     data.LogsBloom,
		Random:        common.BytesToHash(data.PrevRandao),
		Number:        data.BlockNumber,
		GasLimit:      data.GasLimit,
		GasUsed:       data.GasUsed,
		Timestamp:     data.Timestamp,
		ExtraData:     data.ExtraData,
		BaseFeePerGas: new(big.Int).SetBytes(data.BaseFeePerGas),
		BlockHash:     common.BytesToHash(data.BlockHash),
		Transactions:  data.Transactions,
		Withdrawals:   []*ethtypes.Withdrawal{},
		BlobGasUsed:   &data.BlobGasUsed,
		ExcessBlobGas: &data.ExcessBlobGas,
	}

	if v := data.GasRevenue; v != nil {
		res.GasRevenues = append(res.GasRevenues, ethtypes.NewGoatGasRevenue(v.Amount.BigInt()))
	}

	for _, v := range data.AddVoterReq {
		res.AddVoters = append(res.AddVoters, &ethtypes.AddVoter{
			Voter:  common.BytesToAddress(v.Voter),
			Pubkey: common.BytesToHash(v.PubkeyHash),
		})
	}

	for _, v := range data.RmVoterReq {
		res.RemoveVoters = append(res.RemoveVoters, &ethtypes.RemoveVoter{
			Voter: common.BytesToAddress(v.Voter),
		})
	}

	for _, v := range data.WithdrawalReq {
		res.GoatWithdrawals = append(res.GoatWithdrawals, &ethtypes.GoatWithdrawal{
			Id:         v.Id,
			Amount:     v.Amount,
			MaxTxPrice: v.MaxTxPrice,
			Address:    v.Address,
		})
	}

	for _, v := range data.RbfReq {
		res.ReplaceByFees = append(res.ReplaceByFees, &ethtypes.ReplaceByFee{
			Id:         v.Id,
			MaxTxPrice: v.MaxTxPrice,
		})
	}

	for _, v := range data.Cancel1Req {
		res.Cancel1s = append(res.Cancel1s, &ethtypes.Cancel1{Id: v.Id})
	}

	return res
}

func (payload *ExecutionPayload) LogKeyVals() []any {
	return []any{
		"BlockNumber",
		payload.BlockNumber,
		"BlockHash",
		hexutil.Encode(payload.BlockHash),
		"ParentHash",
		hexutil.Encode(payload.ParentHash),
		"FeeRecipient",
		hexutil.Encode(payload.FeeRecipient),
		"StateRoot",
		hexutil.Encode(payload.StateRoot),
		"ReceiptsRoot",
		hexutil.Encode(payload.ReceiptsRoot),
		"LogsBloom",
		hexutil.Encode(payload.LogsBloom),
		"PrevRandao",
		hexutil.Encode(payload.PrevRandao),
		"GasLimit",
		payload.GasLimit,
		"GasUsed",
		payload.GasUsed,
		"Timestamp",
		payload.Timestamp,
		"ExtraData",
		hexutil.Encode(payload.ExtraData),
		"BaseFeePerGas",
		new(big.Int).SetBytes(payload.BaseFeePerGas).String(),
		"len(Transactions)",
		len(payload.Transactions),
		"BeaconRoot",
		hexutil.Encode(payload.BeaconRoot),
		"BlobGasUsed",
		payload.BlobGasUsed,
		"ExcessBlobGas",
		payload.ExcessBlobGas,
		"GasRevenue",
		payload.GasRevenue,
		"len(AddVoters)", len(payload.AddVoterReq),
		"len(RemoveVoters)", len(payload.RmVoterReq),
		"len(BridgeWithdrawals)", len(payload.WithdrawalReq),
		"len(ReplaceByFees)", len(payload.RbfReq),
		"len(Cancel1s)", len(payload.Cancel1Req),
	}
}
