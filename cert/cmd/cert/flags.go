package main

import "github.com/urfave/cli"

var (
	organizationFlag = cli.StringFlag{
		Name:  "organization",
		Usage: "This flag specifies the organization name which will generate the certificate",
		Value: "Dharitri",
	}
	dnsFlag = cli.StringFlag{
		Name:  "dns",
		Usage: "This flag specifies the server's dns for tls connection",
		Value: "localhost",
	}
	availabilityFlag = cli.StringFlag{
		Name:  "availability",
		Usage: "This flag specifies the certificate's availability in days starting from current timestamp",
		Value: "365",
	}
	ipFlag = cli.StringFlag{
		Name:  "ip",
		Usage: "This flag specifies the certificate IP address",
		Value: "127.0.0.1",
	}
)
