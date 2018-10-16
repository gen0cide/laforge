package main

import (
	"github.com/urfave/cli"
)

var (
	exportAsJSON  = false
	exportAsCSV   = false
	exportCommand = cli.Command{
		Name:      "export",
		Usage:     "Allows for various data to be exported from different contexts.",
		UsageText: "laforge export",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "json, j",
				Usage:       "Attempt to export the information in JSON format.",
				Destination: &exportAsJSON,
			},
			cli.BoolFlag{
				Name:        "csv, c",
				Usage:       "Attempt to export the information in CSV format.",
				Destination: &exportAsCSV,
			},
		},
		Subcommands: []cli.Command{
			{
				Name:   "findings",
				Usage:  "Show all findings included in the current environment.",
				Action: exportEnvFindings,
			},
			{
				Name:   "netinfo",
				Usage:  "Export all network information for provisioned hosts in the current environment.",
				Action: exportEnvNetInfo,
			},
		},
	}
)

func exportEnvFindings(c *cli.Context) error {
	return commandNotImplemented(c)
}

func exportEnvNetInfo(c *cli.Context) error {
	return commandNotImplemented(c)
}
