package main

import (
	"github.com/urfave/cli"
)

var (
	infraCommand = cli.Command{
		Name:      "infra",
		Usage:     "Manage infrastructure deployment that has been generated with Laforge.",
		UsageText: "laforge infra",
		Subcommands: []cli.Command{
			{
				Name:            "plan",
				Usage:           "Show the delta between current deployment and the final desired state.",
				Action:          performinfra,
				SkipFlagParsing: true,
			},
			{
				Name:            "status",
				Usage:           "Show the current build's infrastructure status.",
				Action:          performinfra,
				SkipFlagParsing: true,
			},
			{
				Name:            "apply",
				Usage:           "Provision the infrastructure to bring state in line with build blueprint.",
				Action:          performinfra,
				SkipFlagParsing: true,
			},
			{
				Name:            "taint",
				Usage:           "Mark a host for re-provisioning in the laforge infrastructure state.",
				Action:          performinfra,
				SkipFlagParsing: true,
			},
			{
				Name:            "destroy",
				Usage:           "Destroy the builds infrastructure and clean the state.",
				Action:          performinfra,
				SkipFlagParsing: true,
			},
			{
				Name:            "run",
				Usage:           "Run a host provisioner on a specific host within the infrastructure (usually for debugging).",
				Action:          performinfra,
				SkipFlagParsing: true,
			},
		},
	}
)

func performinfra(c *cli.Context) error {
	return commandNotImplemented(c)
}
