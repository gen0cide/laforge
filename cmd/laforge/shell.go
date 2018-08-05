package main

import (
	"github.com/urfave/cli"
)

var (
	shellCommand = cli.Command{
		Name:      "shell",
		Usage:     "launches an interactive SSH or Powershell console on a provisioned host",
		UsageText: "laforge shell",
		Action:    performshell,
	}
)

func performshell(c *cli.Context) error {
	return commandNotImplemented(c)
}
