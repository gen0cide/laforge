package main

import (
	"os"

	"github.com/fatih/color"

	"github.com/gen0cide/laforge/core"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var (
	shellCommand = cli.Command{
		Name:      "shell",
		Usage:     "launches an interactive SSH or Powershell console on a provisioned host",
		UsageText: "laforge shell",
		Action:    performshell,
	}
)

func performshell(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		cliLogger.Errorf("Error encountered during bootstrap: %v", err)
		os.Exit(1)
	}

	err = base.AssertMinContext(core.TeamContext)
	if err != nil {
		cliLogger.Errorf("Must be in a team context to use a shell connector: %v", err)
		os.Exit(1)
	}

	rs, err := core.GetState(base)
	if err != nil {
		panic(err)
	}

	core.SetLogLevel("info")
	cliLogger.Warnf("Environment: %s", color.GreenString("%s", base.Environment.ID))
	cliLogger.Warnf("Builder: %s", color.GreenString("%s", base.Environment.Builder))
	cliLogger.Warnf("Team Number: %s", color.GreenString("%d", base.Team.TeamNumber))
	cliLogger.Infof("Host Information Table")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Host", "Public IP", "Private IP"})

	for _, v := range rs.Hosts {
		table.Append(v.TableInfo())
	}
	table.Render() // Send output

	return nil
}
