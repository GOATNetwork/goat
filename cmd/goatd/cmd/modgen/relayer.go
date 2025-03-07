package modgen

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/relayer/types"
	"github.com/spf13/cobra"
)

func Relayer() *cobra.Command {
	const (
		FlagParamElectingPeriod        = "param.electing_period"
		FlagParamAcceptProposerTimeout = "param.accept_proposer_timeout"

		FlagPubkey  = "key.tx"
		FlagVoteKey = "key.vote"

		FlagKeyOutput = "output"
	)

	cmd := &cobra.Command{
		Use:   "relayer",
		Short: "update relayer module genesis",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			return UpdateModuleGenesis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				period, err := cmd.Flags().GetDuration(FlagParamElectingPeriod)
				if err != nil {
					return err
				}
				genesis.Params.ElectingPeriod = period

				timeout, err := cmd.Flags().GetDuration(FlagParamAcceptProposerTimeout)
				if err != nil {
					return err
				}
				genesis.Params.AcceptProposerTimeout = timeout
				return nil
			})
		},
	}

	addVoter := &cobra.Command{
		Use:   "add-voter",
		Short: "add new voter",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			addrcodec := clientCtx.TxConfig.SigningContext().AddressCodec()

			voteKeyStr, err := cmd.Flags().GetString(FlagVoteKey)
			if err != nil {
				return err
			}

			voteKey, err := DecodeHexOrBase64String(voteKeyStr)
			if err != nil {
				return err
			}

			txKeyStr, err := cmd.Flags().GetString(FlagPubkey)
			if err != nil {
				return err
			}

			txRawKey, err := DecodeHexOrBase64String(txKeyStr)
			if err != nil {
				return err
			}

			txRawKey, err = GetCompressedK256P1Pubkey(txRawKey)
			if err != nil {
				return err
			}

			txKey := &secp256k1.PubKey{Key: txRawKey}
			addrByte := txKey.Address().Bytes()
			addr, err := addrcodec.BytesToString(txKey.Address().Bytes())
			if err != nil {
				return err
			}

			serverCtx.Logger.Info("update genesis", "module", types.ModuleName)
			if err := UpdateModuleGenesis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				votersSet := make(map[string]struct{})
				for _, voter := range genesis.Voters {
					addrStr, err := addrcodec.BytesToString(voter.Address)
					if err != nil {
						return err
					}
					if _, ok := votersSet[addrStr]; ok {
						return fmt.Errorf("voter %x is duplicated in the genesis", voter.Address)
					}
					votersSet[addrStr] = struct{}{}
				}

				if _, ok := votersSet[addr]; ok {
					serverCtx.Logger.Error("relayer already added", "addr", addr)
					return nil
				}

				if genesis.Relayer == nil {
					genesis.Relayer = &types.Relayer{
						Proposer:         addr,
						LastElected:      time.Now().UTC(),
						ProposerAccepted: true,
					}
				} else {
					voters := append(slices.Clone(genesis.Relayer.Voters), addr)
					voters = append(voters, genesis.Relayer.Proposer)
					slices.Sort(voters)
					genesis.Relayer.Proposer = voters[0]
					genesis.Relayer.Voters = voters[1:]
				}

				genesis.Voters = append(genesis.Voters, types.Voter{
					Address: addrByte,
					VoteKey: voteKey,
					Status:  types.VOTER_STATUS_ACTIVATED,
				})
				return genesis.Validate()
			}); err != nil {
				return err
			}

			serverCtx.Logger.Info("update genesis", "module", authtypes.ModuleName)
			// Add the relayer account to auth module to allow it sending tx
			return UpdateModuleGenesis(genesisFile, authtypes.ModuleName, new(authtypes.GenesisState), clientCtx.Codec, func(genesis *authtypes.GenesisState) error {
				baseAccount, err := authtypes.NewBaseAccountWithPubKey(txKey)
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

					if bytes.Equal(acc.GetAddress().Bytes(), addrByte) {
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

	keygen := &cobra.Command{
		Use:   "keygen",
		Short: "create key for relayer voter",
		RunE: func(cmd *cobra.Command, args []string) error {
			type Key struct {
				Type    string `json:"type"`
				Address string `json:"address"`
				TxKey   string `json:"txKey"`
				VoteKey string `json:"voteKey"`
			}

			clientCtx := client.GetClientContextFromCmd(cmd)

			output, err := cmd.Flags().GetString(FlagKeyOutput)
			if err != nil {
				return err
			}

			secretkeys := Key{Type: "SecretKey"}
			publicKeys := Key{Type: "PublicKey"}
			{
				p256k1 := secp256k1.GenPrivKey()
				address := p256k1.PubKey().Address()

				goatAddress, err := clientCtx.TxConfig.SigningContext().AddressCodec().BytesToString(address)
				if err != nil {
					return err
				}

				publicKeys.TxKey = hex.EncodeToString(p256k1.PubKey().Bytes())
				publicKeys.Address = hexutil.Encode(address)

				secretkeys.Address = goatAddress
				secretkeys.TxKey = hex.EncodeToString(p256k1.Bytes())
			}

			{
				secretKey := goatcrypto.GenPrivKey()
				publicKey := new(goatcrypto.PublicKey).From(secretKey)
				publicKeys.VoteKey = hex.EncodeToString(publicKey.Compress())
				secretkeys.VoteKey = hex.EncodeToString(secretKey.Serialize())
			}

			if output == "" || output == "-" {
				if err := WriteJSON(os.Stderr, secretkeys); err != nil {
					return err
				}
			} else {
				if !filepath.IsAbs(output) {
					output = filepath.Join(clientCtx.HomeDir, "relayer", output)
				}

				if _, err := os.Stat(output); err == nil {
					return fmt.Errorf("file %s exists", output)
				}

				if err := os.MkdirAll(filepath.Dir(output), os.ModePerm); err != nil {
					return err
				}

				data, err := json.MarshalIndent(secretkeys, "", "  ")
				if err != nil {
					return err
				}

				if err := os.WriteFile(output, data, 0o400); err != nil {
					return fmt.Errorf("failed to write secret keys: %w", err)
				}
			}
			return PrintJSON(publicKeys)
		},
	}

	cmd.Flags().Duration(FlagParamElectingPeriod, time.Minute*10, "")
	cmd.Flags().Duration(FlagParamAcceptProposerTimeout, time.Minute, "")
	addVoter.Flags().String(FlagPubkey, "", "the voter tx public key(compressed secp256k1)")
	addVoter.Flags().String(FlagVoteKey, "", "the voter vote public key(compressed bls12381 G2)")
	keygen.Flags().String(FlagKeyOutput, "", "the key file name")
	cmd.AddCommand(addVoter, keygen)
	return cmd
}
