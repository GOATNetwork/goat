package keeper_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	keepertest "github.com/goatnetwork/goat/testutil/keeper"
	"github.com/goatnetwork/goat/x/bitcoin/keeper"
	"github.com/goatnetwork/goat/x/bitcoin/mock"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestKeeper(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite
	Ctrl          *gomock.Controller
	Keeper        keeper.Keeper
	Context       sdk.Context
	RelayerKeeper *mock.MockRelayerKeeper
	Param         types.Params
	TestKey       relayer.PublicKey
}

func (suite *KeeperTestSuite) SetupTest() {
	ctl := gomock.NewController(suite.T())
	relayerKeeper := mock.NewMockRelayerKeeper(ctl)

	keeper, ctx, _ := keepertest.BitcoinKeeper(suite.T(), relayerKeeper)

	suite.Keeper = keeper
	suite.Context = ctx
	suite.Param = types.DefaultParams()
	suite.RelayerKeeper = relayerKeeper

	suite.TestKey = relayer.PublicKey{Key: &relayer.PublicKey_Secp256K1{
		Secp256K1: common.Hex2Bytes("0383560def84048edefe637d0119a4428dd12a42765a118b2bf77984057633c50e"),
	}}

	suite.Ctrl = ctl
}

func (suite *KeeperTestSuite) TearDownSuite() {
	suite.Ctrl.Finish()
}

func (suite *KeeperTestSuite) TestNewPubkey() {
	suite.RelayerKeeper.EXPECT().AddNewKey(suite.Context, relayer.EncodePublicKey(&suite.TestKey)).Return(nil)

	suite.Require().NoError(suite.Keeper.NewPubkey(suite.Context, &suite.TestKey))

	pubkey, err := suite.Keeper.Pubkey.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().True(pubkey.Equal(suite.TestKey))
}

func (suite *KeeperTestSuite) TestVerifyDeposit() {
	suite.RelayerKeeper.EXPECT().HasPubkey(suite.Context, relayer.EncodePublicKey(&suite.TestKey)).Return(true, nil)

	evmAddress := common.HexToAddress("0xbC122aEc3EdD80433dfE3c708b2E549B5A7Ab96E")
	blockHash, _ := chainhash.NewHashFromStr("38fb77a25662f9eda5abef8a407ba45e8c3374b5a0724cfa9762f1f9cbf627e2")
	const height = 102
	suite.Require().NoError(suite.Keeper.BlockHashes.Set(suite.Context, 102, blockHash[:]))

	header := common.Hex2Bytes("00000020451119ce15cd42ceb7a00c2ef9843aa613a69f19f7b4fc483f0f28b099c54d1bc8df397f2235b299f7ca89e10f789e598f53dc89789b8a047bc78238ef4bd4daf9f8e466ffff7f2000000000")

	const txIndex = 1
	const txOutput = 1
	const amount = 1e8
	txid, _ := chainhash.NewHashFromStr("9a31c75d3676059c7483d29f12082b4df9e396df5c22612e50fa97b94bbf532c")
	tx := common.Hex2Bytes("0200000001e15e44fc827b0e1a3178b6e07f67e8339faae54e4241e5fa5c1ed61786a84bda0000000000fdffffff020dc74c0001000000225120098ad136e9ed8106af7c1b6b4934011f320b30f6e18871917e0d6fb1bdcb5d1400e1f50500000000220020f7608234b4bc67678cc5498dfe7db5dfda221d3ff669f1d9ee89fbcf14d104f366000000")

	proof := common.Hex2Bytes("4930ac654c3c2e487fcc2106a51ecaaf4188093686dfffcfe880798044aadc02")

	res, err := suite.Keeper.VerifyDeposit(suite.Context,
		map[uint64][]byte{height: header},
		&types.Deposit{
			Version:           0,
			BlockNumber:       height,
			TxIndex:           txIndex,
			NoWitnessTx:       tx,
			OutputIndex:       txOutput,
			IntermediateProof: proof,
			EvmAddress:        evmAddress.Bytes(),
			RelayerPubkey:     &suite.TestKey,
		})
	suite.Require().NoError(err)

	suite.Require().Equal(res, &types.DepositExecReceipt{
		Address: evmAddress.Bytes(),
		Txid:    txid[:],
		Txout:   txOutput,
		Amount:  amount,
	})
}

