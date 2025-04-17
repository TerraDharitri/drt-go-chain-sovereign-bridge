package testscommon

import (
	"context"

	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
)

// ProxyMock mocks Proxy interface
type ProxyMock struct {
	GetAccountCalled       func(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	GetNetworkConfigCalled func(ctx context.Context) (*data.NetworkConfig, error)
	IsInterfaceNilCalled   func() bool
}

// GetAccount mocks the GetAccount method
func (mock *ProxyMock) GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error) {
	if mock.GetAccountCalled != nil {
		return mock.GetAccountCalled(ctx, address)
	}
	return &data.Account{}, nil
}

// GetNetworkConfig mocks the GetNetworkConfig method
func (mock *ProxyMock) GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error) {
	if mock.GetNetworkConfigCalled != nil {
		return mock.GetNetworkConfigCalled(ctx)
	}
	return &data.NetworkConfig{}, nil
}

// IsInterfaceNil -
func (mock *ProxyMock) IsInterfaceNil() bool {
	return mock == nil
}
