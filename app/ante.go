package app

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	txsigning "cosmossdk.io/x/tx/signing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	goattypes "github.com/goatnetwork/goat/x/goat/types"
)

func NewAnteHandler(accKeeper ante.AccountKeeper, relayerKeeper goattypes.RelayerKeeper, signModeHandler *txsigning.HandlerMap) sdk.AnteHandler {
	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		GoatGuardHandler{relayerKeeper},
		ante.NewValidateBasicDecorator(),
		TxTimeoutHeightDecorator{},
		ante.NewSetPubKeyDecorator(accKeeper),
		ante.NewSigVerificationDecorator(accKeeper, signModeHandler),
		ante.NewIncrementSequenceDecorator(accKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...)
}

type GoatGuardHandler struct {
	relayerKeeper goattypes.RelayerKeeper
}

func (ante GoatGuardHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// memo length check
	memoTx, ok := tx.(sdk.TxWithMemo)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	if memoLength := len(memoTx.GetMemo()); memoLength > 0 {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrMemoTooLarge, "no memo required")
	}

	// sig count check
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a sigTx")
	}

	signers, err := sigTx.GetSigners()
	if err != nil {
		return ctx, err
	}

	if len(signers) != 1 {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrTooManySignatures, "signer count %d > 1", len(signers))
	}

	// msg allow list check
	msgs, err := tx.GetMsgsV2()
	if err != nil {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must support MsgV2")
	}

	relayerProposer, err := ante.relayerKeeper.GetCurrentProposer(ctx)
	if err != nil {
		return ctx, nil
	}

	var relayerTxOnly = func(msgName string) error {
		if !strings.HasPrefix(msgName, "goat.bitcoin.") && !strings.HasPrefix(msgName, "goat.relayer.") {
			return errorsmod.Wrapf(sdkerrors.ErrTxDecode, "%s is not a relayer message", msgName)
		}
		if !relayerProposer.Equals(sdk.AccAddress(signers[0])) {
			return errorsmod.Wrapf(sdkerrors.ErrorInvalidSigner, "%s is not current relayer proposer", signers[0])
		}
		return nil
	}

	if !simulate {
		// the message should be belong to goad namespace
		for _, msg := range msgs {
			msgName := msg.ProtoReflect().Descriptor().FullName()
			switch ctx.ExecMode() {
			case sdk.ExecModeCheck, sdk.ExecModeReCheck, sdk.ExecModePrepareProposal: // only accept relayer txs in the mempool
				if err := relayerTxOnly(string(msgName)); err != nil {
					return ctx, err
				}
			case sdk.ExecModeProcessProposal, sdk.ExecModeFinalize:
				if string(msgName) == "goat.goat.v1.MsgNewEthBlock" {
					continue
				}
				if err := relayerTxOnly(string(msgName)); err != nil {
					return ctx, err
				}
			}
		}
	}

	return next(newCtx, tx, simulate)
}

type TxTimeoutHeightDecorator struct{}

func (txh TxTimeoutHeightDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	timeoutTx, ok := tx.(sdk.TxWithTimeoutHeight)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "expected tx to implement TxWithTimeoutHeight")
	}

	timeoutHeight := timeoutTx.GetTimeoutHeight()

	if msgs := tx.GetMsgs(); len(msgs) == 1 {
		if _, ok := msgs[0].(*goattypes.MsgNewEthBlock); ok {
			if timeoutHeight != uint64(ctx.BlockHeight()) {
				return ctx, errorsmod.Wrapf(
					sdkerrors.ErrTxTimeoutHeight, "MsgNewEthBlock timeout height should be current block: %d", ctx.BlockHeight())
			}
		}
	}

	if timeoutHeight > 0 && uint64(ctx.BlockHeight()) > timeoutHeight {
		return ctx, errorsmod.Wrapf(
			sdkerrors.ErrTxTimeoutHeight, "block height: %d, timeout height: %d", ctx.BlockHeight(), timeoutHeight,
		)
	}

	return next(ctx, tx, simulate)
}
