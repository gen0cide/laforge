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

	buildnode, ok := snap.Metastore[path.Join(base.CurrentEnv.Path(), base.CurrentEnv.Builder)]
	if !ok {
		return errors.New("builder was not able to be resolved on the graph")
	}
	// buildmeta, ok := buildnode.(*core.Metadata)
	// if !ok {
	// 	return errors.New("buildnode was not of type *core.Metadata")
	// }
	build, ok := buildnode.Dependency.(*core.Build)
	if !ok {
		return errors.New("build object was not of type *core.Build")
	}

	base.CurrentBuild = build

	// err = snap.Sort()
	// if err != nil {
	// 	panic(err)
	// }

	// arg := c.Args().First()

	// target, ok := snap.Objects[arg]
	// if !ok {
	// 	return errors.New("object was not located on the current graph")
	// }

	// err = graph.WalkRelationship(target, graph.InfiniteDepth, 1, graph.TraverseChildren, func(rel graph.Relationship, distance int) error {
	// 	buf := new(bytes.Buffer)
	// 	for i := 0; i < distance; i++ {
	// 		buf.WriteString("  ")
	// 	}
	// 	buf.WriteString("- ")
	// 	buf.WriteString(rel.GetID())
	// 	fmt.Printf("%s\n", buf.String())
	// 	return nil
	// })

	// snap.Plot(os.Stdout, arg)

	return nil
}
