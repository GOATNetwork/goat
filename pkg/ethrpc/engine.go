package ethrpc

import (
	"context"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
)

func (ec *Client) ForkchoiceUpdatedV3(ctx context.Context, update *engine.ForkchoiceStateV1, params *engine.PayloadAttributes) (engine.ForkChoiceResponse, error) {
	var result engine.ForkChoiceResponse
	err := ec.Client.Client().CallContext(ctx, &result, ForkchoiceUpdatedMethodV3, update, params)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (ec *Client) GetPayloadV3(ctx context.Context, payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error) {
	var result engine.ExecutionPayloadEnvelope
	err := ec.Client.Client().CallContext(ctx, &result, GetPayloadMethodV3, payloadID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *Client) GetPayloadV4(ctx context.Context, payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error) {
	var result engine.ExecutionPayloadEnvelope
	err := ec.Client.Client().CallContext(ctx, &result, GetPayloadMethodV4, payloadID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *Client) NewPayloadV3(ctx context.Context, params *engine.ExecutableData, blobHashes []common.Hash, beaconRoot common.Hash) (*engine.PayloadStatusV1, error) {
	var result engine.PayloadStatusV1
	err := ec.Client.Client().CallContext(ctx, &result, NewPayloadMethodV3, params, blobHashes, beaconRoot)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *Client) NewPayloadV4(ctx context.Context, params *engine.ExecutableData, blobHashes []common.Hash, beaconRoot common.Hash, reqeusts [][]byte) (*engine.PayloadStatusV1, error) {
	var execRequests []hexutil.Bytes
	for i := 0; i < len(reqeusts); i++ {
		execRequests = append(execRequests, reqeusts[i])
	}

	var result engine.PayloadStatusV1
	err := ec.Client.Client().CallContext(ctx, &result, NewPayloadMethodV4, params, blobHashes, beaconRoot, execRequests)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *Client) ExchangeCapabilities(ctx context.Context, caps []string) ([]string, error) {
	var result []string
	err := ec.Client.Client().CallContext(ctx, &result, ExchangeCapabilities, caps)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ec *Client) GetClientVersionV1(ctx context.Context, info engine.ClientVersionV1) ([]engine.ClientVersionV1, error) {
	var result []engine.ClientVersionV1
	err := ec.Client.Client().CallContext(ctx, &result, GetClientVersionV1, info)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ec *Client) GetChainConfig(ctx context.Context) (*params.ChainConfig, error) {
	var result params.ChainConfig
	err := ec.Client.Client().CallContext(ctx, &result, GetChainConfig)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
