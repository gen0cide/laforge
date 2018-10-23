package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gen0cide/laforge"
	"github.com/gen0cide/laforge/agent"
	lfcli "github.com/gen0cide/laforge/core/cli"
	"github.com/kardianos/service"
	"github.com/urfave/cli"
)

var (
	displayBefore = true
	debugOutput   = false
	cliLogger     = lfcli.Logger
	defaultLevel  = "warn"
	verboseOutput = false
	noBanner      = false
	serviceObj    service.Service
)

func init() {
	cli.HelpFlag = cli.BoolFlag{Name: "help, h"}
	cli.VersionFlag = cli.BoolFlag{Name: "version"}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s\n", laforge.Version)
	}
}

func main() {
	app := cli.NewApp()

	app.Writer = color.Output
	app.ErrWriter = color.Output

	cli.AppHelpTemplate = fmt.Sprintf("%s\n%s", strings.Join(laforge.ColorLogo, "\n"), cli.AppHelpTemplate)
	app.Name = "laforge-agent"
	app.Usage = "Endpoint agent for host provisioning"
	app.Description = "Endpoint configuration agent for the Laforge Framework."

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "verbose, v",
			Usage:       "Enables verbose command output",
			Destination: &verboseOutput,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enables low level debug output",
			Destination: &debugOutput,
		},
		cli.StringFlag{
			Name:        "work-dir, w",
			Usage:       "Overrides the home directory for the agent.",
			Value:       agent.AgentHomeDir,
			EnvVar:      "LAFORGE_AGENT_HOME_DIR",
			Destination: &agent.AgentHomeDir,
		},
		cli.StringFlag{
			Name:        "exe-path, e",
			Usage:       "Overrides the location of the agent binary.",
			Value:       agent.ExePath,
			EnvVar:      "LAFORGE_AGENT_EXE_PATH",
			Destination: &agent.ExePath,
		},
	}
	app.Version = laforge.Version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  laforge.AuthorName,
			Email: laforge.AuthorEmail,
		},
	}
	app.Copyright = `(c) 2018 Alex Levinson`
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "start",
			Usage:     "start the laforge-agent system service",
			UsageText: "laforge-agent start",
			Action:    startagent,
		},
		cli.Command{
			Name:      "stop",
			Usage:     "stops the laforge-agent system service",
			UsageText: "laforge-agent stop",
			Action:    stopagent,
		},
		cli.Command{
			Name:      "restart",
			Usage:     "retarts the laforge-agent system service",
			UsageText: "laforge-agent restart",
			Action:    restartagent,
		},
		cli.Command{
			Name:      "install",
			Usage:     "installs the laforge-agent binary as a system service",
			UsageText: "laforge-agent install",
			Action:    installagent,
		},
		cli.Command{
			Name:      "uninstall",
			Usage:     "removes the laforge-agent system service",
			UsageText: "laforge-agent uninstall",
			Action:    uninstallagent,
		},
		cli.Command{
			Name:      "status",
			Usage:     "checks the status of the laforge-agent system service",
			UsageText: "laforge-agent status",
			Action:    agentstatus,
		},
		cli.Command{
			Name:      "run",
			Usage:     "checks the status of the laforge-agent system service",
			UsageText: "laforge-agent run",
			Action:    runagent,
		},
		cli.Command{
			Name:      "serve",
			Usage:     "runs a laforge-agent in the foreground",
			UsageText: "laforge-agent serve",
			Action:    serveagent,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "port, p",
					Usage:       "Overrides the default API port.",
					Value:       agent.ServerPort,
					EnvVar:      "LAFORGE_AGENT_PORT",
					Destination: &agent.ServerPort,
				},
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		if verboseOutput {
			lfcli.SetLogLevel("info")
		}
		if debugOutput {
			lfcli.SetLogLevel("debug")
		}
		svc, err := agent.GetService()
		if err != nil {
			return err
		}
		serviceObj = svc
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		lfcli.Logger.Fatalf("Terminated due to error: %v", err)
	}
}

func startagent(c *cli.Context) error {
	lfcli.Logger.Warnf("Starting Laforge Agent...")
	err := serviceObj.Start()
	if err != nil {
		return err
	}
	lfcli.Logger.Warnf("Laforge Agent Started.")
	return nil
}

func stopagent(c *cli.Context) error {
	lfcli.Logger.Warnf("Stopping Laforge Agent...")
	err := serviceObj.Stop()
	if err != nil {
		return err
	}
	lfcli.Logger.Warnf("Laforge Agent Stopped.")
	return nil
}

func restartagent(c *cli.Context) error {
	lfcli.Logger.Warnf("Restarting Laforge Agent...")
	err := serviceObj.Restart()
	if err != nil {
		return err
	}
	lfcli.Logger.Warnf("Laforge Agent Restarted.")
	return nil
}

func installagent(c *cli.Context) error {
	lfcli.Logger.Warnf("Installing Laforge Agent...")
	err := serviceObj.Install()
	if err != nil {
		return err
	}
	lfcli.Logger.Warnf("Laforge Agent Service Installed.")
	return nil
}

func uninstallagent(c *cli.Context) error {
	lfcli.Logger.Warnf("Uninstalling Laforge Agent...")
	err := serviceObj.Uninstall()
	if err != nil {
		return err
	}
	lfcli.Logger.Warnf("Laforge Agent Service Uninstalled.")
	return nil
}

func agentstatus(c *cli.Context) error {
	lfcli.Logger.Warnf("Uninstalling Laforge Agent...")
	stat, err := serviceObj.Status()
	if err != nil {
		return err
	}
	switch stat {
	case service.StatusUnknown:
		lfcli.Logger.Warnf("Status: UNKNOWN")
	case service.StatusRunning:
		lfcli.Logger.Warnf("Status: RUNNING")
	case service.StatusStopped:
		lfcli.Logger.Warnf("Status: STOPPED")
	}
	return nil
}

func runagent(c *cli.Context) error {
	lfcli.Logger.Warnf("Running Laforge Agent (Service)...")
	err := serviceObj.Run()
	if err != nil {
		return err
	}
	return nil
}

func serveagent(c *cli.Context) error {
	lfcli.Logger.Warnf("Serving In Foreground Laforge Agent...")
	agent.Agent.Serve()
	return nil
}
