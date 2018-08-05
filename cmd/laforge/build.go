package main

import (
	"github.com/urfave/cli"
)

var (
	buildCommand = cli.Command{
		Name:      "build",
		Usage:     "builds environment specific infrastructure configurations",
		UsageText: "laforge build",
		Action:    performbuild,
	}
)

func performbuild(c *cli.Context) error {
	return commandNotImplemented(c)
}
