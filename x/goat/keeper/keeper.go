package keeper

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/goatnetwork/goat/pkg/ethrpc"
	"github.com/goatnetwork/goat/x/goat/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		addressCodec address.Codec
		storeService store.KVStoreService
		logger       log.Logger

		Schema     collections.Schema
		Params     collections.Item[types.Params]
		BeaconRoot collections.Item[[]byte] // the cometbft blockhash
		Block      collections.Item[types.ExecutionPayload]
		ethclient  ethrpc.EngineClient
		txConfig   client.TxConfig
		// this line is used by starport scaffolding # collection/type

		bitcoinKeeper types.BitcoinKeeper
		lockingKeeper types.LockingKeeper
		relayerKeeper types.RelayerKeeper
		accountKeeper types.AccountKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService store.KVStoreService,
	logger log.Logger,

	bitcoinKeeper types.BitcoinKeeper,
	lockingKeeper types.LockingKeeper,
	relayerKeeper types.RelayerKeeper,
	accountKeeper types.AccountKeeper,
	ethclient ethrpc.EngineClient,
	txConfig client.TxConfig,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		storeService: storeService,
		logger:       logger,

		bitcoinKeeper: bitcoinKeeper,
		lockingKeeper: lockingKeeper,
		relayerKeeper: relayerKeeper,
		accountKeeper: accountKeeper,
		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Block:         collections.NewItem(sb, types.BlockKey, "block", codec.CollValue[types.ExecutionPayload](cdc)),
		BeaconRoot:    collections.NewItem(sb, types.ConsHashKey, "consensus_hash", collections.BytesValue),
		ethclient:     ethclient,
		txConfig:      txConfig,
		// this line is used by starport scaffolding # collection/instantiate
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Dequeue(ctx context.Context) ([]hexutil.Bytes, error) {
	btcTxs, err := k.bitcoinKeeper.DequeueBitcoinModuleTx(ctx)
	if err != nil {
		return nil, err
	}
	k.Logger().Debug("dequeue bitcoin module txs", "len", len(btcTxs))

	lockingTxs, err := k.lockingKeeper.DequeueLockingModuleTx(ctx)
	if err != nil {
		return nil, err
	}
	k.Logger().Debug("dequeue locking module txs", "len", len(lockingTxs))

	res := make([]hexutil.Bytes, 0, len(btcTxs)+len(lockingTxs))
	for _, tx := range btcTxs {
		raw, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		res = append(res, raw)
	}

	for _, tx := range lockingTxs {
		raw, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		res = append(res, raw)
	}

	return res, nil
}

// VerifyDequeue verifies if the goat transactions of the new eth block are consistent with expected here
func (k Keeper) VerifyDequeue(ctx context.Context, txRoot []byte, txs [][]byte) error {
	// goat-geth will check the tx root hash
	if len(txRoot) != 33 {
		return errors.New("invalid goat tx root")
	}

	goatTxLen := int(txRoot[0])
	if len(txs) < goatTxLen {
		return errors.New("tx length is less than expected")
	}

	btcTxs, err := k.bitcoinKeeper.DequeueBitcoinModuleTx(ctx)
	if err != nil {
		return err
	}
	k.Logger().Debug("verifying dequeue bitcoin module txs", "len", len(btcTxs))

	if len(txs) < len(btcTxs) {
		return fmt.Errorf("bitcoin module txs length mismatched: len(txs)=%d len(mod)=%d", len(txs), len(btcTxs))
	}

	for idx, tx := range btcTxs {
		raw, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		if !bytes.Equal(raw, txs[idx]) {
			return fmt.Errorf("bridge tx %d bytes mismatched", idx)
		}
		goatTxLen--
	}

	lockingTxs, err := k.lockingKeeper.DequeueLockingModuleTx(ctx)
	if err != nil {
		return err
	}

	k.Logger().Debug("verifing dequeue locking module txs", "len", len(lockingTxs))
	txs = txs[len(btcTxs):]
	if len(txs) < len(lockingTxs) {
		return fmt.Errorf("locking module txs length mismatched: len(txs)=%d len(mod)=%d", len(txs), len(lockingTxs))
	}

	for idx, tx := range lockingTxs {
		raw, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		if !bytes.Equal(raw, txs[idx]) {
			return fmt.Errorf("locking tx %d bytes mismatched", idx)
		}
		goatTxLen--
	}

	if goatTxLen != 0 {
		return errors.New("goat txs length mismatched")
	}
	return nil
}

// Finalized notifies goat-geth to update fork choice state
// if there are any errors, the FinalizeBlock phase will be failed
// we don't use timeout here, validators are responsible for a reliable node
func (k Keeper) Finalized(ctx context.Context) error { // EndBlock phase only!
	block, err := k.Block.Get(ctx)
	if err != nil {
		return err
	}

	k.Logger().Info("notify NewPayload", "number", block.BlockNumber)
	plRes, err := k.ethclient.NewPayloadV3(ctx, types.PayloadToExecutableData(&block),
		[]common.Hash{}, common.BytesToHash(block.BeaconRoot))
	if err != nil {
		return err
	}

	if plRes.Status == engine.INVALID {
		return errors.New("invalid from NewPayloadV3 api")
	}

	// set current block hash to head state and set previous block hash to safe and finalized state
	k.Logger().Info("notify ForkchoiceUpdated",
		"head", hexutil.Encode(block.BlockHash), "finalized", hexutil.Encode(block.ParentHash))
	forkRes, err := k.ethclient.ForkchoiceUpdatedV3(ctx, &engine.ForkchoiceStateV1{
		HeadBlockHash:      common.BytesToHash(block.BlockHash),
		SafeBlockHash:      common.BytesToHash(block.ParentHash),
		FinalizedBlockHash: common.BytesToHash(block.ParentHash),
	}, nil)
	if err != nil {
		return err
	}

	if forkRes.PayloadStatus.Status == engine.INVALID {
		return errors.New("invalid from ForkchoiceUpdatedV3 api")
	}

	return nil
}
