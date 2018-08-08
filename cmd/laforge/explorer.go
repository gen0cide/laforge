package main

import (
	"github.com/gen0cide/laforge/core"
	"github.com/gen0cide/laforge/explorer"
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
	core.SetLogLevel("info")
	base, err := core.Bootstrap()
	if err != nil {
		return err
	}
	return explorer.RenderLaforgeStatusUI(base)
}
