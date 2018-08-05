package main

import (
	"github.com/urfave/cli"
)

var (
	envCommand = cli.Command{
		Name:      "env",
		Usage:     "allows listing and the creation of new environments within a base competition config",
		UsageText: "laforge env [list|create]",
		Action:    performenv,
	}
)

func performenv(c *cli.Context) error {
	return commandNotImplemented(c)
}
