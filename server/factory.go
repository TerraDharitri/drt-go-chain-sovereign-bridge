package server

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/server/cmd/config"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/server/txSender"
)

// CreateSovereignBridgeServer creates a new bridge txs sender grpc server
func CreateSovereignBridgeServer(cfg *config.ServerConfig) (sovereign.BridgeTxSenderServer, error) {
	wallet, err := txSender.LoadWallet(cfg.WalletConfig)
	if err != nil {
		return nil, err
	}

	txSnd, err := txSender.CreateTxSender(wallet, cfg.TxSenderConfig)
	if err != nil {
		return nil, err
	}

	return NewSovereignBridgeTxServer(txSnd)
}
