package keeper_test

import (
	"crypto/rand"
	"math/big"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (suite *KeeperTestSuite) TestDequeueLockingModuleTx() {
	suite.Run("no tx", func() {
		err := suite.Keeper.EthTxQueue.Set(suite.Context, types.EthTxQueue{})
		suite.Require().NoError(err)

		txs, err := suite.Keeper.DequeueLockingModuleTx(suite.Context)
		suite.Require().NoError(err)
		var res []*ethtypes.Transaction
		suite.Require().EqualValues(txs, res)
	})

	suite.Run("full", func() {
		queue := types.EthTxQueue{}
		txs := []*ethtypes.Transaction{}

		nonce := uint64(0)
		for i := uint64(0); i < 16; i++ {
			recipent := make([]byte, 20)
			_, err := rand.Read(recipent)
			suite.Require().NoError(err)

			queue.Rewards = append(queue.Rewards, &types.Reward{
				Id:        i,
				Recipient: recipent,
				Goat:      math.NewIntFromUint64(i),
				Gas:       math.NewIntFromUint64(i),
			})

			txs = append(txs, ethtypes.NewTx(ethtypes.NewGoatTx(
				goattypes.LockingModule,
				goattypes.LockingDistributeRewardAction,
				nonce,
				&goattypes.DistributeRewardTx{
					Id:        i,
					Recipient: common.BytesToAddress(recipent),
					Goat:      new(big.Int).SetUint64(i),
					GasReward: new(big.Int).SetUint64(i),
				},
			)))
			nonce++
		}

		for i := uint64(16); i < 36; i++ {
			recipent := make([]byte, 40)
			_, err := rand.Read(recipent)
			suite.Require().NoError(err)
			queue.Unlocks = append(queue.Unlocks, &types.Unlock{
				Id:        i,
				Recipient: recipent[:20],
				Token:     recipent[20:40],
				Amount:    math.NewIntFromUint64(i),
			})

			txs = append(txs, ethtypes.NewTx(ethtypes.NewGoatTx(
				goattypes.LockingModule,
				goattypes.LockingCompleteUnlockAction,
				nonce,
				&goattypes.CompleteUnlockTx{
					Id:        i,
					Recipient: common.BytesToAddress(recipent[:20]),
					Token:     common.BytesToAddress(recipent[20:]),
					Amount:    new(big.Int).SetUint64(i),
				},
			)))
			nonce++
		}

		err := suite.Keeper.EthTxQueue.Set(suite.Context, queue)
		suite.Require().NoError(err)

		got1, err := suite.Keeper.DequeueLockingModuleTx(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Len(got1, 32)
		for i := 0; i < 32; i++ {
			got, err := got1[i].MarshalBinary()
			suite.Require().NoError(err)
			want, err := txs[i].MarshalBinary()
			suite.Require().NoError(err)
			suite.Require().Equal(got, want, i)
		}
		seq1, err := suite.Keeper.EthTxNonce.Peek(suite.Context)
		suite.Require().NoError(err)
		suite.Require().EqualValues(seq1, 32)
		queue1, err := suite.Keeper.EthTxQueue.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Len(queue1.Rewards, 0)
		suite.Require().Len(queue1.Unlocks, 4)
		suite.Require().Equal(queue1.Unlocks, queue.Unlocks[16:])

		got2, err := suite.Keeper.DequeueLockingModuleTx(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Len(got2, 4)
		for i := 0; i < 4; i++ {
			got, err := got2[i].MarshalBinary()
			suite.Require().NoError(err)
			want, err := txs[i+32].MarshalBinary()
			suite.Require().NoError(err)
			suite.Require().Equal(got, want, i+32)
		}
		seq2, err := suite.Keeper.EthTxNonce.Peek(suite.Context)
		suite.Require().NoError(err)
		suite.Require().EqualValues(seq2, 36)

		queue2, err := suite.Keeper.EthTxQueue.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Len(queue2.Rewards, 0)
		suite.Require().Len(queue2.Unlocks, 0)
	})
}
