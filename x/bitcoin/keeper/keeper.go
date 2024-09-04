package keeper

import (
	"bytes"
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/cosmos/cosmos-sdk/codec"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		addressCodec address.Codec
		storeService store.KVStoreService
		logger       log.Logger
		schema       collections.Schema

		Params         collections.Item[types.Params]
		Pubkey         collections.Item[relayertypes.PublicKey]
		BlockTip       collections.Sequence
		BlockHashes    collections.Map[uint64, []byte]
		Deposited      collections.Map[collections.Pair[[]byte, uint32], uint64]
		EthTxNonce     collections.Sequence
		ExecuableQueue collections.Item[types.ExecuableQueue]
		// this line is used by starport scaffolding # collection/type

		relayerKeeper types.RelayerKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	relayerKeeper types.RelayerKeeper,
) Keeper {

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		storeService: storeService,
		logger:       logger,

		relayerKeeper:  relayerKeeper,
		Params:         collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Pubkey:         collections.NewItem(sb, types.LatestPubkeyKey, "latest_pubkey", codec.CollValue[relayertypes.PublicKey](cdc)),
		BlockTip:       collections.NewSequence(sb, types.LatestHeightKey, "latest_height"),
		BlockHashes:    collections.NewMap(sb, types.BlockHashsKey, "block_hashs", collections.Uint64Key, collections.BytesValue),
		Deposited:      collections.NewMap(sb, types.DepositedKey, "deposited", collections.PairKeyCodec(collections.BytesKey, collections.Uint32Key), collections.Uint64Value),
		EthTxNonce:     collections.NewSequence(sb, types.EthTxNonceKey, "eth_tx_nonce"),
		ExecuableQueue: collections.NewItem(sb, types.ExecuableQueueKey, "queue", codec.CollValue[types.ExecuableQueue](cdc)),
		// this line is used by starport scaffolding # collection/instantiate
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.schema = schema

	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) VerifyDeposit(ctx context.Context, headers map[uint64][]byte, deposit *types.Deposit) (*types.DepositReceipt, error) {
	// check if the pubkey is existed
	hasKey, err := k.relayerKeeper.HasPubkey(ctx, relayertypes.EncodePublicKey(deposit.RelayerPubkey))
	if err != nil {
		return nil, err
	}

	if !hasKey {
		return nil, types.ErrInvalidRequest.Wrap("relayer pubkey not found")
	}

	// check if the block is voted
	blockHash, err := k.BlockHashes.Get(ctx, deposit.BlockNumber)
	if err != nil {
		return nil, err
	}

	rawHeader := headers[deposit.BlockNumber]
	if len(rawHeader) != types.RawBtcHeaderSize {
		return nil, types.ErrInvalidRequest.Wrapf("invalid block header for %d", deposit.BlockNumber)
	}

	if !bytes.Equal(blockHash, goatcrypto.DoubleSHA256Sum(rawHeader)) {
		return nil, types.ErrInvalidRequest.Wrap("inconsistent block hash")
	}

	// check if the tx is valid
	tx, txrd := new(wire.MsgTx), bytes.NewReader(deposit.NoWitnessTx)
	if err := tx.DeserializeNoWitness(txrd); err != nil || txrd.Len() > 0 {
		return nil, types.ErrInvalidRequest.Wrapf("invalid non-witness tx")
	}

	if deposit.OutputIndex >= uint32(len(tx.TxOut)) {
		return nil, types.ErrInvalidRequest.Wrap("output index out of range")
	}

	// check if the deposit is duplicated
	txid := goatcrypto.DoubleSHA256Sum(deposit.NoWitnessTx)
	deposited, err := k.Deposited.Has(ctx, collections.Join(txid, deposit.OutputIndex))
	if err != nil {
		return nil, err
	}
	if deposited {
		return nil, types.ErrInvalidRequest.Wrap("duplicated deposit")
	}

	param, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// check if the deposit script is valid
	txout := tx.TxOut[deposit.OutputIndex]
	if txout.Value < int64(param.MinDepositAmount) {
		return nil, types.ErrInvalidRequest.Wrap("amount too low")
	}

	switch deposit.Version {
	case 0:
		if err := types.VerifyDespositScriptV0(deposit.RelayerPubkey, deposit.EvmAddress, txout.PkScript); err != nil {
			k.logger.Debug("invalid deposit version 0 script", "txid", chainhash.Hash(txid).String(), "txout", deposit.OutputIndex, "err", err.Error())
			return nil, types.ErrInvalidRequest.Wrap("invalid deposit version 0 script")
		}
	case 1:
		if deposit.OutputIndex != 0 || len(tx.TxOut) < 2 {
			return nil, types.ErrInvalidRequest.Wrap("invalid txout index for version 1 deposit")
		}
		if err := types.VerifyDespositScriptV1(deposit.RelayerPubkey,
			param.DepositMagicPrefix, deposit.EvmAddress, txout.PkScript, tx.TxOut[1].PkScript); err != nil {
			k.logger.Debug("invalid deposit version 1 script", "txid", chainhash.Hash(txid).String(), "txout", deposit.OutputIndex, "err", err.Error())
			return nil, types.ErrInvalidRequest.Wrap("invalid deposit version 1 script")
		}
	default:
		return nil, types.ErrInvalidRequest.Wrapf("unknown deposit version")
	}

	// check if the spv is valid
	if !types.VerifyMerkelProof(txid, rawHeader[36:68], deposit.IntermediateProof, deposit.TxIndex) {
		return nil, types.ErrInvalidRequest.Wrap("invalid spv")
	}

	return &types.DepositReceipt{
		Address: deposit.EvmAddress,
		Txid:    txid,
		Txout:   deposit.OutputIndex,
		Amount:  uint64(txout.Value),
	}, nil
}

func (k Keeper) NewPubkey(ctx context.Context, pubkey *relayertypes.PublicKey) error {
	if err := k.relayerKeeper.AddNewKey(ctx, relayertypes.EncodePublicKey(pubkey)); err != nil {
		return err
	}
	return k.Pubkey.Set(ctx, *pubkey)
}

func (k Keeper) DequeueBitcoinModuleTx(ctx context.Context) (txs []*ethtypes.Transaction, err error) {
	queue, err := k.ExecuableQueue.Get(ctx)
	if err != nil {
		return nil, err
	}

	txNonce, err := k.EthTxNonce.Peek(ctx)
	if err != nil {
		return nil, err
	}

	// pop block hash up to 1
	{
		tip, err := k.BlockTip.Peek(ctx)
		if err != nil {
			return nil, err
		}

		if queue.BlockNumber < tip {
			queue.BlockNumber++

			blockHash, err := k.BlockHashes.Get(ctx, queue.BlockNumber)
			if err != nil {
				return nil, err
			}
			txs = append(txs, types.NewBitcoinHashEthTx(txNonce, blockHash))
			txNonce++
		}
	}

	// pop deposit up to 8
	{
		var n int
		for i := 0; n < len(queue.Deposits) && i < 8; i++ {
			deposit := queue.Deposits[i]
			txs = append(txs, deposit.EthTx(txNonce))

			n++
			txNonce++
		}
		queue.Deposits = queue.Deposits[n:]
	}

	if len(txs) > 0 {
		if err := k.ExecuableQueue.Set(ctx, queue); err != nil {
			return nil, err
		}
		if err := k.EthTxNonce.Set(ctx, txNonce); err != nil {
			return nil, err
		}
	}
	return txs, nil
}
