package keeper_test

import (
	"encoding/hex"
	"encoding/json"

	"cosmossdk.io/collections"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/testutil"
	"github.com/goatnetwork/goat/x/bitcoin/keeper"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
)

func (suite *KeeperTestSuite) TestMsgNewDeposits() {
	evmAddress := common.HexToAddress("0xbC122aEc3EdD80433dfE3c708b2E549B5A7Ab96E")
	blockHash, err := chainhash.NewHashFromStr("38fb77a25662f9eda5abef8a407ba45e8c3374b5a0724cfa9762f1f9cbf627e2")
	suite.Require().NoError(err)

	const height = 102
	suite.Require().NoError(suite.Keeper.BlockHashes.Set(suite.Context, 102, blockHash[:]))

	header, err := hex.DecodeString("00000020451119ce15cd42ceb7a00c2ef9843aa613a69f19f7b4fc483f0f28b099c54d1bc8df397f2235b299f7ca89e10f789e598f53dc89789b8a047bc78238ef4bd4daf9f8e466ffff7f2000000000")
	suite.Require().NoError(err)

	const txIndex = 1
	const txOutput = 1
	const amount = 1e8
	txid, err := chainhash.NewHashFromStr("9a31c75d3676059c7483d29f12082b4df9e396df5c22612e50fa97b94bbf532c")
	suite.Require().NoError(err)

	tx, err := hex.DecodeString("0200000001e15e44fc827b0e1a3178b6e07f67e8339faae54e4241e5fa5c1ed61786a84bda0000000000fdffffff020dc74c0001000000225120098ad136e9ed8106af7c1b6b4934011f320b30f6e18871917e0d6fb1bdcb5d1400e1f50500000000220020f7608234b4bc67678cc5498dfe7db5dfda221d3ff669f1d9ee89fbcf14d104f366000000")
	suite.Require().NoError(err)

	proof := common.Hex2Bytes("4930ac654c3c2e487fcc2106a51ecaaf4188093686dfffcfe880798044aadc02")

	headers, err := json.Marshal(map[uint64][]byte{height: header})
	suite.Require().NoError(err)

	err = suite.Keeper.EthTxQueue.Set(suite.Context, types.EthTxQueue{})
	suite.Require().NoError(err)

	msgServer := keeper.NewMsgServerImpl(suite.Keeper)

	req := &types.MsgNewDeposits{
		Proposer:     "goat1xa56637tjn857jyg2plgvhdclzmr4crxzn5xus",
		BlockHeaders: headers,
		Deposits: []*types.Deposit{
			{
				Version:           0,
				BlockNumber:       height,
				TxIndex:           txIndex,
				NoWitnessTx:       tx,
				OutputIndex:       txOutput,
				IntermediateProof: proof,
				EvmAddress:        evmAddress[:],
				RelayerPubkey:     &suite.TestKey,
			},
		},
	}
	suite.RelayerKeeper.EXPECT().HasPubkey(suite.Context, relayertypes.EncodePublicKey(&suite.TestKey)).Return(true, nil)
	suite.RelayerKeeper.EXPECT().VerifyNonProposal(suite.Context, req).Return(nil, nil)
	_, err = msgServer.NewDeposits(suite.Context, req)
	suite.Require().NoError(err)

	queue, err := suite.Keeper.EthTxQueue.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Equal(queue, types.EthTxQueue{
		Deposits: []*types.DepositExecReceipt{
			{Txid: txid[:], Txout: txOutput, Address: evmAddress[:], Amount: amount},
		},
	})

	exist, err := suite.Keeper.Deposited.Has(suite.Context, collections.Join(txid[:], uint32(txOutput)))
	suite.Require().NoError(err)
	suite.Require().True(exist)
}

