package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"

	"github.com/gobwas/glob"

	"github.com/gen0cide/laforge/core"
	lfcli "github.com/gen0cide/laforge/core/cli"
	"github.com/hashicorp/hcl2/hcl"
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

	buildnode, ok := snap.Metastore[path.Join(base.CurrentEnv.Path(), base.CurrentEnv.Builder)]
	if !ok {
		return errors.New("builder was not able to be resolved on the graph")
	}
	build, ok := buildnode.Dependency.(*core.Build)
	if !ok {
		return errors.New("build object was not of type *core.Build")
	}

	base.CurrentBuild = build

	lfcli.SetLogLevel("info")

	pat := c.Args().First()
	if pat != "" {
		g, err := glob.Compile(pat, '/')
		if err != nil {
			return err
		}
		for key, meta := range snap.Metastore {
			if !g.Match(key) {
				continue
			}
			cliLogger.Infof("Parents Of %s (gid=%d) (checksum=%x):", key, meta.GetGID(), meta.Hash())
			for _, x := range meta.Parents() {
				fmt.Printf("  <- (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.Hash()))
			}
			cliLogger.Infof("Children Of %s (gid=%d) (checksum=%x):", key, meta.GetGID(), meta.Hash())
			for _, x := range meta.Children() {
				fmt.Printf("  -> (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.Hash()))
			}
		}
		return nil
	}

	for key, meta := range snap.Metastore {
		cliLogger.Infof("Parents Of %s (gid=%d) (checksum=%x):", key, meta.GetGID(), meta.Hash())
		for _, x := range meta.Parents() {
			fmt.Printf("  <- (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.Hash()))
		}
		cliLogger.Infof("Children Of %s (gid=%d) (%x):", key, meta.GetGID(), meta.Hash())
		for _, x := range meta.Children() {
			fmt.Printf("  -> (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.Hash()))
		}
	}

	return nil
}
