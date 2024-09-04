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
	if accKeeper == nil {
		panic("nil account keeper")
	}

	if relayerKeeper == nil {
		panic("nil relayer keeper")
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		GoatGuardHandler{relayerKeeper},
		ante.NewValidateBasicDecorator(),
		ante.NewSetPubKeyDecorator(accKeeper),
		ante.NewSigVerificationDecorator(accKeeper, signModeHandler),
		ante.NewIncrementSequenceDecorator(accKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...)
}

type GoatGuardHandler struct {
	relayerKeeper goattypes.RelayerKeeper
}

type StdTx interface {
	sdk.TxWithMemo
	sdk.TxWithTimeoutHeight
	authsigning.SigVerifiableTx
}

func (ante GoatGuardHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	stdTx, ok := tx.(StdTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	if len(stdTx.GetMemo()) > 0 {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrMemoTooLarge, "no memo required")
	}

	signers, err := stdTx.GetSigners()
	if err != nil {
		return ctx, err
	}

	if len(signers) != 1 {
		return ctx, errorsmod.Wrapf(sdkerrors.ErrTooManySignatures, "signer count more than 1(%d)", len(signers))
	}

	timeoutHeight := stdTx.GetTimeoutHeight()
	if timeoutHeight > 0 && uint64(ctx.BlockHeight()) > timeoutHeight {
		return ctx, errorsmod.Wrapf(
			sdkerrors.ErrTxTimeoutHeight, "block height: %d, timeout height: %d", ctx.BlockHeight(), timeoutHeight,
		)
	}

	// msg allow list check
	msgs, err := tx.GetMsgsV2()
	if err != nil {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must support MsgV2")
	}

	relayerProposer, err := ante.relayerKeeper.GetCurrentProposer(ctx)
	if err != nil {
		return ctx, err
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
			msgName := string(msg.ProtoReflect().Descriptor().FullName())
			switch ctx.ExecMode() {
			case sdk.ExecModeCheck, sdk.ExecModeReCheck, sdk.ExecModePrepareProposal: // only accept relayer txs in the mempool
				if err := relayerTxOnly(msgName); err != nil {
					return ctx, err
				}
			case sdk.ExecModeProcessProposal, sdk.ExecModeFinalize:
				if (msgName) == "goat.goat.v1.MsgNewEthBlock" {
					if timeoutHeight != uint64(ctx.BlockHeight()) {
						return ctx, errorsmod.Wrapf(
							sdkerrors.ErrTxTimeoutHeight, "MsgNewEthBlock timeout height should be current block: %d", ctx.BlockHeight())
					}
					continue
				}
				if err := relayerTxOnly(msgName); err != nil {
					return ctx, err
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
