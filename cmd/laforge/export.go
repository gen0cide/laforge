package main

import (
	"errors"

	"github.com/gen0cide/laforge/spanner"
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
	s, err := spanner.New(nil, []string(c.Args()), "local-exec", "", false, false)
	if err != nil {
		return err
	}

	err = s.CreateWorkerPool()
	if err != nil {
		return err
	}

	err = s.Do()
	if err != nil {
		return err
	}

	return nil
}

func exportEnvNetInfo(c *cli.Context) error {
	if len(remoteHost) == 0 {
		return errors.New("must provide a target host ID using the -t flag before remote-exec")
	}
	s, err := spanner.New(nil, []string(c.Args()), "remote-exec", remoteHost, false, false)
	if err != nil {
		return err
	}

	err = s.CreateWorkerPool()
	if err != nil {
		return err
	}

	err = s.Do()
	if err != nil {
		return err
	}

	return nil
}
