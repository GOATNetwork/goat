package modgen

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/goatnetwork/goat/x/goat/types"
	"github.com/spf13/cobra"
)

func Goat() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goat eth-genesis-file",
		Short: "update goat module genesis",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			header, err := GetEthGenesisHeaderByFile(args[0])
			if err != nil {
				return err
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			srvConfig := serverCtx.Config.SetRoot(clientCtx.HomeDir)
			serverCtx.Logger.Info("update genesis", "module", types.ModuleName, "geneis", srvConfig.GenesisFile())
			return UpdateModuleGenesis(srvConfig.GenesisFile(), types.ModuleName, new(types.GenesisState), clientCtx.Codec, func(genesis *types.GenesisState) error {
				genesis.BeaconRoot = header.ParentBeaconRoot.Bytes()
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
					BaseFeePerGas: math.NewIntFromBigInt(header.BaseFee),
					BlockHash:     header.Hash().Bytes(),
					Transactions:  nil,
					BeaconRoot:    header.ParentBeaconRoot.Bytes(),
					BlobGasUsed:   *header.BlobGasUsed,
					ExcessBlobGas: *header.ExcessBlobGas,
					Requests:      nil,
				}
				return nil
			})
		},
	}
	return cmd
}
