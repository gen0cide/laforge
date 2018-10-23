package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/emicklei/dot"
	"github.com/fatih/color"

	"github.com/gen0cide/laforge/core"
	lfcli "github.com/gen0cide/laforge/core/cli"
	"github.com/urfave/cli"
)

var (
	shouldgraph  = false
	infraCommand = cli.Command{
		Name:      "infra",
		Usage:     "Manage infrastructure deployment that has been generated with Laforge.",
		UsageText: "laforge infra",
		Subcommands: []cli.Command{
			{
				Name:   "plan",
				Usage:  "Show the delta between current deployment and the final desired state.",
				Action: performplan,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:        "graph",
						Usage:       "graphs the output in dot format",
						Destination: &shouldgraph,
					},
				},
			},
			{
				Name:   "tf",
				Usage:  "Test terraform embedding",
				Action: performtf,
				Flags:  []cli.Flag{},
			},
			{
				Name:            "status",
				Usage:           "Show the current build's infrastructure status.",
				Action:          performinfra,
				SkipFlagParsing: true,
			},
			{
				Name:            "apply",
				Usage:           "Provision the infrastructure to bring state in line with build blueprint.",
				Action:          performapply,
				SkipFlagParsing: true,
			},
			{
				Name:            "taint",
				Usage:           "Mark a host for re-provisioning in the laforge infrastructure state.",
				Action:          performinfra,
				SkipFlagParsing: true,
			},
			{
				Name:            "destroy",
				Usage:           "Destroy the builds infrastructure and clean the state.",
				Action:          performinfra,
				SkipFlagParsing: true,
			},
			{
				Name:            "run",
				Usage:           "Run a host provisioner on a specific host within the infrastructure (usually for debugging).",
				Action:          performinfra,
				SkipFlagParsing: true,
			},
			{
				Name:            "graph",
				Usage:           "Generate a proposed DOT diagram of the target state.",
				Action:          performinfragraph,
				SkipFlagParsing: true,
			},
		},
	}
)

func performplan(c *cli.Context) error {
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

	lfcli.SetLogLevel("info")

	for tid, cmds := range tfcmds {
		cliLogger.Infof("Terraform Commands For Team: %s", tid)
		for _, c := range cmds {
			fmt.Printf("  $ terraform %s\n", c)
		}
	}

	for _, k := range plan.OrderedPriorities {
		cliLogger.Infof("Step #%d:", k)
		for idx, item := range plan.TasksByPriority[k] {
			tcol := ""
			tt := plan.TaskTypes[item]
			switch tt {
			case "MODIFY":
				tcol = color.HiYellowString("[%s]", tt)
			case "DELETE":
				tcol = color.HiRedString("[%s]", tt)
			case "TOUCH":
				tcol = color.HiCyanString("[%s]", tt)
			case "CREATE":
				tcol = color.HiGreenString("[%s]", tt)
			default:
				tcol = "[UNKNOWN]"
			}
			fmt.Printf("%s  %d) %s\n", tcol, idx, item)
		}
	}

	return nil
}

func performinfragraph(c *cli.Context) error {
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

	snap.AltGraph.Remove("root")
	snap.AltGraph.TransitiveReduction()
	nodemap := map[string]dot.Node{}
	g := snap.AltGraph

	di := dot.NewGraph(dot.Directed)
	di.Attr("nodesep", "0.2")
	di.Attr("compound", "true")
	di.Attr("rank", "min")
	di.Attr("rankdir", "LR")
	di.Attr("dpi", "144")
	di.Attr("smoothType", "graph_dist")
	di.Attr("mode", "hier")
	di.Attr("splines", "spline")
	di.Attr("decoreate", "true")
	di.Attr("overlap", "false")
	di.Attr("model", "subset")
	di.Attr("K", "0.6")
	di.Attr("fontname", "Helvetica")

	for _, x := range g.Vertices() {
		id := x.(string)
		nodemap[id] = di.Node(id)
		meta, ok := snap.Metastore[id]
		if !ok {
			panic(fmt.Errorf("could not find dependency for %s", id))
		}
		nodemap[id].Attr("style", meta.Style())
		nodemap[id].Attr("shape", meta.Shape())
		nodemap[id].Attr("height", "0.1")
		nodemap[id].Attr("label", []byte(meta.Label()))
		nodemap[id].Attr("fillcolor", meta.FillColor())
		nodemap[id].Attr("fontname", "Helvetica")
	}

	for _, x := range g.Edges() {
		src := x.Source().(string)
		tar := x.Target().(string)
		nodemap[src].Edge(nodemap[tar])
	}

	graphstring := di.String()
	wat := strings.Replace(graphstring, `"<`, `<<`, -1)
	wat = strings.Replace(wat, `\"`, `"`, -1)
	wat = strings.Replace(wat, `>"`, `>>`, -1)
	fmt.Println(wat)
	return nil
}

func performapply(c *cli.Context) error {
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

	// err = plan.Preflight()
	// if err != nil {
	// 	return err
	// }

	plan.Base = state.Base

	err = plan.SetupTasks()
	if err != nil {
		return err
	}

	diags := plan.Execute()
	if diags.HasErrors() {
		return diags.Err()
	}

	return nil
}

func performtf(c *cli.Context) error {
	return commandNotImplemented(c)
}

func performinfra(c *cli.Context) error {
	return commandNotImplemented(c)
}
