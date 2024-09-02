package modgen

import (
	"fmt"
	"strings"

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
		FlagSafe               = "safe"
		FlagHard               = "hard"
		FlagPubkey             = "pubkey"
		FlagPubkeyType         = "pubkey-type"
		FlagMinDeposit         = "min-deposit"
		FlagDepositMagicPrefix = "deposit-magic-prefix"
		FlagNetworkName        = "network"
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
		Use:   "bitcoin height hash...",
		Short: "init bitcoin module genesis",
		Args:  cobra.RangeArgs(2, 7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			return UpdateGensis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				networkName, err := cmd.Flags().GetString(FlagNetworkName)
				if err != nil {
					return err
				}
				safe, err := cmd.Flags().GetUint32(FlagSafe)
				if err != nil {
					return err
				}

				hard, err := cmd.Flags().GetUint32(FlagHard)
				if err != nil {
					return err
				}

				depositMagicPreifx, err := cmd.Flags().GetString(FlagDepositMagicPrefix)
				if err != nil {
					return err
				}

				minDeposit, err := cmd.Flags().GetUint64(FlagMinDeposit)
				if err != nil {
					return err
				}

				genesis.Params = types.Params{
					NetworkName:           networkName,
					SafeConfirmationBlock: safe,
					HardConfirmationBlock: hard,
					DepositMagicPrefix:    []byte(depositMagicPreifx),
					MinDepositAmount:      minDeposit,
				}

				pubkey, err := cmd.Flags().GetBytesHex(FlagPubkey)
				if err != nil {
					return err
				}

				keyType, err := cmd.Flags().GetString(FlagPubkeyType)
				if err != nil {
					return err
				}

				genesis.Pubkey, err = parsePubkey(pubkey, keyType)
				if err != nil {
					return err
				}

				genesis.StartBlockNumber, err = cast.ToUint64E(args[0])
				if err != nil {
					return fmt.Errorf("invalid height: %s", args[0])
				}

				genesis.BlockHash = genesis.BlockHash[0:]
				for _, hash := range args[1:] {
					r, err := chainhash.NewHashFromStr(hash)
					if err != nil {
						return fmt.Errorf("invalid block hash: %s", hash)
					}
					genesis.BlockHash = append(genesis.BlockHash, r[:])
				}
				return genesis.Validate()
			})
		},
	}

	param := types.DefaultParams()
	cmd.Flags().Uint32(FlagSafe, param.SafeConfirmationBlock, "the safe confirmation number")
	cmd.Flags().Uint32(FlagHard, param.HardConfirmationBlock, "the hard confirmation number")
	cmd.Flags().BytesHex(FlagPubkey, nil, "the initial relayer public key")
	cmd.Flags().String(FlagPubkeyType, "secp256k1", "the public key type [secp256k1,schnorr]")
	cmd.Flags().String(FlagNetworkName, string(param.NetworkName), "the bitcoin network name(mainnet|testnet3|regtest|signet)")
	cmd.Flags().String(FlagDepositMagicPrefix, string(param.DepositMagicPrefix), "the deposit magic prefix for OP_RETURNS")
	cmd.Flags().Uint64(FlagMinDeposit, param.MinDepositAmount, "minimal allowed deposit amount")

	return cmd
}
