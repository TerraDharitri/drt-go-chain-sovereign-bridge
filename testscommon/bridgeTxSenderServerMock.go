package testscommon

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
)

// MockBridgeTxSenderServer is a mock implementation of BridgeTxSenderServer
type MockBridgeTxSenderServer struct {
	SendCalled func(ctx context.Context, req *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error)
	*sovereign.UnimplementedBridgeTxSenderServer
}

// Send is a mock implementation of the Send method
func (m *MockBridgeTxSenderServer) Send(ctx context.Context, req *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error) {
	if m.SendCalled != nil {
		return m.SendCalled(ctx, req)
	}

	return &sovereign.BridgeOperationsResponse{}, nil
}
