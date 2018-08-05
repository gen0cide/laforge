package main

import (
	"github.com/gen0cide/laforge"
	"github.com/gen0cide/laforge/statusui"
	"github.com/urfave/cli"
)

var (
	explorerCommand = cli.Command{
		Name:      "explorer",
		Usage:     "launches an expirimental terminal application for reviewing state configuration",
		UsageText: "laforge explorer",
		Action:    performexplorer,
	}
)

func performexplorer(c *cli.Context) error {
	laforge.SetLogLevel("info")
	base, err := laforge.Bootstrap()
	if err != nil {
		return err
	}
	return statusui.RenderLaforgeStatusUI(base)
}
