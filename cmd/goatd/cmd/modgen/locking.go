package modgen

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/goatnetwork/goat/x/locking/types"
	"github.com/spf13/cobra"
)

func Locking() *cobra.Command {
	const (
		FlagValidatorPubkey = "pubkey"

		FlagTokenAddress   = "token"
		FlagTokenWeight    = "weight"
		FlagTokenThreshold = "threshold"
	)

	cmd := &cobra.Command{
		Use:   "locking",
		Short: "update locking module genesis",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	addValidator := &cobra.Command{
		Use:   "add-validator",
		Short: "append a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

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

				var locking sdktypes.Coins
				var votePower uint64
				for _, token := range genesis.Tokens {
					if !token.Token.Threshold.IsZero() {
						locking = locking.Add(sdktypes.NewCoin(token.Denom, token.Token.Threshold))
						votePower += math.NewIntFromUint64(token.Token.Weight).Mul(token.Token.Threshold).Quo(types.PowerReduction).Uint64()
					}
				}

				genesis.Validators = append(genesis.Validators, types.Validator{
					Pubkey:    pubkeyRaw,
					Power:     votePower,
					Reward:    math.ZeroInt(),
					GasReward: math.ZeroInt(),
					Status:    types.Active,
					Locking:   locking,
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

	addToken := &cobra.Command{
		Use:   "add-token",
		Short: "locking module genesis",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			tokenAddress, err := cmd.Flags().GetBytesHex(FlagTokenAddress)
			if err != nil {
				return err
			}

			weight, err := cmd.Flags().GetUint64(FlagTokenWeight)
			if err != nil {
				return err
			}

			shareStr, err := cmd.Flags().GetString(FlagTokenThreshold)
			if err != nil {
				return err
			}

			var ok bool
			var share = new(big.Int)
			if strings.HasPrefix(shareStr, "0x") {
				_, ok = share.SetString(strings.TrimPrefix(shareStr, "0x"), 16)
			} else {
				_, ok = share.SetString(shareStr, 10)
			}
			if !ok {
				return errors.New("invalid share string")
			}

			serverCtx.Logger.Info("update genesis", "module", types.ModuleName, "geneis", genesisFile)
			if err := UpdateModuleGenesis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				address := hex.EncodeToString(tokenAddress)
				for _, token := range genesis.Tokens {
					if token.Denom == address {
						return nil
					}
				}

				genesis.Tokens = append(genesis.Tokens, &types.TokenGenesis{
					Denom: address,
					Token: types.Token{
						Weight:    weight,
						Threshold: math.NewIntFromBigIntMut(share),
					},
				})
				return nil
			}); err != nil {
				return err
			}
			return nil
		},
	}

	addToken.Flags().BytesHex(FlagTokenAddress, nil, "validator pubkey(compressed secp256k1)")
	addToken.Flags().Uint64(FlagTokenWeight, 0, "validator vote power")
	addToken.Flags().String(FlagTokenThreshold, "", "validator vote power")
	addValidator.Flags().String(FlagValidatorPubkey, "", "validator pubkey(compressed secp256k1)")

	cmd.AddCommand(addToken, addValidator)
	return cmd
}