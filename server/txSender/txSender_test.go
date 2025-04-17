package txSender

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/testscommon"

	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
	"github.com/stretchr/testify/require"
)

const (
	scHeaderVerifierAddress = "drt1qqq"
	scDcdtSafeAddress       = "drt1qqqe"
)

func createArgs() TxSenderArgs {
	return TxSenderArgs{
		Wallet:                  &testscommon.CryptoComponentsHolderMock{},
		Proxy:                   &testscommon.ProxyMock{},
		TxInteractor:            &testscommon.TxInteractorMock{},
		DataFormatter:           &testscommon.DataFormatterMock{},
		TxNonceHandler:          &testscommon.TxNonceSenderHandlerMock{},
		SCHeaderVerifierAddress: scHeaderVerifierAddress,
		SCDcdtSafeAddress:       scDcdtSafeAddress,
	}
}

func TestNewTxSender(t *testing.T) {
	t.Parallel()

	t.Run("nil wallet", func(t *testing.T) {
		args := createArgs()
		args.Wallet = nil

		ts, err := NewTxSender(args)
		require.Nil(t, ts)
		require.Equal(t, errNilWallet, err)
	})
	t.Run("nil proxy", func(t *testing.T) {
		args := createArgs()
		args.Proxy = nil

		ts, err := NewTxSender(args)
		require.Nil(t, ts)
		require.Equal(t, errNilProxy, err)
	})
	t.Run("nil tx interactor", func(t *testing.T) {
		args := createArgs()
		args.TxInteractor = nil

		ts, err := NewTxSender(args)
		require.Nil(t, ts)
		require.Equal(t, errNilTxInteractor, err)
	})
	t.Run("nil data formatter", func(t *testing.T) {
		args := createArgs()
		args.DataFormatter = nil

		ts, err := NewTxSender(args)
		require.Nil(t, ts)
		require.Equal(t, errNilDataFormatter, err)
	})
	t.Run("should work", func(t *testing.T) {
		args := createArgs()

		ts, err := NewTxSender(args)
		require.Nil(t, err)
		require.False(t, ts.IsInterfaceNil())
	})
}

