package main

import (
	"github.com/urfave/cli"
)

var (
	queryCommand = cli.Command{
		Name:      "query",
		Usage:     "gathers information about elements within the configuration state",
		UsageText: "laforge query QUERY",
		Action:    performquery,
	}
)

func performquery(c *cli.Context) error {
	return commandNotImplemented(c)
}
