package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

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

	args := c.Args().Get(0)
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

	err = build.Associate(snap)
	if err != nil {
		panic(err)
	}

	if args == "" {
		data, err := json.MarshalIndent(snap, "", "  ")
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", string(data))

		return nil
	}

	if args != "" {
		statedata, err := ioutil.ReadFile(args)
		if err != nil {
			return err
		}

		var savedSnap core.Snapshot
		err = json.Unmarshal(statedata, &savedSnap)
		if err != nil {
			return err
		}

		err = savedSnap.RebuildGraph()
		if err != nil {
			return err
		}

		plan, err := core.CalculateDelta(&savedSnap, snap)
		if err != nil {
			return err
		}

		for _, k := range plan.OrderedPriorities {
			cliLogger.Infof("Step #%d:", k)
			for idx, item := range plan.TasksByPriority[k] {
				fmt.Printf("  %d) %s\n", idx, item)
			}
		}
	}
	return nil
}
