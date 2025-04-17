package config

import "github.com/TerraDharitri/drt-go-chain-sovereign-bridge/cert"

// ClientConfig holds all grpc client's config
type ClientConfig struct {
	Enabled        bool
	GRPCHost       string
	GRPCPort       string
	CertificateCfg cert.FileCfg
}
