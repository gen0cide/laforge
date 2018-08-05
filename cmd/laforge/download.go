package main

import (
	"github.com/urfave/cli"
)

var (
	downloadCommand = cli.Command{
		Name:      "download",
		Usage:     "downloads a file from a provisioned host",
		UsageText: "laforge download HOST SOURCEFILE DESTFILE",
		Action:    performdownload,
	}
)

func performdownload(c *cli.Context) error {
	return commandNotImplemented(c)
}
