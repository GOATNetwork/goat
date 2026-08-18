package keeper_test

import (
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/consensusfork"
	"github.com/goatnetwork/goat/x/locking/types"
)

// rotationFixture puts one Active validator in the set, keyed by the consensus
// address it was created with, and hands back the key it would rotate to.
type rotationFixture struct {
	id     sdk.ConsAddress
	oldKey *secp256k1.PrivKey
	newKey *secp256k1.PrivKey
	power  uint64
}

func (suite *KeeperTestSuite) setupRotation() rotationFixture {
	oldKey := secp256k1.GenPrivKey()
	newKey := secp256k1.GenPrivKey()
	id := sdk.ConsAddress(oldKey.PubKey().Address())

	validator := types.Validator{
		Pubkey:    oldKey.PubKey().Bytes(),
		Power:     10000,
		Reward:    math.ZeroInt(),
		GasReward: math.ZeroInt(),
		Status:    types.Active,
		Locking: sdk.NewCoins(
			sdk.NewCoin(NativeTokenDenom, math.NewIntFromUint64(1e18)),
		),
	}

	suite.Require().NoError(suite.Keeper.Validators.Set(suite.Context, id, validator))
	suite.Require().NoError(suite.Keeper.ValidatorSet.Set(suite.Context, id, validator.Power))
	suite.Require().NoError(suite.Keeper.PowerRanking.Set(suite.Context,
		collections.Join(validator.Power, id)))
	for _, locking := range validator.Locking {
		suite.Require().NoError(suite.Keeper.Locking.Set(suite.Context,
			collections.Join(locking.Denom, id), locking.Amount))
	}
	suite.Require().NoError(suite.Keeper.Params.Set(suite.Context, suite.Param))

	return rotationFixture{id: id, oldKey: oldKey, newKey: newKey, power: validator.Power}
}

func (f rotationFixture) request(chainID string, key *secp256k1.PrivKey) *goattypes.RotateRequest {
	pubkey := key.PubKey().Bytes()
	msg := types.RotationProofMessage(chainID, f.id, types.KeyTypeSecp256k1, pubkey)
	proof, err := key.Sign(msg)
	if err != nil {
		panic(err)
	}
	return &goattypes.RotateRequest{
		Validator: common.BytesToAddress(f.id),
		KeyType:   types.WireKeyTypeSecp256k1,
		Pubkey:    pubkey,
		Proof:     proof,
	}
}

func (suite *KeeperTestSuite) TestRotateStagesTheNewKey() {
	f := suite.setupRotation()
	ctx := suite.Context.WithBlockHeight(10).WithBlockTime(time.Now().UTC())

	suite.Require().NoError(suite.Keeper.Rotate(ctx,
		[]*goattypes.RotateRequest{f.request(ctx.ChainID(), f.newKey)}))

	updated, err := suite.Keeper.Validators.Get(ctx, f.id)
	suite.Require().NoError(err)
	suite.Require().Equal(f.newKey.PubKey().Bytes(), updated.Pubkey)
	suite.Require().Equal(f.oldKey.PubKey().Bytes(), updated.PrevPubkey)
	suite.Require().EqualValues(10, updated.RotationApplyHeight)
	suite.Require().True(updated.IsRotating())
	// the id is untouched, and so is everything hanging off it
	suite.Require().Equal(f.power, updated.Power)

	// the new consensus address now resolves to the id, and so does the old
	newAddr := sdk.ConsAddress(f.newKey.PubKey().Address())
	got, err := suite.Keeper.ResolveValidatorID(ctx, newAddr)
	suite.Require().NoError(err)
	suite.Require().Equal(f.id, got)

	got, err = suite.Keeper.ResolveValidatorID(ctx, f.id)
	suite.Require().NoError(err)
	suite.Require().Equal(f.id, got)
}

// A rotation leaves the power alone, and updateValidatorSet only emitted an
// update when the power changed. Without an explicit trigger the state would
// hold the new key while CometBFT kept the old one, and the validator would
// miss every block with nothing in the logs to say why.
func (suite *KeeperTestSuite) TestRotateEmitsBothValidatorUpdates() {
	f := suite.setupRotation()
	ctx := suite.Context.WithBlockHeight(10).WithBlockTime(time.Now().UTC())

	suite.Require().NoError(suite.Keeper.Rotate(ctx,
		[]*goattypes.RotateRequest{f.request(ctx.ChainID(), f.newKey)}))

	updates, err := suite.Keeper.EndBlocker(ctx)
	suite.Require().NoError(err)
	suite.Require().Len(updates, 2, "expected the old key removed and the new key added")

	var removed, added *abci.ValidatorUpdate
	for i := range updates {
		if updates[i].Power == 0 {
			removed = &updates[i]
		} else {
			added = &updates[i]
		}
	}
	suite.Require().NotNil(removed)
	suite.Require().NotNil(added)
	suite.Require().Equal(f.oldKey.PubKey().Bytes(), removed.PubKey.GetSecp256K1())
	suite.Require().Equal(f.newKey.PubKey().Bytes(), added.PubKey.GetSecp256K1())
	suite.Require().EqualValues(f.power, added.Power)

	// once emitted the staging fields are cleared, but the retention deadline
	// for the old consensus address is not
	settled, err := suite.Keeper.Validators.Get(ctx, f.id)
	suite.Require().NoError(err)
	suite.Require().False(settled.IsRotating())
	suite.Require().Empty(settled.PrevPubkey)
	suite.Require().False(settled.PrevConsAddrExpiry.IsZero())
}

