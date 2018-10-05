package main

import (
	"github.com/gen0cide/laforge/core"
	"github.com/urfave/cli"
)

var (
	copmetitionCommand = cli.Command{
		Name:      "competition",
		Usage:     "shows current laforge status in an ENV config compatible output",
		UsageText: "laforge status",
		Action:    competitioncommand,
	}
)

func competitioncommand(c *cli.Context) error {
	base, _ := core.Bootstrap()

	core.SetLogLevel("info")
	cliLogger.Infof("Current Context Level\n%s", core.StatusMap(base.GetContext()))
	return nil
}
