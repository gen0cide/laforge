package main

import (
	"github.com/gen0cide/laforge/core"
	"github.com/urfave/cli"
)

var (
	statusCommand = cli.Command{
		Name:      "status",
		Usage:     "shows current laforge status in an ENV config compatible output",
		UsageText: "laforge status",
		Action:    performstatus,
	}
)

func performstatus(c *cli.Context) error {
	base, _ := core.Bootstrap()

	core.SetLogLevel("info")
	cliLogger.Infof("Current Context Level\n%s", core.StatusMap(base.GetContext()))
	return nil
}
