package keeper

import (
	"context"

	"github.com/goatnetwork/goat/x/goat/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) NewEthBlock(ctx context.Context, req *types.MsgNewEthBlock) (*types.MsgNewEthBlockResponse, error) {
	panic("todo")
}
