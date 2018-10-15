package main

import (
	"errors"
	"os"
	"sync"

	"github.com/gen0cide/laforge/core"
	"github.com/gen0cide/laforge/spanner"
	"github.com/urfave/cli"
)

var (
	silent         = false
	writeLog       = false
	remoteHost     = ""
	spannerCommand = cli.Command{
		Name:      "spanner",
		Usage:     "Runs parallel functions across all teams for an environment.",
		UsageText: "laforge spanner",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "log, l",
				Usage:       "Log the outputs of each stream to a file on disk.",
				Destination: &overwrite,
			},
			cli.BoolFlag{
				Name:        "silent, s",
				Usage:       "discards STDOUT, only showing execution status and error messages if they exist.",
				Destination: &silent,
			},
			cli.StringFlag{
				Name:        "target, t",
				Usage:       "Specifies a specific host for the command to be run. (remote-exec only).",
				Destination: &remoteHost,
			},
		},
		Subcommands: []cli.Command{
			{
				Name:            "local-exec",
				Usage:           "Execute a local command inside team context for each team.",
				Action:          spannerLocalExec,
				SkipFlagParsing: true,
			},
			{
				Name:            "remote-exec",
				Usage:           "Run a command on a remote machine in each team's environment.",
				Action:          spannerRemoteExec,
				SkipFlagParsing: true,
			},
			{
				Name:            "terraform",
				Usage:           "Run terraform commands across the environment.",
				Action:          spannerTerraformExec,
				SkipFlagParsing: true,
			},
		},
	}
)

func spannerTerraformExec(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		cliLogger.Errorf("Error encountered during bootstrap: %v", err)
		os.Exit(1)
	}
	args := []string{}
	for _, a := range []string(c.Args()) {
		args = append(args, a)
	}

	err = base.AssertMinContext(core.BuildContext)
	if err != nil {
		cliLogger.Errorf("Must be in a build context to use terraform spanner: %v", err)
		os.Exit(1)
	}
	wg := new(sync.WaitGroup)
	finChan := make(chan bool, 1)

	for _, t := range base.CurrentBuild.Teams {
		wg.Add(1)
		go t.RunTerraformCommand(args, wg)
	}

	go func() {
		wg.Wait()
		close(finChan)
	}()

	select {
	case <-finChan:
		return nil
	}
}

func spannerLocalExec(c *cli.Context) error {
	s, err := spanner.New(nil, []string(c.Args()), "local-exec", "", false, false)
	if err != nil {
		return err
	}

	err = s.CreateWorkerPool()
	if err != nil {
		return err
	}

	err = s.Do()
	if err != nil {
		return err
	}

	return nil
}

func spannerRemoteExec(c *cli.Context) error {
	if len(remoteHost) == 0 {
		return errors.New("must provide a target host ID using the -t flag before remote-exec")
	}
	s, err := spanner.New(nil, []string(c.Args()), "remote-exec", remoteHost, false, false)
	if err != nil {
		return err
	}

	err = s.CreateWorkerPool()
	if err != nil {
		return err
	}

	err = s.Do()
	if err != nil {
		return err
	}

	return nil
}
