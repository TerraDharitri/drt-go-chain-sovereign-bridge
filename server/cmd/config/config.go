package config

import (
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/cert"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/server/txSender"
)

// ServerConfig holds necessary config for the grpc server
type ServerConfig struct {
	GRPCPort          string
	TxSenderConfig    txSender.TxSenderConfig
	WalletConfig      txSender.WalletConfig
	CertificateConfig cert.FileCfg
}
