package main

import (
	"os"

	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/urfave/cli"

	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/cert"
)

var log = logger.GetOrCreate("cert")

func main() {

	app := cli.NewApp()
	app.Name = "Certificate generator"
	app.Usage = "Generate certificate (.crt + .pem) for grpc tls connection between server and client.\n" +
		"->Certificate Generation: To enable secure communication, generate a certificate pair containing a .crt (certificate) " +
		"and a .pem (private key) for both the server and the sovereign nodes (clients). This will facilitate the encryption and " +
		"authentication required for the gRPC TLS connection.\n" +
		"->Authentication of Clients: The server, acting as the hot wallet binary, should authenticate and validate the sovereign nodes (clients) " +
		"attempting to connect. Only trusted clients with the matching certificate will be granted access to interact with the hot wallet binary.\n" +
		"->Ensuring Secure Transactions: Utilize the certificate-based authentication mechanism to ensure that only authorized sovereign nodes can access the hot wallet binary. " +
		"This step is crucial in maintaining the integrity and security of transactions being sent from the sovereign shards to the main chain.\n" +
		"->Ongoing Security Measures: Regularly review and update the certificate mechanism to maintain security. This includes renewal of certificates, " +
		"implementing security best practices, and promptly revoking access for compromised or unauthorized clients."
	app.Action = generateCertificate
	app.Flags = []cli.Flag{
		organizationFlag,
		dnsFlag,
		availabilityFlag,
		ipFlag,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

}

func generateCertificate(ctx *cli.Context) error {
	organization := ctx.GlobalString(organizationFlag.Name)
	dns := ctx.GlobalString(dnsFlag.Name)
	availability := ctx.GlobalInt64(availabilityFlag.Name)
	ipAddress := ctx.GlobalString(ipFlag.Name)

	err := cert.GenerateCertFiles(cert.CertificateCfg{
		CertCfg: cert.CertCfg{
			Organization: organization,
			DNSName:      dns,
			Availability: availability,
			IPAddress:    ipAddress,
		},
		CertFileCfg: cert.FileCfg{
			CertFile: "certificate.crt",
			PkFile:   "private_key.pem",
		},
	})
	if err != nil {
		return err
	}

	log.Info("generated certificate files successfully")
	return nil
}
