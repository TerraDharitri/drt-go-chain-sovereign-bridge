package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"

	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/cert"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/client"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/client/config"
)

var log = logger.GetOrCreate("client-tx-sender")

const (
	envGRPCHost   = "GRPC_HOST"
	envGRPCPort   = "GRPC_PORT"
	envCertFile   = "CERT_FILE"
	envCertPkFile = "CERT_PK_FILE"
)

func main() {
	app := cli.NewApp()
	app.Name = "GRPC Client"
	app.Action = startClient
	app.Flags = []cli.Flag{
		logLevel,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startClient(ctx *cli.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	err = initializeLogger(ctx)
	if err != nil {
		return err
	}

	log.Info("starting client...")

	bridgeClient, err := client.CreateClient(cfg)
	if err != nil {
		return err
	}

	defer func() {
		err = bridgeClient.Close()
		log.LogIfError(err)
	}()

	return sendData(bridgeClient)
}

func sendData(bridgeClient client.ClientHandler) error {
	txHashes := make(map[string]struct{})
	mut := sync.RWMutex{}

	numBridgeOps := 5
	expectedNumBridgeTxs := 3 * numBridgeOps
	wg := sync.WaitGroup{}
	wg.Add(expectedNumBridgeTxs)

	for i := 0; i < numBridgeOps; i++ {
		hash := []byte(fmt.Sprintf("hash_%d", i))
		log.Info("sending data", "hash", hash)

		go func() {
			res, errSend := bridgeClient.Send(context.Background(), &sovereign.BridgeOperations{
				Data: []*sovereign.BridgeOutGoingData{
					{
						Hash: hash,
						OutGoingOperations: []*sovereign.OutGoingOperation{
							{
								Hash: []byte("opHash1"),
								Data: []byte("bridgeOp1"),
							},
							{
								Hash: []byte("opHash2"),
								Data: []byte("bridgeOp2"),
							},
						},
						AggregatedSignature: []byte("aggregatedSig"),
						LeaderSignature:     []byte("leaderSig"),
					},
				},
			})
			if errSend != nil {
				log.Error("error sending bridge data", "error", errSend)
				wg.Done()
				return
			}

			addTxHashes(res.TxHashes, txHashes, &mut, &wg)
		}()
	}

	wg.Wait()

	numSentTxs := len(txHashes)
	if numSentTxs != expectedNumBridgeTxs {
		return fmt.Errorf("did not send all txs; expected num send txs: %d, received: %d",
			expectedNumBridgeTxs, numSentTxs)
	}

	return nil
}

func addTxHashes(txHashes []string, txHashesMap map[string]struct{}, mut *sync.RWMutex, wg *sync.WaitGroup) {
	for _, txHash := range txHashes {
		log.Info("received", "tx hash", txHash)

		mut.Lock()
		txHashesMap[txHash] = struct{}{}
		mut.Unlock()

		wg.Done()
	}
}

func loadConfig() (*config.ClientConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	grpcHost := os.Getenv(envGRPCHost)
	grpcPort := os.Getenv(envGRPCPort)
	certFile := os.Getenv(envCertFile)
	certPkFile := os.Getenv(envCertPkFile)

	log.Info("loaded config", "grpc host", grpcHost)
	log.Info("loaded config", "grpc port", grpcPort)

	log.Info("loaded config", "certificate file", certFile)
	log.Info("loaded config", "certificate pk", certPkFile)

	return &config.ClientConfig{
		Enabled:  true,
		GRPCHost: grpcHost,
		GRPCPort: grpcPort,
		CertificateCfg: cert.FileCfg{
			CertFile: certFile,
			PkFile:   certPkFile,
		},
	}, nil
}

func initializeLogger(ctx *cli.Context) error {
	logLevelFlagValue := ctx.GlobalString(logLevel.Name)
	return logger.SetLogLevel(logLevelFlagValue)
}
