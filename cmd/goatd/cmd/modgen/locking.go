package modgen

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"cosmossdk.io/math"
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
		FlagTotalReward     = "total-reward"
		FlagUnlockDuration  = "unlock-duration"
		FlagExitingDuration = "exit-duration"

		FlagDowntimeJailDuration    = "downtime-jail-duration"
		FlagMaxValidators           = "max-validators"
		FlagSignedBlocksWindow      = "signed-blocks-window"
		FlagMaxMissedPerWindow      = "max-missed-per-window"
		FlagSlashFractionDoubleSign = "slash-double-sign"
		FlagSlashFractionDowntime   = "slash-down"
		FlagHalvingInterval         = "halving-interval"
		FlagInitialBlockReward      = "initial-block-reward"

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

			rewardStr, err := cmd.Flags().GetString(FlagTotalReward)
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

			cmd.Println(hex.EncodeToString(pubkeyRaw))

			pubkeyRaw, err = GetCompressedK256P1Pubkey(pubkeyRaw)
			if err != nil {
				return err
			}

			cmd.Println(hex.EncodeToString(pubkeyRaw))

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

	showValidator := &cobra.Command{
		Use:   "show-validator",
		Short: "Shows this node's validator address",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			privValidator := privval.LoadFilePVEmptyState(
				serverCtx.Config.PrivValidatorKeyFile(),
				serverCtx.Config.PrivValidatorStateFile(),
			)
			cmd.Println(hexutil.Encode(privValidator.GetAddress()))
			return nil
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

			pvKey := privval.LoadFilePVEmptyState(config.PrivValidatorKeyFile(), "").Key
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

			exitingDuration, err := cmd.Flags().GetDuration(FlagExitingDuration)
			if err != nil {
				return err
			}

			maxValidators, err := cmd.Flags().GetInt64(FlagMaxValidators)
			if err != nil {
				return err
			}

			signedBlocksWindow, err := cmd.Flags().GetInt64(FlagSignedBlocksWindow)
			if err != nil {
				return err
			}
			maxMissedPerWindow, err := cmd.Flags().GetInt64(FlagMaxMissedPerWindow)
			if err != nil {
				return err
			}

			slashFractionDoubleSignStr, err := cmd.Flags().GetString(FlagSlashFractionDoubleSign)
			if err != nil {
				return err
			}
			slashFractionDoubleSign, err := math.LegacyNewDecFromStr(slashFractionDoubleSignStr)
			if err != nil {
				return err
			}

			slashFractionDowntimeStr, err := cmd.Flags().GetString(FlagSlashFractionDowntime)
			if err != nil {
				return err
			}
			slashFractionDowntime, err := math.LegacyNewDecFromStr(slashFractionDowntimeStr)
			if err != nil {
				return err
			}

			halvingInterval, err := cmd.Flags().GetInt64(FlagHalvingInterval)
			if err != nil {
				return err
			}

			initialBlockReward, err := cmd.Flags().GetInt64(FlagInitialBlockReward)
			if err != nil {
				return err
			}

			serverCtx.Logger.Info("update param", "module", types.ModuleName, "geneis", genesisFile)
			return UpdateModuleGenesis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				genesis.Params.UnlockDuration = unlockDuration
				genesis.Params.ExitingDuration = exitingDuration

				genesis.Params.MaxValidators = maxValidators
				genesis.Params.SignedBlocksWindow = signedBlocksWindow
				genesis.Params.MaxMissedPerWindow = maxMissedPerWindow

				genesis.Params.SlashFractionDoubleSign = slashFractionDoubleSign
				genesis.Params.SlashFractionDowntime = slashFractionDowntime

				genesis.Params.HalvingInterval = halvingInterval
				genesis.Params.InitialBlockReward = initialBlockReward
				return genesis.Params.Validate()
			})
		},
	}

	cmd.Flags().String(FlagTotalReward, "", "total reward amount in genesis")

	defaultParam := types.DefaultParams()
	setParam.Flags().Duration(FlagUnlockDuration, defaultParam.UnlockDuration, "waiting time for partial unlock")
	setParam.Flags().Duration(FlagExitingDuration, defaultParam.ExitingDuration, "waiting time for exit unlock")
	setParam.Flags().Duration(FlagDowntimeJailDuration, defaultParam.DowntimeJailDuration, "jail duration for downgraded validators")
	setParam.Flags().Int64(FlagMaxValidators, defaultParam.MaxValidators, "max number in the validator set")
	setParam.Flags().Int64(FlagSignedBlocksWindow, defaultParam.SignedBlocksWindow, "sgined block window")
	setParam.Flags().Int64(FlagMaxMissedPerWindow, defaultParam.MaxMissedPerWindow, "max missed block per block window")
	setParam.Flags().String(FlagSlashFractionDoubleSign, defaultParam.SlashFractionDoubleSign.String(), "frachtion for double sign")
	setParam.Flags().String(FlagSlashFractionDowntime, defaultParam.SlashFractionDowntime.String(), "frachtion for down")
	setParam.Flags().Int64(FlagHalvingInterval, defaultParam.HalvingInterval, "halving block interval")
	setParam.Flags().Int64(FlagInitialBlockReward, defaultParam.InitialBlockReward, "reward for the first block")

	addToken.Flags().String(FlagTokenAddress, "", "token address")
	addToken.Flags().Uint64(FlagTokenWeight, 0, "validator vote power")
	addToken.Flags().String(FlagTokenThreshold, "", "validator vote power")
	addValidator.Flags().String(FlagValidatorPubkey, "", "validator pubkey(secp256k1)")

	sign.Flags().Uint64(FlagEthChainID, 48815, "the goat-geth chain id")
	sign.Flags().String(FlagOwner, "", "the validator owner")
	cmd.AddCommand(addToken, addValidator, showValidator, sign, setParam)
	return cmd
}
