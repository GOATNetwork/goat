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
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/goat/types"
	"golang.org/x/sync/errgroup"
)

const maxTxLen = 16

func (k Keeper) PrepareProposalHandler(
	txpool mempool.Mempool,
	txVerifier baseapp.ProposalTxVerifier,
	keyProvider cryptotypes.PrivKey,
	txConfig client.TxConfig,
) sdk.PrepareProposalHandler {
	if keyProvider == nil {
		panic("no eth block signer provided")
	}

	return func(sdkctx sdk.Context, rpp *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		var ethTx []byte

		eg := new(errgroup.Group)

		// build eth block msg
		eg.Go(func() (err error) {
			ethTx, err = k.createEthBlockProposal(sdkctx, keyProvider, txConfig, rpp)
			return
		})

		// select relayer message from mempool
		var memTxs [][]byte
		eg.Go(func() error {
			iterator := txpool.Select(sdkctx, rpp.Txs)
			for iterator != nil {
				memTx := iterator.Tx()
				txBytes, err := txVerifier.PrepareProposalVerifyTx(memTx)
				if err != nil {
					k.Logger().Info("Remove tx from mempool", "reason", err.Error())

					err := txpool.Remove(memTx)
					if err != nil && !errors.Is(err, mempool.ErrTxNotFound) {
						return err
					}
					iterator = iterator.Next()
					continue
				}
				memTxs = append(memTxs, txBytes)
				if len(memTxs)+1 >= maxTxLen {
					return nil
				}
				iterator = iterator.Next()
			}
			return nil
		})

		if err := eg.Wait(); err != nil {
			return nil, err
		}
		return &abci.ResponsePrepareProposal{Txs: append([][]byte{ethTx}, memTxs...)}, nil
	}
}

func (k Keeper) createEthBlockProposal(sdkctx sdk.Context, keyProvider cryptotypes.PrivKey, txConfig client.TxConfig, rpp *abci.RequestPrepareProposal) ([]byte, error) {
	validatorAddr, err := k.addressCodec.BytesToString(rpp.ProposerAddress)
	if err != nil {
		return nil, err
	}

	validatorAcc := k.accountKeeper.GetAccount(sdkctx, rpp.ProposerAddress)
	if validatorAcc == nil {
		return nil, fmt.Errorf("nil validator account: %s", validatorAddr)
	}

	if !bytes.Equal(validatorAcc.GetPubKey().Bytes(), keyProvider.PubKey().Bytes()) {
		return nil, fmt.Errorf("validator pubkey mismatched: expected %x got %x", validatorAcc.GetPubKey().Bytes(), keyProvider.PubKey().Bytes())
	}

	parentBlock, err := k.Block.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	beaconBlock, err := k.BeaconRoot.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	goatTxs, err := k.Dequeue(sdkctx)
	if err != nil {
		return nil, err
	}

	tmctx, cancel := context.WithTimeout(sdkctx, time.Second*2)
	defer cancel()

	// yeah, we have a proposer based random number
	var random common.Hash
	_, _ = rand.Read(random[:])

	beaconRoot := common.BytesToHash(beaconBlock)
	forkChoiceResp, err := k.ethclient.ForkchoiceUpdatedV3(tmctx,
		&engine.ForkchoiceStateV1{HeadBlockHash: common.BytesToHash(parentBlock.BlockHash)},
		&engine.PayloadAttributes{
			// Why we use system timestamp instead of CometBFT timestamp?
			// CometBFT 0.38 uses last vote median time by all voters
			// The proposer-based timestamps (PBTS) will be enabled on CometBFT 1.0
			// We can use the consensus layer timestamp after migrating to 1.0
			Timestamp:             uint64(time.Now().UTC().Unix()),
			Random:                random,
			SuggestedFeeRecipient: common.BytesToAddress(rpp.ProposerAddress),
			Withdrawals:           ethtypes.Withdrawals{},
			BeaconRoot:            &beaconRoot,
			GoatTxs:               goatTxs,
		},
	)
	if err != nil {
		return nil, err
	}

	if status := forkChoiceResp.PayloadStatus; status.Status != engine.VALID {
		if status.ValidationError != nil {
			return nil, fmt.Errorf("failed to build goat-geth txs: %s", *status.ValidationError)
		}
		return nil, errors.New("failed to build goat-geth txs")
	}

	if forkChoiceResp.PayloadID == nil {
		return nil, errors.New("got nil payloadId")
	}

	// Note: the waiting duration is for the payload building starting instead of finishing
	<-time.After(time.Millisecond * 50)
	envelope, err := k.ethclient.GetPayloadV4(tmctx, *forkChoiceResp.PayloadID)
	if err != nil {
		return nil, err
	}

	txBuilder := txConfig.NewTxBuilder()
	txBuilder.SetGasLimit(1e8)
	txBuilder.SetTimeoutHeight(uint64(rpp.Height))

	payload := types.ExecutableDataToPayload(envelope.ExecutionPayload, beaconBlock, envelope.Requests)
	k.Logger().Info("Propose new executable payload", payload.LogKeyVals()...)
	if err := txBuilder.SetMsgs(&types.MsgNewEthBlock{Proposer: validatorAddr, Payload: payload}); err != nil {
		return nil, err
	}

	sigMode := signing.SignMode(txConfig.SignModeHandler().DefaultMode())
	if err := txBuilder.SetSignatures(signing.SignatureV2{
		PubKey:   validatorAcc.GetPubKey(),
		Data:     &signing.SingleSignatureData{SignMode: sigMode},
		Sequence: validatorAcc.GetSequence(),
	}); err != nil {
		return nil, err
	}

	sigs, err := tx.SignWithPrivKey(sdkctx, sigMode, xauthsigning.SignerData{
		Address:       validatorAddr,
		ChainID:       sdkctx.ChainID(),
		AccountNumber: validatorAcc.GetAccountNumber(),
		Sequence:      validatorAcc.GetSequence(),
		PubKey:        validatorAcc.GetPubKey(),
	}, txBuilder, keyProvider, txConfig, validatorAcc.GetSequence())
	if err != nil {
		return nil, err
	}

	if err := txBuilder.SetSignatures(sigs); err != nil {
		return nil, err
	}

	ethTx, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	return ethTx, nil
}

