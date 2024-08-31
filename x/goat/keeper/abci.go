package keeper

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/goatnetwork/goat/x/goat/types"
	"golang.org/x/sync/errgroup"
)

func (k Keeper) PrepareProposalHandler(txpool mempool.Mempool, txVerifier baseapp.ProposalTxVerifier) sdk.PrepareProposalHandler {
	signerAdapter := mempool.NewDefaultSignerExtractionAdapter()

	txSelector := func(proposer sdk.AccAddress, seq uint64, memTx sdk.Tx) ([]byte, error) {
		signerData, err := signerAdapter.GetSigners(memTx)
		if err != nil {
			return nil, err
		}

		if len(signerData) != 1 {
			return nil, errors.New("Not a relayer tx")
		}

		signer := signerData[0]
		if !proposer.Equals(signer.Signer) {
			return nil, errors.New("Not a relayer proposer")
		}

		if signer.Sequence != seq {
			return nil, errors.New("Not next relayer proposer sequence")
		}

		for _, msg := range memTx.GetMsgs() {
			switch msg.(type) {
			case *types.MsgNewEthBlock:
				return nil, errors.New("MsgNewEthBlock should not be placed in the mempool")
			}
		}
		return txVerifier.PrepareProposalVerifyTx(memTx)
	}

	return func(sdkctx sdk.Context, rpp *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		var ethBlockTxByte []byte

		var eg = new(errgroup.Group)

		// build eth block msg
		eg.Go(func() error {
			validator := k.accountKeeper.GetAccount(sdkctx, rpp.ProposerAddress)
			if validator == nil {
				return errors.New("nil validator account")
			}

			block, err := k.Block.Get(sdkctx)
			if err != nil {
				return err
			}

			blockHash, err := k.BeaconRoot.Get(sdkctx)
			if err != nil {
				return err
			}

			goatTxs, err := k.Dequeue(sdkctx)
			if err != nil {
				return err
			}

			tmctx, cancel := context.WithTimeout(sdkctx, time.Second*2)
			defer cancel()

			var random common.Hash
			_, _ = rand.Read(random[:])

			beaconRoot := common.BytesToHash(blockHash)
			forkChoiceResp, err := k.ethclient.ForkchoiceUpdatedV3(tmctx,
				&engine.ForkchoiceStateV1{HeadBlockHash: common.BytesToHash(block.BlockHash)},
				&engine.PayloadAttributes{
					// Note: proposer-based timestamps (PBTS) is not enable on cometbft 0.38
					Timestamp:             uint64(time.Now().UTC().Unix()),
					Random:                random,
					SuggestedFeeRecipient: common.BytesToAddress(rpp.ProposerAddress),
					BeaconRoot:            &beaconRoot,
					GoatTxs:               goatTxs,
				},
			)
			if err != nil {
				return err
			}

			if status := forkChoiceResp.PayloadStatus; status.Status != engine.VALID {
				if status.ValidationError != nil {
					return fmt.Errorf("failed to build goat-geth txs: %s", *status.ValidationError)
				}
				return errors.New("failed to build goat-geth txs")
			}

			envelope, err := k.ethclient.GetPayloadV3(tmctx, *forkChoiceResp.PayloadID)
			if err != nil {
				return err
			}

			proposer, err := k.addressCodec.BytesToString(rpp.ProposerAddress)
			if err != nil {
				return err
			}

			txBuilder := k.txConfig.NewTxBuilder()
			if err := txBuilder.SetMsgs(&types.MsgNewEthBlock{
				Proposer: proposer,
				Payload:  types.ExecutableDataToPayload(envelope.ExecutionPayload),
			}); err != nil {
				return err
			}
			txBuilder.SetGasLimit(1e8)
			txBuilder.SetTimeoutHeight(uint64(rpp.Height))

			sigMode := signing.SignMode(k.txConfig.SignModeHandler().DefaultMode())
			if err := txBuilder.SetSignatures(signing.SignatureV2{
				PubKey:   validator.GetPubKey(),
				Data:     &signing.SingleSignatureData{SignMode: sigMode},
				Sequence: validator.GetSequence(),
			}); err != nil {
				return err
			}

			sigsV2, err := tx.SignWithPrivKey(sdkctx, sigMode, xauthsigning.SignerData{
				ChainID:       sdkctx.BlockHeader().ChainID,
				AccountNumber: validator.GetAccountNumber(),
				Sequence:      validator.GetSequence(),
			}, txBuilder, k.PrivKey, k.txConfig, validator.GetSequence())
			if err != nil {
				return err
			}

			if err := txBuilder.SetSignatures(sigsV2); err != nil {
				return err
			}

			ethBlockTxByte, err = k.txConfig.TxEncoder()(txBuilder.GetTx())
			if err != nil {
				return err
			}

			return nil
		})

		// select relayer message from mempool
		var memTxs [][]byte
		eg.Go(func() error {
			proposer, err := k.relayerKeeper.GetCurrentProposer(sdkctx)
			if err != nil {
				return err
			}

			sequence, err := k.accountKeeper.GetSequence(sdkctx, proposer)
			if err != nil {
				return err
			}

			iterator := txpool.Select(sdkctx, rpp.Txs)
			for iterator != nil {
				memTx := iterator.Tx()
				txBytes, err := txSelector(proposer, sequence, memTx)
				if err != nil {
					k.Logger().Debug("Remove mempool tx", "reason", err.Error())

					err := txpool.Remove(memTx)
					if err != nil && !errors.Is(err, mempool.ErrTxNotFound) {
						return err
					}
					iterator = iterator.Next()
					continue
				}
				memTxs = append(memTxs, txBytes)
				iterator = iterator.Next()
				sequence++
			}

			return nil
		})

		if err := eg.Wait(); err != nil {
			return nil, err
		}
		return &abci.ResponsePrepareProposal{Txs: append([][]byte{ethBlockTxByte}, memTxs...)}, nil
	}
}