func (suite *KeeperTestSuite) TestMsgNewBlockHashes() {
	const parentHeight = 100
	err := suite.Keeper.BlockTip.Set(suite.Context, parentHeight)
	suite.Require().NoError(err)

	req := &types.MsgNewBlockHashes{
		Proposer:         "goat1xa56637tjn857jyg2plgvhdclzmr4crxzn5xus",
		Vote:             &relayertypes.Votes{Signature: make([]byte, goatcrypto.SignatureLength)},
		StartBlockNumber: 101,
		BlockHash: [][]byte{
			common.Hex2Bytes("9973a71cbd41470951a10bb09236ce5f535edbccecf4c3593ca4755f3c4fe7d3"),
			common.Hex2Bytes("dea491aae4c5cc43a7b5c39559580adde5ef10995bdd2d53e722a1674d78806f"),
		},
	}

	sequence := uint64(100)
	suite.RelayerKeeper.EXPECT().VerifyProposal(suite.Context, req).Return(sequence, nil)
	suite.RelayerKeeper.EXPECT().SetProposalSeq(suite.Context, sequence+1)
	suite.RelayerKeeper.EXPECT().UpdateRandao(suite.Context, req)

	msgServer := keeper.NewMsgServerImpl(suite.Keeper)

	_, err = msgServer.NewBlockHashes(suite.Context, req)
	suite.Require().NoError(err)

	bhash, err := suite.Keeper.BlockHashes.Get(suite.Context, req.StartBlockNumber)
	suite.Require().NoError(err)
	suite.Require().Equal(bhash, req.BlockHash[0])

	bhash, err = suite.Keeper.BlockHashes.Get(suite.Context, req.StartBlockNumber+1)
	suite.Require().NoError(err)
	suite.Require().Equal(bhash, req.BlockHash[1])

	height, err := suite.Keeper.BlockTip.Peek(suite.Context)
	suite.Require().NoError(err)
	suite.Require().EqualValues(height, parentHeight+2)
}

func (suite *KeeperTestSuite) TestMsgNewPubkey() {
	req := &types.MsgNewPubkey{
		Proposer: "goat1xa56637tjn857jyg2plgvhdclzmr4crxzn5xus",
		Vote:     &relayertypes.Votes{Signature: make([]byte, goatcrypto.SignatureLength)},
		Pubkey: &relayertypes.PublicKey{Key: &relayertypes.PublicKey_Secp256K1{
			Secp256K1: common.Hex2Bytes("03a466deef30f68c03ad54f89f9deb8284f0529aad1c095985a015be27daec20c6"),
		}},
	}

	sequence := uint64(100)

	rawKey := relayertypes.EncodePublicKey(req.Pubkey)
	suite.RelayerKeeper.EXPECT().HasPubkey(suite.Context, rawKey).Return(false, nil)
	suite.RelayerKeeper.EXPECT().VerifyProposal(suite.Context, req).Times(2).Return(sequence, nil)
	suite.RelayerKeeper.EXPECT().AddNewKey(suite.Context, rawKey)
	suite.RelayerKeeper.EXPECT().SetProposalSeq(suite.Context, sequence+1)
	suite.RelayerKeeper.EXPECT().UpdateRandao(suite.Context, req)
	suite.RelayerKeeper.EXPECT().HasPubkey(suite.Context, rawKey).Return(true, nil)

	msgServer := keeper.NewMsgServerImpl(suite.Keeper)
	_, err := msgServer.NewPubkey(suite.Context, req)
	suite.Require().NoError(err)

	_, err = msgServer.NewPubkey(suite.Context, req)
	suite.Require().ErrorContains(err, "the key already existed")
}

