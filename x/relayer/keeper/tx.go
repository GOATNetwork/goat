package keeper

import (
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

// func (k msgServer) OnBoarding(ctx context.Context, req *types.MsgOnBoardingRequest) (*types.MsgOnBoardingResponse, error) {
// 	panic("todo")
// }
