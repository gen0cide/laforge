package main

import (
	"github.com/fatih/color"
	"github.com/gen0cide/laforge"
	"github.com/urfave/cli"
)

var (
	overwrite = false

	initCommand = cli.Command{
		Name:      "init",
		Usage:     "creates a base configuration in the current working directory",
		UsageText: "laforge init [OPTIONS]",
		Action:    performinit,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "force, f",
				Usage:       "force removes and deletes any conflicting directories (dangerous)",
				Destination: &overwrite,
			},
		},
	}
)

func performinit(c *cli.Context) error {
	base, err := laforge.Bootstrap()
	if err == nil && !overwrite {
		cliLogger.Errorf("Cannot initialize a competition repository - you are inside a competition base!\n\t%20s%s\n\t%20s%s", "Base Directory = ", color.HiWhiteString(base.BaseRoot), color.HiYellowString("Current Directory = "), color.HiWhiteString(base.CurrDir))
		return nil
	}
	newErr := base.InitializeBaseDirectory(overwrite)
	if newErr != nil {
		return newErr
	}
	laforge.SetLogLevel("info")
	cliLogger.Infof("Successfully initialized base competition repository.")
	return nil
}
