package server

import (
	"context"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/testscommon"
	"github.com/stretchr/testify/require"
)

func TestNewSovereignBridgeTxServer(t *testing.T) {
	t.Parallel()

	t.Run("nil tx sender", func(t *testing.T) {
		bridgeServer, err := NewSovereignBridgeTxServer(nil)
		require.Equal(t, errNilTxSender, err)
		require.Nil(t, bridgeServer)
	})
	t.Run("should work", func(t *testing.T) {
		bridgeServer, err := NewSovereignBridgeTxServer(&testscommon.TxSenderMock{})
		require.Nil(t, err)
		require.False(t, bridgeServer.IsInterfaceNil())
	})
}

func TestServer_Send(t *testing.T) {
	t.Parallel()

	expectedTxHashes := []string{"txHash"}
	expectedBridgeOps := &sovereign.BridgeOperations{
		Data: []*sovereign.BridgeOutGoingData{
			{
				Hash: []byte("hash"),
			},
		},
	}
	txSender := &testscommon.TxSenderMock{
		SendTxsCalled: func(ctx context.Context, data *sovereign.BridgeOperations) ([]string, error) {
			require.Equal(t, expectedBridgeOps, data)
			return expectedTxHashes, nil
		},
	}

	bridgeServer, _ := NewSovereignBridgeTxServer(txSender)
	res, err := bridgeServer.Send(context.Background(), expectedBridgeOps)
	require.Nil(t, err)
	require.Equal(t, &sovereign.BridgeOperationsResponse{
		TxHashes: expectedTxHashes,
	}, res)
}
