package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gen0cide/laforge"
	"github.com/urfave/cli"
)

var (
	displayBefore = true
	debugOutput   = false
	cliLogger     = laforge.Logger
	defaultLevel  = "warn"
	verboseOutput = false
	noBanner      = false
)

func init() {
	cli.HelpFlag = cli.BoolFlag{Name: "help, h"}
	cli.VersionFlag = cli.BoolFlag{Name: "version"}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s", laforge.Version)
	}
}

func main() {
	app := cli.NewApp()

	app.Writer = color.Output
	app.ErrWriter = color.Output
	app.Name = "laforge"

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
		cli.BoolFlag{
			Name:        "no-banner, n",
			Usage:       "Disables the ASCII text banner",
			Destination: &noBanner,
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

	app.Before = func(c *cli.Context) error {
		if verboseOutput {
			laforge.SetLogLevel("info")
		}
		if debugOutput {
			laforge.SetLogLevel("debug")
		}
		if !noBanner {
			laforge.PrintLogo()
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		laforge.Logger.Fatalf("Terminated due to error: %v", err)
	}
}

func commandNotImplemented(c *cli.Context) error {
	return fmt.Errorf("%s command not implemented", c.Command.FullName())
}