func (suite *KeeperTestSuite) TestDequeueBitcoinModuleTx() {
	height := uint64(0)

	blockHashes := [][]byte{
		common.Hex2Bytes("0f9188f13cb7b2c71f2a335e3a4fc328bf5beb436012afca590b1a11466e2206"), // 0
		common.Hex2Bytes("2527aeed8a8f4ab5163ca5b4099e4ebc65cb4085a8462fc664f2115e1937a01a"), // 1
		common.Hex2Bytes("778d0758ef2720140a09d96b59071aa87bbcfadefd3063cd35c32f1a0bdcc537"), // 2
		common.Hex2Bytes("26195e235e742bcf886ff695b050cc7054ca42777f8a206be1547b9a2daf14ab"), // 3
	}

	deposits := []*types.DepositExecReceipt{
		{
			Txid:    common.Hex2Bytes("2f6af73e06798b2caaf4c355e0cdf2b8581667462c2f26e1ee27c3ab07e4a05e"),
			Txout:   0,
			Address: common.Hex2Bytes("dc90965f6ba338ec181e652fef3f2f26804ed823"),
			Amount:  1,
		},
		{
			Txid:    common.Hex2Bytes("91d474a2711d757cd0ea3ebdc2d48d4c81fd32cd10b69fcac0a3ef4bddab51b7"),
			Txout:   1,
			Address: common.Hex2Bytes("1e7966a5c7a99bc1eca451914546ca60c8721cbf"),
			Amount:  1,
		},
		{
			Txid:    common.Hex2Bytes("1df17c545c48730cea429d49b717b69ab800b4f43c1a29d718d39a5c3bd37dc1"),
			Txout:   2,
			Address: common.Hex2Bytes("0d23b5dfe32bfe402b632ac5528ecc888bc74df7"),
			Amount:  1,
		},
		{
			Txid:    common.Hex2Bytes("1378b2024e96e724c6bdd4ab7d94e478d85820156aedd528b2422bf7ad2ac800"),
			Txout:   3,
			Address: common.Hex2Bytes("d82d17523d36fcfbd3617fcd729c55e526233210"),
			Amount:  1,
		},
		{
			Txid:    common.Hex2Bytes("be37ce9d3f9c1384df919326745ef775058c055723039db1fef530f248ebab91"),
			Txout:   4,
			Address: common.Hex2Bytes("bdbec1dedac40cd9df7bb7aa200291bf1eee3da3"),
			Amount:  1,
		},
		{
			Txid:    common.Hex2Bytes("134eca9c59da98c9e7d0fbbf8da36ac54f4d3c794b17a8b98c183714b7752394"),
			Txout:   5,
			Address: common.Hex2Bytes("cf3ffe9bcdeb250694fd8ca57a98e092e4e7bc9d"),
			Amount:  1,
		},
		{
			Txid:    common.Hex2Bytes("934a26463bfb0c839d1b4f8480717cab135ca71d23e27fe91f587c0fbe5b2ea7"),
			Txout:   6,
			Address: common.Hex2Bytes("f9c4420f2d9b248d516ab591f148a15f11edfe6a"),
			Amount:  1,
		},
		{
			Txid:    common.Hex2Bytes("7a5d364885d47b67900b3aad37ef523bcd4447e13a376e256c08ef05547527ee"),
			Txout:   7,
			Address: common.Hex2Bytes("6482daed5fd4395b141418315a26fa8735b0ad11"),
			Amount:  1,
		},
		{
			Txid:    common.Hex2Bytes("44b8877466e9d14860ed8c172f12d82e92ca21fbad89f445b87dd786539651d0"),
			Txout:   8,
			Address: common.Hex2Bytes("f59d04d7c30c0e3a25407dd6e7a0475ef0fb19f6"),
			Amount:  1,
		},
	}

	rejected := []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8}

	paid := []*types.WithdrawalExecReceipt{
		{
			Id: 9,
			Receipt: &types.WithdrawalReceipt{
				Txid:   common.Hex2Bytes("439f4b18923c16fdb2e9e9cc703268e034491e06123740ff4fe249ee210e72e4"),
				Txout:  0,
				Amount: 1,
			},
		},
	}

	for idx, v := range blockHashes {
		if idx != 0 {
			height++
		}
		suite.Require().NoError(
			suite.Keeper.BlockHashes.Set(suite.Context, uint64(idx), v),
		)
	}

	suite.Require().NoError(suite.Keeper.BlockTip.Set(suite.Context, height))

	err := suite.Keeper.EthTxQueue.Set(suite.Context, types.EthTxQueue{
		BlockNumber:         0,
		Deposits:            deposits,
		PaidWithdrawals:     paid,
		RejectedWithdrawals: rejected,
	})
	suite.Require().NoError(err)

	{
		txsGot, err := suite.Keeper.DequeueBitcoinModuleTx(suite.Context)
		suite.Require().NoError(err)

		txsWant := []*ethtypes.Transaction{
			types.NewBitcoinHashEthTx(0, blockHashes[1]),
			deposits[0].EthTx(1),
			deposits[1].EthTx(2),
			deposits[2].EthTx(3),
			deposits[3].EthTx(4),
			deposits[4].EthTx(5),
			deposits[5].EthTx(6),
			deposits[6].EthTx(7),
			deposits[7].EthTx(8),
			paid[0].EthTx(9),
			types.NewRejectEthTx(rejected[0], 10),
			types.NewRejectEthTx(rejected[1], 11),
			types.NewRejectEthTx(rejected[2], 12),
			types.NewRejectEthTx(rejected[3], 13),
			types.NewRejectEthTx(rejected[4], 14),
			types.NewRejectEthTx(rejected[5], 15),
			types.NewRejectEthTx(rejected[6], 16),
		}

		suite.Require().Equal(len(txsGot), len(txsWant))

		for i, tx := range txsGot {
			got, err := tx.MarshalBinary()
			suite.Require().NoError(err)
			want, err := txsWant[i].MarshalBinary()
			suite.Require().NoError(err)
			suite.Require().Equal(got, want)
		}

		newNonce, err := suite.Keeper.EthTxNonce.Peek(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(newNonce, uint64(17))

		newQueue, err := suite.Keeper.EthTxQueue.Get(suite.Context)
		suite.Require().NoError(err)

		suite.Require().Equal(newQueue, types.EthTxQueue{
			BlockNumber:         1,
			Deposits:            deposits[8:],
			PaidWithdrawals:     nil,
			RejectedWithdrawals: rejected[7:],
		})
	}

	{
		txsGot, err := suite.Keeper.DequeueBitcoinModuleTx(suite.Context)
		suite.Require().NoError(err)

		txsWant := []*ethtypes.Transaction{
			types.NewBitcoinHashEthTx(17, blockHashes[2]),
			deposits[8].EthTx(18),
			types.NewRejectEthTx(rejected[7], 19),
			types.NewRejectEthTx(rejected[8], 20),
		}
		suite.Require().Equal(len(txsGot), len(txsWant))

		for i := 0; i < len(txsGot); i++ {
			got, err := txsGot[i].MarshalBinary()
			suite.Require().NoError(err)
			want, err := txsWant[i].MarshalBinary()
			suite.Require().NoError(err)
			suite.Require().Equal(got, want)
		}

		newNonce, err := suite.Keeper.EthTxNonce.Peek(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(newNonce, uint64(21))

		newQueue, err := suite.Keeper.EthTxQueue.Get(suite.Context)
		suite.Require().NoError(err)

		suite.Require().Equal(newQueue, types.EthTxQueue{
			BlockNumber:         2,
			Deposits:            nil,
			PaidWithdrawals:     nil,
			RejectedWithdrawals: nil,
		})
	}

	{
		txsGot, err := suite.Keeper.DequeueBitcoinModuleTx(suite.Context)
		suite.Require().NoError(err)

		txsWant := []*ethtypes.Transaction{
			types.NewBitcoinHashEthTx(21, blockHashes[3]),
		}
		suite.Require().Equal(len(txsGot), len(txsWant))

		for i := 0; i < len(txsGot); i++ {
			got, err := txsGot[i].MarshalBinary()
			suite.Require().NoError(err)
			want, err := txsWant[i].MarshalBinary()
			suite.Require().NoError(err)
			suite.Require().Equal(got, want)
		}

		newNonce, err := suite.Keeper.EthTxNonce.Peek(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(newNonce, uint64(22))

		newQueue, err := suite.Keeper.EthTxQueue.Get(suite.Context)
		suite.Require().NoError(err)

		suite.Require().Equal(newQueue, types.EthTxQueue{
			BlockNumber:         3,
			Deposits:            nil,
			PaidWithdrawals:     nil,
			RejectedWithdrawals: nil,
		})
	}

	paid2 := []*types.WithdrawalExecReceipt{
		{
			Id: 9,
			Receipt: &types.WithdrawalReceipt{
				Txid:   common.Hex2Bytes("439f4b18923c16fdb2e9e9cc703268e034491e06123740ff4fe249ee210e72e4"),
				Txout:  0,
				Amount: 1,
			},
		},
		{
			Id: 10,
			Receipt: &types.WithdrawalReceipt{
				Txid:   common.Hex2Bytes("a581e6b25fd3b06d36afbc4a59a9466fb8357ccc6c690d61e32fa26a26ece88f"),
				Txout:  0,
				Amount: 1,
			},
		},
		{
			Id: 11,
			Receipt: &types.WithdrawalReceipt{
				Txid:   common.Hex2Bytes("6593c3e2836908ffbe9fa27238629dcf609baeef1c9a3521c1522aa56c163b37"),
				Txout:  0,
				Amount: 1,
			},
		},
	}

	rejected2 := []uint64{11, 12}

	{
		txsGot, err := suite.Keeper.DequeueBitcoinModuleTx(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(len(txsGot), 0)

		newNonce, err := suite.Keeper.EthTxNonce.Peek(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(newNonce, uint64(22))

		newQueue, err := suite.Keeper.EthTxQueue.Get(suite.Context)
		suite.Require().NoError(err)

		suite.Require().Equal(newQueue, types.EthTxQueue{
			BlockNumber:         3,
			Deposits:            nil,
			PaidWithdrawals:     nil,
			RejectedWithdrawals: nil,
		})

		newQueue.PaidWithdrawals = paid2
		newQueue.RejectedWithdrawals = rejected2
		suite.Require().NoError(suite.Keeper.EthTxQueue.Set(suite.Context, newQueue))
	}

	{
		txsGot, err := suite.Keeper.DequeueBitcoinModuleTx(suite.Context)
		suite.Require().NoError(err)

		txsWant := []*ethtypes.Transaction{
			paid2[0].EthTx(22),
			paid2[1].EthTx(23),
			paid2[2].EthTx(24),
			types.NewRejectEthTx(rejected2[0], 25),
			types.NewRejectEthTx(rejected2[1], 26),
		}

		suite.Require().Equal(len(txsGot), len(txsWant))

		for i := 0; i < len(txsGot); i++ {
			got, err := txsGot[i].MarshalBinary()
			suite.Require().NoError(err)
			want, err := txsWant[i].MarshalBinary()
			suite.Require().NoError(err)
			suite.Require().Equal(got, want)
		}

		newNonce, err := suite.Keeper.EthTxNonce.Peek(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Equal(newNonce, uint64(27))

		newQueue, err := suite.Keeper.EthTxQueue.Get(suite.Context)
		suite.Require().NoError(err)

		suite.Require().Equal(newQueue, types.EthTxQueue{
			BlockNumber:         3,
			Deposits:            nil,
			PaidWithdrawals:     nil,
			RejectedWithdrawals: nil,
		})
	}
}

func (suite *KeeperTestSuite) TestProcessBridgeRequest() {
	withdrawals1 := []*ethtypes.GoatWithdrawal{
		{
			Id:         0,
			Amount:     1,
			MaxTxPrice: 1,
			Address:    "bc1qaplz7zlds5hwkr9wz62hrmdvlrljly0c2uzye09rnht8venzy7rsz245xr",
		},
		{
			Id:         1,
			Amount:     1,
			MaxTxPrice: 1,
			Address:    "17yhJ5DME9Fu3wVjVoVfP4jKxjrc9WRyaB",
		},
		{
			Id:         2,
			Amount:     1,
			MaxTxPrice: 1,
			Address:    "bcrt1q268jv3p5gcs8a0xf2pgty8lv9a87ufy38nxwclt88txf4ptzzzaqwa2hp7",
		},
	}

	suite.Require().NoError(suite.Keeper.EthTxQueue.Set(suite.Context, types.EthTxQueue{}))
	err := suite.Keeper.ProcessBridgeRequest(suite.Context, types.ExecRequests{Withdrawals: withdrawals1})
	suite.Require().NoError(err)
	suite.Require().Equal(len(suite.Context.EventManager().Events()), 1)

	has0, err := suite.Keeper.Withdrawals.Has(suite.Context, 0)
	suite.Require().NoError(err)
	suite.Require().False(has0)
	has1, err := suite.Keeper.Withdrawals.Has(suite.Context, 1)
	suite.Require().NoError(err)
	suite.Require().False(has1)

	queue, err := suite.Keeper.EthTxQueue.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(queue, types.EthTxQueue{
		RejectedWithdrawals: []uint64{0, 1},
	})

	wd2, err := suite.Keeper.Withdrawals.Get(suite.Context, 2)
	suite.Require().NoError(err)
	suite.Require().Equal(wd2, types.Withdrawal{
		Address:       withdrawals1[2].Address,
		RequestAmount: withdrawals1[2].Amount,
		MaxTxPrice:    withdrawals1[2].MaxTxPrice,
		OutputScript:  common.Hex2Bytes("0020568f26443446207ebcc95050b21fec2f4fee24913cccec7d673acc9a856210ba"),
		Status:        types.WITHDRAWAL_STATUS_PENDING,
	})

	withdrawals2 := []*ethtypes.GoatWithdrawal{
		{
			Id:         3,
			Amount:     1,
			MaxTxPrice: 1,
			Address:    "bcrt1q5h7tn8l6euv670gzjp7s9nlcadmcmttprrylyj",
		},
	}

	rbf1 := []*ethtypes.ReplaceByFee{
		{Id: 2, MaxTxPrice: 2},
	}

	err = suite.Keeper.ProcessBridgeRequest(suite.Context, types.ExecRequests{Withdrawals: withdrawals2, RBFs: rbf1})
	suite.Require().NoError(err)

	wd2, err = suite.Keeper.Withdrawals.Get(suite.Context, 2)
	suite.Require().NoError(err)
	suite.Require().Equal(wd2, types.Withdrawal{
		Address:       withdrawals1[2].Address,
		RequestAmount: withdrawals1[2].Amount,
		MaxTxPrice:    rbf1[0].MaxTxPrice,
		OutputScript:  common.Hex2Bytes("0020568f26443446207ebcc95050b21fec2f4fee24913cccec7d673acc9a856210ba"),
		Status:        types.WITHDRAWAL_STATUS_PENDING,
	})

	wd3, err := suite.Keeper.Withdrawals.Get(suite.Context, 3)
	suite.Require().NoError(err)
	suite.Require().Equal(wd3, types.Withdrawal{
		Address:       withdrawals2[0].Address,
		RequestAmount: withdrawals2[0].Amount,
		MaxTxPrice:    withdrawals2[0].MaxTxPrice,
		OutputScript:  common.Hex2Bytes("0014a5fcb99ffacf19af3d02907d02cff8eb778dad61"),
		Status:        types.WITHDRAWAL_STATUS_PENDING,
	})

	wd2.Status = types.WITHDRAWAL_STATUS_PROCESSING
	suite.Require().NoError(suite.Keeper.Withdrawals.Set(suite.Context, 2, wd2))

	rbf2 := []*ethtypes.ReplaceByFee{
		{Id: 2, MaxTxPrice: 3},
	}

	err = suite.Keeper.ProcessBridgeRequest(suite.Context, types.ExecRequests{RBFs: rbf2})
	suite.Require().NoError(err)

	wd2, err = suite.Keeper.Withdrawals.Get(suite.Context, 2)
	suite.Require().NoError(err)
	suite.Require().Equal(wd2, types.Withdrawal{
		Address:       withdrawals1[2].Address,
		RequestAmount: withdrawals1[2].Amount,
		MaxTxPrice:    rbf1[0].MaxTxPrice,
		OutputScript:  common.Hex2Bytes("0020568f26443446207ebcc95050b21fec2f4fee24913cccec7d673acc9a856210ba"),
		Status:        types.WITHDRAWAL_STATUS_PROCESSING,
	})

	cancel1 := []*ethtypes.Cancel1{
		{Id: 2},
		{Id: 3},
	}

	err = suite.Keeper.ProcessBridgeRequest(suite.Context, types.ExecRequests{Cancel1s: cancel1})
	suite.Require().NoError(err)

	wd2, err = suite.Keeper.Withdrawals.Get(suite.Context, 2)
	suite.Require().NoError(err)
	suite.Require().Equal(wd2, types.Withdrawal{
		Address:       withdrawals1[2].Address,
		RequestAmount: withdrawals1[2].Amount,
		MaxTxPrice:    rbf1[0].MaxTxPrice,
		OutputScript:  common.Hex2Bytes("0020568f26443446207ebcc95050b21fec2f4fee24913cccec7d673acc9a856210ba"),
		Status:        types.WITHDRAWAL_STATUS_PROCESSING,
	})

	wd3, err = suite.Keeper.Withdrawals.Get(suite.Context, 3)
	suite.Require().NoError(err)
	suite.Require().Equal(wd3, types.Withdrawal{
		Address:       withdrawals2[0].Address,
		RequestAmount: withdrawals2[0].Amount,
		MaxTxPrice:    withdrawals2[0].MaxTxPrice,
		OutputScript:  common.Hex2Bytes("0014a5fcb99ffacf19af3d02907d02cff8eb778dad61"),
		Status:        types.WITHDRAWAL_STATUS_CANCELING,
	})
}
