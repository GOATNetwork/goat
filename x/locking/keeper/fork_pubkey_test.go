package keeper_test

import (
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/mldsa65"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/goatnetwork/goat/x/consensusfork"
	"github.com/goatnetwork/goat/x/locking/types"
)

// Adding a key type to the whitelist has no transaction path: GOAT registers
// no x/gov and the ante handler admits nothing that could carry a consensus
// parameter update, so a hardcoded height is the only carrier.
func (suite *KeeperTestSuite) TestPubKeyTypesWidenedAtTheFork() {
	forkHeight := consensusfork.PubKeyTypesForkHeight["unitest"]
	suite.Require().Equal([]string{types.KeyTypeSecp256k1}, suite.ConsensusParams.PubKeyTypes())

	ctx := suite.Context.WithChainID("unitest").
		WithBlockHeight(forkHeight - 1).
		WithBlockTime(time.Now().UTC())
	suite.Require().NoError(suite.Keeper.UpdateForkParams(ctx))
	suite.Require().Equal([]string{types.KeyTypeSecp256k1}, suite.ConsensusParams.PubKeyTypes())

	ctx = ctx.WithBlockHeight(forkHeight)
	suite.Require().NoError(suite.Keeper.UpdateForkParams(ctx))
	suite.Require().Equal(
		[]string{types.KeyTypeSecp256k1, types.KeyTypeMlDsa65},
		suite.ConsensusParams.PubKeyTypes(),
		"secp256k1 has to survive: the validators already in the set still use it",
	)

	// running it again must not append a second time
	suite.Require().NoError(suite.Keeper.UpdateForkParams(ctx))
	suite.Require().Len(suite.ConsensusParams.PubKeyTypes(), 2)
}

// A validator update naming a key type that is off the whitelist is rejected
// by CometBFT, and a rejected update fails the block. Refuse the rotation
// while refusing is still free.
func (suite *KeeperTestSuite) TestRotateToMlDsaBeforeTheWhitelistOpens() {
	f := suite.setupRotation()
	newKey, keyErr := mldsa65.GenPrivKey()
	suite.Require().NoError(keyErr)
	ctx := suite.Context.WithBlockHeight(10).WithBlockTime(time.Now().UTC())

	req := mlDsaRotateRequest(f, ctx.ChainID(), newKey)

	before, err := suite.Keeper.Validators.Get(ctx, f.id)
	suite.Require().NoError(err)

	suite.Require().NoError(suite.Keeper.Rotate(ctx, []*goattypes.RotateRequest{req}))

	after, err := suite.Keeper.Validators.Get(ctx, f.id)
	suite.Require().NoError(err)
	suite.Require().Equal(before, after, "state must be untouched")

	updates, err := suite.Keeper.EndBlocker(ctx)
	suite.Require().NoError(err)
	suite.Require().Empty(updates)
}

// What the rotation mechanism and the widened key-type whitelist exist for: a
// validator moves from secp256k1 to ML-DSA-65 without leaving the set.
func (suite *KeeperTestSuite) TestRotateSecp256k1ToMlDsa() {
	f := suite.setupRotation()
	newKey, keyErr := mldsa65.GenPrivKey()
	suite.Require().NoError(keyErr)

	ctx := suite.Context.WithChainID("unitest").
		WithBlockHeight(consensusfork.PubKeyTypesForkHeight["unitest"]).
		WithBlockTime(time.Now().UTC())
	suite.Require().NoError(suite.Keeper.UpdateForkParams(ctx))

	ctx = ctx.WithBlockHeight(10)
	suite.Require().NoError(suite.Keeper.Rotate(ctx,
		[]*goattypes.RotateRequest{mlDsaRotateRequest(f, ctx.ChainID(), newKey)}))

	updated, err := suite.Keeper.Validators.Get(ctx, f.id)
	suite.Require().NoError(err)
	suite.Require().Equal(types.KeyTypeMlDsa65, updated.KeyType)
	suite.Require().Equal(newKey.PubKey().Bytes(), updated.Pubkey)
	suite.Require().Equal(types.KeyTypeSecp256k1, types.KeyTypeName(updated.PrevKeyType))
	suite.Require().Equal(f.power, updated.Power, "the stake does not move")

	updates, err := suite.Keeper.EndBlocker(ctx)
	suite.Require().NoError(err)
	suite.Require().Len(updates, 2)

	var removed, added = updates[0], updates[1]
	if removed.Power != 0 {
		removed, added = added, removed
	}
	// the old key leaves as a secp256k1 entry, the new one arrives as ml_dsa_65
	suite.Require().Equal(f.oldKey.PubKey().Bytes(), removed.PubKey.GetSecp256K1())
	suite.Require().EqualValues(0, removed.Power)
	suite.Require().Equal(newKey.PubKey().Bytes(), added.PubKey.GetMldsa65())
	suite.Require().EqualValues(f.power, added.Power)

	// and the ML-DSA consensus address resolves back to the unchanged id
	id, err := suite.Keeper.ResolveValidatorID(ctx, sdk.ConsAddress(newKey.PubKey().Address()))
	suite.Require().NoError(err)
	suite.Require().Equal(f.id, id)
}

func mlDsaRotateRequest(f rotationFixture, chainID string, key mldsa65.PrivKey) *goattypes.RotateRequest {
	pubkey := key.PubKey().Bytes()
	proof, err := key.Sign(
		types.RotationProofMessage(chainID, f.id, types.KeyTypeMlDsa65, pubkey))
	if err != nil {
		panic(err)
	}
	return &goattypes.RotateRequest{
		Validator: common.BytesToAddress(f.id),
		KeyType:   types.WireKeyTypeMlDsa65,
		Pubkey:    pubkey,
		Proof:     proof,
	}
}
