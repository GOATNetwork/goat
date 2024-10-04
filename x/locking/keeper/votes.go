package keeper

import (
	"context"

	"cosmossdk.io/core/comet"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleVoteInfo(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	for _, voteInfo := range sdkctx.VoteInfos() {
		address := sdktypes.ConsAddress(voteInfo.Validator.Address)
		signed := comet.BlockIDFlag(voteInfo.BlockIdFlag)
		if err := k.handleVoteInfo(sdkctx, address, signed); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) handleVoteInfo(ctx context.Context, address sdktypes.ConsAddress, signed comet.BlockIDFlag) error {
	// sdkctx := sdktypes.UnwrapSDKContext(ctx)
	// todo
	return nil
}
