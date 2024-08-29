package ethrpc

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	*ethclient.Client
}

// DialContext connects a client to the given URL with context.
func DialContext(ctx context.Context, rawurl string, jwt []byte) (*Client, error) {
	var opts []rpc.ClientOption
	if len(jwt) == 32 {
		opts = append(opts, rpc.WithHTTPAuth(node.NewJWTAuth([32]byte(jwt))))
	}

	client, err := rpc.DialOptions(ctx, rawurl, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{ethclient.NewClient(client)}, nil
}
