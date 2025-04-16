package client

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
)

type client struct {
	bridgeClient sovereign.BridgeTxSenderClient
	conn         GRPCConn
}

// NewClient creates a wrapper over the grpc client connection and tx sender
func NewClient(bridgeClient sovereign.BridgeTxSenderClient, conn GRPCConn) (*client, error) {
	if conn == nil {
		return nil, errNilClientConnection
	}
	if bridgeClient == nil {
		return nil, nil
	}

	return &client{
		conn:         conn,
		bridgeClient: bridgeClient,
	}, nil
}

// Send sends bridge operations to the server
func (c *client) Send(ctx context.Context, data *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error) {
	return c.bridgeClient.Send(ctx, data)
}

// Close closes internal grpc connection
func (c *client) Close() error {
	return c.conn.Close()
}

// IsInterfaceNil checks if the underlying pointer is nil
func (c *client) IsInterfaceNil() bool {
	return c == nil
}
