package main

import (
	"os"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/core"
	"github.com/urfave/cli"
)

var (
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
		},
	}
)

func performbuild(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		cliLogger.Errorf("Error encountered during bootstrap: %v", err)
		os.Exit(1)
	}

	bldr, err := builder.New(base, overwrite)
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
