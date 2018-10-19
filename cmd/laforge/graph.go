package main

import (
	"errors"
	"os"
	"path"

	"github.com/gen0cide/laforge/core"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/urfave/cli"
)

var (
	graphCommand = cli.Command{
		Name:      "graph",
		Usage:     "Generates graphviz diagrams of the state graph within the current build.",
		UsageText: "laforge graph [ROOT]",
		Action:    performgraph,
	}
)

func performgraph(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		if _, ok := err.(hcl.Diagnostics); ok {
			return errors.New("aborted due to parsing error")
		}
		cliLogger.Errorf("Error encountered during bootstrap: %v", err)
		os.Exit(1)
	}

	err = base.AssertMinContext(core.BuildContext)
	if err != nil {
		cliLogger.Errorf("Must be in a team context to use this command: %v", err)
		os.Exit(1)
	}

	snap, err := core.NewSnapshotFromEnv(base.CurrentEnv)
	if err != nil {
		return err
	}

	build, ok := snap.Objects[path.Join(base.CurrentEnv.Path(), base.CurrentEnv.Builder)].(*core.Build)
	if !ok {
		return errors.New("builder was not able to resolve object of type Build")
	}
	base.CurrentBuild = build

	err = snap.Sort()
	if err != nil {
		panic(err)
	}

	arg := c.Args().First()
	snap.Plot(os.Stdout, arg)

	return nil
}
