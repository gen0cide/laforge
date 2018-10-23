package main

import (
	"errors"

	"github.com/gen0cide/laforge/core"
	"github.com/urfave/cli"
)

var (
	queryCommand = cli.Command{
		Name:      "query",
		Usage:     "gathers information about elements within the configuration state",
		UsageText: "laforge query QUERY",
		Action:    performquery,
	}
)

func performquery(c *cli.Context) error {
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
