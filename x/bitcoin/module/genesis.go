package bitcoin

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/bitcoin/keeper"
	"github.com/goatnetwork/goat/x/bitcoin/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := genState.Validate(); err != nil {
		panic(err)
	}

	if err := k.Params.Set(ctx, genState.Params); err != nil {
		panic(err)
	}

	if err := k.BlockTip.Set(ctx, genState.BlockTip); err != nil {
		panic(err)
	}

	for idx, hash := range genState.BlockHashes {
		if len(hash) != sha256.Size {
			panic(fmt.Sprintf("invalid block hash length: %x", hash))
		}
		if genState.BlockTip < uint64(idx) {
			panic("invalid block hash length for block tip")
		}
		if err := k.BlockHashes.Set(ctx, genState.BlockTip-uint64(idx), hash); err != nil {
			panic(err)
		}
	}

	if err := k.EthTxNonce.Set(ctx, genState.EthTxNonce); err != nil {
		panic(err)
	}

	k.MustHasKey(ctx, genState.Pubkey)
	if err := k.Pubkey.Set(ctx, *genState.Pubkey); err != nil {
		panic(err)
	}

	if err := k.EthTxQueue.Set(ctx, genState.GetEthTxQueue()); err != nil {
		panic(err)
	}

	// Deposits
	for _, item := range genState.Deposits {
		if err := k.Deposited.Set(ctx, collections.Join(item.Txid, item.Txout), item.Amount); err != nil {
			panic(err)
		}
	}

	// withdrawals
	for _, item := range genState.Withdrawals {
		if err := k.Withdrawals.Set(ctx, item.Id, item.Withdrawal); err != nil {
			panic(err)
		}
	}

	// processing
	for _, item := range genState.Processing {
		if err := k.Processing.Set(ctx, item.Id, item.Processing); err != nil {
			panic(err)
		}
	}
	if err := k.ProcessID.Set(ctx, genState.LatestProcessId); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := new(types.GenesisState)

	var err error
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	genesis.BlockTip, err = k.BlockTip.Peek(ctx)
	if err != nil {
		panic(err)
	}

	pubkey, err := k.Pubkey.Get(ctx)
	if err != nil {
		panic(err)
	}
	genesis.Pubkey = &pubkey

	// BlockHashes
	{
		for i := genesis.BlockTip + 1; i > 0; i-- {
			hash, err := k.BlockHashes.Get(ctx, i-1)
			if err != nil {
				if errors.Is(err, collections.ErrNotFound) {
					break
				}
				panic(err)
			}
			genesis.BlockHashes = append(genesis.BlockHashes, hash)
		}
	}

	{
		genesis.EthTxNonce, err = k.EthTxNonce.Peek(ctx)
		if err != nil {
			panic(err)
		}

		genesis.EthTxQueue, err = k.EthTxQueue.Get(ctx)
		if err != nil {
			panic(err)
		}
	}

	// deposited
	{
		iter, err := k.Deposited.Iterate(ctx, nil)
		if err != nil {
			panic(err)
		}
		defer iter.Close()

		for ; iter.Valid(); iter.Next() {
			kv, err := iter.KeyValue()
			if err != nil {
				panic(err)
			}

			genesis.Deposits = append(genesis.Deposits, types.DepositGenesis{
				Txid:   kv.Key.K1(),
				Txout:  kv.Key.K2(),
				Amount: kv.Value,
			})
		}
	}

	// withdrawals
	{
		iter, err := k.Withdrawals.Iterate(ctx, (&collections.Range[uint64]{}).Descending())
		if err != nil {
			panic(err)
		}
		defer iter.Close()

		for ; iter.Valid(); iter.Next() {
			kv, err := iter.KeyValue()
			if err != nil {
				panic(err)
			}

			genesis.Withdrawals = append(genesis.Withdrawals, types.WithdrawalGenesis{
				Id:         kv.Key,
				Withdrawal: kv.Value,
			})
		}
	}

	// processing
	{
		genesis.LatestProcessId, err = k.ProcessID.Peek(ctx)
		if err != nil {
			panic(err)
		}

		iter, err := k.Processing.Iterate(ctx, (&collections.Range[uint64]{}).Descending())
		if err != nil {
			panic(err)
		}
		defer iter.Close()

		for ; iter.Valid(); iter.Next() {
			kv, err := iter.KeyValue()
			if err != nil {
				panic(err)
			}
			genesis.Processing = append(genesis.Processing, types.ProcessingGenesis{
				Id:         kv.Key,
				Processing: kv.Value,
			})
		}
	}

	return genesis
}
