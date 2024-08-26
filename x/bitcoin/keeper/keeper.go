package keeper

import (
	"bytes"
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/btcsuite/btcd/wire"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/goatnetwork/goat/pkg/utils"
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

		Params        collections.Item[types.Params]
		LatestPubkey  collections.Item[relayertypes.PublicKey]
		LastestHeight collections.Sequence
		BlockHashs    collections.Map[uint64, []byte]
		Deposited     collections.Map[collections.Pair[[]byte, uint32], int64]
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

		relayerKeeper: relayerKeeper,
		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		LatestPubkey:  collections.NewItem(sb, types.LatestPubkeyKey, "latest_pubkey", codec.CollValue[relayertypes.PublicKey](cdc)),
		LastestHeight: collections.NewSequence(sb, types.LatestHeightKey, "latest_height"),
		BlockHashs:    collections.NewMap(sb, types.BlockHashsKey, "block_hashs", collections.Uint64Key, collections.BytesValue),
		Deposited:     collections.NewMap(sb, types.Depositedkey, "deposited", collections.PairKeyCodec(collections.BytesKey, collections.Uint32Key), collections.Int64Value),
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

func (k Keeper) VerifyNewDeposit(ctx context.Context, deposit *types.Deposit) error {
	// check if the pubkey is existed
	hasKey, err := k.relayerKeeper.HasPubkey(ctx, relayertypes.EncodePublicKey(deposit.RelayerPubkey))
	if err != nil {
		return err
	}

	if !hasKey {
		return types.ErrInvalidRequest.Wrap("relayer pubkey not found")
	}

	// check if the block is voted
	blockHash, err := k.BlockHashs.Get(ctx, deposit.BlockNumber)
	if err != nil {
		return err
	}
	if !bytes.Equal(blockHash, utils.DoubleSHA256Sum(deposit.BlockHeader)) {
		return types.ErrInvalidRequest.Wrapf("incorrect block hash, expected %x", blockHash)
	}

	// check if the tx is valid
	tx, txrd := new(wire.MsgTx), bytes.NewReader(deposit.NoWitnessTx)
	if err := tx.DeserializeNoWitness(txrd); err != nil || txrd.Len() > 0 {
		return types.ErrInvalidRequest.Wrapf("invalid non-witness tx")
	}

	if deposit.OutputIndex >= uint32(len(tx.TxOut)) {
		return types.ErrInvalidRequest.Wrap("out of range for the output")
	}

	// check if the deposit is done
	txid := utils.DoubleSHA256Sum(deposit.NoWitnessTx)
	deposited, err := k.Deposited.Has(ctx, collections.Join(txid, deposit.OutputIndex))
	if err != nil {
		return err
	}
	if deposited {
		return types.ErrInvalidRequest.Wrap("duplicated deposit")
	}

	// check if the deposit script is valid
	txout := tx.TxOut[deposit.OutputIndex]
	isValidDepositScript, _ := types.ValidateDespositTxOut(deposit.RelayerPubkey, deposit.EvmAddress, txout.PkScript)
	if !isValidDepositScript {
		return types.ErrInvalidRequest.Wrap("invalid txout script")
	}

	// check if the spv is valid
	if !types.VerifyMerkelProof(txid, deposit.BlockHeader[36:68], deposit.IntermediateProof, deposit.TxIndex) {
		return types.ErrInvalidRequest.Wrap("invalid spv")
	}

	return nil
}
