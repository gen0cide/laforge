package main

import (
	"github.com/urfave/cli"
)

var (
	serveCommand = cli.Command{
		Name:      "serve",
		Usage:     "starts an HTTP server that can be used to serve files for local provisioning",
		UsageText: "laforge serve",
		Action:    performserve,
	}
)

func performserve(c *cli.Context) error {
	return commandNotImplemented(c)
}