func (suite *KeeperTestSuite) TestMsgWithdrawal() {
	withdrawals := []string{
		"mqEgATpvdzTNpdpCRAomE9nCH8cw7Sp4R3",
		"2Mu3CdYeLs1Jb7ywDwfgtE58u5z1yseD2KN",
		"bcrt1qy728d54p6ftlpnwvfpjpkdne6sg3saq4qzezpx",
		"bcrt1q8kk55x5qcwf97a9y3apfam9yx3q2mc9wu02ak6yquua5wjtpvvwsggh3r2",
		"bcrt1qpjzz885yglnu8dr3kfnnvc27kakyw4w9h9z3vdmftrchurdyvj8srq7d9m",
	}

	msgServer := keeper.NewMsgServerImpl(suite.Keeper)

	fristTxid, err := chainhash.NewHashFromStr("884879222670f204932d5dfc275bfb11bdc38a93b884f6c54501ff9e721de411")
	suite.Require().NoError(err)
	secondTxid, err := chainhash.NewHashFromStr("445ede02feadfca4ff92ad46b640b2bd6d8f60c14c09efe7659ef7c5e6c26192")
	suite.Require().NoError(err)

	expected, events := []types.Withdrawal{}, []sdktypes.Event{}
	// process it
	{
		req := &types.MsgProcessWithdrawal{
			Proposer:    "goat1xa56637tjn857jyg2plgvhdclzmr4crxzn5xus",
			Vote:        &relayertypes.Votes{Signature: make([]byte, goatcrypto.SignatureLength)},
			NoWitnessTx: common.Hex2Bytes("02000000012e0e3e521ac999cfc292a78aaeb31fe19dfb7867c660ae5560537370d55fdf0e0000000000ffffffff06e8030000000000001976a9146a9d23174484d7ba74f7bc2a64ed102b4846267588ace80300000000000017a91413aa207651e0f3724cbe6134f54675aa2d5cdbf987e803000000000000160014279476d2a1d257f0cdcc48641b3679d411187415e8030000000000002200203dad4a1a80c3925f74a48f429eeca43440ade0aee3d5db6880e73b474961631de8030000000000002200200c84239e8447e7c3b471b26736615eb76c4755c5b94516376958f17e0da4648f90c9f505000000001600145b029559baaea5e928e8e2774e9e2350a5fc9c2d00000000"),
			TxFee:       1000,
		}

		for idx, address := range withdrawals {
			status := types.WITHDRAWAL_STATUS_PENDING
			if idx%2 == 0 {
				status = types.WITHDRAWAL_STATUS_CANCELING
			}
			wd := types.Withdrawal{
				Address:       address,
				RequestAmount: 1e3,
				MaxTxPrice:    50,
				Status:        status,
			}
			err = suite.Keeper.Withdrawals.Set(suite.Context, uint64(idx), wd)
			req.Id = append(req.Id, uint64(idx))
			suite.Require().NoError(err)
			expected = append(expected, wd)
		}

		err := suite.Keeper.Pubkey.Set(suite.Context, relayertypes.PublicKey{Key: &relayertypes.PublicKey_Secp256K1{
			Secp256K1: common.Hex2Bytes("037e7bee29c1956152e308d1310823295d720b4cef9e1118726eb1705ffc5a4701"),
		}})
		suite.Require().NoError(err)

		sequence := uint64(100)
		suite.RelayerKeeper.EXPECT().VerifyProposal(suite.Context, req).Return(sequence, nil)
		suite.RelayerKeeper.EXPECT().SetProposalSeq(suite.Context, sequence+1)
		suite.RelayerKeeper.EXPECT().UpdateRandao(suite.Context, req)

		_, err = msgServer.ProcessWithdrawal(suite.Context, req)
		suite.Require().NoError(err)

		for idx, wd := range expected {
			wd.Status = types.WITHDRAWAL_STATUS_PROCESSING
			wd.Receipt = &types.WithdrawalReceipt{Txid: fristTxid[:], Txout: uint32(idx), Amount: 1e3}
			withdrawal, err := suite.Keeper.Withdrawals.Get(suite.Context, uint64(idx))
			suite.Require().NoError(err)
			suite.Require().Equal(withdrawal, wd)
		}

		processing, err := suite.Keeper.Processing.Get(suite.Context, 0)
		suite.Require().NoError(err)
		suite.Require().Equal(processing, types.Processing{
			Txid:        [][]byte{fristTxid.CloneBytes()},
			Output:      []types.TxOuptut{{Values: []uint64{1e3, 1e3, 1e3, 1e3, 1e3}}},
			Withdrawals: []uint64{0, 1, 2, 3, 4},
			Fee:         1000,
		})
		latestPid, err := suite.Keeper.ProcessID.Peek(suite.Context)
		suite.Require().NoError(err)
		suite.Require().EqualValues(latestPid, 1)
		suite.Require().Len(suite.Context.EventManager().ABCIEvents(), 2)
		events = append(events, sdktypes.NewEvent(types.EventTypeWithdrawalProcessing,
			sdktypes.NewAttribute("pid", "0"), sdktypes.NewAttribute("txid", fristTxid.String())))
		events = append(events, sdktypes.NewEvent(relayertypes.EventFinalizedProposal, sdktypes.NewAttribute("sequence", "100")))
	}

	// replace it
	{
		req := &types.MsgReplaceWithdrawal{
			Proposer:       "goat1xa56637tjn857jyg2plgvhdclzmr4crxzn5xus",
			Vote:           &relayertypes.Votes{Signature: make([]byte, goatcrypto.SignatureLength)},
			Pid:            0,
			NewNoWitnessTx: common.Hex2Bytes("02000000012e0e3e521ac999cfc292a78aaeb31fe19dfb7867c660ae5560537370d55fdf0e0000000000ffffffff06b0030000000000001976a9146a9d23174484d7ba74f7bc2a64ed102b4846267588aca00300000000000017a91413aa207651e0f3724cbe6134f54675aa2d5cdbf987e303000000000000160014279476d2a1d257f0cdcc48641b3679d41118741599030000000000002200203dad4a1a80c3925f74a48f429eeca43440ade0aee3d5db6880e73b474961631de3030000000000002200200c84239e8447e7c3b471b26736615eb76c4755c5b94516376958f17e0da4648f10270000000000001600145b029559baaea5e928e8e2774e9e2350a5fc9c2d00000000"),
			NewTxFee:       1100,
		}

		sequence := uint64(101)
		suite.RelayerKeeper.EXPECT().VerifyProposal(suite.Context, req).Return(sequence, nil)
		suite.RelayerKeeper.EXPECT().SetProposalSeq(suite.Context, sequence+1)
		suite.RelayerKeeper.EXPECT().UpdateRandao(suite.Context, req)

		_, err = msgServer.ReplaceWithdrawal(suite.Context, req)
		suite.Require().NoError(err)

		processing, err := suite.Keeper.Processing.Get(suite.Context, 0)
		suite.Require().NoError(err)

		values := []uint64{944, 928, 995, 921, 995}
		suite.Require().Equal(processing, types.Processing{
			Txid:        [][]byte{fristTxid.CloneBytes(), secondTxid.CloneBytes()},
			Output:      []types.TxOuptut{{Values: []uint64{1e3, 1e3, 1e3, 1e3, 1e3}}, {Values: values}},
			Withdrawals: []uint64{0, 1, 2, 3, 4},
			Fee:         1100,
		})

		for idx, wd := range expected {
			wd.Status = types.WITHDRAWAL_STATUS_PROCESSING
			wd.Receipt = &types.WithdrawalReceipt{Txid: secondTxid[:], Txout: uint32(idx), Amount: values[idx]}
			withdrawal, err := suite.Keeper.Withdrawals.Get(suite.Context, uint64(idx))
			suite.Require().NoError(err)
			suite.Require().Equal(withdrawal, wd)
		}
		suite.Require().Len(suite.Context.EventManager().ABCIEvents(), 4)

		events = append(events, sdktypes.NewEvent(types.EventTypeWithdrawalRelayerReplace,
			sdktypes.NewAttribute("pid", "0"), sdktypes.NewAttribute("txid", secondTxid.String())))
		events = append(events, sdktypes.NewEvent(relayertypes.EventFinalizedProposal, sdktypes.NewAttribute("sequence", "101")))
	}

	// finalize it but the first tx is confirmed
	{
		err = suite.Keeper.EthTxQueue.Set(suite.Context, types.EthTxQueue{})
		suite.Require().NoError(err)

		const height = 200

		blockHash, err := chainhash.NewHashFromStr("72e60f220f81fa7a154a686b93732f727e805eacd22abf69bc02349b93ab748c")
		suite.Require().NoError(err)

		suite.Require().NoError(suite.Keeper.BlockHashes.Set(suite.Context, height, blockHash[:]))

		header, err := hex.DecodeString("00000020a4ed2c96b81609cb808079d7138e45fb02fa4f7f5e411ad5af345bd85a925a0a36438a7ba17e98e54c0abdc5ae80af0a7159be79dc877975bffe4f7a08fc53023f480a67ffff7f2007000000")
		suite.Require().NoError(err)

		proof, err := hex.DecodeString("68540c8ef2ca5d5c2bc271b869767b25471c47b13f50f8d393c081ed8d646e1d")
		suite.Require().NoError(err)

		req := &types.MsgFinalizeWithdrawal{
			Proposer:          "goat1xa56637tjn857jyg2plgvhdclzmr4crxzn5xus",
			Pid:               0,
			Txid:              fristTxid.CloneBytes(),
			BlockNumber:       height,
			TxIndex:           1,
			IntermediateProof: proof,
			BlockHeader:       header,
		}
		suite.RelayerKeeper.EXPECT().VerifyNonProposal(suite.Context, req).Return(nil, nil)
		_, err = msgServer.FinalizeWithdrawal(suite.Context, req)
		suite.Require().NoError(err)

		queue, err := suite.Keeper.EthTxQueue.Get(suite.Context)
		suite.Require().NoError(err)
		suite.Require().Len(queue.PaidWithdrawals, len(expected))

		for idx, wd := range expected {
			wd.Status = types.WITHDRAWAL_STATUS_PAID
			wd.Receipt = &types.WithdrawalReceipt{Txid: fristTxid.CloneBytes(), Txout: uint32(idx), Amount: 1e3}
			withdrawal, err := suite.Keeper.Withdrawals.Get(suite.Context, uint64(idx))
			suite.Require().NoError(err)
			suite.Require().Equal(withdrawal, wd)
			suite.Require().Equal(&types.WithdrawalExecReceipt{
				Id:      uint64(idx),
				Receipt: wd.Receipt,
			}, queue.PaidWithdrawals[idx])
		}

		noProcessing, err := suite.Keeper.Processing.Has(suite.Context, 0)
		suite.Require().NoError(err)
		suite.Require().False(noProcessing)
		suite.Require().Len(suite.Context.EventManager().ABCIEvents(), 5)
		events = append(events, sdktypes.NewEvent(types.EventTypeWithdrawalFinalized,
			sdktypes.NewAttribute("pid", "0"),
			sdktypes.NewAttribute("txid", fristTxid.String())))
	}

	testutil.EventEquals(suite.T(), events, suite.Context.EventManager().Events())
}

