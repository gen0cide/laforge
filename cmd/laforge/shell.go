package main

import (
	"errors"
	"os"

	"github.com/fatih/color"
	"github.com/gen0cide/laforge/core"
	lfcli "github.com/gen0cide/laforge/core/cli"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var (
	listHosts    = false
	shellCommand = cli.Command{
		Name:      "shell",
		Usage:     "launches an interactive SSH or Powershell console on a provisioned host",
		UsageText: "laforge shell",
		Action:    performshell,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "list, l",
				Usage:       "Shows a table of hosts available to be connected to",
				Destination: &listHosts,
			},
		},
	}
)

func performshell(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		if _, ok := err.(hcl.Diagnostics); ok {
			return errors.New("aborted due to parsing error")
		}
		cliLogger.Errorf("Error encountered during bootstrap: %v", err)
		os.Exit(1)
	}

	err = base.AssertMinContext(core.TeamContext)
	if err != nil {
		cliLogger.Errorf("Must be in a team context to use a shell connector: %v", err)
		os.Exit(1)
	}

	if listHosts {
		lfcli.SetLogLevel("info")
		cliLogger.Warnf("Environment: %s", color.GreenString("%s", base.CurrentEnv.ID))
		cliLogger.Warnf("Builder: %s", color.GreenString("%s", base.CurrentBuild.ID))
		cliLogger.Warnf("Team Number: %s", color.GreenString("%d", base.CurrentTeam.TeamNumber))
		cliLogger.Infof("Host Information Table")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Network", "Hostname", "Internal IP", "Public IP"})

		for _, net := range base.CurrentTeam.ProvisionedNetworks {
			for hid, host := range net.ProvisionedHosts {
				if host.Conn == nil {
					continue
				}
				if host.Conn.RemoteAddr == core.NullIP {
					continue
				}
				table.Append([]string{hid, net.Base(), host.Host.Hostname, host.Conn.LocalAddr, host.Conn.RemoteAddr})
			}
		}
		table.Render()

		return nil
	}

	target := c.Args().Get(0)
	if len(target) < 1 {
		cliLogger.Errorf("No host name was specified as an argument!")
		os.Exit(1)
	}

	provisionedHost, found := base.ProvisionedHosts[target]
	if !found || (provisionedHost != nil && provisionedHost.Conn != nil && provisionedHost.Conn.Active == false) {
		cliLogger.Errorf("Host %s is currently not active in this team's environment", target)
		os.Exit(1)
	}

	return provisionedHost.Conn.RemoteShell()
}
