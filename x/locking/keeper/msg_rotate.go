package keeper

import (
	"bytes"
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (k *Keeper) Rotate(ctx context.Context, req []*goattypes.RotateRequest) error {
	for _, rotate := range req {
		if err := k.rotateValidator(ctx, rotate); err != nil {
			return err
		}
	}
	return nil
}

// rotateValidator records a validator's intent to change its consensus public
// key. The change is staged here and only reaches CometBFT in the EndBlocker of
// the same height, so that the removal of the old key and the addition of the
// new one leave together.
//
// Every rejection here is soft. This request came from an EVM log, and the
// execution layer cannot check any of the things checked below: it has no
// ml-dsa precompile, no view of the validator set, and no view of the consensus
// address index. Returning an error would fail ProcessLockingRequest, which
// fails MsgNewEthBlock, which is the block itself. The only errors returned are
// storage errors.
func (k Keeper) rotateValidator(ctx context.Context, req *goattypes.RotateRequest) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	id := sdktypes.ConsAddress(req.Validator.Bytes())
	name := types.ValidatorName(id)

	validator, err := k.Validators.Get(sdkctx, id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			k.Logger().Warn("Rotate: unknown validator", "validator", name)
			return nil
		}
		return err
	}

	if validator.IsRotating() {
		k.Logger().Warn("Rotate: rotation already in flight", "validator", name,
			"apply_height", validator.RotationApplyHeight)
		return nil
	}

	switch validator.Status {
	case types.Active, types.Pending, types.Downgrade:
	default:
		k.Logger().Warn("Rotate: validator status does not allow rotation",
			"validator", name, "status", validator.Status)
		return nil
	}

	keyType, err := types.KeyTypeFromRequest(req.KeyType)
	if err != nil {
		k.Logger().Warn("Rotate: unsupported key type", "validator", name, "error", err)
		return nil
	}

	newAddr, err := types.ConsAddress(keyType, req.Pubkey)
	if err != nil {
		k.Logger().Warn("Rotate: invalid pubkey", "validator", name, "error", err)
		return nil
	}

	if types.SameConsensusKey(&validator, keyType, req.Pubkey) {
		k.Logger().Warn("Rotate: new key equals the current one", "validator", name)
		return nil
	}

	// The proof is the only evidence that whoever sent the EVM transaction
	// actually holds the new private key. create checks the equivalent with
	// ECDSA.recover on chain; there is no precompile to do that here.
	if err := types.VerifyRotationProof(sdkctx.ChainID(), id, keyType, req.Pubkey, req.Proof); err != nil {
		k.Logger().Warn("Rotate: invalid proof of possession", "validator", name, "error", err)
		return nil
	}

	// Refuse to point two validators at one consensus address. Whoever holds
	// it first keeps it.
	taken, err := k.ConsAddrIndex.Has(sdkctx, newAddr)
	if err != nil {
		return err
	}
	// A validator id is itself a consensus address that ResolveValidatorID
	// hands back unchanged, so taking one over would silently redirect that
	// validator's votes and evidence to us. Our own id is fine: the index entry
	// would just be the identity mapping the fallback already provides.
	if !taken && !bytes.Equal(newAddr, id) {
		taken, err = k.Validators.Has(sdkctx, newAddr)
		if err != nil {
			return err
		}
	}
	if taken {
		k.Logger().Warn("Rotate: consensus address already in use", "validator", name,
			"address", types.ValidatorName(newAddr))
		return nil
	}

	param, err := k.Params.Get(sdkctx)
	if err != nil {
		return err
	}

	validator.PrevPubkey, validator.PrevKeyType = validator.Pubkey, validator.KeyType
	validator.Pubkey, validator.KeyType = req.Pubkey, keyType
	validator.RotationApplyHeight = sdkctx.BlockHeight()
	// The old consensus address has to stay resolvable for as long as evidence
	// naming it can still be slashed. Binding this to ExitingDuration rather
	// than to a snapshot of the evidence window means a later parameter change
	// moves the deadline with it.
	validator.PrevConsAddrExpiry = sdkctx.BlockTime().Add(param.ExitingDuration)

	if err := k.ConsAddrIndex.Set(sdkctx, newAddr, id); err != nil {
		return err
	}
	if err := k.Validators.Set(sdkctx, id, validator); err != nil {
		return err
	}

	k.Logger().Info("Rotate", "validator", name, "key_type", keyType,
		"new_address", types.ValidatorName(newAddr), "apply_height", validator.RotationApplyHeight)
	sdkctx.EventManager().EmitEvent(types.ValidatorRotatedEvent(id, newAddr))
	return nil
}

// pruneConsAddrIndex drops the index entries for consensus addresses whose
// retention has run out. Dropping one early halts the chain, so this only runs
// once the expiry has passed.
func (k Keeper) pruneConsAddrIndex(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	now := sdkctx.BlockTime()

	iter, err := k.ConsAddrIndex.Iterate(sdkctx, nil)
	if err != nil {
		return err
	}
	defer iter.Close()

	type entry struct{ consAddr, id sdktypes.ConsAddress }
	var stale []entry

	for ; iter.Valid(); iter.Next() {
		kv, err := iter.KeyValue()
		if err != nil {
			return err
		}
		// the entry for the key currently in effect is never stale
		validator, err := k.Validators.Get(sdkctx, kv.Value)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				stale = append(stale, entry{kv.Key, kv.Value})
				continue
			}
			return err
		}
		current, err := types.ConsAddress(validator.KeyType, validator.Pubkey)
		if err != nil {
			return err
		}
		if bytes.Equal(current, kv.Key) {
			continue
		}
		if validator.PrevConsAddrExpiry.After(now) {
			continue
		}
		stale = append(stale, entry{kv.Key, kv.Value})
	}

	for _, e := range stale {
		if err := k.ConsAddrIndex.Remove(sdkctx, e.consAddr); err != nil {
			return err
		}
		k.Logger().Info("Drop stale consensus address", "address", types.ValidatorName(e.consAddr),
			"validator", types.ValidatorName(e.id))
	}
	return nil
}

// settleRotations clears the staging fields of every rotation whose validator
// update has just been emitted. PrevConsAddrExpiry is deliberately left alone:
// it governs how long the old consensus address stays resolvable, which has
// nothing to do with whether CometBFT has been told yet.
func (k Keeper) settleRotations(ctx context.Context) error {
	sdkctx := sdktypes.UnwrapSDKContext(ctx)
	height := sdkctx.BlockHeight()

	iter, err := k.Validators.Iterate(sdkctx, nil)
	if err != nil {
		return err
	}
	defer iter.Close()

	type settled struct {
		id        sdktypes.ConsAddress
		validator types.Validator
	}
	var done []settled

	for ; iter.Valid(); iter.Next() {
		kv, err := iter.KeyValue()
		if err != nil {
			return err
		}
		if !kv.Value.IsRotating() || kv.Value.RotationApplyHeight != height {
			continue
		}
		v := kv.Value
		v.PrevPubkey, v.PrevKeyType = nil, ""
		v.RotationApplyHeight = 0
		done = append(done, settled{kv.Key, v})
	}

	for _, s := range done {
		if err := k.Validators.Set(sdkctx, s.id, s.validator); err != nil {
			return err
		}
	}
	return nil
}
