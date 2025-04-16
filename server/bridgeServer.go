package server

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
)

var log = logger.GetOrCreate("server")

type server struct {
	txSender TxSender
	*sovereign.UnimplementedBridgeTxSenderServer
}

// NewSovereignBridgeTxServer creates a new sovereign bridge operations server. This server receives bridge data operations from
// sovereign nodes and sends transactions to main chain.
func NewSovereignBridgeTxServer(txSender TxSender) (*server, error) {
	if check.IfNil(txSender) {
		return nil, errNilTxSender
	}

	return &server{
		txSender: txSender,
	}, nil
}

// Send should handle receiving data bridge operations from sovereign shard and forward transactions to main chain
func (s *server) Send(ctx context.Context, data *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error) {
	hashes, err := s.txSender.SendTxs(ctx, data)
	if err != nil {
		return nil, err
	}

	logTxHashes(hashes)

	return &sovereign.BridgeOperationsResponse{
		TxHashes: hashes,
	}, nil
}

func logTxHashes(hashes []string) {
	for _, hash := range hashes {
		log.Info("sent tx", "hash", hash)
	}
}

// IsInterfaceNil checks if the underlying pointer is nil
func (s *server) IsInterfaceNil() bool {
	return s == nil
}
