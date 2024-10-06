package keeper

import (
	"bytes"
	"context"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k *Keeper) Create(ctx context.Context, req []*ethtypes.CreateValidator) error {
	for _, create := range req {
		if err := k.createValidator(ctx, create); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) createValidator(ctx context.Context, req *ethtypes.CreateValidator) error {
	uncomp, err := ethcrypto.UnmarshalPubkey(append([]byte{0x04}, req.Pubkey[:]...))
	if err != nil {
		return err
	}

	// check if the address provided is consistent with caculation
	pubkey := &secp256k1.PubKey{Key: ethcrypto.CompressPubkey(uncomp)}
	address := sdktypes.ConsAddress(pubkey.Address())
	if !bytes.Equal(address, req.Validator.Bytes()) {
		return types.ErrInvalid.Wrapf("invalid address for pubkey %x: expect %x but got %x",
			req.Pubkey[:], req.Validator.Bytes(), address.Bytes())
	}

	// check if the validator address exists
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	exists, err := k.Validators.Has(sdkctx, address)
	if err != nil {
		return err
	}
	if exists {
		return types.ErrInvalid.Wrapf("validator %x has been created", address.Bytes())
	}

	// check if account for the address exists
	if acc := sdktypes.AccAddress(address); k.accountKeeper.HasAccount(ctx, acc) {
		return types.ErrInvalid.Wrapf("account %x has been created", address.Bytes())
	} else {
		acc := k.accountKeeper.NewAccountWithAddress(ctx, acc)
		if err := acc.SetPubKey(&secp256k1.PubKey{Key: pubkey.Key}); err != nil {
			return types.ErrInvalid.Wrapf("unable to set pubkey")
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
		Status:    types.ValidatorStatus_Pending,
	}
	if err := k.Validators.Set(sdkctx, address, validator); err != nil {
		return err
	}

	return nil
}
