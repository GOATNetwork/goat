package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/goatnetwork/goat/x/locking/types"
)

func (suite *KeeperTestSuite) TestCreates() {
	var creates []*goattypes.CreateRequest

	for idx, validator := range suite.Validator {
		pubkey, err := ethcrypto.DecompressPubkey(validator.Pubkey)
		suite.Require().NoError(err)
		suite.Require().NotNil(pubkey)

		uncompressed := ethcrypto.FromECDSAPub(pubkey)
		suite.Require().Len(uncompressed, 65)

		account, err := authtypes.NewBaseAccountWithPubKey(&secp256k1.PubKey{Key: validator.Pubkey})
		suite.Require().NoError(err)

		creates = append(creates, &goattypes.CreateRequest{
			Validator: common.BytesToAddress(account.GetAddress()),
			Pubkey:    [64]byte(uncompressed[1:]),
		})

		if idx < 2 {
			suite.Account.EXPECT().HasAccount(suite.Context, account.GetAddress()).Return(false)
			suite.Account.EXPECT().NewAccountWithAddress(suite.Context, account.GetAddress()).Return(account)
			suite.Account.EXPECT().SetAccount(suite.Context, account)
		} else {
			suite.Account.EXPECT().HasAccount(suite.Context, account.GetAddress()).Return(true)
		}

		if idx == 0 {
			creates = append(creates, &goattypes.CreateRequest{
				Validator: common.BytesToAddress(account.GetAddress()),
				Pubkey:    [64]byte(uncompressed[1:]),
			})
		}
	}
	err := suite.Keeper.Create(suite.Context, creates)
	suite.Require().NoError(err)

	for idx, address := range suite.Address {
		validator, err := suite.Keeper.Validators.Get(suite.Context, address)
		suite.Require().NoError(err)

		status := types.Inactive
		if idx < 2 {
			status = types.Pending
		}

		suite.Require().Equal(validator, types.Validator{
			Pubkey:    suite.Validator[idx].Pubkey,
			Reward:    math.ZeroInt(),
			GasReward: math.ZeroInt(),
			Status:    status,
		})
	}
}
