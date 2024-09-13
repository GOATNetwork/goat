package keeper

import (
	"bytes"
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	goattypes "github.com/goatnetwork/goat/x/goat/types"
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
		Withdrawals    collections.Map[uint64, types.Withdrawal]
		Processing     collections.Map[[]byte, types.WithdrawalIds] // processing withdrawal(a pair of txid and withdrawal id list)
		ExecuableQueue collections.Item[types.ExecuableQueue]

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

		Withdrawals: collections.NewMap(sb, types.WithdrawalKey, "withdrawals", collections.Uint64Key, codec.CollValue[types.Withdrawal](cdc)),
		Processing:  collections.NewMap(sb, types.ProcessingKey, "processings", collections.BytesKey, codec.CollValue[types.WithdrawalIds](cdc)),
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

func (k Keeper) VerifyDeposit(ctx context.Context, headers map[uint64][]byte, deposit *types.Deposit) (*types.DepositExecReceipt, error) {
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

	// coinbase confirmation check
	if deposit.TxIndex == 0 {
		tip, err := k.BlockTip.Peek(ctx)
		if err != nil {
			return nil, err
		}
		if tip < deposit.BlockNumber+100 {
			return nil, types.ErrInvalidRequest.Wrap("coinbase tx should be confirmed more than 100 blocks")
		}
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
			k.logger.Debug("invalid deposit version 0 script", "txid", types.BtcTxid(txid), "txout", deposit.OutputIndex, "err", err.Error())
			return nil, types.ErrInvalidRequest.Wrap("invalid deposit version 0 script")
		}
	case 1:
		if deposit.OutputIndex != 0 || len(tx.TxOut) < 2 {
			return nil, types.ErrInvalidRequest.Wrap("invalid txout index for version 1 deposit")
		}
		if err := types.VerifyDespositScriptV1(deposit.RelayerPubkey,
			param.DepositMagicPrefix, deposit.EvmAddress, txout.PkScript, tx.TxOut[1].PkScript); err != nil {
			k.logger.Debug("invalid deposit version 1 script", "txid", types.BtcTxid(txid), "txout", deposit.OutputIndex, "err", err.Error())
			return nil, types.ErrInvalidRequest.Wrap("invalid deposit version 1 script")
		}
	default:
		return nil, types.ErrInvalidRequest.Wrapf("unknown deposit version")
	}

	// check if the spv is valid
	if !types.VerifyMerkelProof(txid, rawHeader[36:68], deposit.IntermediateProof, deposit.TxIndex) {
		return nil, types.ErrInvalidRequest.Wrap("invalid spv")
	}

	return &types.DepositExecReceipt{
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

func (k Keeper) DequeueBitcoinModuleTx(ctx context.Context) (txs []*ethtypes.Transaction, err error) {
	const (
		MaxDeposit    = 8
		MaxWithdrawal = 8
	)

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
		for i := 0; i < len(queue.Deposits) && n < MaxDeposit; i++ {
			deposit := queue.Deposits[i]
			txs = append(txs, deposit.EthTx(txNonce))

			n++
			txNonce++
		}
		queue.Deposits = queue.Deposits[n:]
	}

	// pop paid and reject withdrwal up to 8
	{
		var n int
		for i := 0; i < len(queue.PaidWithdrawals) && n < MaxWithdrawal; i++ {
			paid := queue.PaidWithdrawals[i]
			txs = append(txs, paid.EthTx(txNonce))

			n++
			txNonce++
		}

		for i := 0; n < len(queue.RejectedWithdrawals) && n < MaxWithdrawal; i++ {
			txs = append(txs, types.NewRejectEthTx(queue.RejectedWithdrawals[i], txNonce))
			n++
			txNonce++
		}
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

func (k Keeper) ProcessBridgeRequest(ctx context.Context, withdrawals []*goattypes.WithdrawalReq, rbf []*goattypes.ReplaceByFeeReq, cancel1 []*goattypes.Cancel1Req) (sdktypes.Events, error) {
	reqLens := len(withdrawals) + len(rbf) + len(cancel1)
	if reqLens == 0 {
		return nil, nil
	}

	param, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	chaincfg := param.ChainConfig.ToBtcdParam()

	queue, err := k.ExecuableQueue.Get(ctx)
	if err != nil {
		return nil, err
	}

	var rejecting bool

	events := make(sdktypes.Events, 0, reqLens)
	for _, v := range withdrawals {
		// reject if we have an invalid address
		addr, err := btcutil.DecodeAddress(v.Address, chaincfg)
		if err != nil {
			queue.RejectedWithdrawals = append(queue.RejectedWithdrawals, v.Id)
			rejecting = true
			continue
		}

		script, err := txscript.PayToAddrScript(addr)
		if err != nil {
			queue.RejectedWithdrawals = append(queue.RejectedWithdrawals, v.Id)
			rejecting = true
			continue
		}

		if err := k.Withdrawals.Set(ctx, v.Id, types.Withdrawal{
			Address:       v.Address,
			RequestAmount: v.Amount,
			MaxTxPrice:    v.MaxTxPrice,
			OutputScript:  script,
			Status:        types.WITHDRAWAL_STATUS_PENDING,
		}); err != nil {
			return nil, err
		}
		events = append(events, types.NewWithdrawalRequestEvent(v.Address, v.Id, v.MaxTxPrice, v.Amount))
	}

	if rejecting {
		if err := k.ExecuableQueue.Set(ctx, queue); err != nil {
			return nil, err
		}
	}

	for _, v := range rbf {
		withdrawal, err := k.Withdrawals.Get(ctx, v.Id)
		if err != nil {
			return nil, err
		}
		if withdrawal.Status != types.WITHDRAWAL_STATUS_PENDING {
			continue
		}
		withdrawal.MaxTxPrice = v.MaxTxPrice
		if err := k.Withdrawals.Set(ctx, v.Id, withdrawal); err != nil {
			return nil, err
		}
		events = append(events, types.NewWithdrawalReplaceEvent(v.Id, v.MaxTxPrice))
	}

	for _, v := range cancel1 {
		withdrawal, err := k.Withdrawals.Get(ctx, v.Id)
		if err != nil {
			return nil, err
		}

		if withdrawal.Status != types.WITHDRAWAL_STATUS_PENDING {
			continue
		}

		withdrawal.Status = types.WITHDRAWAL_STATUS_CANCELING
		if err := k.Withdrawals.Set(ctx, v.Id, withdrawal); err != nil {
			return nil, err
		}
		events = append(events, types.NewWithdrawalCancellationEvent(v.Id))
	}
	return events, nil
}
