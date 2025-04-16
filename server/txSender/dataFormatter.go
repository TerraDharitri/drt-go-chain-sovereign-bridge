package txSender

import (
	"bytes"
	"encoding/hex"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	"github.com/TerraDharitri/drt-go-chain-core/hashing"
)

const (
	registerBridgeOpsPrefix = "registerBridgeOps"
	executeBridgeOpsPrefix  = "executeBridgeOps"
)

type dataFormatter struct {
	hasher hashing.Hasher
}

// NewDataFormatter creates a sovereign bridge tx data formatter
func NewDataFormatter(hasher hashing.Hasher) (*dataFormatter, error) {
	if check.IfNil(hasher) {
		return nil, core.ErrNilHasher
	}

	return &dataFormatter{
		hasher: hasher,
	}, nil
}

// CreateTxsData creates txs data for bridge operations
func (df *dataFormatter) CreateTxsData(data *sovereign.BridgeOperations) [][]byte {
	txsData := make([][]byte, 0)
	if data == nil {
		return txsData
	}

	for _, bridgeData := range data.Data {
		log.Debug("creating tx data", "bridge op hash", bridgeData.Hash, "no. of operations", len(bridgeData.OutGoingOperations))

		registerBridgeOpData := df.createRegisterBridgeOperationsData(bridgeData)
		if len(registerBridgeOpData) != 0 {
			txsData = append(txsData, registerBridgeOpData)
		}

		txsData = append(txsData, createBridgeOperationsData(bridgeData.Hash, bridgeData.OutGoingOperations)...)
	}

	return txsData
}

func (df *dataFormatter) createRegisterBridgeOperationsData(bridgeData *sovereign.BridgeOutGoingData) []byte {
	hashes := make([]byte, 0)
	hashesHexEncodedArgs := make([]byte, 0)
	for _, operation := range bridgeData.OutGoingOperations {
		hashesHexEncodedArgs = append(hashesHexEncodedArgs, "@"+hex.EncodeToString(operation.Hash)...)
		hashes = append(hashes, operation.Hash...)
	}

	// unconfirmed operation, should not register it, only resend it
	computedHashOfHashes := df.hasher.Compute(string(hashes))
	if !bytes.Equal(bridgeData.Hash, computedHashOfHashes) {
		return nil
	}

	registerBridgeOpData := []byte(registerBridgeOpsPrefix +
		"@" + hex.EncodeToString(bridgeData.AggregatedSignature) +
		"@" + hex.EncodeToString(bridgeData.Hash))

	return append(registerBridgeOpData, hashesHexEncodedArgs...)
}

func createBridgeOperationsData(hashOfHashes []byte, outGoingOperations []*sovereign.OutGoingOperation) [][]byte {
	executeBridgeOpsTxData := make([][]byte, 0)
	for _, operation := range outGoingOperations {
		bridgeOpTxData := []byte(
			executeBridgeOpsPrefix +
				"@" + hex.EncodeToString(hashOfHashes) +
				"@" + hex.EncodeToString(operation.Data))

		executeBridgeOpsTxData = append(executeBridgeOpsTxData, bridgeOpTxData)
	}

	return executeBridgeOpsTxData
}

// IsInterfaceNil checks if the underlying pointer is nil
func (df *dataFormatter) IsInterfaceNil() bool {
	return df == nil
}
