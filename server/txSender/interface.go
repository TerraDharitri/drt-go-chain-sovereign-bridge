package txSender

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
)

// TxInteractor defines a tx interactor with dharitri blockchain
type TxInteractor interface {
	ApplyUserSignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
	IsInterfaceNil() bool
}

// Proxy defines the proxy to interact with Dharitri blockchain
type Proxy interface {
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	IsInterfaceNil() bool
}

// DataFormatter should format txs data for bridge operations
type DataFormatter interface {
	CreateTxsData(data *sovereign.BridgeOperations) [][]byte
	IsInterfaceNil() bool
}

// TxNonceSenderHandler should handle nonce management and tx interactions
type TxNonceSenderHandler interface {
	ApplyNonceAndGasPrice(ctx context.Context, txs ...*transaction.FrontendTransaction) error
	SendTransactions(ctx context.Context, txs ...*transaction.FrontendTransaction) ([]string, error)
	IsInterfaceNil() bool
}
