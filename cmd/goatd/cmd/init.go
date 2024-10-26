package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tmsecp256k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/libs/tempfile"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

//go:embed genesis
var genesis embed.FS

func allChainID() string {
	dirEntry, _ := genesis.ReadDir("genesis")
	var res []string
	for _, item := range dirEntry {
		res = append(res, strings.TrimSuffix(item.Name(), ".json"))
	}
	return strings.Join(res, ",")
}

// initializeNodeFiles creates private validator and p2p configuration files if they doesn't exist
func initializeNodeFiles(cmd *cobra.Command) error {
	serverCtx := server.GetServerContextFromCmd(cmd)

	if _, err := p2p.LoadOrGenNodeKey(serverCtx.Config.NodeKeyFile()); err != nil {
		return err
	}

	pvKeyFile := serverCtx.Config.PrivValidatorKeyFile()
	if _, pvKeyFileErr := os.Stat(pvKeyFile); os.IsNotExist(pvKeyFileErr) {
		pvStateFile := serverCtx.Config.PrivValidatorStateFile()
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

	if cmd.Name() != "start" {
		return nil
	}

	serverCtx.Viper.Set("minimum-gas-prices", "0gas")
	genFile := serverCtx.Config.GenesisFile()
	if _, genFileErr := os.Stat(genFile); os.IsNotExist(genFileErr) {
		chainID := serverCtx.Viper.GetString(flags.FlagChainID)
		if chainID == "" {
			return errors.New("no chain id")
		}
		jsonBytes, err := genesis.ReadFile(fmt.Sprintf("genesis/%s.json", chainID))
		if err != nil {
			return fmt.Errorf("genesis not found for chain id %s", chainID)
		}
		return tempfile.WriteFileAtomic(genFile, jsonBytes, 0o600)
	}
	return nil
}
