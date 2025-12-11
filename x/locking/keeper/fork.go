package keeper

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) UpdateForkParams(sdkctx sdktypes.Context) error {
	if osakaHeight, ok := types.OsakaForkHeight[sdkctx.ChainID()]; ok && sdkctx.BlockHeight() == osakaHeight {
		param, err := k.Params.Get(sdkctx)
		if err != nil {
			return err
		}
		// 2.186 Goat per block
		param.InitialBlockReward = 2186000000000000000
		return k.Params.Set(sdkctx, param)
	}
	return nil
}
