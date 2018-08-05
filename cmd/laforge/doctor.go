package main

import (
	"github.com/urfave/cli"
)

var (
	doctorCommand = cli.Command{
		Name:      "doctor",
		Usage:     "checks that dependency applications are installed properly",
		UsageText: "laforge doctor",
		Action:    performdoctor,
	}
)

func performdoctor(c *cli.Context) error {
	return commandNotImplemented(c)
}
