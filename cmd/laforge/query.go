package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/yourbasic/graph"

	"github.com/fatih/color"

	"github.com/gobwas/glob"

	"github.com/gen0cide/laforge/core"
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

	build, ok := snap.Objects[path.Join(base.CurrentEnv.Path(), base.CurrentEnv.Builder)].(*core.Build)
	if !ok {
		return errors.New("builder was not able to resolve object of type Build")
	}
	base.CurrentBuild = build

	err = snap.Sort()
	if err != nil {
		panic(err)
	}

	core.SetLogLevel("info")

	pat := c.Args().First()
	if pat != "" {
		g, err := glob.Compile(pat, '/')
		if err != nil {
			return err
		}
		mut := graph.Copy(graph.Transpose(snap.Graph))
		for key, meta := range snap.Metadata {
			if !g.Match(key) {
				continue
			}
			mid := snap.ObjectToGID[meta.Dependency]
			cliLogger.Infof("Parents Of %s (gid=%d) (%x):", key, mid, meta.Checksum)
			for _, x := range meta.ParentIDs {
				pm := snap.Metadata[x]
				pid := snap.ObjectToGID[pm.Dependency]
				dacost := mut.Cost(mid, pid)
				fmt.Printf("  <- (pid=%d / cost=%d) %s (checksum=%s)\n", snap.ObjectToGID[pm.Dependency], dacost, color.YellowString(x), color.CyanString("%x", pm.Checksum))
			}
			cliLogger.Infof("Children Of %s (%x):", key, meta.Checksum)
			for _, x := range meta.ChildrenIDs {
				pm := snap.Metadata[x]
				fmt.Printf("  -> (%d) %s (checksum=%s)\n", snap.ObjectToGID[pm.Dependency], color.YellowString(x), color.CyanString("%x", pm.Checksum))
			}
		}
		return nil
	}

	for key, meta := range snap.Metadata {
		cliLogger.Infof("Parents Of %s (%x):", key, meta.Checksum)
		for _, x := range meta.ParentIDs {
			pm := snap.Metadata[x]
			fmt.Printf("  <- (%d) %s (checksum=%s)\n", snap.ObjectToGID[pm.Dependency], color.YellowString(x), color.CyanString("%x", pm.Checksum))
		}
		cliLogger.Infof("Children Of %s (%x):", key, meta.Checksum)
		for _, x := range meta.ChildrenIDs {
			pm := snap.Metadata[x]
			fmt.Printf("  -> (%d) %s (checksum=%s)\n", snap.ObjectToGID[pm.Dependency], color.YellowString(x), color.CyanString("%x", pm.Checksum))
		}
	}

	// for i := 0; i < oplen; i++ {
	// 	for depidx, dep := range ops[i] {
	// 		fmt.Printf("[%d] - (%d) %s\n", i, depidx, dep.Path())
	// 	}
	// }

	// fmt.Println(snap.Graph.String())

	// pp.Println(graph.Sort(snap.Graph))

	return nil
}
