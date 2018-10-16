package main

import (
	"errors"
	"os"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/core"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/urfave/cli"
)

var (
	updateConfig = false
	buildCommand = cli.Command{
		Name:      "build",
		Usage:     "builds environment specific infrastructure configurations",
		UsageText: "laforge build [OPTIONS]",
		Action:    performbuild,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "force, f",
				Usage:       "force removes and deletes any conflicting directories (dangerous)",
				Destination: &overwrite,
			},
			cli.BoolFlag{
				Name:        "update, u",
				Usage:       "Updates a build directory (expirimental)",
				Destination: &updateConfig,
			},
		},
	}
)

func performbuild(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		if _, ok := err.(hcl.Diagnostics); ok {
			return errors.New("aborted due to parsing error")
		}
		cliLogger.Errorf("Error encountered during bootstrap: %v", err)
		os.Exit(1)
	}

	bldr, err := builder.New(base, overwrite, updateConfig)
	if err != nil {
		cliLogger.Errorf("Error encountered initializing builder:\n%v", err)
		os.Exit(1)
	}

	cliLogger.Infof("Build directory initialized")

	err = bldr.Do()

	if err != nil {
		cliLogger.Errorf("Error encountered initializing builder:\n%v", err)
		os.Exit(1)
	}

	return nil
}
