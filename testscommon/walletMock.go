package testscommon

import (
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"

	"github.com/TerraDharitri/drt-go-sdk/core"
)

// CryptoComponentsHolderMock mocks CryptoComponentsHolder interface
type CryptoComponentsHolderMock struct {
	GetPublicKeyCalled      func() crypto.PublicKey
	GetPrivateKeyCalled     func() crypto.PrivateKey
	GetBech32Called         func() string
	GetAddressHandlerCalled func() core.AddressHandler
	IsInterfaceNilCalled    func() bool
}

// GetPublicKey mocks the GetPublicKey method
func (mock *CryptoComponentsHolderMock) GetPublicKey() crypto.PublicKey {
	if mock.GetPublicKeyCalled != nil {
		return mock.GetPublicKeyCalled()
	}
	return nil // Return appropriate default value if needed
}

// GetPrivateKey mocks the GetPrivateKey method
func (mock *CryptoComponentsHolderMock) GetPrivateKey() crypto.PrivateKey {
	if mock.GetPrivateKeyCalled != nil {
		return mock.GetPrivateKeyCalled()
	}
	return nil // Return appropriate default value if needed
}

// GetBech32 mocks the GetBech32 method
func (mock *CryptoComponentsHolderMock) GetBech32() string {
	if mock.GetBech32Called != nil {
		return mock.GetBech32Called()
	}
	return "" // Return appropriate default value if needed
}

// GetAddressHandler mocks the GetAddressHandler method
func (mock *CryptoComponentsHolderMock) GetAddressHandler() core.AddressHandler {
	if mock.GetAddressHandlerCalled != nil {
		return mock.GetAddressHandlerCalled()
	}
	return nil // Return appropriate default value if needed
}

// IsInterfaceNil -
func (mock *CryptoComponentsHolderMock) IsInterfaceNil() bool {
	return mock == nil
}
