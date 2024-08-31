package subgen

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/goatnetwork/goat/x/goat/types"
	"github.com/spf13/cobra"
)

func Goat() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goat ethHeaderJsonFile",
		Short: "update goat module genesis",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			var header ethtypes.Header
			if err := json.Unmarshal(data, &header); err != nil {
				return err
			}

			if header.BaseFee == nil || header.WithdrawalsHash == nil {
				return fmt.Errorf("shanghai upgrade should be activated")
			}

			if *header.WithdrawalsHash != ethtypes.EmptyWithdrawalsHash {
				return fmt.Errorf("No withdrawals required")
			}

			if header.BlobGasUsed == nil || header.ExcessBlobGas == nil || header.ParentBeaconRoot == nil {
				return fmt.Errorf("cancun upgrade should be activated")
			}

			if *header.BlobGasUsed != 0 {
				return fmt.Errorf("No blob txes required")
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			genesisFile := config.GenesisFile()

			return UpdateGensis(genesisFile, types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				genesis.EthBlock = types.ExecutionPayload{
					ParentHash:    header.ParentHash.Bytes(),
					FeeRecipient:  header.Coinbase.Bytes(),
					StateRoot:     header.Root.Bytes(),
					ReceiptsRoot:  header.ReceiptHash.Bytes(),
					LogsBloom:     header.Bloom.Bytes(),
					PrevRandao:    header.Difficulty.Bytes(),
					BlockNumber:   header.Number.Uint64(),
					GasLimit:      header.GasLimit,
					GasUsed:       header.GasUsed,
					Timestamp:     header.Time,
					ExtraData:     header.Extra,
					BaseFeePerGas: header.BaseFee.Bytes(),
					BlockHash:     header.Hash().Bytes(),
					Transactions:  nil,
					BeaconRoot:    header.ParentBeaconRoot.Bytes(),
					BlobGasUsed:   *header.BlobGasUsed,
					ExcessBlobGas: *header.ExcessBlobGas,
				}
				return nil
			})
		},
	}
	return cmd
}
