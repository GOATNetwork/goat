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
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
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
		ethclient  *ethrpc.Client
		txConfig   client.TxConfig
		// this line is used by starport scaffolding # collection/type

		bitcoinKeeper types.BitcoinKeeper
		lockingKeeper types.LockingKeeper
		relayerKeeper types.RelayerKeeper
		accountKeeper types.AccountKeeper
		PrivKey       cryptotypes.PrivKey
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
	ethclient *ethrpc.Client,
	txConfig client.TxConfig,
	privKey cryptotypes.PrivKey,
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
		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Block:         collections.NewItem(sb, types.BlockKey, "block", codec.CollValue[types.ExecutionPayload](cdc)),
		BeaconRoot:    collections.NewItem(sb, types.ConsHashKey, "consensus_hash", collections.BytesValue),
		ethclient:     ethclient,
		txConfig:      txConfig,
		PrivKey:       privKey,
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

	lockingTxs, err := k.lockingKeeper.DequeueLockingModuleTx(ctx)
	if err != nil {
		return nil, err
	}

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

	if len(txs) < len(btcTxs) {
		return errors.New("bridge txs length mismatched")
	}

	for idx, tx := range btcTxs {
		raw, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		if !bytes.Equal(raw, txs[idx]) {
			return errors.New("bridge tx bytes mismatched")
		}
		goatTxLen--
	}

	lockingTxs, err := k.lockingKeeper.DequeueLockingModuleTx(ctx)
	if err != nil {
		return err
	}

	txs = txs[len(btcTxs):]
	if len(txs) < len(lockingTxs) {
		return errors.New("locking txs length mismatched")
	}

	for idx, tx := range lockingTxs {
		raw, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		if !bytes.Equal(raw, txs[idx]) {
			return errors.New("locking tx bytes mismatched")
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
	k.Logger().Info("notify ForkchoiceUpdated", "head", block.BlockHash, "finalized", block.ParentHash)
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
