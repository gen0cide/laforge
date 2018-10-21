package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/gen0cide/laforge/core"
	lfcli "github.com/gen0cide/laforge/core/cli"
	"github.com/hashicorp/hcl2/hcl"
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
		},
	}
)

func performplan(c *cli.Context) error {
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

	err = build.Associate(snap)
	if err != nil {
		return err
	}

	state := core.NewState()

	dbfile := filepath.Join(base.CurrentBuild.Dir, "build.db")

	_, err = os.Stat(dbfile)

	if err != nil {
		return err
	}

	state.SetCurrent(snap)

	err = state.Open(dbfile)
	if err != nil {
		return err
	}

	defer state.DB.Close()

	_, err = state.LoadSnapshotFromDB()
	if err != nil {
		return err
	}

	plan, err := core.CalculateDelta(state.Persisted, state.Current)
	if err != nil {
		return err
	}

	tfcmds, err := core.CalculateTerraformNeeds(plan)
	if err != nil {
		return err
	}

	// pp.Println(tfcmds)
	// pp.Println(plan)

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
			case "UPDATE":
				tcol = color.HiYellowString("[%s]", tt)
			case "DESTROY":
				tcol = color.HiRedString("[%s]", tt)
			case "REFRESH":
				tcol = color.HiCyanString("[%s]", tt)
			case "MODIFY":
				tcol = color.HiGreenString("[%s]", tt)
			default:
				tcol = "[UNKNOWN]"
			}
			fmt.Printf("%s  %d) %s\n", tcol, idx, item)
		}
	}

	return nil
}

func performapply(c *cli.Context) error {
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

	// lfcli.SetLogLevel("info")

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
		return err
	}

	state := core.NewState()

	dbfile := filepath.Join(base.CurrentBuild.Dir, "build.db")

	_, err = os.Stat(dbfile)

	if err != nil {
		return err
	}

	state.SetCurrent(snap)

	err = state.Open(dbfile)
	if err != nil {
		return err
	}

	defer state.DB.Close()

	_, err = state.LoadSnapshotFromDB()
	if err != nil {
		return err
	}

	plan, err := core.CalculateDelta(state.Persisted, state.Current)
	if err != nil {
		return err
	}

	state.Plan = plan

	lfcli.SetLogLevel("info")

	err = plan.Preflight()
	if err != nil {
		return err
	}

	return nil
}

// type LFUI struct{}

// func (l *LFUI) Ask(s string) (string, error) {
// 	return "", nil
// }

// func (l *LFUI) AskSecret(s string) (string, error) {
// 	return "", nil
// }

// func (l *LFUI) Output(s string) {
// 	fmt.Println(s)
// }

// func (l *LFUI) Info(s string) {
// 	cliLogger.Info(s)
// }

// func (l *LFUI) Error(s string) {
// 	cliLogger.Error(s)
// }

// func (l *LFUI) Warn(s string) {
// 	cliLogger.Warn(s)
// }

// func performtf(c *cli.Context) error {
// 	lfcli.SetLogLevel("debug")
// 	lfui := LFUI{}
// 	tfcmd := &command.PlanCommand{
// 		Ui: lfui,
// 	}

// 	conf, err := tfcmd.Config(c.Args().Get(0))
// 	if err != nil {
// 		panic(err)
// 	}

// 	ret := tfcmd.Run([]string{})
// 	cliLogger.Infof("output was %d", ret)
// 	return nil
// }

func performtf(c *cli.Context) error {
	return commandNotImplemented(c)
}

func performinfra(c *cli.Context) error {
	return commandNotImplemented(c)
}
