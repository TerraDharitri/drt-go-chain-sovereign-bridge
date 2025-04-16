package testscommon

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
)

// TxNonceSenderHandlerMock mocks TxNonceSenderHandler interface
type TxNonceSenderHandlerMock struct {
	ApplyNonceAndGasPriceCalled func(ctx context.Context, txs ...*transaction.FrontendTransaction) error
	SendTransactionsCalled      func(ctx context.Context, txs ...*transaction.FrontendTransaction) ([]string, error)
}

// ApplyNonceAndGasPrice mocks the ApplyNonceAndGasPrice method
func (mock *TxNonceSenderHandlerMock) ApplyNonceAndGasPrice(ctx context.Context, txs ...*transaction.FrontendTransaction) error {
	if mock.ApplyNonceAndGasPriceCalled != nil {
		return mock.ApplyNonceAndGasPriceCalled(ctx, txs...)
	}
	return nil
}

// SendTransactions mocks the SendTransaction method
func (mock *TxNonceSenderHandlerMock) SendTransactions(ctx context.Context, txs ...*transaction.FrontendTransaction) ([]string, error) {
	if mock.SendTransactionsCalled != nil {
		return mock.SendTransactionsCalled(ctx, txs...)
	}
	return make([]string, 0), nil
}

// IsInterfaceNil -
func (mock *TxNonceSenderHandlerMock) IsInterfaceNil() bool {
	return mock == nil
}
