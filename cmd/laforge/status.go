package main

import (
	"fmt"

	"github.com/gen0cide/laforge/core"
	lfcli "github.com/gen0cide/laforge/core/cli"
	"github.com/urfave/cli"
)

var (
	short         = false
	statusCommand = cli.Command{
		Name:      "status",
		Usage:     "shows current laforge status in an ENV config compatible output",
		UsageText: "laforge status",
		Action:    performstatus,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "short, s",
				Usage:       "displays a short name of the context without formatting",
				Destination: &short,
			},
		},
	}
)

func performstatus(c *cli.Context) error {
	base := &core.Laforge{}
	err := base.InitializeContext()
	if err != nil {
		return err
	}
	if short {
		switch base.GetContext() {
		case core.GlobalContext:
			fmt.Printf("global")
		case core.BaseContext:
			fmt.Printf("base")
		case core.EnvContext:
			fmt.Printf("env")
		case core.BuildContext:
			fmt.Printf("build")
		case core.TeamContext:
			fmt.Printf("team")
		}
		return nil
	}
	lfcli.SetLogLevel("info")
	cliLogger.Infof("Current Context Level\n%s", core.StatusMap(base.GetContext()))
	return nil
}
