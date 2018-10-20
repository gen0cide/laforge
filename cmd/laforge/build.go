package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/core"
	lfcli "github.com/gen0cide/laforge/core/cli"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/urfave/cli"
)

var (
	updateConfig = false
	buildCommand = cli.Command{
		Name:      "build",
		Usage:     "builds environment specific infrastructure configurations",
		UsageText: "laforge build [OPTIONS]",
		Action:    performbuild,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "force, f",
				Usage:       "force removes and deletes any conflicting directories (dangerous)",
				Destination: &overwrite,
			},
			cli.BoolFlag{
				Name:        "update, u",
				Usage:       "Updates a build directory (expirimental)",
				Destination: &updateConfig,
			},
		},
	}
)

func performbuild(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		if _, ok := err.(hcl.Diagnostics); ok {
			return errors.New("aborted due to parsing error")
		}
		cliLogger.Errorf("Error encountered during bootstrap: %v", err)
		os.Exit(1)
	}

	bldr, err := builder.New(base, overwrite, updateConfig)
	if err != nil {
		cliLogger.Errorf("Error encountered initializing builder:\n%v", err)
		os.Exit(1)
	}

	cliLogger.Infof("Build directory initialized")

	err = bldr.Do()

	if err != nil {
		cliLogger.Errorf("Error encountered initializing builder:\n%v", err)
		os.Exit(1)
	}

	dbfile := filepath.Join(base.CurrentBuild.Dir, "build.db")
	_, err = os.Stat(dbfile)
	if err == nil || !os.IsNotExist(err) {
		return err
	}

	os.Chdir(base.CurrentBuild.Dir)
	base, err = core.Bootstrap()
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

	err = build.Associate(snap)
	if err != nil {
		panic(err)
	}

	state := core.NewState()
	err = state.Open(dbfile)
	if err != nil {
		return err
	}

	defer state.DB.Close()

	err = state.PersistSnapshot(snap)
	if err != nil {
		return err
	}

	cliLogger.Infof("State DB has been persisted to %s", dbfile)

	return nil
}