func (suite *KeeperTestSuite) TestRotateSoftFailures() {
	chainID := suite.Context.ChainID()

	cases := map[string]func(f rotationFixture) *goattypes.RotateRequest{
		"unknown validator": func(f rotationFixture) *goattypes.RotateRequest {
			req := f.request(chainID, f.newKey)
			req.Validator = common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
			return req
		},
		"proof signed by the wrong key": func(f rotationFixture) *goattypes.RotateRequest {
			req := f.request(chainID, f.newKey)
			bad, err := f.oldKey.Sign(types.RotationProofMessage(chainID, f.id, types.KeyTypeSecp256k1, req.Pubkey))
			suite.Require().NoError(err)
			req.Proof = bad
			return req
		},
		"proof for another chain": func(f rotationFixture) *goattypes.RotateRequest {
			req := f.request("some-other-chain", f.newKey)
			return req
		},
		"unknown key type": func(f rotationFixture) *goattypes.RotateRequest {
			req := f.request(chainID, f.newKey)
			req.KeyType = 99
			return req
		},
		"malformed pubkey": func(f rotationFixture) *goattypes.RotateRequest {
			req := f.request(chainID, f.newKey)
			req.Pubkey = req.Pubkey[:16]
			return req
		},
		"rotating to the key already in use": func(f rotationFixture) *goattypes.RotateRequest {
			return f.request(chainID, f.oldKey)
		},
	}

	for name, build := range cases {
		suite.Run(name, func() {
			suite.SetupTest()
			f := suite.setupRotation()
			ctx := suite.Context.WithBlockHeight(10).WithBlockTime(time.Now().UTC())

			before, err := suite.Keeper.Validators.Get(ctx, f.id)
			suite.Require().NoError(err)

			// soft: no error, so MsgNewEthBlock and therefore the block survives
			suite.Require().NoError(suite.Keeper.Rotate(ctx, []*goattypes.RotateRequest{build(f)}))

			after, err := suite.Keeper.Validators.Get(ctx, f.id)
			suite.Require().NoError(err)
			suite.Require().Equal(before, after, "state must be untouched")

			updates, err := suite.Keeper.EndBlocker(ctx)
			suite.Require().NoError(err)
			suite.Require().Empty(updates, "a rejected rotation must not reach CometBFT")
		})
	}
}

func (suite *KeeperTestSuite) TestRotateRejectsASecondRotationInFlight() {
	f := suite.setupRotation()
	ctx := suite.Context.WithBlockHeight(10).WithBlockTime(time.Now().UTC())

	suite.Require().NoError(suite.Keeper.Rotate(ctx,
		[]*goattypes.RotateRequest{f.request(ctx.ChainID(), f.newKey)}))

	third := secp256k1.GenPrivKey()
	suite.Require().NoError(suite.Keeper.Rotate(ctx,
		[]*goattypes.RotateRequest{f.request(ctx.ChainID(), third)}))

	updated, err := suite.Keeper.Validators.Get(ctx, f.id)
	suite.Require().NoError(err)
	suite.Require().Equal(f.newKey.PubKey().Bytes(), updated.Pubkey, "the first rotation stands")

	has, err := suite.Keeper.ConsAddrIndex.Has(ctx, sdk.ConsAddress(third.PubKey().Address()))
	suite.Require().NoError(err)
	suite.Require().False(has)
}

// The old consensus address has to stay resolvable for as long as evidence
// naming it can still be slashed, and not one block longer than necessary.
func (suite *KeeperTestSuite) TestConsAddrIndexRetention() {
	f := suite.setupRotation()
	now := time.Now().UTC()
	ctx := suite.Context.WithBlockHeight(10).WithBlockTime(now)

	suite.Require().NoError(suite.Keeper.Rotate(ctx,
		[]*goattypes.RotateRequest{f.request(ctx.ChainID(), f.newKey)}))
	_, err := suite.Keeper.EndBlocker(ctx)
	suite.Require().NoError(err)

	newAddr := sdk.ConsAddress(f.newKey.PubKey().Address())

	// still inside the window: the old address resolves and the new one does too
	ctx = ctx.WithBlockHeight(11).WithBlockTime(now.Add(suite.Param.ExitingDuration - time.Minute))
	_, err = suite.Keeper.EndBlocker(ctx)
	suite.Require().NoError(err)

	has, err := suite.Keeper.ConsAddrIndex.Has(ctx, newAddr)
	suite.Require().NoError(err)
	suite.Require().True(has, "the address in effect is never pruned")

	// past the window: the entry for the address no longer in effect goes away
	ctx = ctx.WithBlockHeight(12).WithBlockTime(now.Add(suite.Param.ExitingDuration + time.Minute))
	_, err = suite.Keeper.EndBlocker(ctx)
	suite.Require().NoError(err)

	has, err = suite.Keeper.ConsAddrIndex.Has(ctx, newAddr)
	suite.Require().NoError(err)
	suite.Require().True(has, "the address in effect is still never pruned")

	id, err := suite.Keeper.ResolveValidatorID(ctx, newAddr)
	suite.Require().NoError(err)
	suite.Require().Equal(f.id, id)
}

