package main

import (
	"github.com/urfave/cli"
)

var (
	uploadCommand = cli.Command{
		Name:      "upload",
		Usage:     "uploads a file to a provisioned host",
		UsageText: "laforge upload HOST SOURCEFILE DESTFILE",
		Action:    performupload,
	}
)

func performupload(c *cli.Context) error {
	return commandNotImplemented(c)
}
