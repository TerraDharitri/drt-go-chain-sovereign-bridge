package txSender

import (
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/hashing/factory"
	"github.com/TerraDharitri/drt-go-sdk/blockchain"
	"github.com/TerraDharitri/drt-go-sdk/blockchain/cryptoProvider"
	"github.com/TerraDharitri/drt-go-sdk/builders"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/interactors"
	"github.com/TerraDharitri/drt-go-sdk/interactors/nonceHandlerV3"
)

// CreateTxSender creates a new transactions sender
func CreateTxSender(wallet core.CryptoComponentsHolder, cfg TxSenderConfig) (*txSender, error) {
	args := blockchain.ArgsProxy{
		ProxyURL:            cfg.Proxy,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
	proxy, err := blockchain.NewProxy(args)
	if err != nil {
		return nil, err
	}

	nonceHandler, err := nonceHandlerV3.NewNonceTransactionHandlerV3(nonceHandlerV3.ArgsNonceTransactionsHandlerV3{
		Proxy:          proxy,
		IntervalToSend: time.Millisecond * time.Duration(cfg.IntervalToSend),
	})
	if err != nil {
		return nil, err
	}

	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		return nil, err
	}

	ti, err := interactors.NewTransactionInteractor(proxy, txBuilder)
	if err != nil {
		return nil, err
	}

	hasher, err := factory.NewHasher(cfg.Hasher)
	if err != nil {
		return nil, err
	}

	dtaFormatter, err := NewDataFormatter(hasher)
	if err != nil {
		return nil, err
	}

	return NewTxSender(TxSenderArgs{
		Wallet:                  wallet,
		Proxy:                   proxy,
		TxInteractor:            ti,
		TxNonceHandler:          nonceHandler,
		DataFormatter:           dtaFormatter,
		SCHeaderVerifierAddress: cfg.HeaderVerifierSCAddress,
		SCDcdtSafeAddress:       cfg.DcdtSafeSCAddress,
	})
}
