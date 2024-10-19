package ethrpc

import (
	"context"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

type EngineClient interface {
	ForkchoiceUpdatedV3(ctx context.Context, update *engine.ForkchoiceStateV1, params *engine.PayloadAttributes) (engine.ForkChoiceResponse, error)
	GetPayloadV4(ctx context.Context, payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error)
	NewPayloadV4(ctx context.Context, params *engine.ExecutableData, versionedHashes []common.Hash, beaconRoot common.Hash, requests [][]byte) (*engine.PayloadStatusV1, error)
	ExchangeCapabilities(ctx context.Context, caps []string) ([]string, error)
	GetClientVersionV1(ctx context.Context, info engine.ClientVersionV1) ([]engine.ClientVersionV1, error)
	GetChainConfig(ctx context.Context) (*params.ChainConfig, error)
}
