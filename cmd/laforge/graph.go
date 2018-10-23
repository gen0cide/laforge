package main

import (
	"errors"

	"github.com/gen0cide/laforge/core"
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
	state, err := core.BootstrapWithState(true)
	if err != nil {
		return err
	}
	if state == nil {
		return errors.New("cannot proceed with a nil state")
	}

	plan, err := state.CalculateDelta()
	if err != nil {
		return err
	}
	tfcmds, err := core.CalculateTerraformNeeds(plan)
	if err != nil {
		return err
	}

	_ = tfcmds

	snap := state.Current
	defer state.DB.Close()
	_ = snap

	return nil
}
