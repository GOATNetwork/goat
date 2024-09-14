package modgen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	relayer "github.com/goatnetwork/goat/x/relayer/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func Bitcoin() *cobra.Command {
	const (
		FlagPubkey             = "pubkey"
		FlagPubkeyType         = "pubkey-type"
		FlagMinDeposit         = "min-deposit"
		FlagNetworkName        = "network"
		FlagConfirmationNumber = "confirmation-number"
	)

	parsePubkey := func(raw []byte, typ string) (*relayer.PublicKey, error) {
		var key relayer.PublicKey
		switch strings.ToLower(typ) {
		case "secp256k1":
			key.Key = &relayer.PublicKey_Secp256K1{Secp256K1: raw}
		case "schnorr":
			key.Key = &relayer.PublicKey_Schnorr{Schnorr: raw}
		default:
			return nil, fmt.Errorf("unknown key type: %s", typ)
		}
		return &key, nil
	}

	var cmd = &cobra.Command{
		Use:   "bitcoin [height] [hash...]",
		Short: "init bitcoin module genesis",
		Args:  cobra.RangeArgs(0, 7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			serverCtx.Logger.Info("update genesis", "module", types.ModuleName, "geneis", genesisFile)

			rawPubkey, err := cmd.Flags().GetBytesHex(FlagPubkey)
			if err != nil {
				return err
			}

			keyType, err := cmd.Flags().GetString(FlagPubkeyType)
			if err != nil {
				return err
			}

			newPubkey, err := parsePubkey(rawPubkey, keyType)
			if err != nil {
				return err
			}

			if err := UpdateModuleGenesis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				networkName, err := cmd.Flags().GetString(FlagNetworkName)
				if err != nil {
					return err
				}

				network, ok := types.BitcoinNetworks[networkName]
				if !ok {
					return fmt.Errorf("unknown bitcoin network: %s", networkName)
				}

				confirmationNumber, err := cmd.Flags().GetUint32(FlagConfirmationNumber)
				if err != nil {
					return err
				}

				if networkName == chaincfg.MainNetParams.Name {
					if confirmationNumber < 6 {
						return errors.New("confirmation number can't be less than 6")
					}
				}

				minDeposit, err := cmd.Flags().GetUint64(FlagMinDeposit)
				if err != nil {
					return err
				}

				genesis.Params = types.Params{
					ChainConfig: &types.ChainConfig{
						NetworkName:          network.Name,
						PubkeyHashAddrPrefix: uint32(network.PubKeyHashAddrID),
						ScriptHashAddrPrefix: uint32(network.ScriptHashAddrID),
						Bech32Hrp:            network.Bech32HRPSegwit,
					},
					ConfirmationNumber: confirmationNumber,
					DepositMagicPrefix: []byte(types.DepositMagicPreifxs[networkName]),
					MinDepositAmount:   minDeposit,
				}

				genesis.Pubkey = newPubkey

				if len(args) > 0 {
					start, err := cast.ToUint64E(args[0])
					if err != nil {
						return fmt.Errorf("invalid height: %s", args[0])
					}

					if genesis.Queue == nil {
						genesis.Queue = &types.ExecuableQueue{}
					}

					genesis.BlockHashes = make(map[uint64][]byte)
					for idx, hash := range args[1:] {
						if idx != 0 {
							start++
						}
						r, err := chainhash.NewHashFromStr(strings.TrimPrefix(hash, "0x"))
						if err != nil {
							return fmt.Errorf("invalid block hash: %s", hash)
						}
						genesis.BlockHashes[start] = r[:]
					}
					genesis.BlockTip = start
					genesis.Queue.BlockNumber = start
				}
				return genesis.Validate()
			}); err != nil {
				panic(err)
			}

			serverCtx.Logger.Info("update genesis", "module", relayer.ModuleName, "geneis", genesisFile)
			return UpdateModuleGenesis(genesisFile, relayer.ModuleName, new(relayer.GenesisState), clientCtx.Codec, func(state *relayer.GenesisState) error {
				for _, item := range state.Pubkeys {
					if item.Equal(newPubkey) {
						return nil
					}
				}
				state.Pubkeys = append(state.Pubkeys, newPubkey)
				return nil
			})
		},
	}

	param := types.DefaultParams()
	cmd.Flags().Uint32(FlagConfirmationNumber, param.ConfirmationNumber, "the confirmation number")
	cmd.Flags().BytesHex(FlagPubkey, nil, "the initial relayer public key")
	cmd.Flags().String(FlagPubkeyType, "secp256k1", "the public key type [secp256k1,schnorr]")
	cmd.Flags().String(FlagNetworkName, param.ChainConfig.NetworkName, "the bitcoin network name(mainnet|testnet3|regtest|signet)")
	cmd.Flags().Uint64(FlagMinDeposit, param.MinDepositAmount, "minimal allowed deposit amount")

	return cmd
}
