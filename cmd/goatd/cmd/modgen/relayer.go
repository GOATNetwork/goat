package modgen

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/goatnetwork/goat/x/relayer/types"
	"github.com/spf13/cobra"
)

func Relayer() *cobra.Command {
	const (
		FlagParamElectingPeriod = "param.electing_period"
		FlagThreshold           = "threshold"
		FlagPubkey              = "key.tx"
		FlagVoteKey             = "key.vote"
	)

	cmd := &cobra.Command{
		Use:   "relayer",
		Short: "update relayer module genesis",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			return UpdateGensis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				period, err := cmd.Flags().GetDuration(FlagParamElectingPeriod)
				if err != nil {
					return err
				}
				genesis.Params.ElectingPeriod = period
				return nil
			})
		},
	}

	appendVoter := &cobra.Command{
		Use:   "append address",
		Short: "append new voter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			addrcodec := clientCtx.TxConfig.SigningContext().AddressCodec()

			addr := args[0]
			addrByte, err := addrcodec.StringToBytes(addr)
			if err != nil {
				return fmt.Errorf("invalid address: %s", addr)
			}

			voteKey, err := cmd.Flags().GetBytesHex(FlagVoteKey)
			if err != nil {
				return err
			}

			txRawKey, err := cmd.Flags().GetBytesHex(FlagPubkey)
			if err != nil {
				return err
			}

			threshold, err := cmd.Flags().GetUint64(FlagThreshold)
			if err != nil {
				return err
			}

			if len(txRawKey) != secp256k1.PubKeySize || (txRawKey[0] != 2 && txRawKey[0] != 3) {
				return errors.New("not a valid secp256k1 compressed key")
			}

			txKey := &secp256k1.PubKey{Key: txRawKey}
			if txKeyAddr := txKey.Address().Bytes(); !bytes.Equal(txKeyAddr, addrByte) {
				addr, _ := addrcodec.BytesToString(txKeyAddr)
				return fmt.Errorf("address and public key not matched, expected address %s", addr)
			}

			if err := UpdateGensis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				if threshold != 0 {
					genesis.Threshold = threshold
				}
				if _, ok := genesis.Voters[addr]; ok {
					return fmt.Errorf("%s already added", addr)
				}
				genesis.Voters[addr] = &types.Voter{VoteKey: voteKey}
				return genesis.Validate()
			}); err != nil {
				return err
			}

			// Add the relayer account to auth module to allow it sending tx
			return UpdateGensis(genesisFile, authtypes.ModuleName, new(authtypes.GenesisState), clientCtx.Codec, func(genesis *authtypes.GenesisState) error {
				baseAccount, err := authtypes.NewBaseAccountWithPubKey(txKey)
				if err != nil {
					return err
				}

				if err := genesis.UnpackInterfaces(clientCtx.Codec); err != nil {
					return err
				}

				for _, v := range genesis.GetAccounts() {
					var acc authtypes.BaseAccount
					if err := clientCtx.Codec.UnpackAny(v, &acc); err != nil {
						return err
					}

					if acc.Address == addr {
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

	cmd.Flags().Duration(FlagParamElectingPeriod, time.Minute*10, "")

	appendVoter.Flags().Uint64(FlagThreshold, 0, "voter threshold")
	appendVoter.Flags().BytesHex(FlagPubkey, nil, "the voter tx public key(compressed secp256k1)")
	appendVoter.Flags().BytesHex(FlagVoteKey, nil, "the voter vote public key(compressed bls12381 G2)")
	cmd.AddCommand(appendVoter)
	return cmd
}
