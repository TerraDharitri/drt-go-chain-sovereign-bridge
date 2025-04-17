package txSender

import (
	"fmt"
	"strings"

	"github.com/TerraDharitri/drt-go-chain-crypto/signing"
	"github.com/TerraDharitri/drt-go-chain-crypto/signing/ed25519"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-sdk/blockchain/cryptoProvider"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/interactors"
)

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
	log    = logger.GetOrCreate("drt-go-chain-sovereign-bridge")
)

const (
	json = "json"
	pem  = "pem"
)

// LoadWallet loads a wallet using provided config
func LoadWallet(cfg WalletConfig) (core.CryptoComponentsHolder, error) {
	var privateKey []byte
	var err error

	w := interactors.NewWallet()
	walletType := getWalletType(cfg.Path)
	switch walletType {
	case pem:
		privateKey, err = w.LoadPrivateKeyFromPemFile(cfg.Path)
	case json:
		privateKey, err = w.LoadPrivateKeyFromJsonFile(cfg.Path, cfg.Password)
	default:
		return nil, fmt.Errorf("%w: %s, acceptable:%s, %s", errInvalidWalletType, walletType, pem, json)
	}

	if err != nil {
		return nil, err
	}

	_, err = w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
}

func getWalletType(walletPath string) string {
	tokens := strings.Split(walletPath, ".")
	if len(tokens) < 2 {
		return ""
	}

	return tokens[len(tokens)-1]
}
