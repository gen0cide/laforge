package main

import (
	"github.com/urfave/cli"
)

var (
	statusCommand = cli.Command{
		Name:      "status",
		Usage:     "shows current laforge status in an ENV config compatible output",
		UsageText: "laforge status",
		Action:    performstatus,
	}
)

func performstatus(c *cli.Context) error {
	return commandNotImplemented(c)
}
