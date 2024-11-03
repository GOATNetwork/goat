package cmd

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
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
var genesisFiles embed.FS

func chainList() string {
	dirEntry, _ := genesisFiles.ReadDir("genesis")
	var res []string
	for _, item := range dirEntry {
		res = append(res, strings.TrimSuffix(item.Name(), ".json"))
	}
	return strings.Join(res, ",")
}

var bootnodes = map[string][]string{}

// initializeNodeFiles creates private validator and p2p configuration files if they doesn't exist
func initializeNodeFiles(cmd *cobra.Command) error {
	serverCtx := server.GetServerContextFromCmd(cmd)

	// add default p2p key
	if _, err := p2p.LoadOrGenNodeKey(serverCtx.Config.NodeKeyFile()); err != nil {
		return err
	}

	// add validator private key
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

	// we have gas system, but cosmos requires a non-empty gas value
	serverCtx.Viper.Set(server.FlagMinGasPrices, "0gas")

	chainID, _ := cmd.Flags().GetString(flags.FlagChainID)

	// insert genesis file the file doesn't exist
	genFile := serverCtx.Config.GenesisFile()
	if _, genFileErr := os.Stat(genFile); os.IsNotExist(genFileErr) {
		if chainID == "" {
			return errors.New("no chain id")
		}
		jsonBytes, err := genesisFiles.ReadFile(fmt.Sprintf("genesis/%s.json", chainID))
		if err != nil {
			return fmt.Errorf("genesis not found for chain id %s", chainID)
		}
		if err := tempfile.WriteFileAtomic(genFile, jsonBytes, 0o644); err != nil {
			return err
		}
	}

	// add bootnodes if not provided
	if chainID != "" && serverCtx.Viper.GetString(FlagPersistentPeers) == "" {
		if bootnode, ok := bootnodes[chainID]; ok {
			serverCtx.Viper.Set(FlagPersistentPeers, strings.Join(bootnode, ","))
		}
	}

	preset, _ := cmd.Flags().GetString(FlagGoatPreset)
	presets := strings.Split(preset, ",")
	if slices.Contains(presets, "bootnode") {
		if serverCtx.Viper.GetString(FlagExternalIP) == "" {
			if ip, err := getPublicIP(); err != nil {
				serverCtx.Logger.Warn("Failed to fetch external public IP", "err", err.Error())
			} else {
				serverCtx.Logger.Info("Set external public IP", "ip", ip)
				serverCtx.Viper.Set(FlagP2PListener, "tcp://0.0.0.0:26656")
				serverCtx.Viper.Set(FlagExternalIP, ip+":26656")
			}
		}
		serverCtx.Viper.Set("p2p.max_num_inbound_peers", 200)
		serverCtx.Viper.Set("p2p.max_num_outbound_peers", 200)
	}

	if slices.Contains(presets, "rpc") {
		serverCtx.Viper.Set(flags.FlagGRPC, "0.0.0.0:9090")
		serverCtx.Viper.Set(server.FlagAPIEnable, true)
	}

	if slices.Contains(presets, "regtest") {
		serverCtx.Viper.Set(FlagP2PPex, false)
	}

	return nil
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
