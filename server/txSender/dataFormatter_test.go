package txSender

import (
	"encoding/hex"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/testscommon"
	"github.com/stretchr/testify/require"
)

func TestNewDataFormatter(t *testing.T) {
	t.Parallel()

	t.Run("nil hasher, should fail", func(t *testing.T) {
		df, err := NewDataFormatter(nil)
		require.Equal(t, core.ErrNilHasher, err)
		require.Nil(t, df)
	})

	t.Run("should work", func(t *testing.T) {
		df, err := NewDataFormatter(&testscommon.HasherMock{})
		require.Nil(t, err)
		require.False(t, df.IsInterfaceNil())
	})
}

func TestDataFormatter_CreateTxsData(t *testing.T) {
	t.Parallel()

	t.Run("nil input, should return empty result", func(t *testing.T) {
		df, _ := NewDataFormatter(&testscommon.HasherMock{})
		require.Empty(t, df.CreateTxsData(nil))
	})

	t.Run("empty data, should return empty result", func(t *testing.T) {
		df, _ := NewDataFormatter(&testscommon.HasherMock{})
		require.Empty(t, df.CreateTxsData(&sovereign.BridgeOperations{Data: nil}))
	})

	t.Run("should work to create register and execute txs data", func(t *testing.T) {
		bridgeDataHash1 := []byte("bridgeDataHash1")
		bridgeDataHash2 := []byte("bridgeDataHash2")

		aggregatedSig1 := []byte("aggregatedSig1")
		aggregatedSig2 := []byte("aggregatedSig2")

		leaderSig1 := []byte("leaderSig1")
		leaderSig2 := []byte("leaderSig2")

		opHash1 := []byte("outGoingOpHash1")
		opHash2 := []byte("outGoingOpHash2")
		opHash3 := []byte("outGoingOpHash3")

		bridgeDataOp1 := []byte("bridgeDataOp1")
		bridgeDataOp2 := []byte("bridgeDataOp2")
		bridgeDataOp3 := []byte("bridgeDataOp3")

		bridgeOps := &sovereign.BridgeOperations{
			Data: []*sovereign.BridgeOutGoingData{
				{
					Hash: bridgeDataHash1,
					OutGoingOperations: []*sovereign.OutGoingOperation{
						{
							Hash: opHash1,
							Data: bridgeDataOp1,
						},
						{
							Hash: opHash2,
							Data: bridgeDataOp2,
						},
					},
					AggregatedSignature: aggregatedSig1,
					LeaderSignature:     leaderSig1,
				},
				{
					Hash: bridgeDataHash2,
					OutGoingOperations: []*sovereign.OutGoingOperation{
						{
							Hash: opHash3,
							Data: bridgeDataOp3,
						},
					},
					AggregatedSignature: aggregatedSig2,
					LeaderSignature:     leaderSig2,
				},
			},
		}

		computeHashCt := 0
		hasher := &testscommon.HasherMock{
			ComputeCalled: func(s string) []byte {
				defer func() {
					computeHashCt++
				}()

				switch computeHashCt {
				case 0:
					require.Equal(t, string(append(opHash1, opHash2...)), s)
					return bridgeDataHash1
				case 1:
					require.Equal(t, string(opHash3), s)
					return bridgeDataHash2
				default:
					require.Fail(t, "should have not compute another hash")
				}

				return nil
			},
		}
		df, _ := NewDataFormatter(hasher)

		registerOp1 := []byte(
			registerBridgeOpsPrefix +
				"@" + hex.EncodeToString(aggregatedSig1) +
				"@" + hex.EncodeToString(bridgeDataHash1) +
				"@" + hex.EncodeToString(opHash1) +
				"@" + hex.EncodeToString(opHash2))
		execOp1 := []byte(executeBridgeOpsPrefix +
			"@" + hex.EncodeToString(bridgeDataHash1) +
			"@" + hex.EncodeToString(bridgeDataOp1))
		execOp2 := []byte(executeBridgeOpsPrefix +
			"@" + hex.EncodeToString(bridgeDataHash1) +
			"@" + hex.EncodeToString(bridgeDataOp2))

		registerOp2 := []byte(
			registerBridgeOpsPrefix +
				"@" + hex.EncodeToString(aggregatedSig2) +
				"@" + hex.EncodeToString(bridgeDataHash2) +
				"@" + hex.EncodeToString(opHash3))
		execOp3 := []byte(executeBridgeOpsPrefix +
			"@" + hex.EncodeToString(bridgeDataHash2) +
			"@" + hex.EncodeToString(bridgeDataOp3))

		expectedTxsData := [][]byte{
			registerOp1,
			execOp1,
			execOp2,
			registerOp2,
			execOp3,
		}

		txsData := df.CreateTxsData(bridgeOps)
		require.Equal(t, expectedTxsData, txsData)
		require.Equal(t, computeHashCt, 2)
	})

	t.Run("computed hash != received hash, should only format execute operations, without register", func(t *testing.T) {
		bridgeDataHash1 := []byte("bridgeDataHash1")
		aggregatedSig1 := []byte("aggregatedSig1")
		leaderSig1 := []byte("leaderSig1")

		opHash1 := []byte("outGoingOpHash1")
		opHash2 := []byte("outGoingOpHash2")

		bridgeDataOp1 := []byte("bridgeDataOp1")
		bridgeDataOp2 := []byte("bridgeDataOp2")

		bridgeOps := &sovereign.BridgeOperations{
			Data: []*sovereign.BridgeOutGoingData{
				{
					Hash: bridgeDataHash1,
					OutGoingOperations: []*sovereign.OutGoingOperation{
						{
							Hash: opHash1,
							Data: bridgeDataOp1,
						},
						{
							Hash: opHash2,
							Data: bridgeDataOp2,
						},
					},
					AggregatedSignature: aggregatedSig1,
					LeaderSignature:     leaderSig1,
				},
			},
		}

		computeHashCt := 0
		hasher := &testscommon.HasherMock{
			ComputeCalled: func(s string) []byte {
				defer func() {
					computeHashCt++
				}()

				switch computeHashCt {
				case 0:
					require.Equal(t, string(append(opHash1, opHash2...)), s)
					return []byte("another hash")
				default:
					require.Fail(t, "should have not compute another hash")
				}

				return nil
			},
		}
		df, _ := NewDataFormatter(hasher)

		execOp1 := []byte(executeBridgeOpsPrefix +
			"@" + hex.EncodeToString(bridgeDataHash1) +
			"@" + hex.EncodeToString(bridgeDataOp1))
		execOp2 := []byte(executeBridgeOpsPrefix +
			"@" + hex.EncodeToString(bridgeDataHash1) +
			"@" + hex.EncodeToString(bridgeDataOp2))

		expectedTxsData := [][]byte{
			execOp1,
			execOp2,
		}

		txsData := df.CreateTxsData(bridgeOps)
		require.Equal(t, expectedTxsData, txsData)
		require.Equal(t, computeHashCt, 1)
	})
}
