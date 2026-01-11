package lightclient

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/goatnetwork/goat/app"
	"github.com/goatnetwork/goat/cmd/goatd/cmd/goatflags"
	"github.com/goatnetwork/goat/pkg/ethrpc"
	"github.com/goatnetwork/goat/x/goat/types"
	"github.com/spf13/cobra"
)

const (
	FlagInterval = "interval"
)

func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "light-client",
		Short:   `Start a light client to sync geth blocks from a trust goatd full node`,
		Example: `goatd light-client --goat-geth <path-to-goat-geth-ipc> --chain-id <goat-chain-id> --node <trust-node-jsonrpc-url>`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			basectx, cancel := signal.NotifyContext(cmd.Context(),
				syscall.SIGTERM, syscall.SIGTERM)
			defer cancel()

			client, err := newLightClient(basectx, cmd)
			if err != nil {
				return err
			}
			defer client.Stop()
			return client.Start(basectx)
		},
	}

	cmd.Flags().String(flags.FlagNode, "tcp://localhost:26657", "<host>:<port> to CometBFT RPC interface for this chain")
	cmd.Flags().String(flags.FlagChainID, "", "The network chain ID")
	cmd.Flags().Duration(FlagInterval, time.Second*3, "interval for fetching new block")
	return cmd
}

type lightClient struct {
	context      client.Context
	chainID      string
	engineClient *ethrpc.Client
	engineConfig *params.ChainConfig
	ticker       *time.Timer
	interval     time.Duration
	logger       log.Logger
}

func newLightClient(basectx context.Context, cmd *cobra.Command) (*lightClient, error) {
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return nil, err
	}
	clientCtx = clientCtx.WithCmdContext(basectx)

	tickerNum, err := cmd.Flags().GetDuration(FlagInterval)
	if err != nil {
		return nil, err
	}

	chainID, err := cmd.Flags().GetString(flags.FlagChainID)
	if err != nil {
		return nil, err
	}
	if chainID == "" {
		return nil, errors.New("chain id is required")
	}
	endpoint, err := cmd.Flags().GetString(goatflags.GoatGeth)
	if err != nil {
		return nil, err
	}

	logger := server.GetServerContextFromCmd(cmd).Logger
	engineClient, engineConfig, err := app.ConnectEngineClient(basectx, logger, endpoint)
	if err != nil {
		return nil, err
	}
	return &lightClient{
		context:      clientCtx,
		chainID:      chainID,
		engineClient: engineClient,
		engineConfig: engineConfig,
		interval:     tickerNum,
		ticker:       time.NewTimer(0),
		logger:       logger,
	}, nil
}

func (client *lightClient) Start(basectx context.Context) error {
	client.logger.Info("light client started",
		"interval", client.interval.String(), "node", client.context.NodeURI, "chain-id", client.chainID)
	for {
		select {
		case <-basectx.Done():
			return nil
		case <-client.ticker.C:
			if err := client.run(basectx); err != nil {
				client.logger.Error("syncing error", "error", err)
			}
			client.ticker.Reset(client.interval)
		}
	}
}

func (client *lightClient) Stop() {
	client.ticker.Stop()
	client.engineClient.Close()
}

func (client *lightClient) run(ctx context.Context) error {
	node, err := client.context.GetNode()
	if err != nil {
		return err
	}

	// get the latest block from cometbft
	status, err := node.Status(ctx)
	if err != nil {
		return err
	}
	height := status.SyncInfo.LatestBlockHeight

	blockRes, err := node.Block(ctx, &height)
	if err != nil {
		return err
	}
	if blockRes == nil {
		return fmt.Errorf("block %d not found", height)
	}

	if got := blockRes.Block.ChainID; client.chainID != got {
		return fmt.Errorf("goat chain id mismatch: expected %s, got %s", client.chainID, got)
	}

	for i, raw := range blockRes.Block.Txs {
		if i != 0 {
			break
		}

		txResp, err := node.Tx(ctx, raw.Hash(), true)
		if err != nil {
			return err
		}

		if txResp.TxResult.Code != 0 {
			return nil
		}

		tx, err := client.context.TxConfig.TxDecoder()(txResp.Tx)
		if err != nil {
			return err
		}

		gethHeight, err := client.engineClient.BlockNumber(ctx)
		if err != nil {
			return err
		}

		var execData *engine.ExecutableData
		var beaconRoot common.Hash
		var requests [][]byte
		for _, msg := range tx.GetMsgs() {
			switch v := msg.(type) {
			case *types.MsgNewEthBlock:
				if v.Payload.BlockNumber <= gethHeight {
					return fmt.Errorf("block number %d is not greater than goat-geth height %d", v.Payload.BlockNumber, gethHeight)
				}
				execData = types.PayloadToExecutableData(v.Payload)
				beaconRoot = common.BytesToHash(v.Payload.BeaconRoot)
				requests = v.Payload.Requests
			default:
				return errors.New("unknown tx type")
			}
		}

		client.logger.Info("Notify NewPayload", "number", execData.Number)
		response, err := client.engineClient.NewPayloadV4(ctx, execData,
			[]common.Hash{}, beaconRoot, requests)
		if err != nil {
			return err
		}

		if response.Status == engine.INVALID {
			return errors.New("invalid from NewPayloadV4 api")
		}

		// set current block hash to head state and set previous block hash to safe and finalized state
		client.logger.Info("Notify ForkChoiceUpdated",
			"head", execData.BlockHash.Hex(), "finalized", execData.ParentHash.Hex())
		forkRes, err := client.engineClient.ForkchoiceUpdatedV3(ctx, &engine.ForkchoiceStateV1{
			HeadBlockHash: execData.BlockHash,
			SafeBlockHash: execData.ParentHash, FinalizedBlockHash: execData.ParentHash,
		}, nil)
		if err != nil {
			return err
		}
		if forkRes.PayloadStatus.Status == engine.INVALID {
			return errors.New("invalid from ForkchoiceUpdatedV3 api")
		}
	}

	return nil
}
