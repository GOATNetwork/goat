package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/relayer/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) NewVoter(ctx context.Context, req *types.MsgNewVoterRequest) (*types.MsgNewVoterResponse, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	if err := req.Validate(); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	relayer, err := k.VerifyNonProposal(sdkctx, req)
	if err != nil {
		return nil, err
	}

	addrRaw := sdktypes.AccAddress(goatcrypto.Hash160Sum(req.VoterTxKey))
	addrStr, err := k.AddrCodec.BytesToString(addrRaw)
	if err != nil {
		return nil, err
	}

	voter, err := k.Voters.Get(sdkctx, addrStr)
	if err != nil {
		return nil, err
	}

	if voter.Status != types.VOTER_STATUS_PENDING {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "not a pending voter")
	}

	if !bytes.Equal(goatcrypto.SHA256Sum(req.VoterBlsKey), voter.VoteKey) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "vote key hash not match")
	}

	reqMsg := types.NewOnBoardingVoterRequest(voter.Height, addrRaw, voter.VoteKey)
	sigMsg := types.VoteSignDoc(
		reqMsg.MethodName(), sdkctx.ChainID(), req.Proposer, 0 /* sequence */, relayer.GetEpoch(), reqMsg.SignDoc())

	if !ethcrypto.VerifySignature(req.VoterTxKey, sigMsg, req.VoterTxKeyProof) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "false tx key proof")
	}

	if !goatcrypto.Verify(req.VoterBlsKey, sigMsg, req.VoterBlsKeyProof) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "false vote key proof")
	}

	// add account
	hasAccount := k.accountKeeper.HasAccount(sdkctx, addrRaw)
	if !hasAccount {
		acc := k.accountKeeper.NewAccountWithAddress(sdkctx, addrRaw)
		if err := acc.SetPubKey(&secp256k1.PubKey{Key: req.VoterTxKey}); err != nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "failed to set voter's account pubkey: %s", err.Error())
		}
		k.accountKeeper.SetAccount(sdkctx, acc)
	}

	queue, err := k.Queue.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	voter.VoteKey = req.VoterBlsKey
	if hasAccount {
		voter.Status = types.VOTER_STATUS_OFF_BOARDING
		queue.OffBoarding = append(queue.OffBoarding, addrStr)
		sdkctx.EventManager().EmitEvent(types.RemovingVoterEvent(addrStr))
	} else {
		voter.Status = types.VOTER_STATUS_ON_BOARDING
		queue.OnBoarding = append(queue.OnBoarding, addrStr)
		sdkctx.EventManager().EmitEvent(types.VoterBoardedEvent(relayer.GetProposer(), addrStr))
	}

	if err := k.Queue.Set(sdkctx, queue); err != nil {
		return nil, err
	}

	if err := k.Voters.Set(sdkctx, addrStr, voter); err != nil {
		return nil, err
	}
	return &types.MsgNewVoterResponse{}, nil
}

func (k msgServer) AcceptProposer(ctx context.Context, req *types.MsgAcceptProposerRequest) (*types.MsgAcceptProposerResponse, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	relayer, err := k.Relayer.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	if relayer.Proposer != req.Proposer {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "not the current proposer")
	}

	if relayer.ProposerAccepted {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "proposer has been accepted")
	}

	if relayer.Epoch != req.Epoch {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid epoch: expected %d", relayer.Epoch)
	}

	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	if sdkctx.BlockTime().Sub(relayer.LastElected) > param.AcceptProposerTimeout {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "timeout to accept proposer role")
	}

	relayer.ProposerAccepted = true
	if err := k.Relayer.Set(sdkctx, relayer); err != nil {
		return nil, err
	}

	k.Logger().Info("New proposer is accepted", "epoch", relayer.Epoch, "proposer", relayer.Proposer)
	sdkctx.EventManager().EmitEvent(types.AcceptedProposerEvent(relayer.Proposer, relayer.Epoch))
	return &types.MsgAcceptProposerResponse{}, nil
}
