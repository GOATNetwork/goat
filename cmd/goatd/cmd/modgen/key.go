package modgen

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/goatnetwork/goat/x/bitcoin/types"
	"github.com/spf13/cobra"
)

func NewKey() *cobra.Command {
	const (
		FlagTxKey   = "tx"
		FlagVoteKey = "vote"
		FlagNetwork = "network"
	)

	cmd := &cobra.Command{
		Use: "keygen",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			isTxKey, err := cmd.Flags().GetBool(FlagTxKey)
			if err != nil {
				return err
			}

			isVoteKey, err := cmd.Flags().GetBool(FlagVoteKey)
			if err != nil {
				return err
			}

			networkName, err := cmd.Flags().GetString(FlagNetwork)
			if err != nil {
				return err
			}

			network, ok := types.BitcoinNetworks[networkName]
			if !ok {
				return fmt.Errorf("unknown bitcoin network: %s", networkName)
			}

			if isTxKey {
				key := secp256k1.GenPrivKey()

				fmt.Println("secp256k1 prvkey", hex.EncodeToString(key.Bytes()))

				rawAddress := key.PubKey().Address()
				goatAddr, err := clientCtx.TxConfig.SigningContext().AddressCodec().BytesToString(rawAddress)
				if err != nil {
					return err
				}

				pubkey, err := ethcrypto.DecompressPubkey(key.PubKey().Bytes())
				if err != nil {
					return err
				}

				fmt.Println("secp256k1 compressed pubkey", hex.EncodeToString(key.PubKey().Bytes()))
				fmt.Println("uncompressed pubkey", hex.EncodeToString(ethcrypto.FromECDSAPub(pubkey)))
				fmt.Println()

				btcAddr, err := btcutil.NewAddressWitnessPubKeyHash(rawAddress, network)
				if err != nil {
					return err
				}

				fmt.Println("goat address", goatAddr)
				fmt.Println("goat address bytes", hex.EncodeToString(rawAddress))
				fmt.Println()
				fmt.Println("eth address", ethcrypto.PubkeyToAddress(*pubkey).String())
				fmt.Println("btc address(p2wpkh)", btcAddr.EncodeAddress())
			}

			if isVoteKey {
				secretKey := goatcrypto.GenPrivKey()
				publicKey := new(goatcrypto.PublicKey).From(secretKey)

				fmt.Println()
				fmt.Println("bls12-381 prvkey", hex.EncodeToString(secretKey.Serialize()))
				fmt.Println("bls12-381 pubkey", hex.EncodeToString(publicKey.Compress()))
			}
			return nil
		},
	}

	cmd.Flags().Bool(FlagTxKey, false, "create secp256k1 key")
	cmd.Flags().Bool(FlagVoteKey, false, "create bls12-381 key")
	cmd.Flags().String(FlagNetwork, "regtest", "bitcoin network name")

	return cmd
}
