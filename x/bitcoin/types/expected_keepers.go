package types

import (
	"context"

	relayer "github.com/goatnetwork/goat/x/relayer/types"
)

//go:generate mockgen -source=expected_keepers.go -destination=../mock/keeper.go -package=mock

type RelayerKeeper interface {
	VerifyProposal(ctx context.Context, req relayer.IVoteMsg, verifyFn ...func(sigdoc []byte) error) (uint64, error)
	VerifyNonProposal(ctx context.Context, req relayer.INonVoteMsg) (relayer.IRelayer, error)
	UpdateRandao(ctx context.Context, req relayer.IVoteMsg) error
	HasPubkey(ctx context.Context, raw []byte) (bool, error)
	AddNewKey(ctx context.Context, raw []byte) error
	SetProposalSeq(ctx context.Context, seq uint64) error
}