func (k Keeper) ProcessProposalHandler(txVerifier baseapp.ProposalTxVerifier) sdk.ProcessProposalHandler {
	return func(sdkctx sdk.Context, rpp *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		if l := len(rpp.Txs); l == 0 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no transactions")
		} else if l > maxTxLen {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "too many transactions")
		}

		for txIdx, rawTx := range rpp.Txs {
			tx, err := txVerifier.ProcessProposalVerifyTx(rawTx)
			if err != nil {
				return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid transaction: index %d: %s", txIdx, err)
			}

			msgs := tx.GetMsgs()

			// the first tx should be for MsgNewEthBlock
			if txIdx == 0 {
				if len(msgs) != 1 {
					return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid MsgNewEthBlock message")
				}
				ethBlock, ok := msgs[0].(*types.MsgNewEthBlock)
				if !ok {
					return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "the first tx should be MsgNewEthBlock")
				}
				if err := k.verifyEthBlockProposal(sdkctx, ethBlock); err != nil {
					return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid MsgNewEthBlock: %s", err.Error())
				}
				continue
			}

			// the rest txs should be belong to goat package but must not be MsgNewEthBlock
			// the GoatGuardHandler checks the first case is matched
			// and we check the second case here
			for _, msg := range msgs {
				if _, ok := msg.(*types.MsgNewEthBlock); ok {
					return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "MsgNewEthBlock should be first tx in block")
				}
			}
		}
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

func (k Keeper) verifyEthBlockProposal(sdkctx sdk.Context, msg *types.MsgNewEthBlock) error {
	payload := msg.Payload
	if payload == nil {
		return errors.New("empty payload")
	}

	k.Logger().Info("Verify new executable payload", payload.LogKeyVals()...)
	eg, egctx := errgroup.WithContext(sdkctx)
	eg.Go(func() error {
		proposer, err := k.addressCodec.StringToBytes(msg.Proposer)
		if err != nil {
			return err
		}

		if expect := sdkctx.CometInfo().GetProposerAddress(); !bytes.Equal(proposer, expect) {
			return fmt.Errorf("invalid MsgNewEthBlock proposer: expect %x got %x", expect, proposer)
		}

		if !bytes.Equal(proposer, payload.FeeRecipient) {
			return errors.New("fee recipient mismatched")
		}

		// we don't use cometbft timestamp
		// refer to the note in the PrepareProposalHandler for the details
		if systemTime := uint64(time.Now().UTC().Unix()); payload.Timestamp > systemTime {
			return errors.New("invalid MsgNewEthBlock timestamp")
		}

		block, err := k.Block.Get(sdkctx)
		if err != nil {
			return err
		}

		if !bytes.Equal(block.BlockHash, payload.ParentHash) {
			return fmt.Errorf("incorrect parent block number: expected %x got %x", block.BlockHash, payload.ParentHash)
		}

		if block.BlockNumber+1 != payload.BlockNumber {
			return fmt.Errorf("incorrect parent block hash: expected %d got %d", block.BlockNumber+1, payload.BlockNumber)
		}

		if _, _, lockingReqs, err := goattypes.DecodeRequests(payload.Requests, false); err != nil {
			return fmt.Errorf("invalid goat requests: %w", err)
		} else if len(lockingReqs.Gas) != 1 {
			return errors.New("gas revenue request length is not 1")
		}

		beaconRoot, err := k.BeaconRoot.Get(sdkctx)
		if err != nil {
			return err
		}

		if !bytes.Equal(beaconRoot, payload.BeaconRoot) {
			return fmt.Errorf("refer to incorrect beacon root: expected %x got %x", beaconRoot, payload.BeaconRoot)
		}

		if err := k.VerifyDequeue(sdkctx, payload.ExtraData, payload.Transactions); err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		res, err := k.ethclient.NewPayloadV4(egctx, types.PayloadToExecutableData(payload),
			[]common.Hash{}, common.BytesToHash(payload.BeaconRoot), payload.Requests)
		if err != nil {
			return err
		}

		if res.Status != engine.VALID {
			if res.ValidationError != nil {
				return fmt.Errorf("NewPayloadV4 non-VALID status(%s): %s", res.Status, *res.ValidationError)
			}
			return fmt.Errorf("NewPayloadV4 non-VALID status: %s", res.Status)
		}
		return nil
	})

	return eg.Wait()
}
