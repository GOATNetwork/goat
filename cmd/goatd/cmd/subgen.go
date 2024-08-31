package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/goatnetwork/goat/app"
	"github.com/goatnetwork/goat/cmd/goatd/cmd/subgen"
	"github.com/spf13/cobra"
)

func SubgenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subgen",
		Short: "update module genesis",
		RunE:  client.ValidateCmd,
	}
	cmd.PersistentFlags().String(flags.FlagHome, app.DefaultNodeHome, "node's home directory")
	cmd.AddCommand(
		subgen.Bitcoin(),
		subgen.Relayer(),
		subgen.Goat(),
	)
	return cmd
}
