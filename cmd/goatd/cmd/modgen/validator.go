package modgen

import (
	"bytes"
	"fmt"

	cmt256k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
)

func Validator() *cobra.Command {
	const (
		FlagValidatorPubkey = "pubkey"
		FlagValidatorPower  = "power"
	)

	cmd := &cobra.Command{
		Use:   "validator name",
		Short: "append a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			validatorName := args[0]
			votePower, err := cmd.Flags().GetInt64(FlagValidatorPower)
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

			serverCtx.Logger.Info("update genesis", "module", "cometbft", "geneis", genesisFile)
			if err := UpdateGenesis(genesisFile, func(state *genutiltypes.AppGenesis) error {
				pubkey := cmt256k1.PubKey(pubkeyRaw)
				address := pubkey.Address()
				for _, validator := range state.Consensus.Validators {
					if validator.Name == validatorName {
						return fmt.Errorf("validator %s has been added", validatorName)
					}
					if bytes.Equal(validator.Address.Bytes(), address.Bytes()) {
						return fmt.Errorf("conflict pubkey with validator %s", validator.Name)
					}
				}

				state.Consensus.Validators = append(state.Consensus.Validators,
					cmttypes.GenesisValidator{
						Name:    validatorName,
						PubKey:  pubkey,
						Address: address,
						Power:   votePower,
					},
				)
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
	cmd.Flags().Int64(FlagValidatorPower, 1, "validator vote power")

	return cmd
}
