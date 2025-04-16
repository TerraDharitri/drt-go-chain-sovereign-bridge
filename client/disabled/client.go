package disabled

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
)

type client struct{}

// NewDisabledClient creates a new instance of disabled client
func NewDisabledClient() *client {
	return &client{}
}

// Send does nothing and returns empty bridge operations response
func (c *client) Send(_ context.Context, _ *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error) {
	return &sovereign.BridgeOperationsResponse{}, nil
}

// Close returns no error
func (c *client) Close() error {
	return nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (c *client) IsInterfaceNil() bool {
	return c == nil
}
