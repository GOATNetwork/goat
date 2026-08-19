package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	cmtcrypto "github.com/cometbft/cometbft/crypto"
	cmtmldsa65 "github.com/cometbft/cometbft/crypto/mldsa65"
	cmtsecp256k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/privval"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/goatnetwork/goat/x/locking/types"
	"github.com/spf13/cobra"
)

const (
	flagValidator = "validator"
	flagKeyType   = "key-type"
	flagOut       = "out"
	flagFromKey   = "from-key"
)

// RotateProofCommand produces the four arguments of Rotator.rotate.
//
// The proof has to be signed by the key being rotated to, and for ml_dsa_65
// nothing outside this repository can produce it: the EVM has no precompile
// and the message layout is ours. Without this command every operator would be
// copying key handling code out of the runbook.
func RotateProofCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate-proof",
		Short: "Generate the arguments for a consensus key rotation",
		Long: `Generate the arguments for a consensus key rotation.

By default a new consensus key is generated and written as a
priv_validator_key.json for the node to swap in at the height the consensus
layer reports. Pass --from-key to sign with a key you already hold instead, in
which case nothing is written.

The proof is verified before anything is printed, so whatever this prints is
something the consensus layer accepts.

The validator id is the address the validator was created with, which is not
the current consensus address once it has rotated.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			validator, err := cmd.Flags().GetString(flagValidator)
			if err != nil {
				return err
			}
			if !common.IsHexAddress(validator) {
				return fmt.Errorf("--%s must be a hex address, have %q", flagValidator, validator)
			}
			id := sdktypes.ConsAddress(common.HexToAddress(validator).Bytes())

			chainID, err := resolveChainID(cmd)
			if err != nil {
				return err
			}

			priv, out, err := rotationKey(cmd)
			if err != nil {
				return err
			}
			keyType := priv.Type()
			wire, err := types.RequestKeyType(keyType)
			if err != nil {
				return err
			}

			pubkey := priv.PubKey().Bytes()
			proof, err := priv.Sign(types.RotationProofMessage(chainID, id, keyType, pubkey))
			if err != nil {
				return err
			}
			// never hand out something the consensus layer would reject
			if err := types.VerifyRotationProof(chainID, id, keyType, pubkey, proof); err != nil {
				return fmt.Errorf("generated proof does not verify: %w", err)
			}

			if out != "" {
				// only the key file: priv_validator_state.json records what
				// this node has signed and must survive the swap untouched
				if err := os.MkdirAll(filepath.Dir(out), 0o700); err != nil {
					return err
				}
				privval.NewFilePV(priv, out, "").Key.Save()
			}

			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "chain_id      %s\n", chainID)
			fmt.Fprintf(w, "validator     %s\n", common.HexToAddress(validator))
			fmt.Fprintf(w, "key_type      %d   (%s)\n", wire, keyType)
			fmt.Fprintf(w, "pubkey        0x%s\n", hex.EncodeToString(pubkey))
			fmt.Fprintf(w, "proof         0x%s\n", hex.EncodeToString(proof))
			fmt.Fprintf(w, "new cons addr %X\n", priv.PubKey().Address())
			if out != "" {
				fmt.Fprintf(w, "key written   %s\n", out)
			}
			return nil
		},
	}

	cmd.Flags().String(flagValidator, "", "validator id: the address it was created with")
	cmd.Flags().String(flags.FlagChainID, "", "consensus chain id; read from the genesis file under --home when unset")
	cmd.Flags().String(flagKeyType, types.KeyTypeMlDsa65,
		fmt.Sprintf("key type to generate, %s or %s", types.KeyTypeMlDsa65, types.KeyTypeSecp256k1))
	cmd.Flags().String(flagOut, "new_priv_validator_key.json", "where to write the generated key")
	cmd.Flags().String(flagFromKey, "", "sign with this priv_validator_key.json instead of generating one")
	if err := cmd.MarkFlagRequired(flagValidator); err != nil {
		panic(err)
	}
	return cmd
}

// rotationKey returns the key to rotate to, and where to write it, or the
// empty string when the caller brought their own key.
func rotationKey(cmd *cobra.Command) (cmtcrypto.PrivKey, string, error) {
	from, err := cmd.Flags().GetString(flagFromKey)
	if err != nil {
		return nil, "", err
	}
	if from != "" {
		// LoadFilePVEmptyState leaves priv_validator_state.json alone
		pv := privval.LoadFilePVEmptyState(from, "")
		return pv.Key.PrivKey, "", nil
	}

	out, err := cmd.Flags().GetString(flagOut)
	if err != nil {
		return nil, "", err
	}
	if out == "" {
		return nil, "", fmt.Errorf("--%s must not be empty", flagOut)
	}
	// a consensus key is not something to overwrite by accident
	if _, err := os.Stat(out); err == nil {
		return nil, "", fmt.Errorf("%s already exists, refusing to overwrite it", out)
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, "", err
	}

	keyType, err := cmd.Flags().GetString(flagKeyType)
	if err != nil {
		return nil, "", err
	}
	switch types.KeyTypeName(keyType) {
	case types.KeyTypeMlDsa65:
		key, err := cmtmldsa65.GenPrivKey()
		if err != nil {
			return nil, "", err
		}
		return key, out, nil
	case types.KeyTypeSecp256k1:
		return cmtsecp256k1.GenPrivKey(), out, nil
	default:
		return nil, "", fmt.Errorf("unsupported consensus key type %q", keyType)
	}
}

// resolveChainID prefers the flag and falls back to the genesis file, because
// a proof signed for the wrong chain is rejected with nothing to point at it.
func resolveChainID(cmd *cobra.Command) (string, error) {
	chainID, err := cmd.Flags().GetString(flags.FlagChainID)
	if err != nil {
		return "", err
	}
	if chainID != "" {
		return chainID, nil
	}

	home, err := cmd.Flags().GetString(flags.FlagHome)
	if err != nil || home == "" {
		return "", fmt.Errorf("--%s is required when there is no genesis file to read it from", flags.FlagChainID)
	}
	genesis, err := genutiltypes.AppGenesisFromFile(filepath.Join(home, "config", "genesis.json"))
	if err != nil {
		return "", fmt.Errorf("--%s is required: %w", flags.FlagChainID, err)
	}
	return genesis.ChainID, nil
}
