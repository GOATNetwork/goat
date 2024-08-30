package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func ExecutableDataToPayload(data *engine.ExecutableData) *ExecutionPayload {
	var BlobGasUsed uint64
	if data.BlobGasUsed != nil {
		BlobGasUsed = *data.BlobGasUsed
	}

	var ExcessBlobGas uint64
	if data.ExcessBlobGas != nil {
		ExcessBlobGas = *data.ExcessBlobGas
	}

	return &ExecutionPayload{
		ParentHash:    data.ParentHash.Bytes(),
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
		BlobGasUsed:   BlobGasUsed,
		ExcessBlobGas: ExcessBlobGas,
	}
}

func PayloadToExecutableData(data *ExecutionPayload, propser []byte) *engine.ExecutableData {
	return &engine.ExecutableData{
		ParentHash:    common.BytesToHash(data.ParentHash),
		FeeRecipient:  common.BytesToAddress(propser),
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
}
