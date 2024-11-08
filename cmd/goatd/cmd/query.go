package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	goatcrypto "github.com/goatnetwork/goat/pkg/crypto"
	"github.com/spf13/cobra"
)

func QueryMsgsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "msgs [height]",
		Short: "Query transactions and its results for a block by height",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse block height: %w", err)
			}

			node, err := clientCtx.GetNode()
			if err != nil {
				return err
			}

			block, err := node.Block(cmd.Context(), &height)
			if err != nil {
				return err
			}

			if block == nil {
				return fmt.Errorf("block %d not found", height)
			}

			// there is no easy way to marshal an array of proto msg
			results := make([]json.RawMessage, 0, len(block.Block.Txs))
			for _, raw := range block.Block.Txs {
				txHash := hex.EncodeToString(goatcrypto.SHA256Sum(raw))
				txResp, err := authtx.QueryTx(clientCtx, txHash)
				if err != nil {
					return err
				}
				tx, err := clientCtx.Codec.MarshalJSON(txResp)
				if err != nil {
					return err
				}
				results = append(results, tx)
			}
			final, err := json.Marshal(results)
			if err != nil {
				return err
			}
			return clientCtx.PrintRaw(final)
		},
	}
	return cmd
}
