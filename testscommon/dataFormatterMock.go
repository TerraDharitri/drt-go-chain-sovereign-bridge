package testscommon

import "github.com/TerraDharitri/drt-go-chain-core/data/sovereign"

// DataFormatterMock mocks DataFormatter interface
type DataFormatterMock struct {
	CreateTxsDataCalled func(data *sovereign.BridgeOperations) [][]byte
}

// CreateTxsData mocks the CreateTxsData method
func (mock *DataFormatterMock) CreateTxsData(data *sovereign.BridgeOperations) [][]byte {
	if mock.CreateTxsDataCalled != nil {
		return mock.CreateTxsDataCalled(data)
	}
	return nil
}

// IsInterfaceNil -
func (mock *DataFormatterMock) IsInterfaceNil() bool {
	return mock == nil
}
