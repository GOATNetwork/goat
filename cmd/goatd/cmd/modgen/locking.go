package modgen

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"cosmossdk.io/math"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/privval"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	"github.com/goatnetwork/goat/x/locking/types"
	"github.com/spf13/cobra"
)

func Locking() *cobra.Command {
	const (
		FlagInitialReward  = "init-reward"
		FlagUnlockDuration = "unlock-duration"
		FlagExitDuration   = "exit-duration"

		FlagValidatorPubkey = "pubkey"

		FlagTokenAddress   = "token"
		FlagTokenWeight    = "weight"
		FlagTokenThreshold = "threshold"

		FlagEthChainID = "eth-chain-id"
		FlagOwner      = "owner"
	)

	cmd := &cobra.Command{
		Use:   "locking",
		Short: "update locking module genesis",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			rewardStr, err := cmd.Flags().GetString(FlagInitialReward)
			if err != nil {
				return err
			}

			if rewardStr == "" {
				return fmt.Errorf("no reward param provided")
			}

			reward := math.ZeroInt()
			switch {
			case strings.HasSuffix(rewardStr, "ether"):
				i, ok := new(big.Int).SetString(strings.TrimSuffix(rewardStr, "ether"), 10)
				if !ok {
					return fmt.Errorf("invalid reward string: %s", rewardStr)
				}
				i.Mul(i, big.NewInt(1e18))
				reward = math.NewIntFromBigIntMut(i)
			case strings.HasPrefix(rewardStr, "0x"):
				i, ok := new(big.Int).SetString(strings.TrimPrefix(rewardStr, "0x"), 16)
				if !ok {
					return fmt.Errorf("invalid reward string: %s", rewardStr)
				}
				reward = math.NewIntFromBigIntMut(i)
			default:
				i, ok := new(big.Int).SetString(rewardStr, 10)
				if !ok {
					return fmt.Errorf("invalid reward string: %s", rewardStr)
				}
				reward = math.NewIntFromBigIntMut(i)
			}

			serverCtx.Logger.Info("update param", "module", types.ModuleName, "geneis", genesisFile)
			return UpdateModuleGenesis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				genesis.RewardPool.Remain = reward
				return nil
			})
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

			fmt.Println(hex.EncodeToString(pubkeyRaw))

			pubkeyRaw, err = GetCompressedK256P1Pubkey(pubkeyRaw)
			if err != nil {
				return err
			}

			fmt.Println(hex.EncodeToString(pubkeyRaw))

			serverCtx.Logger.Info("adding validator", "module", types.ModuleName, "geneis", genesisFile)
			if err := UpdateModuleGenesis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				if len(genesis.Tokens) == 0 {
					return errors.New("no token setting found")
				}

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

				if votePower == 0 {
					return errors.New("no threshold setting found")
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

			serverCtx.Logger.Info("adding account", "module", authtypes.ModuleName, "geneis", genesisFile)
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

			tokenAddress, err := cmd.Flags().GetString(FlagTokenAddress)
			if err != nil {
				return err
			}

			tokenByte, err := bitcointypes.DecodeEthAddress(tokenAddress)
			if err != nil {
				return err
			}

			tokenDenom := types.TokenDenom(common.BytesToAddress(tokenByte))

			weight, err := cmd.Flags().GetUint64(FlagTokenWeight)
			if err != nil {
				return err
			}

			shareStr, err := cmd.Flags().GetString(FlagTokenThreshold)
			if err != nil {
				return err
			}

			var ok bool
			share := new(big.Int)
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
				for _, token := range genesis.Tokens {
					if token.Denom == tokenDenom {
						token.Token = types.Token{
							Weight:    weight,
							Threshold: math.NewIntFromBigIntMut(share),
						}
						return nil
					}
				}

				genesis.Tokens = append(genesis.Tokens, &types.TokenGenesis{
					Denom: tokenDenom,
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

	sign := &cobra.Command{
		Use:   "sign",
		Short: "get signature for current validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)

			ownerStr, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			chainID, err := cmd.Flags().GetUint64(FlagEthChainID)
			if err != nil {
				return err
			}

			ownerByte, err := bitcointypes.DecodeEthAddress(ownerStr)
			if err != nil {
				return err
			}

			keyJSONBytes, err := os.ReadFile(config.PrivValidatorKeyFile())
			if err != nil {
				panic(err)
			}

			var pvKey privval.FilePVKey
			err = cmtjson.Unmarshal(keyJSONBytes, &pvKey)
			if err != nil {
				return err
			}

			prvkey, err := ethcrypto.ToECDSA(pvKey.PrivKey.Bytes())
			if err != nil {
				return err
			}

			msgHash := getValidatorSignMsg(chainID, ownerByte, pvKey.Address.Bytes())
			sig, err := ethcrypto.Sign(msgHash, prvkey)
			if err != nil {
				return err
			}

			pubkey := make([]byte, 64)
			prvkey.X.FillBytes(pubkey[:32])
			prvkey.Y.FillBytes(pubkey[32:])

			return PrintJSON(map[string]string{
				"owner":     hexutil.Encode(ownerByte),
				"pubkey":    hexutil.Encode(pubkey),
				"signature": hexutil.Encode(sig),
				"validator": hexutil.Encode(pvKey.Address.Bytes()),
			})
		},
	}

	setParam := &cobra.Command{
		Use:   "param",
		Short: "update param",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			unlockDuration, err := cmd.Flags().GetDuration(FlagUnlockDuration)
			if err != nil {
				return err
			}

			exitDuration, err := cmd.Flags().GetDuration(FlagExitDuration)
			if err != nil {
				return err
			}

			serverCtx.Logger.Info("update param", "module", types.ModuleName, "geneis", genesisFile)
			return UpdateModuleGenesis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				genesis.Params.ExitingDuration = exitDuration
				genesis.Params.UnlockDuration = unlockDuration
				return genesis.Params.Validate()
			})
		},
	}

	cmd.Flags().String(FlagInitialReward, "", "remain reward amount")

	defaultParam := types.DefaultParams()
	setParam.Flags().Duration(FlagUnlockDuration, defaultParam.UnlockDuration, "unlock duation")
	setParam.Flags().Duration(FlagExitDuration, defaultParam.ExitingDuration, "exit duation")

	addToken.Flags().String(FlagTokenAddress, "", "token address")
	addToken.Flags().Uint64(FlagTokenWeight, 0, "validator vote power")
	addToken.Flags().String(FlagTokenThreshold, "", "validator vote power")
	addValidator.Flags().String(FlagValidatorPubkey, "", "validator pubkey(secp256k1)")

	sign.Flags().Uint64(FlagEthChainID, 48815, "the goat-geth chain id")
	sign.Flags().String(FlagOwner, "", "the validator owner")
	cmd.AddCommand(addToken, addValidator, sign, setParam)
	return cmd
}
