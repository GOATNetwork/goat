package cmd

import (
	"bufio"
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
	"cosmossdk.io/math/unsafe"
	cfg "github.com/cometbft/cometbft/config"
	tmsecp256k1 "github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/go-bip39"
	"github.com/goatnetwork/goat/app"
	"github.com/goatnetwork/goat/cmd/goatd/cmd/modgen"
	"github.com/spf13/cobra"
)

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(mbm module.BasicManager) *cobra.Command {
	const (
		// FlagOverwrite defines a flag to overwrite an existing genesis JSON file.
		FlagOverwrite = "overwrite"

		// FlagSeed defines a flag to initialize the private validator key from a specific seed.
		FlagRecover = "recover"

		// FlagDefaultBondDenom defines the default denom to use in the genesis file.
		FlagDefaultBondDenom = "default-denom"

		FlagGenesisTime = "genesis-time"
	)

	type printInfo struct {
		Moniker   string `json:"moniker" yaml:"moniker"`
		ChainID   string `json:"chain_id" yaml:"chain_id"`
		NodeID    string `json:"node_id" yaml:"node_id"`
		Validator string `json:"validator" yaml:"validator"`
	}

	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			genesisTime, _ := cmd.Flags().GetString(FlagGenesisTime)

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			switch {
			case chainID != "":
			case clientCtx.ChainID != "":
				chainID = clientCtx.ChainID
			default:
				chainID = fmt.Sprintf("test-chain-%v", unsafe.Str(6))
			}

			// Get bip39 mnemonic
			var mnemonic string
			isRecover, _ := cmd.Flags().GetBool(FlagRecover)
			if isRecover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				value, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				mnemonic = value
				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			// Get initial height
			initHeight, _ := cmd.Flags().GetInt64(flags.FlagInitHeight)
			if initHeight < 1 {
				initHeight = 1
			}

			nodeID, validatorKey, err := initializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(FlagOverwrite)
			defaultDenom, _ := cmd.Flags().GetString(FlagDefaultBondDenom)

			// use os.Stat to check if the file exists
			_, err = os.Stat(genFile)
			if !overwrite && !os.IsNotExist(err) {
				fmt.Printf("genesis.json file already exists: %v", genFile)
				return nil
			}

			// Overwrites the SDK default denom for side-effects
			if defaultDenom != "" {
				sdk.DefaultBondDenom = defaultDenom
			}
			appGenState := mbm.DefaultGenesis(cdc)

			appState, err := json.MarshalIndent(appGenState, "", " ")
			if err != nil {
				return errorsmod.Wrap(err, "Failed to marshal default genesis state")
			}

			appGenesis := &types.AppGenesis{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				appGenesis, err = types.AppGenesisFromFile(genFile)
				if err != nil {
					return errorsmod.Wrap(err, "Failed to read genesis doc from file")
				}
			}

			if genesisTime != "" {
				switch {
				case strings.HasPrefix(genesisTime, "+"):
					du, err := time.ParseDuration(genesisTime[1:])
					if err != nil {
						return fmt.Errorf("invalid duation string %s: %w", genesisTime, err)
					}
					appGenesis.GenesisTime = time.Now().Add(du).Round(0).UTC()
				case regexp.MustCompile(`^[0-9]+$`).MatchString(genesisTime):
					unix, err := strconv.ParseInt(genesisTime, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid unix timestamp %s: %w", genesisTime, err)
					}
					appGenesis.GenesisTime = time.Unix(unix, 0).Round(0).UTC()
				default:
					parsed, err := time.Parse(time.RFC3339, genesisTime)
					if err != nil {
						return fmt.Errorf("invalid RFC3339 timestamp %s: %w", genesisTime, err)
					}
					appGenesis.GenesisTime = parsed.Round(0).UTC()
				}
			}

			consensusParam := cmttypes.DefaultConsensusParams()
			consensusParam.Validator.PubKeyTypes = []string{cmttypes.ABCIPubKeyTypeSecp256k1}
			consensusParam.Block.MaxBytes = 50 * 1024 * 124
			consensusParam.Block.MaxGas = -1

			appGenesis.AppName = version.AppName
			appGenesis.AppVersion = version.Version
			appGenesis.ChainID = chainID
			appGenesis.AppState = appState
			appGenesis.InitialHeight = initHeight
			appGenesis.Consensus = &types.ConsensusGenesis{
				Validators: nil,
				Params:     consensusParam,
			}

			if err = genutil.ExportGenesisFile(appGenesis, genFile); err != nil {
				return errorsmod.Wrap(err, "Failed to export genesis file")
			}

			toPrint := printInfo{
				Moniker:   config.Moniker,
				ChainID:   chainID,
				NodeID:    nodeID,
				Validator: validatorKey.Address().String(),
			}

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return modgen.PrintJSON(toPrint)
		},
	}

	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().Bool(FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(FlagGenesisTime, "", "genesis time(rfc3399/unix number/duration(e.g. +1h))")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(FlagDefaultBondDenom, "goat", "genesis file default denomination, if left blank default value is 'stake'")
	cmd.Flags().Int64(flags.FlagInitHeight, 1, "specify the initial block height at genesis")

	return cmd
}

// initializeNodeValidatorFilesFromMnemonic creates private validator and p2p configuration files using the given mnemonic.
// If no valid mnemonic is given, a random one will be used instead.
func initializeNodeValidatorFilesFromMnemonic(config *cfg.Config, mnemonic string) (nodeID string, valPubKey cryptotypes.PubKey, err error) {
	if len(mnemonic) > 0 && !bip39.IsMnemonicValid(mnemonic) {
		return "", nil, fmt.Errorf("invalid mnemonic")
	}
	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return "", nil, err
	}

	nodeID = string(nodeKey.ID())

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := os.MkdirAll(filepath.Dir(pvKeyFile), 0o777); err != nil {
		return "", nil, fmt.Errorf("could not create directory %q: %w", filepath.Dir(pvKeyFile), err)
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := os.MkdirAll(filepath.Dir(pvStateFile), 0o777); err != nil {
		return "", nil, fmt.Errorf("could not create directory %q: %w", filepath.Dir(pvStateFile), err)
	}

	_, pvKeyFileErr := os.Stat(pvKeyFile)
	_, pvStateFileErr := os.Stat(pvStateFile)

	var filePV *privval.FilePV
	switch {
	case pvKeyFileErr == nil && pvStateFileErr == nil:
		filePV = privval.LoadFilePV(pvKeyFile, pvStateFile)
	case len(mnemonic) == 0:
		filePV = privval.NewFilePV(tmsecp256k1.GenPrivKey(), pvKeyFile, pvStateFile)
		filePV.Save()
	default:
		privKey := tmsecp256k1.GenPrivKeySecp256k1(bip39.NewSeed(mnemonic, ""))
		filePV = privval.NewFilePV(privKey, pvKeyFile, pvStateFile)
		filePV.Save()
	}

	tmValPubKey, err := filePV.GetPubKey()
	if err != nil {
		return "", nil, err
	}

	valPubKey, err = cryptocodec.FromCmtPubKeyInterface(tmValPubKey)
	if err != nil {
		return "", nil, err
	}

	return nodeID, valPubKey, nil
}
