package types

import (
	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func ExecutableDataToPayload(data *engine.ExecutableData, beaconRoot []byte, execRequests [][]byte) *ExecutionPayload {
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
		BaseFeePerGas: math.NewIntFromBigInt(data.BaseFeePerGas),
		BlockHash:     data.BlockHash.Bytes(),
		Transactions:  data.Transactions,
		BeaconRoot:    beaconRoot,
		BlobGasUsed:   BlobGasUsed,
		ExcessBlobGas: ExcessBlobGas,
		Requests:      execRequests,
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
		BaseFeePerGas: data.BaseFeePerGas.BigInt(),
		BlockHash:     common.BytesToHash(data.BlockHash),
		Transactions:  data.Transactions,
		Withdrawals:   []*ethtypes.Withdrawal{},
		BlobGasUsed:   &data.BlobGasUsed,
		ExcessBlobGas: &data.ExcessBlobGas,
	}

	return res
}

func (payload *ExecutionPayload) LogKeyVals() []any {
	return []any{
		"BlockNumber", payload.BlockNumber,
		"BlockHash", hexutil.Encode(payload.BlockHash),
		"Timestamp", payload.Timestamp,
		"len(Transactions)", len(payload.Transactions),
	}
}
