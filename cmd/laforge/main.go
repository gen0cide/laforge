package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gen0cide/laforge"
	"github.com/gen0cide/laforge/core"
	"github.com/urfave/cli"
)

var (
	displayBefore = true
	debugOutput   = false
	cliLogger     = core.Logger
	defaultLevel  = "warn"
	verboseOutput = false
	noBanner      = false
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
	app.Name = "laforge"
	app.Usage = "Distributed competition development and automation"
	app.Description = "Security competition infrastructure automation framework"

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
		configureCommand,
		initCommand,
		statusCommand,
		dumpCommand,
		buildCommand,
		envCommand,
		queryCommand,
		serveCommand,
		shellCommand,
		debugCommand,
		downloadCommand,
		exampleCommand,
		depsCommand,
		spannerCommand,
		infraCommand,
		fmtCommand,
		graphCommand,
	}

	app.Before = func(c *cli.Context) error {
		if verboseOutput {
			core.SetLogLevel("info")
		}
		if debugOutput {
			core.SetLogLevel("debug")
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		core.Logger.Fatalf("Terminated due to error: %v", err)
	}
}

func commandNotImplemented(c *cli.Context) error {
	return fmt.Errorf("%s command not implemented", c.Command.FullName())
}
