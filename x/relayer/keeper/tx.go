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

	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return nil, err
	}

	if relayer.Proposer != req.Proposer {
		return nil, types.ErrInvalid.Wrapf("not the current proposer")
	}

	address := sdktypes.AccAddress(
		(&secp256k1.PubKey{Key: req.VoterTxKey}).Address().Bytes())
	voter, err := k.Voters.Get(ctx, address)
	if err != nil {
		return nil, err
	}

	if voter.Status != types.Pending {
		return nil, types.ErrInvalid.Wrapf("not a pending voter")
	}

	if !bytes.Equal(goatcrypto.SHA256Sum(req.VoterBlsKey), voter.VoteKey) {
		return nil, types.ErrInvalid.Wrapf("vote key hash not match")
	}

	sigMsg := req.SignDoc(sdkctx.ChainID(), voter.Height, address, voter.VoteKey)
	if !ethcrypto.VerifySignature(req.VoterTxKey, sigMsg, req.VoterTxKeyProof) {
		return nil, types.ErrInvalid.Wrapf("false tx key proof")
	}

	blsPubKey := new(goatcrypto.PublicKey).Uncompress(req.VoterBlsKey)
	if !goatcrypto.Verify(blsPubKey, sigMsg, req.VoterBlsKeyProof) {
		return nil, types.ErrInvalid.Wrapf("false vote key proof")
	}

	voter.VoteKey = req.VoterBlsKey
	voter.Status = types.Activated
	if err := k.Voters.Set(ctx, address, voter); err != nil {
		return nil, err
	}

	return &types.MsgNewVoterResponse{}, nil
}

func (k msgServer) AcceptProposer(ctx context.Context, req *types.MsgAcceptProposerRequest) (*types.MsgAcceptProposerResponse, error) {
	param, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if param.AcceptProposerTimeout == 0 {
		return &types.MsgAcceptProposerResponse{}, nil
	}

	relayer, err := k.Relayer.Get(ctx)
	if err != nil {
		return nil, err
	}

	if relayer.Proposer != req.Proposer {
		return nil, types.ErrInvalid.Wrapf("not the current proposer")
	}

	if relayer.ProposerAccepted {
		return nil, types.ErrInvalid.Wrapf("proposer has been accepted")
	}

	if relayer.Version != req.Version {
		return nil, types.ErrInvalid.Wrapf("invalid version: expected: %d", relayer.Version)
	}

	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	if sdkctx.BlockTime().Sub(relayer.LastElected) > param.AcceptProposerTimeout {
		return nil, types.ErrInvalid.Wrapf("timeout to accept proposer role")
	}

	relayer.ProposerAccepted = true

	if err := k.Relayer.Set(ctx, relayer); err != nil {
		return nil, err
	}

	return &types.MsgAcceptProposerResponse{}, nil
}
