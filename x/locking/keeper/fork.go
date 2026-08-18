package keeper

import (
	"context"
	"slices"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/goatnetwork/goat/x/consensusfork"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k Keeper) UpdateForkParams(sdkctx sdktypes.Context) error {
	if osakaHeight, ok := consensusfork.OsakaForkHeight[sdkctx.ChainID()]; ok && sdkctx.BlockHeight() == osakaHeight {
		param, err := k.Params.Get(sdkctx)
		if err != nil {
			return err
		}
		// 2.186 Goat per block
		param.InitialBlockReward = 2186000000000000000
		if err := k.Params.Set(sdkctx, param); err != nil {
			return err
		}
	}

	if height, ok := consensusfork.PubKeyTypesForkHeight[sdkctx.ChainID()]; ok && sdkctx.BlockHeight() == height {
		if err := k.widenPubKeyTypes(sdkctx); err != nil {
			return err
		}
	}
	return nil
}

// widenPubKeyTypes adds ml_dsa_65 to the validator pub_key_types whitelist.
//
// CometBFT rejects a validator update naming a key type that is not on the
// list, and a rejected update fails the block, so this has to be in place
// before the first rotation to a post quantum key. Adding a type does not
// touch the validators already in the set: existing entries keep their
// secp256k1 keys and the two coexist for as long as it takes.
//
// There is no transaction that could do this. GOAT registers neither x/gov nor
// x/upgrade, and app/ante.go admits only bitcoin, relayer and MsgNewEthBlock
// messages, so cosmos.consensus.v1.MsgUpdateParams cannot even enter the
// mempool. A hardcoded height is the only carrier.
func (k Keeper) widenPubKeyTypes(sdkctx sdktypes.Context) error {
	params, err := k.consensusKeeper.Get(sdkctx)
	if err != nil {
		return err
	}
	if params.Validator == nil {
		params.Validator = &cmtproto.ValidatorParams{}
	}
	if slices.Contains(params.Validator.PubKeyTypes, types.KeyTypeMlDsa65) {
		return nil
	}

	params.Validator.PubKeyTypes = append(
		slices.Clone(params.Validator.PubKeyTypes), types.KeyTypeMlDsa65)

	// baseapp hands whatever the store holds back to CometBFT in every
	// FinalizeBlock response, so writing here is what publishes the change
	if err := k.consensusKeeper.Set(sdkctx, params); err != nil {
		return err
	}

	k.Logger().Info("Widen the consensus pubkey whitelist",
		"pub_key_types", params.Validator.PubKeyTypes, "height", sdkctx.BlockHeight())
	return nil
}

// PubKeyTypeAllowed reports whether a consensus key type may be handed to
// CometBFT yet.
func (k Keeper) PubKeyTypeAllowed(ctx context.Context, keyType string) (bool, error) {
	params, err := k.consensusKeeper.Get(ctx)
	if err != nil {
		return false, err
	}
	if params.Validator == nil {
		return false, nil
	}
	return slices.Contains(params.Validator.PubKeyTypes, types.KeyTypeName(keyType)), nil
}
