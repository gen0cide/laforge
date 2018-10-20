package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"

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

	// snap := core.NewEmptySnapshot()
	// build := base.CurrentEnv.CreateBuild()
	// err = build.CreateTeams()
	// if err != nil {
	// 	panic(err)
	// }

	// envMeta := snap.AddNode(base.CurrentEnv)
	// buildMeta := snap.AddNode(build)

	// for _, n := range :=
	buildnode, ok := snap.Metastore[path.Join(base.CurrentEnv.Path(), base.CurrentEnv.Builder)]
	if !ok {
		return errors.New("builder was not able to be resolved on the graph")
	}
	build, ok := buildnode.Dependency.(*core.Build)
	if !ok {
		return errors.New("build object was not of type *core.Build")
	}

	base.CurrentBuild = build

	err = snap.Sort()
	if err != nil {
		panic(err)
	}

	if args == "" {
		// data, err := json.MarshalIndent(snap.Graph, "", "  ")
		// if err != nil {
		// 	return err
		// }

		// spew.Dump(dag.StronglyConnected(&snap.Graph.Graph))

		// fmt.Printf("%s\n", snap.Graph.String())

		err = build.Associate(snap)
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(snap.Graph, "", "  ")
		if err != nil {
			return err
		}

		// for edgeid, edge := range snap.Graph.Edges() {
		// 	cliLogger.Infof("[%d] %s -> %s", edgeid, edge.Source().(*core.Metadata).GetID(), edge.Target().(*core.Metadata).GetID())
		// }

		data = snap.Graph.Dot(nil)
		// fmt.Printf("%s\n", string(data))
		fmt.Printf("%s\n", string(data))

		return nil
	}

	if args != "" {
		// pat := c.Args().Get(1)
		// if pat == "" {
		// 	return errors.New("must provide a comparison pattern")
		// }
		// g, err := glob.Compile(pat, '/')
		// if err != nil {
		// 	return err
		// }

		// for key, meta := range snap.Objects {
		// 	if !g.Match(key) {
		// 		continue
		// 	}
		// 	cliLogger.Infof("*** INLINE OBJECT *** Parents Of %s (gid=%d) (checksum=%x):", key, meta.GetGID(), meta.Hash())
		// 	for _, x := range meta.Parents() {
		// 		fmt.Printf("  <- (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.Hash()))
		// 	}

		// 	cliLogger.Infof("*** INLINE OBJECT *** Children Of %s (gid=%d) (checksum=%x):", key, meta.GetGID(), meta.Hash())
		// 	for _, x := range meta.Children() {
		// 		fmt.Printf("  -> (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.Hash()))
		// 	}
		// }

		statedata, err := ioutil.ReadFile(args)
		if err != nil {
			return err
		}

		savedSnap := core.NewEmptySnapshot()
		err = json.Unmarshal(statedata, savedSnap)
		if err != nil {
			return err
		}

		err = savedSnap.RebuildGraph()
		if err != nil {
			return err
		}

		// cliLogger.Infof("Current Checksum: %x", snap.Hash())
		// cliLogger.Infof("Persisted Checksum: %x", savedSnap.Hash())

		curRoot, err := snap.Graph.Root()
		if err != nil {
			panic(err)
		}

		savedRoot, err := savedSnap.Graph.Root()
		if err != nil {
			panic(err)
		}

		curWorld, err := snap.Graph.Descendents(curRoot)
		if err != nil {
			panic(err)
		}

		savedWorld, err := savedSnap.Graph.Descendents(savedRoot)
		if err != nil {
			panic(err)
		}

		_ = curWorld
		_ = savedWorld

		ddata := snap.Graph.Dot(&dag.DotOpts{})
		if err != nil {
			panic(err)
		}

		fmt.Printf("\n%s\n", string(ddata))
		os.Exit(0)

		tainted := []dag.Vertex{}
		vertMap := map[string]bool{}

		cb1 := func(v dag.Vertex) tfdiags.Diagnostics {
			meta, ok := v.(*core.Metadata)
			if !ok {
				return tfdiags.Diagnostics{tfdiags.SimpleWarning(fmt.Sprintf("vertex %v was not of type *core.Metadata", v))}
			}
			newmeta, ok := snap.Metastore[meta.ID]
			if !ok {
				// this is getting deleted
				tainted = append(tainted, v)
				return nil
			}
			if newmeta.Checksum != meta.Checksum {
				// cliLogger.Warnf(fmt.Sprintf("** UPDATE ** %s", meta.ID))
				// return fmt.Errorf("** UPDATE ** %s", meta.ID)
				tainted = append(tainted, v)
				vertMap[newmeta.ID] = true
			}

			return nil
		}

		cb2 := func(v dag.Vertex) tfdiags.Diagnostics {
			meta, ok := v.(*core.Metadata)
			if !ok {
				return tfdiags.Diagnostics{tfdiags.SimpleWarning(fmt.Sprintf("vertex %v was not of type *core.Metadata", v))}
			}
			newmeta, ok := savedSnap.Metastore[meta.ID]
			if !ok {
				cliLogger.Warnf("** ADDITION ** %s", meta.ID)
				// return fmt.Errorf("** ADDITION ** %s", meta.ID)
				return nil
			}
			if !savedSnap.Graph.HasVertex(newmeta) {
				cliLogger.Warnf("** UPDATE ** %s", newmeta.ID)
			}
			return nil
		}

		savedSnap.Graph.SetDebugWriter(ioutil.Discard)
		log.SetOutput(ioutil.Discard)
		savedSnap.Graph.Walk(cb1)
		edgeRemovals := []dag.Edge{}
		vertRemovals := []dag.Vertex{}
		newtaints := []dag.Vertex{}
		taintSet := &dag.Set{}
		for _, t := range tainted {
			dset, err := savedSnap.Graph.Descendents(t)
			if err != nil {
				panic(err)
			}
			lits := dset.Filter(func(i interface{}) bool {
				imet := i.(*core.Metadata)
				if imet.IsGlobalType() && !vertMap[imet.ID] {
					return false
				}
				return true
			})

			for _, x := range dag.AsVertexList(lits) {
				newtaints = append(newtaints, x)
				taintSet.Add(x)
			}
		}

		spew.Dump(newtaints)

		savedSnap.Graph.DepthFirstWalk(newtaints, func(v dag.Vertex, depth int) error {
			if v == savedRoot {
				return nil
			}
			met := v.(*core.Metadata)
			edgesFucked := 0
			isPV := false
			if met.TypeByPath() == core.LFTypeProvisionedHost {
				isPV = true
			}
			_ = isPV
			for _, e := range savedSnap.Graph.EdgesFrom(v) {
				tv := e.Target()
				tvmet := tv.(*core.Metadata)

				if !taintSet.Include(tv) && tvmet.IsGlobalType() && !taintSet.Include(met) && met.IsGlobalType() {
					continue
				}
				edgeRemovals = append(edgeRemovals, e)
				edgesFucked++
			}
			if edgesFucked > 0 {
				vertRemovals = append(vertRemovals, v)
			}
			return nil
		})
		for _, e := range edgeRemovals {
			savedSnap.Graph.RemoveEdge(e)
		}
		for _, v := range vertRemovals {
			savedSnap.Graph.Remove(v)
		}

		// for _, x := range snap.Graph.Edges() {
		// 	cliLogger.Infof("Edge: %v", x.Hashcode())
		// 	cliLogger.Infof("  (source) %s", x.Source().Name)
		// 	cliLogger.Infof("  (target) %s", x.Target().Name)
		// }

		snap.Graph.Walk(cb2)

		// pat := c.Args().Get(1)
		// if pat != "" {
		// 	cliLogger.Infof("%s -> Current Checksum: %x", pat, snap.Metastore[pat].Hash())
		// 	cliLogger.Infof("%s -> Persisted Checksum: %x", pat, savedSnap.Metastore[pat].Hash())
		// }
		// mapper.Update(savedSnap.Graph)
		// mapper.Update(snap.Graph)
		// for key, meta := range savedSnap.Objects {
		// 	if !g.Match(key) {
		// 		continue
		// 	}
		// 	cliLogger.Infof("~~~ PERSISTED OBJECT ~~~ Parents Of %s (gid=%d) (checksum=%x):", key, meta.GetGID(), meta.Hash())
		// 	for _, x := range meta.Parents() {
		// 		fmt.Printf("  <- (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.Hash()))
		// 	}
		// 	cliLogger.Infof("~~~ PERSISTED OBJECT ~~~ Children Of %s (gid=%d) (checksum=%x):", key, meta.GetGID(), meta.Hash())
		// 	for _, x := range meta.Children() {
		// 		fmt.Printf("  -> (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.Hash()))
		// 	}
		// }

	}

	return nil
}