func (suite *KeeperTestSuite) TestMsgApproveCancellation() {
	withdrawals := []string{
		"mqEgATpvdzTNpdpCRAomE9nCH8cw7Sp4R3",
		"2Mu3CdYeLs1Jb7ywDwfgtE58u5z1yseD2KN",
		"bcrt1qy728d54p6ftlpnwvfpjpkdne6sg3saq4qzezpx",
		"bcrt1q8kk55x5qcwf97a9y3apfam9yx3q2mc9wu02ak6yquua5wjtpvvwsggh3r2",
		"bcrt1qpjzz885yglnu8dr3kfnnvc27kakyw4w9h9z3vdmftrchurdyvj8srq7d9m",
	}

	err := suite.Keeper.EthTxQueue.Set(suite.Context, types.EthTxQueue{})
	suite.Require().NoError(err)

	msgServer := keeper.NewMsgServerImpl(suite.Keeper)

	req := &types.MsgApproveCancellation{
		Proposer: "goat1xa56637tjn857jyg2plgvhdclzmr4crxzn5xus",
	}

	expected := []types.Withdrawal{}
	for idx, address := range withdrawals {
		wd := types.Withdrawal{
			Address:       address,
			RequestAmount: 1e3,
			MaxTxPrice:    5,
			Status:        types.WITHDRAWAL_STATUS_CANCELING,
		}
		err = suite.Keeper.Withdrawals.Set(suite.Context, uint64(idx), wd)
		req.Id = append(req.Id, uint64(idx))
		suite.Require().NoError(err)
		expected = append(expected, wd)
	}

	suite.RelayerKeeper.EXPECT().VerifyNonProposal(suite.Context, req).Return(nil, nil)

	_, err = msgServer.ApproveCancellation(suite.Context, req)
	suite.Require().NoError(err)

	queue, err := suite.Keeper.EthTxQueue.Get(suite.Context)
	suite.Require().NoError(err)
	suite.Require().Len(queue.RejectedWithdrawals, len(expected))

	for idx, wd := range expected {
		wd.Status = types.WITHDRAWAL_STATUS_CANCELED
		withdrawal, err := suite.Keeper.Withdrawals.Get(suite.Context, uint64(idx))
		suite.Require().NoError(err)
		suite.Require().Equal(withdrawal, wd)
		suite.Require().EqualValues(idx, queue.RejectedWithdrawals[idx])
	}
}

