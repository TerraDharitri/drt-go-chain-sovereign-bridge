package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/cert"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/server"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/server/cmd/config"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/server/txSender"

	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/core/closing"
	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	"github.com/TerraDharitri/drt-go-chain-core/marshal"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain-logger/file"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logger.GetOrCreate("sov-bridge-sender")

const (
	retrialTimeServe = 1
	logsPath         = "logs"
	logsPrefix       = "sov-bridge-sender"
	logLifeSpanMb    = 1024   //# 1GB
	logLifeSpanSec   = 432000 // 5 days
)

const (
	envGRPCPort             = "GRPC_PORT"
	envWallet               = "WALLET_PATH"
	envPassword             = "WALLET_PASSWORD"
	envHeaderVerifierSCAddr = "HEADER_VERIFIER_SC_ADDRESS"
	envDcdtSafeSCAddr       = "DCDT_SAFE_SC_ADDRESS"
	envDharitriProxy        = "DHARITRI_PROXY"
	envIntervalToSend       = "INTERVAL_TO_SEND"
	envCertFile             = "CERT_FILE"
	envCertPkFile           = "CERT_PK_FILE"
	envHasher               = "HASHER"
)

func main() {
	app := cli.NewApp()
	app.Name = "Sovereign bridge tx server"
	app.Action = startServer
	app.Flags = []cli.Flag{
		logLevel,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startServer(ctx *cli.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	logFile, err := initializeLogger(ctx)
	if err != nil {
		return err
	}

	tlsConfig, err := cert.LoadTLSServerConfig(cfg.CertificateConfig)
	if err != nil {
		return err
	}

	tlsCredentials := credentials.NewTLS(tlsConfig)
	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)
	bridgeServer, err := server.CreateSovereignBridgeServer(cfg)
	if err != nil {
		return err
	}

	sovereign.RegisterBridgeTxSenderServer(grpcServer, bridgeServer)
	log.Info("starting server...")

	ginHandler, err := server.NewGinHandler(&marshal.GogoProtoMarshalizer{})
	if err != nil {
		return err
	}

	serverHandler, err := server.NewServerHandler(ginHandler, grpcServer)
	if err != nil {
		return err
	}

	go func() {
		for {
			err = http.ListenAndServeTLS(
				fmt.Sprintf(":%s", cfg.GRPCPort),
				cfg.CertificateConfig.CertFile,
				cfg.CertificateConfig.PkFile,
				serverHandler,
			)
			if err != nil {
				log.Error("sovereign bridge tx sender: ListenAndServeTLS", "error", err)
				time.Sleep(retrialTimeServe * time.Second)
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interrupt
	log.Info("closing app at user's signal")

	grpcServer.Stop()

	if !check.IfNilReflect(logFile) {
		err = logFile.Close()
		log.LogIfError(err)
	}

	return nil
}

func loadConfig() (*config.ServerConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	grpcPort := os.Getenv(envGRPCPort)
	walletPath := os.Getenv(envWallet)
	walletPassword := os.Getenv(envPassword)
	headerVerifierSCAddress := os.Getenv(envHeaderVerifierSCAddr)
	dcdtSafeSCAddress := os.Getenv(envDcdtSafeSCAddr)
	proxy := os.Getenv(envDharitriProxy)
	intervalToSendStr := os.Getenv(envIntervalToSend)
	certFile := os.Getenv(envCertFile)
	certPkFile := os.Getenv(envCertPkFile)
	hasher := os.Getenv(envHasher)

	intervalToSend, err := strconv.Atoi(intervalToSendStr)
	if err != nil {
		return nil, err
	}

	log.Info("loaded config", "grpc port", grpcPort)
	log.Info("loaded config", "headerVerifierSCAddress", headerVerifierSCAddress)
	log.Info("loaded config", "dcdtSafeSCAddress", dcdtSafeSCAddress)
	log.Info("loaded config", "proxy", proxy)
	log.Info("loaded config", "intervalToSend", intervalToSend)
	log.Info("loaded config", "hasher", hasher)

	log.Info("loaded config", "certificate file", certFile)
	log.Info("loaded config", "certificate pk", certPkFile)

	return &config.ServerConfig{
		GRPCPort: grpcPort,
		WalletConfig: txSender.WalletConfig{
			Path:     walletPath,
			Password: walletPassword,
		},
		TxSenderConfig: txSender.TxSenderConfig{
			HeaderVerifierSCAddress: headerVerifierSCAddress,
			DcdtSafeSCAddress:       dcdtSafeSCAddress,
			Proxy:                   proxy,
			IntervalToSend:          intervalToSend,
			Hasher:                  hasher,
		},
		CertificateConfig: cert.FileCfg{
			CertFile: certFile,
			PkFile:   certPkFile,
		},
	}, nil
}

func initializeLogger(ctx *cli.Context) (closing.Closer, error) {
	logLevelFlagValue := ctx.GlobalString(logLevel.Name)
	err := logger.SetLogLevel(logLevelFlagValue)
	if err != nil {
		return nil, err
	}

	withLogFile := ctx.GlobalBool(logSaveFile.Name)
	if !withLogFile {
		return nil, nil
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.LogIfError(err)
		workingDir = ""
	}

	fileLogging, err := file.NewFileLogging(file.ArgsFileLogging{
		WorkingDir:      workingDir,
		DefaultLogsPath: logsPath,
		LogFilePrefix:   logsPrefix,
	})
	if err != nil {
		return nil, fmt.Errorf("%w creating a log file", err)
	}

	err = fileLogging.ChangeFileLifeSpan(
		time.Second*time.Duration(logLifeSpanSec),
		uint64(logLifeSpanMb),
	)
	if err != nil {
		return nil, err
	}

	disableAnsi := ctx.GlobalBool(disableAnsiColor.Name)
	err = removeANSIColorsForLoggerIfNeeded(disableAnsi)
	if err != nil {
		return nil, err
	}

	return fileLogging, nil
}

func removeANSIColorsForLoggerIfNeeded(disableAnsi bool) error {
	if !disableAnsi {
		return nil
	}

	err := logger.RemoveLogObserver(os.Stdout)
	if err != nil {
		return err
	}

	return logger.AddLogObserver(os.Stdout, &logger.PlainFormatter{})
}
