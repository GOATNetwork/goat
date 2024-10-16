package modgen

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
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

		FlagDepositTxid    = "txid"
		FlagDepositTxout   = "txout"
		FlagDepositSatoshi = "satoshi"
	)

	parsePubkey := func(raw []byte, typ string) (*relayer.PublicKey, error) {
		var key relayer.PublicKey
		switch strings.ToLower(typ) {
		case types.Secp256K1Name:
			key.Key = &relayer.PublicKey_Secp256K1{Secp256K1: raw}
		case types.SchnorrName:
			key.Key = &relayer.PublicKey_Schnorr{Schnorr: raw}
		default:
			return nil, fmt.Errorf("unknown key type: %s", typ)
		}
		return &key, nil
	}

	cmd := &cobra.Command{
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
					NetworkName:        network.Name,
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

					genesis.BlockHashes = make([][]byte, 0)
					for idx, hash := range args[1:] {
						if idx != 0 {
							start++
						}
						r, err := chainhash.NewHashFromStr(hash)
						if err != nil {
							return fmt.Errorf("invalid block hash: %s", hash)
						}
						genesis.BlockHashes = append(genesis.BlockHashes, r.CloneBytes())
					}
					slices.Reverse(genesis.BlockHashes) // tip is first
					genesis.BlockTip = start
					genesis.EthTxQueue.BlockNumber = start
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

	addDeposit := &cobra.Command{
		Use:   "add-deposit",
		Short: "genesis validator deposits",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			satoshi, err := cmd.Flags().GetUint64(FlagDepositSatoshi)
			if err != nil {
				return err
			}

			if satoshi == 0 {
				return errors.New("deposit value is 0")
			}

			// big endian
			txHash, err := cmd.Flags().GetString(FlagDepositTxid)
			if err != nil {
				return err
			}

			txid, err := chainhash.NewHashFromStr(txHash)
			if err != nil {
				return err
			}

			txout, err := cmd.Flags().GetUint32(FlagDepositTxout)
			if err != nil {
				return err
			}

			return UpdateModuleGenesis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				for _, item := range genesis.Deposits {
					if bytes.Equal(item.Txid, txid[:]) && item.Txout == txout {
						return nil
					}
				}
				genesis.Deposits = append(genesis.Deposits, types.DepositGenesis{
					Txid:   txid[:],
					Txout:  txout,
					Amount: satoshi,
				})
				return nil
			})
		},
	}

	param := types.DefaultParams()
	cmd.Flags().Uint32(FlagConfirmationNumber, param.ConfirmationNumber, "the confirmation number")
	cmd.Flags().BytesHex(FlagPubkey, nil, "the initial relayer public key")
	cmd.Flags().String(FlagPubkeyType, types.Secp256K1Name, "the public key type [secp256k1,schnorr]")
	cmd.Flags().String(FlagNetworkName, param.NetworkName, "the bitcoin network name(mainnet|testnet3|regtest|signet)")
	cmd.Flags().Uint64(FlagMinDeposit, param.MinDepositAmount, "minimal allowed deposit amount")

	addDeposit.Flags().Uint64(FlagDepositSatoshi, 0, "deposit amount in satoshi")
	addDeposit.Flags().String(FlagDepositTxid, "", "deposit txid")
	addDeposit.Flags().Uint32(FlagDepositTxout, 0, "deposit txout")

	cmd.AddCommand(addDeposit)
	return cmd
}
