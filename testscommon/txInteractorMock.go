package testscommon

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk/core"
)

// TxInteractorMock mocks TxInteractor interface
type TxInteractorMock struct {
	ApplyUserSignatureCalled func(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
	IsInterfaceNilCalled     func() bool
}

// ApplyUserSignature -
func (mock *TxInteractorMock) ApplyUserSignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
	if mock.ApplyUserSignatureCalled != nil {
		return mock.ApplyUserSignatureCalled(cryptoHolder, tx)
	}
	return nil
}

// IsInterfaceNil -
func (mock *TxInteractorMock) IsInterfaceNil() bool {
	return mock == nil
}
