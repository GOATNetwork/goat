package keeper

import (
	"bytes"
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"github.com/btcsuite/btcd/wire"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

		Params      collections.Item[types.Params]
		Pubkey      collections.Item[relayertypes.PublicKey]
		BlockTip    collections.Sequence
		BlockHashes collections.Map[uint64, []byte]
		Deposited   collections.Map[collections.Pair[[]byte, uint32], uint64]
		EthTxNonce  collections.Sequence
		Withdrawals collections.Map[uint64, types.Withdrawal]
		// processing withdrawal(a pair of pid and the details)
		ProcessID  collections.Sequence
		Processing collections.Map[uint64, types.Processing]
		EthTxQueue collections.Item[types.EthTxQueue]

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

		relayerKeeper: relayerKeeper,
		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Pubkey:        collections.NewItem(sb, types.LatestPubkeyKey, "latest_pubkey", codec.CollValue[relayertypes.PublicKey](cdc)),
		BlockTip:      collections.NewSequence(sb, types.LatestHeightKey, "latest_height"),
		BlockHashes:   collections.NewMap(sb, types.BlockHashsKey, "block_hashs", collections.Uint64Key, collections.BytesValue),
		Deposited:     collections.NewMap(sb, types.DepositedKey, "deposited", collections.PairKeyCodec(collections.BytesKey, collections.Uint32Key), collections.Uint64Value),
		EthTxNonce:    collections.NewSequence(sb, types.EthTxNonceKey, "eth_tx_nonce"),
		EthTxQueue:    collections.NewItem(sb, types.EthTxQueueKey, "eth_tx_queue", codec.CollValue[types.EthTxQueue](cdc)),
		Withdrawals:   collections.NewMap(sb, types.WithdrawalKey, "withdrawals", collections.Uint64Key, codec.CollValue[types.Withdrawal](cdc)),
		Processing:    collections.NewMap(sb, types.ProcessingKey, "processings", collections.Uint64Key, codec.CollValue[types.Processing](cdc)),
		ProcessID:     collections.NewSequence(sb, types.ProcessIDKey, "process_id"),
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

func (k Keeper) VerifyDeposit(ctx context.Context, headers map[uint64][]byte, deposit *types.Deposit) (*types.DepositExecReceipt, error) {
	// check if the pubkey is existed
	hasKey, err := k.relayerKeeper.HasPubkey(ctx, relayertypes.EncodePublicKey(deposit.RelayerPubkey))
	if err != nil {
		return nil, err
	}

	if !hasKey {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "relayer pubkey not found")
	}

	// check if the block is voted
	blockHash, err := k.BlockHashes.Get(ctx, deposit.BlockNumber)
	if err != nil {
		return nil, err
	}

	// coinbase confirmation check
	if deposit.TxIndex == 0 {
		tip, err := k.BlockTip.Peek(ctx)
		if err != nil {
			return nil, err
		}
		if tip < deposit.BlockNumber+100 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "coinbase tx should be confirmed more than 100 blocks")
		}
	}

	rawHeader := headers[deposit.BlockNumber]
	if len(rawHeader) != types.RawBtcHeaderSize {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block header for %d", deposit.BlockNumber)
	}

	if !bytes.Equal(blockHash, goatcrypto.DoubleSHA256Sum(rawHeader)) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "inconsistent block hash")
	}

	// check if the tx is valid
	tx, txrd := new(wire.MsgTx), bytes.NewReader(deposit.NoWitnessTx)
	if err := tx.DeserializeNoWitness(txrd); err != nil || txrd.Len() > 0 {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid non-witness tx")
	}

	if deposit.OutputIndex >= uint32(len(tx.TxOut)) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "output index out of range")
	}

	// check if the deposit is duplicated
	txid := goatcrypto.DoubleSHA256Sum(deposit.NoWitnessTx)
	deposited, err := k.Deposited.Has(ctx, collections.Join(txid, deposit.OutputIndex))
	if err != nil {
		return nil, err
	}
	if deposited {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "duplicated deposit")
	}

	param, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// check if the deposit script is valid
	txOut := tx.TxOut[deposit.OutputIndex]
	txAmount := uint64(txOut.Value)
	if txAmount < param.MinDepositAmount {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount too low")
	}

	switch deposit.Version {
	case 0:
		if err := types.VerifyDespositScriptV0(deposit.RelayerPubkey, deposit.EvmAddress, txOut.PkScript); err != nil {
			k.logger.Warn("invalid deposit version 0 script", "txid", types.BtcTxid(txid), "txout", deposit.OutputIndex, "err", err.Error())
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid deposit version 0 script")
		}
	case 1:
		if deposit.OutputIndex != 0 || len(tx.TxOut) < 2 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid txout index for version 1 deposit")
		}
		if err := types.VerifyDespositScriptV1(deposit.RelayerPubkey,
			param.DepositMagicPrefix, deposit.EvmAddress, txOut.PkScript, tx.TxOut[1].PkScript); err != nil {
			k.logger.Warn("invalid deposit version 1 script", "txid", types.BtcTxid(txid), "txout", deposit.OutputIndex, "err", err.Error())
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid deposit version 1 script")
		}
	default:
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unknown deposit version")
	}

	// check if the spv is valid
	if !types.VerifyMerkelProof(txid, rawHeader[36:68], deposit.IntermediateProof, deposit.TxIndex) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid spv")
	}

	tax := uint64(0)
	if param.DepositTaxRate > 0 && txAmount > types.MaxTaxBP {
		tax = txAmount / types.MaxTaxBP * param.DepositTaxRate
		// 0 represents no limit
		if param.MaxDepositTax > 0 && tax > param.MaxDepositTax {
			tax = param.MaxDepositTax
		}
		txAmount -= tax
	}

	return &types.DepositExecReceipt{
		Address: deposit.EvmAddress,
		Txid:    txid,
		Txout:   deposit.OutputIndex,
		Amount:  txAmount,
		Tax:     tax,
	}, nil
}

func (k Keeper) NewPubkey(ctx context.Context, pubkey *relayertypes.PublicKey) error {
	if err := k.relayerKeeper.AddNewKey(ctx, relayertypes.EncodePublicKey(pubkey)); err != nil {
		return err
	}
	return k.Pubkey.Set(ctx, *pubkey)
}

// MustHasKey is for genesis check
func (k Keeper) MustHasKey(ctx context.Context, pubkey *relayertypes.PublicKey) {
	hasKey, err := k.relayerKeeper.HasPubkey(ctx, relayertypes.EncodePublicKey(pubkey))
	if err != nil {
		panic(err)
	}
	if !hasKey {
		panic(fmt.Sprintf("the pubkey %x doesn't exists", relayertypes.EncodePublicKey(pubkey)))
	}
}
