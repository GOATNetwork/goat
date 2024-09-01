package keeper

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/goatnetwork/goat/x/goat/types"
	"golang.org/x/sync/errgroup"
)

func (k Keeper) PrepareProposalHandler(txpool mempool.Mempool, txVerifier baseapp.ProposalTxVerifier) sdk.PrepareProposalHandler {
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
					// Why we use system timestamp instead of CometBFT timestamp?
					// CometBFT 0.38 uses last vote median time by all voters
					// The proposer-based timestamps (PBTS) will be enabled on CometBFT 1.0
					// We can use the consensus layer timestamp after migrating to 1.0
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
			iterator := txpool.Select(sdkctx, rpp.Txs)
			for iterator != nil {
				memTx := iterator.Tx()
				txBytes, err := txVerifier.PrepareProposalVerifyTx(memTx)
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
			return nil, errors.New("empty block")
		}

		for txIdx, rawTx := range rpp.Txs {
			tx, err := txVerifier.ProcessProposalVerifyTx(rawTx)
			if err != nil {
				return nil, errors.New("invalid message")
			}

			msgs := tx.GetMsgs()

			// the first tx should be for MsgNewEthBlock
			if txIdx == 0 {
				if len(msgs) != 1 {
					return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "invalid MsgNewEthBlock message")
				}
				ethBlock, ok := msgs[0].(*types.MsgNewEthBlock)
				if !ok {
					return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "the first tx should be MsgNewEthBlock")
				}
				if err := k.verifyEthBlockProposal(sdkctx, rpp.ProposerAddress, ethBlock); err != nil {
					return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "invalid MsgNewEthBlock: %s", err.Error())
				}
				continue
			}

			// the rest txs should be belong to goat package but must not be MsgNewEthBlock
			// the GoatGuardHandler checks the first case is matched
			// and we check the second case here
			for _, msg := range msgs {
				switch msg.(type) {
				case *types.MsgNewEthBlock:
					return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "MsgNewEthBlock should be first tx in block")
				}
			}
		}
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

func (k Keeper) verifyEthBlockProposal(sdkctx sdk.Context, cometProposer []byte, ethBlock *types.MsgNewEthBlock) error {
	proposer, err := k.addressCodec.StringToBytes(ethBlock.Proposer)
	if err != nil {
		return err
	}

	if !bytes.Equal(proposer, cometProposer) {
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
	systemTime, cometTime := uint64(time.Now().UTC().Unix()), uint64(sdkctx.BlockTime().UTC().Unix())
	if payload.Timestamp > systemTime || payload.Timestamp < cometTime {
		return errors.New("invalid MsgNewEthBlock timestamp")
	}

	block, err := k.Block.Get(sdkctx)
	if err != nil {
		return err
	}

	if !bytes.Equal(block.ParentHash, payload.ParentHash) || block.BlockNumber+1 != payload.BlockNumber {
		return errors.New("refer to incorrect parent block")
	}

	if payload.BlobGasUsed > 0 {
		return errors.New("blob tx type is not activated")
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
