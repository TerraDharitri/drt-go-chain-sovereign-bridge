package server

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
)

// TxSender defines a tx sender for bridge operations
type TxSender interface {
	SendTxs(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error)
	IsInterfaceNil() bool
}
