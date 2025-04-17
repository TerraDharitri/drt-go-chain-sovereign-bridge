package txSender

import (
	"context"
	"strings"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	coreTx "github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
)

// TxSenderArgs holds args to create a new tx sender
type TxSenderArgs struct {
	Wallet                  core.CryptoComponentsHolder
	Proxy                   Proxy
	TxInteractor            TxInteractor
	TxNonceHandler          TxNonceSenderHandler
	DataFormatter           DataFormatter
	SCHeaderVerifierAddress string
	SCDcdtSafeAddress       string
}

type txSender struct {
	wallet                  core.CryptoComponentsHolder
	netConfigs              *data.NetworkConfig
	txInteractor            TxInteractor
	txNonceHandler          TxNonceSenderHandler
	dataFormatter           DataFormatter
	scHeaderVerifierAddress string
	scDcdtSafeAddress       string
}

// NewTxSender creates a new tx sender
func NewTxSender(args TxSenderArgs) (*txSender, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	networkConfig, err := args.Proxy.GetNetworkConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &txSender{
		wallet:                  args.Wallet,
		netConfigs:              networkConfig,
		txInteractor:            args.TxInteractor,
		txNonceHandler:          args.TxNonceHandler,
		dataFormatter:           args.DataFormatter,
		scHeaderVerifierAddress: args.SCHeaderVerifierAddress,
		scDcdtSafeAddress:       args.SCDcdtSafeAddress,
	}, nil
}

func checkArgs(args TxSenderArgs) error {
	if check.IfNil(args.Wallet) {
		return errNilWallet
	}
	if check.IfNil(args.Proxy) {
		return errNilProxy
	}
	if check.IfNil(args.TxInteractor) {
		return errNilTxInteractor
	}
	if check.IfNil(args.DataFormatter) {
		return errNilDataFormatter
	}
	if check.IfNil(args.TxNonceHandler) {
		return errNilNonceHandler
	}
	if len(args.SCHeaderVerifierAddress) == 0 {
		return errNoHeaderVerifierSCAddress
	}
	if len(args.SCDcdtSafeAddress) == 0 {
		return errNoDcdtSafeSCAddress
	}

	return nil
}

// SendTxs should send bridge data operation txs
func (ts *txSender) SendTxs(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
	if len(data.Data) == 0 {
		return make([]string, 0), nil
	}

	return ts.createAndSendTxs(ctx, data)
}

func (ts *txSender) createAndSendTxs(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
	txHashes := make([]string, 0)
	txsData := ts.dataFormatter.CreateTxsData(data)

	for _, txData := range txsData {
		var tx *coreTx.FrontendTransaction

		switch {
		case strings.HasPrefix(string(txData), registerBridgeOpsPrefix):
			tx = &coreTx.FrontendTransaction{
				Value:    "0",
				Receiver: ts.scHeaderVerifierAddress,
				Sender:   ts.wallet.GetBech32(),
				GasPrice: ts.netConfigs.MinGasPrice,
				GasLimit: 50_000_000, // todo
				Data:     txData,
				ChainID:  ts.netConfigs.ChainID,
				Version:  ts.netConfigs.MinTransactionVersion,
			}
		case strings.HasPrefix(string(txData), executeBridgeOpsPrefix):
			tx = &coreTx.FrontendTransaction{
				Value:    "0",
				Receiver: ts.scDcdtSafeAddress,
				Sender:   ts.wallet.GetBech32(),
				GasPrice: ts.netConfigs.MinGasPrice,
				GasLimit: 50_000_000, // todo
				Data:     txData,
				ChainID:  ts.netConfigs.ChainID,
				Version:  ts.netConfigs.MinTransactionVersion,
			}
		default:
			log.Error("invalid tx data received", "data", string(txData))
			continue
		}

		err := ts.txNonceHandler.ApplyNonceAndGasPrice(ctx, tx)
		if err != nil {
			return nil, err
		}

		err = ts.txInteractor.ApplyUserSignature(ts.wallet, tx)
		if err != nil {
			return nil, err
		}

		hash, err := ts.txNonceHandler.SendTransactions(ctx, tx)
		if err != nil {
			log.Error("failed to send tx", "error", err, "nonce", tx.Nonce)
			return nil, err
		}

		txHashes = append(txHashes, hash...)
	}

	return txHashes, nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (ts *txSender) IsInterfaceNil() bool {
	return ts == nil
}
