package modgen

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/spf13/cobra"
)

func NewKey() *cobra.Command {
	const (
		FlagTxKey   = "tx"
		FlagVoteKey = "vote"
	)

	cmd := &cobra.Command{
		Use: "keygen",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			isTxKey, err := cmd.Flags().GetBool(FlagTxKey)
			if err != nil {
				return err
			}

			isVoteKey, err := cmd.Flags().GetBool(FlagVoteKey)
			if err != nil {
				return err
			}

			if isTxKey {
				key := secp256k1.GenPrivKey()
				serverCtx.Logger.Info(
					"secp256k1",
					"prvkey", hex.EncodeToString(key.Bytes()),
					"pubkey", hex.EncodeToString(key.PubKey().Bytes()),
				)
				address, err := clientCtx.TxConfig.SigningContext().AddressCodec().BytesToString(key.PubKey().Address())
				if err != nil {
					return err
				}

				pubkey, err := ethcrypto.DecompressPubkey(key.PubKey().Bytes())
				if err != nil {
					return err
				}
				serverCtx.Logger.Info(
					"address",
					"goat", address,
					"eth", ethcrypto.PubkeyToAddress(*pubkey).String(),
				)
			}

			if isVoteKey {
				secretKey := goatcrypto.GenPrivKey()
				publicKey := new(goatcrypto.PublicKey).From(secretKey)

				serverCtx.Logger.Info(
					"bls12-381",
					"prvkey", hex.EncodeToString(secretKey.Serialize()),
					"pubkey", hex.EncodeToString(publicKey.Compress()),
				)
			}
			return nil
		},
	}

	cmd.Flags().Bool(FlagTxKey, false, "create secp256k1 key")
	cmd.Flags().Bool(FlagVoteKey, false, "create bls12-381 key")

	return cmd
}
