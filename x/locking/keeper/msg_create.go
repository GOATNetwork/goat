package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k *Keeper) Create(ctx context.Context, req []*goattypes.CreateRequest) error {
	for _, create := range req {
		if err := k.createValidator(ctx, create); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) createValidator(ctx context.Context, req *goattypes.CreateRequest) error {
	// check if the address provided is consistent with caculation
	pubkey := &secp256k1.PubKey{Key: goatcrypto.CompressP256k1Pubkey(req.Pubkey)}
	address := sdktypes.ConsAddress(pubkey.Address())
	if !bytes.Equal(address, req.Validator.Bytes()) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid address for pubkey %x: expect %x but got %x",
			req.Pubkey[:], req.Validator.Bytes(), address.Bytes())
	}

	// check if the validator address exists
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	exists, err := k.Validators.Has(sdkctx, address)
	if err != nil {
		return err
	}
	if exists {
		k.Logger().Warn("validator %x has been created", address.Bytes())
		return nil
	}

	// check if account for the validator exists
	acc := sdktypes.AccAddress(address)
	hasAccount := k.accountKeeper.HasAccount(ctx, acc)
	if !hasAccount {
		acc := k.accountKeeper.NewAccountWithAddress(ctx, acc)
		if err := acc.SetPubKey(&secp256k1.PubKey{Key: pubkey.Key}); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrLogic, "unable to set pubkey")
		}
		k.accountKeeper.SetAccount(ctx, acc)
	}

	// save
	validator := types.Validator{
		Pubkey:    pubkey.Key,
		Power:     0,
		Locking:   nil,
		Reward:    math.ZeroInt(),
		GasReward: math.ZeroInt(),
		Status:    types.Pending,
	}

	// Don't allow conflict
	// case for a relayer voter is using the same account
	if hasAccount {
		k.Logger().Warn("validator account %x has been created", address.Bytes())
		validator.Status = types.Inactive
	}

	if err := k.Validators.Set(sdkctx, address, validator); err != nil {
		return err
	}

	k.Logger().Info("Create", "address", types.ValidatorName(address))
	return nil
}
