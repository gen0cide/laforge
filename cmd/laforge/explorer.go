package main

import (
	"github.com/urfave/cli"
)

var (
	explorerCommand = cli.Command{
		Name:      "explorer",
		Usage:     "launches an expirimental terminal application for reviewing state configuration",
		UsageText: "laforge explorer",
		Action:    performexplorer,
	}
)

func performexplorer(c *cli.Context) error {
	return commandNotImplemented(c)
}
