package modgen

import (
	"bytes"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/goatnetwork/goat/x/locking/types"
	"github.com/spf13/cobra"
)

func Validator() *cobra.Command {
	const (
		FlagValidatorPubkey = "pubkey"
		FlagValidatorPower  = "power"
	)

	cmd := &cobra.Command{
		Use:   "validator",
		Short: "append a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			votePower, err := cmd.Flags().GetUint64(FlagValidatorPower)
			if err != nil {
				return err
			}

			pubkeyStr, err := cmd.Flags().GetString(FlagValidatorPubkey)
			if err != nil {
				return err
			}

			pubkeyRaw, err := DecodeHexOrBase64String(pubkeyStr)
			if err != nil {
				return err
			}

			if err := IsValidSecp256Pubkey(pubkeyRaw); err != nil {
				return err
			}

			serverCtx.Logger.Info("update genesis", "module", types.ModuleName, "geneis", genesisFile)
			if err := UpdateModuleGenesis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				for _, v := range genesis.GetValidators() {
					if bytes.Equal(v.Pubkey, pubkeyRaw) {
						return nil
					}
				}
				genesis.Validators = append(genesis.Validators, types.Validator{
					Pubkey:    pubkeyRaw,
					Power:     votePower,
					Reward:    math.ZeroInt(),
					GasReward: math.ZeroInt(),
					Status:    types.ValidatorStatus_Active,
				})
				return nil
			}); err != nil {
				return err
			}

			serverCtx.Logger.Info("update genesis", "module", authtypes.ModuleName, "geneis", genesisFile)
			// Add the validator account to auth module to allow it sending tx
			return UpdateModuleGenesis(genesisFile, authtypes.ModuleName, new(authtypes.GenesisState), clientCtx.Codec, func(genesis *authtypes.GenesisState) error {
				pubkey := &secp256k1.PubKey{Key: pubkeyRaw}
				baseAccount, err := authtypes.NewBaseAccountWithPubKey(pubkey)
				if err != nil {
					return err
				}

				if err := genesis.UnpackInterfaces(clientCtx.Codec); err != nil {
					return err
				}

				for _, v := range genesis.GetAccounts() {
					var acc authtypes.GenesisAccount
					if err := clientCtx.Codec.UnpackAny(v, &acc); err != nil {
						return err
					}

					if bytes.Equal(acc.GetAddress().Bytes(), pubkey.Address()) {
						return nil
					}
				}

				if err := baseAccount.SetAccountNumber(uint64(len(genesis.GetAccounts()))); err != nil {
					return err
				}

				baseAccountAny, err := codectypes.NewAnyWithValue(baseAccount)
				if err != nil {
					return err
				}
				genesis.Accounts = append(genesis.Accounts, baseAccountAny)
				return nil
			})
		},
	}

	cmd.Flags().String(FlagValidatorPubkey, "", "validator pubkey(compressed secp256k1)")
	cmd.Flags().Uint64(FlagValidatorPower, 1, "validator vote power")

	return cmd
}
