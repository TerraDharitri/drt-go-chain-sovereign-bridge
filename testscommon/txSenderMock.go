package testscommon

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
)

// TxSenderMock mocks TxSender interface
type TxSenderMock struct {
	SendTxsCalled func(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error)
}

// SendTxs mocks the SendTxs method
func (mock *TxSenderMock) SendTxs(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
	if mock.SendTxsCalled != nil {
		return mock.SendTxsCalled(ctx, data)
	}
	return nil, nil // Return appropriate default values if needed
}

// IsInterfaceNil mocks the IsInterfaceNil method
func (mock *TxSenderMock) IsInterfaceNil() bool {
	return mock == nil
}
