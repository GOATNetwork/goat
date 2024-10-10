package keeper_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
)

func (suite *KeeperTestSuite) TestDequeue() {
	txs := []*ethtypes.Transaction{
		types.NewTx(ethtypes.NewGoatTx(
			goattypes.BirdgeModule,
			goattypes.BridgeDepoitAction,
			0,
			&goattypes.DepositTx{
				Txid:   common.HexToHash("0x344fb824c793fc370a38577eea12aba8842cb0516cf52099911a36c0c36f11ee"),
				TxOut:  0,
				Target: common.HexToAddress("0xe896f4afff6c2424819aa493b1724fc11851dc54"),
				Amount: big.NewInt(10),
			},
		)),
		types.NewTx(
			types.NewGoatTx(
				goattypes.LockingModule,
				goattypes.LockingDistributeRewardAction,
				0,
				&goattypes.DistributeRewardTx{
					Id:        0,
					Recipient: common.HexToAddress("0x8743c8f569103715dce9d16394185a8c8dc721ec"),
					Goat:      big.NewInt(10),
					GasReward: big.NewInt(1),
				},
			),
		),
	}

	bytes := make([][]byte, len(txs))
	for i := 0; i < len(txs); i++ {
		tx, err := txs[i].MarshalBinary()
		suite.Require().NoError(err)
		bytes[i] = tx
	}

	suite.Bitcoin.EXPECT().DequeueBitcoinModuleTx(suite.Context).Return(txs[:1], nil)
	suite.Locking.EXPECT().DequeueLockingModuleTx(suite.Context).Return(txs[1:], nil)

	res, err := suite.Keeper.Dequeue(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(res, bytes)
}

func (suite *KeeperTestSuite) TestVerifyDequeue() {
	txs := []*ethtypes.Transaction{
		types.NewTx(ethtypes.NewGoatTx(
			goattypes.BirdgeModule,
			goattypes.BridgeDepoitAction,
			0,
			&goattypes.DepositTx{
				Txid:   common.HexToHash("0x344fb824c793fc370a38577eea12aba8842cb0516cf52099911a36c0c36f11ee"),
				TxOut:  0,
				Target: common.HexToAddress("0xe896f4afff6c2424819aa493b1724fc11851dc54"),
				Amount: big.NewInt(10),
			},
		)),
		types.NewTx(
			types.NewGoatTx(
				goattypes.LockingModule,
				goattypes.LockingDistributeRewardAction,
				0,
				&goattypes.DistributeRewardTx{
					Id:        0,
					Recipient: common.HexToAddress("0x8743c8f569103715dce9d16394185a8c8dc721ec"),
					Goat:      big.NewInt(10),
					GasReward: big.NewInt(1),
				},
			),
		),
	}

	bytes := make([][]byte, len(txs))
	for i := 0; i < len(txs); i++ {
		tx, err := txs[i].MarshalBinary()
		suite.Require().NoError(err)
		bytes[i] = tx
	}

	suite.Bitcoin.EXPECT().DequeueBitcoinModuleTx(suite.Context).Return(txs[:1], nil)
	suite.Locking.EXPECT().DequeueLockingModuleTx(suite.Context).Return(txs[1:], nil)

	err := suite.Keeper.VerifyDequeue(suite.Context, common.Hex2Bytes("0290d76198d64802e418eafb8d6b011bae8bd790af09a49fc80b835aeed8924e48"), bytes)
	suite.Require().NoError(err)
}