func TestTxSender_SendTxs(t *testing.T) {
	t.Parallel()

	expectedCtx := context.Background()
	expectedNonce := 0
	expectedDataIdx := 0
	expectedTxHashes := []string{"txHash1", "txHash2", "txHash3"}
	expectedTxsData := [][]byte{
		[]byte(registerBridgeOpsPrefix + "txData1"),
		[]byte(executeBridgeOpsPrefix + "txData2"),
		[]byte(executeBridgeOpsPrefix + "txData3"),
	}
	expectedTxsReceiver := []string{
		scHeaderVerifierAddress,
		scDcdtSafeAddress,
		scDcdtSafeAddress,
	}
	expectedSigs := []string{"sig1", "sig2", "sig3"}
	expectedBridgeData := &sovereign.BridgeOperations{
		Data: []*sovereign.BridgeOutGoingData{
			{
				Hash: []byte("bridgeDataHash1"),
			},
		},
	}
	expectedNetworkConfig := &data.NetworkConfig{
		MinGasPrice:           5000,
		ChainID:               "1",
		MinTransactionVersion: 2,
	}

	args := createArgs()
	args.Proxy = &testscommon.ProxyMock{
		GetNetworkConfigCalled: func(ctx context.Context) (*data.NetworkConfig, error) {
			require.Equal(t, expectedCtx, ctx)
			return expectedNetworkConfig, nil
		},
	}
	args.DataFormatter = &testscommon.DataFormatterMock{
		CreateTxsDataCalled: func(data *sovereign.BridgeOperations) [][]byte {
			require.Equal(t, expectedBridgeData, data)
			return expectedTxsData
		},
	}
	args.TxInteractor = &testscommon.TxInteractorMock{
		ApplyUserSignatureCalled: func(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
			tx.Signature = expectedSigs[expectedDataIdx]
			return nil
		},
	}
	args.TxNonceHandler = &testscommon.TxNonceSenderHandlerMock{
		ApplyNonceAndGasPriceCalled: func(ctx context.Context, txs ...*transaction.FrontendTransaction) error {
			require.Len(t, txs, 1) // we update transactions one at a time
			require.Equal(t, &transaction.FrontendTransaction{
				Nonce:    0,
				Value:    "0",
				Receiver: expectedTxsReceiver[expectedDataIdx],
				Sender:   args.Wallet.GetBech32(),
				GasPrice: expectedNetworkConfig.MinGasPrice,
				GasLimit: 50_000_000,
				Data:     expectedTxsData[expectedDataIdx],
				ChainID:  expectedNetworkConfig.ChainID,
				Version:  expectedNetworkConfig.MinTransactionVersion,
			}, txs[0])

			expectedNonce++
			txs[0].Nonce = uint64(expectedNonce)
			return nil
		},
		SendTransactionsCalled: func(ctx context.Context, txs ...*transaction.FrontendTransaction) ([]string, error) {
			defer func() {
				expectedDataIdx++
			}()

			require.Equal(t, expectedCtx, ctx)
			require.Len(t, txs, 1) // we send transactions one at a time
			require.Equal(t, &transaction.FrontendTransaction{
				Nonce:     uint64(expectedNonce),
				Value:     "0",
				Receiver:  expectedTxsReceiver[expectedDataIdx],
				Sender:    args.Wallet.GetBech32(),
				GasPrice:  expectedNetworkConfig.MinGasPrice,
				GasLimit:  50_000_000,
				Data:      expectedTxsData[expectedDataIdx],
				Signature: expectedSigs[expectedDataIdx],
				ChainID:   expectedNetworkConfig.ChainID,
				Version:   expectedNetworkConfig.MinTransactionVersion,
			}, txs[0])

			return []string{expectedTxHashes[expectedDataIdx]}, nil
		},
	}

	ts, _ := NewTxSender(args)
	txHashes, err := ts.SendTxs(expectedCtx, expectedBridgeData)
	require.Nil(t, err)
	require.Equal(t, expectedTxHashes, txHashes)
	require.Equal(t, 3, expectedNonce)
	require.Equal(t, 3, expectedDataIdx)
}

func TestTxSender_SendTxsConcurrently(t *testing.T) {
	t.Parallel()

	args := createArgs()

	expectedTxHashes := []string{"hash"}
	numTxsToSend := 1000
	numSentTxs := 0

	mut := sync.RWMutex{}
	wg := sync.WaitGroup{}
	wg.Add(numTxsToSend)

	args.DataFormatter = &testscommon.DataFormatterMock{
		CreateTxsDataCalled: func(data *sovereign.BridgeOperations) [][]byte {
			return [][]byte{[]byte(executeBridgeOpsPrefix + "txData")}
		},
	}
	args.TxNonceHandler = &testscommon.TxNonceSenderHandlerMock{
		SendTransactionsCalled: func(ctx context.Context, txs ...*transaction.FrontendTransaction) ([]string, error) {
			defer func() {
				numSentTxs++
				wg.Done()
				mut.Unlock()
			}()

			mut.Lock()
			return []string{"hash"}, nil
		},
	}

	ts, _ := NewTxSender(args)

	for i := 0; i < numTxsToSend; i++ {
		go func(idx int) {
			txHashes, err := ts.SendTxs(context.Background(), &sovereign.BridgeOperations{
				Data: []*sovereign.BridgeOutGoingData{
					{
						Hash: []byte(fmt.Sprintf("hash%d", idx)),
					},
				},
			})
			require.Nil(t, err)
			require.Equal(t, expectedTxHashes, txHashes)
		}(i)
	}

	wg.Wait()
	require.Equal(t, numTxsToSend, numSentTxs)
}