// The id is itself a consensus address that ResolveValidatorID hands back
// unchanged, so the first rotation never leaves anything stale behind. Only a
// second rotation strands an address: the one from the first rotation stops
// being in effect while evidence naming it is still slashable.
func (suite *KeeperTestSuite) TestConsAddrIndexPrunesAStrandedAddress() {
	f := suite.setupRotation()
	now := time.Now().UTC()

	ctx := suite.Context.WithBlockHeight(10).WithBlockTime(now)
	suite.Require().NoError(suite.Keeper.Rotate(ctx,
		[]*goattypes.RotateRequest{f.request(ctx.ChainID(), f.newKey)}))
	_, err := suite.Keeper.EndBlocker(ctx)
	suite.Require().NoError(err)
	middle := sdk.ConsAddress(f.newKey.PubKey().Address())

	third := secp256k1.GenPrivKey()
	ctx = ctx.WithBlockHeight(11).WithBlockTime(now.Add(time.Minute))
	suite.Require().NoError(suite.Keeper.Rotate(ctx,
		[]*goattypes.RotateRequest{f.request(ctx.ChainID(), third)}))
	_, err = suite.Keeper.EndBlocker(ctx)
	suite.Require().NoError(err)
	latest := sdk.ConsAddress(third.PubKey().Address())

	// inside the window all three addresses lead to the same validator, which
	// is what keeps evidence against any of them slashable
	ctx = ctx.WithBlockHeight(12).WithBlockTime(now.Add(suite.Param.ExitingDuration))
	_, err = suite.Keeper.EndBlocker(ctx)
	suite.Require().NoError(err)
	for _, addr := range []sdk.ConsAddress{f.id, middle, latest} {
		id, err := suite.Keeper.ResolveValidatorID(ctx, addr)
		suite.Require().NoError(err)
		suite.Require().Equal(f.id, id)
	}

	// past it the stranded one goes, and only that one
	ctx = ctx.WithBlockHeight(13).WithBlockTime(now.Add(suite.Param.ExitingDuration + time.Hour))
	_, err = suite.Keeper.EndBlocker(ctx)
	suite.Require().NoError(err)

	has, err := suite.Keeper.ConsAddrIndex.Has(ctx, middle)
	suite.Require().NoError(err)
	suite.Require().False(has, "the stranded address should have been pruned")

	has, err = suite.Keeper.ConsAddrIndex.Has(ctx, latest)
	suite.Require().NoError(err)
	suite.Require().True(has, "the address in effect must survive")

	id, err := suite.Keeper.ResolveValidatorID(ctx, f.id)
	suite.Require().NoError(err)
	suite.Require().Equal(f.id, id, "the id always resolves to itself")
}

// Nothing can produce a rotation before the execution layer carries a Locking
// contract that has the function, but the consensus layer pins its own
// behavior to a height rather than trusting that.
func (suite *KeeperTestSuite) TestRotateIgnoredBeforeTheFork() {
	f := suite.setupRotation()
	ctx := suite.Context.WithChainID("unitest").
		WithBlockHeight(consensusfork.RotationForkHeight["unitest"] - 1).
		WithBlockTime(time.Now().UTC())

	before, err := suite.Keeper.Validators.Get(ctx, f.id)
	suite.Require().NoError(err)

	suite.Require().NoError(suite.Keeper.Rotate(ctx,
		[]*goattypes.RotateRequest{f.request(ctx.ChainID(), f.newKey)}))

	after, err := suite.Keeper.Validators.Get(ctx, f.id)
	suite.Require().NoError(err)
	suite.Require().Equal(before, after)

	// and it works from the fork height on
	ctx = ctx.WithBlockHeight(consensusfork.RotationForkHeight["unitest"])
	suite.Require().NoError(suite.Keeper.Rotate(ctx,
		[]*goattypes.RotateRequest{f.request(ctx.ChainID(), f.newKey)}))

	rotated, err := suite.Keeper.Validators.Get(ctx, f.id)
	suite.Require().NoError(err)
	suite.Require().True(rotated.IsRotating())
}