func (suite *KeeperTestSuite) TestNewConsolidation() {
	err := suite.Keeper.Pubkey.Set(suite.Context, relayertypes.PublicKey{Key: &relayertypes.PublicKey_Secp256K1{
		Secp256K1: common.Hex2Bytes("037e7bee29c1956152e308d1310823295d720b4cef9e1118726eb1705ffc5a4701"),
	}})
	suite.Require().NoError(err)

	msgServer := keeper.NewMsgServerImpl(suite.Keeper)
	req := &types.MsgNewConsolidation{
		Proposer:    "goat1xa56637tjn857jyg2plgvhdclzmr4crxzn5xus",
		NoWitnessTx: common.Hex2Bytes("020000000a59470436f741658f24f25ec84e1af99a1355173b36e19abcd1d8d9079435be470000000000ffffffff59470436f741658f24f25ec84e1af99a1355173b36e19abcd1d8d9079435be470100000000ffffffff59470436f741658f24f25ec84e1af99a1355173b36e19abcd1d8d9079435be470200000000ffffffff59470436f741658f24f25ec84e1af99a1355173b36e19abcd1d8d9079435be470300000000ffffffff59470436f741658f24f25ec84e1af99a1355173b36e19abcd1d8d9079435be470400000000ffffffff59470436f741658f24f25ec84e1af99a1355173b36e19abcd1d8d9079435be470500000000ffffffff59470436f741658f24f25ec84e1af99a1355173b36e19abcd1d8d9079435be470600000000ffffffff59470436f741658f24f25ec84e1af99a1355173b36e19abcd1d8d9079435be470700000000ffffffff59470436f741658f24f25ec84e1af99a1355173b36e19abcd1d8d9079435be470800000000ffffffff59470436f741658f24f25ec84e1af99a1355173b36e19abcd1d8d9079435be470900000000ffffffff0100e1f505000000001600145b029559baaea5e928e8e2774e9e2350a5fc9c2d00000000"),
		Vote:        &relayertypes.Votes{Signature: make([]byte, goatcrypto.SignatureLength)},
	}

	sequence := uint64(100)
	suite.RelayerKeeper.EXPECT().VerifyProposal(suite.Context, req).Return(sequence, nil)
	suite.RelayerKeeper.EXPECT().SetProposalSeq(suite.Context, sequence+1)
	suite.RelayerKeeper.EXPECT().UpdateRandao(suite.Context, req)

	_, err = msgServer.NewConsolidation(suite.Context, req)
	suite.Require().NoError(err)

	txid, err := chainhash.NewHashFromStr("1ed1c224d55d41aa03d0baa14f370f627c39ac21970d92ef87694abf60c9b76d")
	suite.Require().NoError(err)

	events := suite.Context.EventManager().ABCIEvents()
	suite.Require().Equal(len(events), 2)

	event0 := types.NewConsolidationEvent(txid[:])
	suite.Require().Equal(events[0].Type, event0.Type)
	suite.Require().Equal(len(events[0].Attributes), len(event0.Attributes))
	suite.Require().Equal(events[0].Attributes[0].Key, event0.Attributes[0].Key)
	suite.Require().Equal(events[0].Attributes[0].Value, event0.Attributes[0].Value)

	event1 := relayertypes.FinalizedProposalEvent(sequence)
	suite.Require().Equal(events[1].Type, event1.Type)
	suite.Require().Equal(len(events[1].Attributes), len(event1.Attributes))
	suite.Require().Equal(events[1].Attributes[0].Key, event1.Attributes[0].Key)
	suite.Require().Equal(events[1].Attributes[0].Value, event1.Attributes[0].Value)
}
