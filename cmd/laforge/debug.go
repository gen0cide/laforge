package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gen0cide/laforge/core"
	lfcli "github.com/gen0cide/laforge/core/cli"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/urfave/cli"
)

var (
	debugCommand = cli.Command{
		Name:      "debug",
		Usage:     "Allows for expirimental debugging features that only laforge developers would find useful",
		UsageText: "laforge debug ...",
		Action:    performdebug,
	}
)

func performdebug(c *cli.Context) error {
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

	lfcli.SetLogLevel("info")

	snap, err := core.NewSnapshotFromEnv(base.CurrentEnv)
	if err != nil {
		return err
	}

	buildnode, ok := snap.Metastore[path.Join(base.CurrentEnv.Path(), base.CurrentEnv.Builder)]
	if !ok {
		return errors.New("builder was not able to be resolved on the graph")
	}
	build, ok := buildnode.Dependency.(*core.Build)
	if !ok {
		return errors.New("build object was not of type *core.Build")
	}

	base.CurrentBuild = build

	state := core.NewState()
	state.Base = base

	lfcli.SetLogLevel("info")

	err = build.Associate(snap)
	if err != nil {
		panic(err)
	}

	dbfile := filepath.Join(base.CurrentBuild.Dir, "build.db")

	err = state.Open(dbfile)
	if err != nil {
		return err
	}

	defer state.DB.Close()

	state.SetCurrent(snap)

	_, err = state.LoadSnapshotFromDB()
	if err != nil {
		return err
	}

	plan, err := state.CalculateDelta()
	if err != nil {
		return err
	}

	for _, k := range plan.OrderedPriorities {
		cliLogger.Infof("Step #%d:", k)
		for idx, item := range plan.TasksByPriority[k] {
			fmt.Printf("  %d) %s\n", idx, item)
		}
	}

	return nil
}
