package types

import (
	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	lockingtypes "github.com/goatnetwork/goat/x/locking/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
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
		BaseFeePerGas: math.NewIntFromBigInt(data.BaseFeePerGas),
		BlockHash:     data.BlockHash.Bytes(),
		Transactions:  data.Transactions,
		BeaconRoot:    beaconRoot,
		BlobGasUsed:   BlobGasUsed,
		ExcessBlobGas: ExcessBlobGas,
		Requests:      nil, // todo
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
		"ParentHash", hexutil.Encode(payload.ParentHash),
		"FeeRecipient", hexutil.Encode(payload.FeeRecipient),
		"StateRoot", hexutil.Encode(payload.StateRoot),
		"ReceiptsRoot", hexutil.Encode(payload.ReceiptsRoot),
		"LogsBloom", hexutil.Encode(payload.LogsBloom),
		"PrevRandao", hexutil.Encode(payload.PrevRandao),
		"GasLimit", payload.GasLimit,
		"GasUsed", payload.GasUsed,
		"Timestamp", payload.Timestamp,
		"BeaconRoot", hexutil.Encode(payload.BeaconRoot),
		"ExtraData", hexutil.Encode(payload.ExtraData),
		"BaseFeePerGas", payload.BaseFeePerGas.BigInt(),
		"BlobGasUsed", payload.BlobGasUsed,
		"ExcessBlobGas", payload.ExcessBlobGas,
		"len(Transactions)", len(payload.Transactions),
		"len(Requests)", len(payload.Requests),
	}
}

func (payload *ExecutionPayload) DecodeGoatRequests() (
	bridge bitcointypes.ExecRequests,
	relayer relayertypes.ExecRequests,
	locking lockingtypes.ExecRequests,
	err error,
) {
	for _, raw := range payload.Requests {
		req := new(ethtypes.Request)
		err = req.UnmarshalBinary(raw)
		if err != nil {
			return
		}

		switch v := req.Inner().(type) {
		case *ethtypes.GasRevenue:
			locking.GasRevenues = append(locking.GasRevenues, v)
		case *ethtypes.AddVoter:
			relayer.AddVoters = append(relayer.AddVoters, v)
		case *ethtypes.RemoveVoter:
			relayer.RemoveVoters = append(relayer.RemoveVoters, v)
		case *ethtypes.GoatWithdrawal:
			bridge.Withdrawals = append(bridge.Withdrawals, v)
		case *ethtypes.ReplaceByFee:
			bridge.RBFs = append(bridge.RBFs, v)
		case *ethtypes.Cancel1:
			bridge.Cancel1s = append(bridge.Cancel1s, v)
		case *ethtypes.CreateValidator:
			locking.Creates = append(locking.Creates, v)
		case *ethtypes.GoatLock:
			locking.Locks = append(locking.Locks, v)
		case *ethtypes.GoatUnlock:
			locking.Unlocks = append(locking.Unlocks, v)
		case *ethtypes.GoatClaimReward:
			locking.Claims = append(locking.Claims, v)
		case *ethtypes.GoatGrant:
			locking.Grants = append(locking.Grants, v)
		case *ethtypes.UpdateTokenThreshold:
			locking.UpdateThresholds = append(locking.UpdateThresholds, v)
		case *ethtypes.UpdateTokenWeight:
			locking.UpdateWeights = append(locking.UpdateWeights, v)
		}
	}
	return
}
