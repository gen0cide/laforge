package main

import (
	"github.com/gen0cide/laforge/core"
	"github.com/k0kubun/pp"
	"github.com/urfave/cli"
)

var (
	dumpCommand = cli.Command{
		Name:      "dump",
		Usage:     "dumps the current configuration state in a pretty printed output",
		UsageText: "laforge dump",
		Action:    performdump,
	}
)

func performdump(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		return err
	}
	pp.Println(base)
	return nil
}
