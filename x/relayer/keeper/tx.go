package keeper

import (
	"bytes"
	"context"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
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
		return nil, types.ErrInvalid.Wrapf("invalid request: %s", err.Error())
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
		return nil, types.ErrInvalid.Wrapf("not a pending voter")
	}

	if !bytes.Equal(goatcrypto.SHA256Sum(req.VoterBlsKey), voter.VoteKey) {
		return nil, types.ErrInvalid.Wrapf("vote key hash not match")
	}

	reqMsg := types.NewOnBoardingVoterRequest(voter.Height, addrRaw, voter.VoteKey)
	sigMsg := types.VoteSignDoc(
		reqMsg.MethodName(), sdkctx.ChainID(), req.Proposer, 0 /* sequence */, relayer.GetEpoch(), reqMsg.SignDoc())

	if !ethcrypto.VerifySignature(req.VoterTxKey, sigMsg, req.VoterTxKeyProof) {
		return nil, types.ErrInvalid.Wrapf("false tx key proof")
	}

	if !goatcrypto.Verify(req.VoterBlsKey, sigMsg, req.VoterBlsKeyProof) {
		return nil, types.ErrInvalid.Wrapf("false vote key proof")
	}

	// add account
	if !k.accountKeeper.HasAccount(sdkctx, addrRaw) {
		acc := k.accountKeeper.NewAccountWithAddress(sdkctx, addrRaw)
		if err := acc.SetPubKey(&secp256k1.PubKey{Key: req.VoterTxKey}); err != nil {
			return nil, types.ErrInvalid.Wrapf("unable to set pubkey")
		}
		k.accountKeeper.SetAccount(sdkctx, acc)
	}

	voter.VoteKey = req.VoterBlsKey
	voter.Status = types.VOTER_STATUS_ON_BOARDING
	if err := k.Voters.Set(sdkctx, addrStr, voter); err != nil {
		return nil, err
	}

	queue, err := k.Queue.Get(sdkctx)
	if err != nil {
		return nil, err
	}
	queue.OnBoarding = append(queue.OnBoarding, addrStr)

	if err := k.Queue.Set(sdkctx, queue); err != nil {
		return nil, err
	}

	sdkctx.EventManager().EmitEvent(types.VoterBoardedEvent(relayer.GetProposer(), addrStr))
	return &types.MsgNewVoterResponse{}, nil
}

func (k msgServer) AcceptProposer(ctx context.Context, req *types.MsgAcceptProposerRequest) (*types.MsgAcceptProposerResponse, error) {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)

	relayer, err := k.Relayer.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	if relayer.Proposer != req.Proposer {
		return nil, types.ErrInvalid.Wrapf("not the current proposer")
	}

	if relayer.ProposerAccepted {
		return nil, types.ErrInvalid.Wrapf("proposer has been accepted")
	}

	if relayer.Epoch != req.Epoch {
		return nil, types.ErrInvalid.Wrapf("invalid epoch: expected %d", relayer.Epoch)
	}

	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return nil, err
	}

	if sdkctx.BlockTime().Sub(relayer.LastElected) > param.AcceptProposerTimeout {
		return nil, types.ErrInvalid.Wrapf("timeout to accept proposer role")
	}

	relayer.ProposerAccepted = true
	if err := k.Relayer.Set(sdkctx, relayer); err != nil {
		return nil, err
	}

	k.Logger().Info("New proposer is accepted", "epoch", relayer.Epoch, "proposer", relayer.Proposer)
	sdkctx.EventManager().EmitEvent(types.AcceptedProposerEvent(relayer.Proposer, relayer.Epoch))
	return &types.MsgAcceptProposerResponse{}, nil
}
