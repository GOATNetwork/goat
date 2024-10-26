package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	cfg "github.com/cometbft/cometbft/config"
	tmsecp256k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/goatnetwork/goat/app"
	"github.com/spf13/cobra"
)

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(mbm module.BasicManager) *cobra.Command {
	const (
		// FlagOverwrite defines a flag to overwrite an existing genesis JSON file.
		FlagOverwrite = "overwrite"

		// FlagGenesisTime defines a flag to set genesis time
		FlagGenesisTime = "genesis-time"
	)

	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize the default geneis",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			serverCtx.Config.SetRoot(clientCtx.HomeDir)
			serverCtx.Config.Moniker = args[0]

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID == "" {
				return errors.New("no chain id provided")
			}

			genFile := serverCtx.Config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(FlagOverwrite)

			// use os.Stat to check if the file exists
			if _, err := os.Stat(genFile); !overwrite && !os.IsNotExist(err) {
				fmt.Println("genesis.json file already exists", genFile)
				return nil
			}

			appState, err := json.MarshalIndent(mbm.DefaultGenesis(clientCtx.Codec), "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshal default genesis state: %w", err)
			}

			consensusParam := cmttypes.DefaultConsensusParams()
			consensusParam.Validator.PubKeyTypes = []string{cmttypes.ABCIPubKeyTypeSecp256k1}
			consensusParam.Block.MaxBytes = 50 * 1024 * 124
			consensusParam.Block.MaxGas = -1

			appGenesis := &types.AppGenesis{
				AppName:       version.AppName,
				AppVersion:    version.Version,
				ChainID:       chainID,
				AppState:      appState,
				InitialHeight: 1,
				Consensus: &types.ConsensusGenesis{
					Validators: nil,
					Params:     consensusParam,
				},
			}

			if gtime, _ := cmd.Flags().GetString(FlagGenesisTime); gtime != "" {
				switch {
				case strings.HasPrefix(gtime, "+"):
					du, err := time.ParseDuration(gtime[1:])
					if err != nil {
						return fmt.Errorf("invalid duation string %s: %w", gtime, err)
					}
					appGenesis.GenesisTime = time.Now().Add(du).Round(0).UTC()
				case regexp.MustCompile(`^[0-9]+$`).MatchString(gtime):
					unix, err := strconv.ParseInt(gtime, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid unix timestamp %s: %w", gtime, err)
					}
					appGenesis.GenesisTime = time.Unix(unix, 0).Round(0).UTC()
				default:
					parsed, err := time.Parse(time.RFC3339, gtime)
					if err != nil {
						return fmt.Errorf("invalid RFC3339 timestamp %s: %w", gtime, err)
					}
					appGenesis.GenesisTime = parsed.Round(0).UTC()
				}
			}

			if err = genutil.ExportGenesisFile(appGenesis, genFile); err != nil {
				return errorsmod.Wrap(err, "Failed to export genesis file")
			}
			cfg.WriteConfigFile(filepath.Join(serverCtx.Config.RootDir, "config", "config.toml"), serverCtx.Config)
			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(FlagGenesisTime, "", "genesis time(rfc3399/unix number/duration(e.g. +1h))")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	return cmd
}

// initializeNodeFiles creates private validator and p2p configuration files if they doesn't exist
func initializeNodeFiles(config *cfg.Config) error {
	if _, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile()); err != nil {
		return err
	}

	pvKeyFile, pvStateFile := config.PrivValidatorKeyFile(), config.PrivValidatorStateFile()
	if _, pvKeyFileErr := os.Stat(pvKeyFile); os.IsNotExist(pvKeyFileErr) {
		filePV := privval.NewFilePV(tmsecp256k1.GenPrivKey(), pvKeyFile, pvStateFile)
		if err := os.MkdirAll(filepath.Dir(pvKeyFile), 0o777); err != nil {
			return fmt.Errorf("could not create directory %q: %w", filepath.Dir(pvKeyFile), err)
		}
		filePV.Key.Save()

		if _, pvStateFileErr := os.Stat(pvStateFile); os.IsNotExist(pvStateFileErr) {
			if err := os.MkdirAll(filepath.Dir(pvStateFile), 0o777); err != nil {
				return fmt.Errorf("could not create directory %q: %w", filepath.Dir(pvStateFile), err)
			}
			filePV.LastSignState.Save()
		}
	}
	return nil
}
