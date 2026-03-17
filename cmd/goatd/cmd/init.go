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

var bootnodes = map[string][]string{
	"goat-testnet3": {
		"20443bfea09e1eab6fa4a8634e89dd6f6f12c5cf@35.82.170.238:26656",
		"cd23f3c3d21d27fc7d1119ee478819378765966a@100.22.4.171:26656",
		"4e1cd16b9377da92f10cde336de8b549c4be6499@52.36.155.137:26656",

		// Deprecated bootnodes, will be removed later
		"997f925d3d4947483c9339ec2765ed8c825ace32@3.222.213.223:26656",
		"9106b59e244eb8bf4dbedcd03b56e30790278765@54.68.179.184:26656",
		"c99c2abe0886a3c82c12f611477ce22fe178186a@52.32.82.160:26656",
	},
	"goat-mainnet": {
		"5b7503cd5e819ebb400a9068a72b05f7514b0a44@54.69.121.173:26656",
		"e1728a836481154bca8012871fee692fba4f9f9d@44.235.88.193:26656",
		"9f0dcce3c0fd0f2f292271a4734718b6ac07e720@44.254.4.159:26656",

		// Deprecated bootnodes, will be removed later
		// ZKM
		"84b041c3800b67319d378fbb8d2f83e2c686e738@3.16.248.103:26656",
		// Goat
		"445da2a86ee97f08154464fab98d90a3fef08a8f@18.220.78.155:26656",
		// Metis
		"f96b6429528bf55c57b9ba37790233c2ffefb4e8@3.147.156.156:26656",
	},
}

// initializeNodeFiles creates private validator and p2p configuration files if they doesn't exist
func initializeNodeFiles(cmd *cobra.Command, regtest bool) error {
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

	// we don't have gas system, but cosmos requires a non-empty gas value
	serverCtx.Viper.Set(server.FlagMinGasPrices, "0gas")

	chainID, _ := cmd.Flags().GetString(flags.FlagChainID)

	// insert genesis file if it doesn't exist
	genFile := serverCtx.Config.GenesisFile()
	if _, genFileErr := os.Stat(genFile); os.IsNotExist(genFileErr) {
		if regtest {
			return errors.New("genesis file should be existed for regtest")
		}
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

	// Note: cometbft doesn't use the viper to get the config

	// add bootnodes if not provided
	if !regtest && chainID != "" && serverCtx.Viper.GetString(FlagPersistentPeers) == "" {
		if bootnode, ok := bootnodes[chainID]; ok {
			serverCtx.Logger.Info("Set persistent peers", "chainID", chainID)
			serverCtx.Config.P2P.PersistentPeers = strings.Join(bootnode, ",")
		}
	}

	preset, _ := cmd.Flags().GetString(FlagGoatPreset)
	presets := strings.Split(preset, ",")
	if !regtest && slices.Contains(presets, "bootnode") {
		if ip, err := getPublicIP(); err != nil {
			serverCtx.Logger.Warn("Failed to fetch external public IP", "err", err.Error())
		} else {
			serverCtx.Logger.Info("Set external public IP", "ip", ip)
			serverCtx.Config.P2P.ListenAddress = "tcp://0.0.0.0:26656"
			serverCtx.Config.P2P.ExternalAddress = ip + ":26656"
		}
		serverCtx.Config.P2P.MaxNumInboundPeers = 200
		serverCtx.Config.P2P.MaxNumOutboundPeers = 200
		serverCtx.Config.P2P.AllowDuplicateIP = true
	}

	if regtest || slices.Contains(presets, "rpc") {
		serverCtx.Config.RPC.ListenAddress = "tcp://0.0.0.0:26657"
		serverCtx.Config.RPC.CORSAllowedOrigins = []string{"*"}
		serverCtx.Viper.Set("grpc.address", "0.0.0.0:9090")
		serverCtx.Viper.Set("api.enable", true)
		serverCtx.Viper.Set("api.address", "tcp://0.0.0.0:1317")
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