func (k Keeper) ProcessProposalHandler(txVerifier baseapp.ProposalTxVerifier) sdk.ProcessProposalHandler {
	return func(sdkctx sdk.Context, rpp *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		if len(rpp.Txs) == 0 {
			return nil, errors.New("No transactions")
		}

		for txIdx, rawTx := range rpp.Txs {
			tx, err := txVerifier.ProcessProposalVerifyTx(rawTx)
			if err != nil {
				return nil, errors.New("invalid message")
			}

			if txIdx == 0 {
				msgs := tx.GetMsgs()
				if len(msgs) != 1 {
					return nil, errors.New("invalid MsgNewEthBlock message")
				}
				ethBlock, ok := msgs[0].(*types.MsgNewEthBlock)
				if !ok {
					return nil, errors.New("the first tx should be MsgNewEthBlock")
				}
				if err := k.verifyEthBlockProposal(sdkctx, rpp.ProposerAddress, ethBlock); err != nil {
					return nil, err
				}
				continue
			}

			// todo: using msg url to check if the tx is allowed
			// we should deny many message types like x/bank and x/staking
			for _, msg := range tx.GetMsgs() {
				switch msg.(type) {
				case *types.MsgNewEthBlock:
					return nil, errors.New("MsgNewEthBlock should be first tx in block")
				}
			}
		}
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

func (k Keeper) verifyEthBlockProposal(sdkctx context.Context, expectProposer []byte, ethBlock *types.MsgNewEthBlock) error {
	proposer, err := k.addressCodec.StringToBytes(ethBlock.Proposer)
	if err != nil {
		return err
	}

	if !bytes.Equal(proposer, expectProposer) {
		return errors.New("invalid MsgNewEthBlock proposer")
	}

	payload := ethBlock.Payload
	if !bytes.Equal(proposer, ethBlock.Payload.FeeRecipient) {
		return errors.New("fee recipient mismatched")
	}

	if payload.BlobGasUsed > 0 {
		return errors.New("blob tx is not allowed")
	}

	// we don't use cometbft timestamp
	// refer to the note in the PrepareProposalHandler for the details
	if now := uint64(time.Now().UTC().Unix()); payload.Timestamp > now {
		return errors.New("invalid MsgNewEthBlock timestamp")
	}

	block, err := k.Block.Get(sdkctx)
	if err != nil {
		return err
	}

	if !bytes.Equal(block.ParentHash, payload.ParentHash) || block.BlockNumber+1 != payload.BlockNumber {
		return errors.New("refer to incorrect parent block")
	}

	beaconRoot, err := k.BeaconRoot.Get(sdkctx)
	if err != nil {
		return err
	}

	if !bytes.Equal(beaconRoot, payload.BeaconRoot) {
		return errors.New("refer to incorrect beacon root")
	}

	if err := k.VerifyDequeue(sdkctx, payload.Transactions); err != nil {
		return err
	}

	res, err := k.ethclient.NewPayloadV3(sdkctx,
		types.PayloadToExecutableData(&payload), nil, common.BytesToHash(beaconRoot))
	if err != nil {
		return err
	}

	if res.Status != engine.VALID {
		if res.ValidationError != nil {
			return fmt.Errorf("NewPayloadV3 non-VALID status(%s): %s", res.Status, *res.ValidationError)
		}
		return fmt.Errorf("NewPayloadV3 non-VALID status: %s", res.Status)
	}
	return nil
}
