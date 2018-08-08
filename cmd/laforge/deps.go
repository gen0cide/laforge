package main

import (
	"fmt"

	"github.com/gen0cide/laforge/core"
	"github.com/urfave/cli"
)

var (
	depsCommand = cli.Command{
		Name:      "deps",
		Usage:     "prints a tree of laforge dependencies and their load preference",
		UsageText: "laforge deps",
		Action:    performdeps,
	}
)

func performdeps(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		return err
	}
	core.SetLogLevel("info")
	cliLogger.Infof("== Dependency Graph ==")
	fmt.Println(base.DependencyGraph.String())
	return nil
}
