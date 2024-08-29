package keeper

import (
	"bytes"
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/btcsuite/btcd/chaincfg"
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
		BlockHeight    collections.Sequence
		BlockHashes    collections.Map[uint64, []byte]
		Deposited      collections.Map[collections.Pair[[]byte, uint32], int64]
		ExecuableQueue collections.Item[types.ExecuableQueue]
		BtcChainConfig *chaincfg.Params
		// this line is used by starport scaffolding # collection/type

		relayerKeeper types.RelayerKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	btcConfig *chaincfg.Params,

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
		BtcChainConfig: btcConfig,
		Pubkey:         collections.NewItem(sb, types.LatestPubkeyKey, "latest_pubkey", codec.CollValue[relayertypes.PublicKey](cdc)),
		BlockHeight:    collections.NewSequence(sb, types.LatestHeightKey, "latest_height"),
		BlockHashes:    collections.NewMap(sb, types.BlockHashsKey, "block_hashs", collections.Uint64Key, collections.BytesValue),
		Deposited:      collections.NewMap(sb, types.DepositedKey, "deposited", collections.PairKeyCodec(collections.BytesKey, collections.Uint32Key), collections.Int64Value),
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

func (k Keeper) NewDeposit(ctx context.Context, deposit *types.Deposit) (*types.ExecuableDeposit, error) {
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
	if !bytes.Equal(blockHash, goatcrypto.DoubleSHA256Sum(deposit.BlockHeader)) {
		return nil, types.ErrInvalidRequest.Wrapf("incorrect block hash, expected %x", blockHash)
	}

	// check if the tx is valid
	tx, txrd := new(wire.MsgTx), bytes.NewReader(deposit.NoWitnessTx)
	if err := tx.DeserializeNoWitness(txrd); err != nil || txrd.Len() > 0 {
		return nil, types.ErrInvalidRequest.Wrapf("invalid non-witness tx")
	}

	if deposit.OutputIndex >= uint32(len(tx.TxOut)) {
		return nil, types.ErrInvalidRequest.Wrap("out of range for outputs")
	}

	// check if the deposit is done
	txid := goatcrypto.DoubleSHA256Sum(deposit.NoWitnessTx)
	deposited, err := k.Deposited.Has(ctx, collections.Join(txid, deposit.OutputIndex))
	if err != nil {
		return nil, err
	}
	if deposited {
		return nil, types.ErrInvalidRequest.Wrap("duplicated deposit")
	}

	// check if the deposit script is valid
	txout := tx.TxOut[deposit.OutputIndex]
	if txout.Value <= 0 {
		return nil, types.ErrInvalidRequest.Wrap("invalid txout amount")
	}

	if err := types.VerifyDespositScript(deposit.RelayerPubkey, deposit.EvmAddress, txout.PkScript); err != nil {
		return nil, types.ErrInvalidRequest.Wrapf("invalid txout script: %s", err.Error())
	}

	// check if the spv is valid
	if !types.VerifyMerkelProof(txid, deposit.BlockHeader[36:68], deposit.IntermediateProof, deposit.TxIndex) {
		return nil, types.ErrInvalidRequest.Wrap("invalid spv")
	}

	return &types.ExecuableDeposit{
		Address: deposit.EvmAddress,
		Txid:    txid,
		Txout:   deposit.OutputIndex,
		Amount:  math.NewInt(txout.Value).Mul(types.Satoshi),
	}, nil
}

func (k Keeper) NewPubkey(ctx context.Context, pubkey *relayertypes.PublicKey) error {
	if err := k.relayerKeeper.AddNewKey(ctx, relayertypes.EncodePublicKey(pubkey)); err != nil {
		return err
	}
	return k.Pubkey.Set(ctx, *pubkey)
}

func (k Keeper) DequeueBitcoinModuleTx() []*ethtypes.Transaction {
	return nil
}
