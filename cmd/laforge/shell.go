package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/gen0cide/laforge/core"
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
		cliLogger.Errorf("Error encountered during bootstrap: %v", err)
		os.Exit(1)
	}

	err = base.AssertMinContext(core.TeamContext)
	if err != nil {
		cliLogger.Errorf("Must be in a team context to use a shell connector: %v", err)
		os.Exit(1)
	}

	if listHosts {
		core.SetLogLevel("info")
		cliLogger.Warnf("Environment: %s", color.GreenString("%s", base.CurrentEnv.ID))
		cliLogger.Warnf("Builder: %s", color.GreenString("%s", base.CurrentBuild.ID))
		cliLogger.Warnf("Team Number: %s", color.GreenString("%d", base.CurrentTeam.TeamNumber))
		cliLogger.Infof("Host Information Table")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Hostname", "Public IP"})

		for hid, ph := range base.CurrentTeam.Hosts {
			table.Append([]string{hid, ph.Host.Hostname, ph.RemoteAddr})
		}
		table.Render()

		return nil
	}

	target := c.Args().Get(0)
	if len(target) < 1 {
		cliLogger.Errorf("No host name was specified as an argument!")
		os.Exit(1)
	}

	_, found := base.CurrentEnv.IncludedHosts[target]
	if !found {
		cliLogger.Errorf("Host %s was not found in this environment!", target)
		os.Exit(1)
	}

	provisionedHost, found := base.CurrentTeam.Hosts[target]
	if !found || (provisionedHost != nil && provisionedHost.Active == false) {
		cliLogger.Errorf("Host %s is currently not active in this team's environment", target)
		os.Exit(1)
	}

	return provisionedHost.RemoteShell()

	// if host.IsWindows() {
	// 	s := shells.WinRM{}
	// 	s.SetIO(os.Stdout, os.Stderr, os.Stdin)
	// 	err = s.SetConfig(provisionedHost.WinRMAuthConfig)
	// 	if err != nil {
	// 		cliLogger.Errorf("Error applying configuration: %v", err)
	// 		os.Exit(1)
	// 	}
	// 	err = s.LaunchInteractiveShell()
	// 	if err != nil {
	// 		cliLogger.Errorf("interactive shell error: %v", err)
	// 		os.Exit(1)
	// 	}
	// } else {
	// 	s := shells.SSH{}
	// 	s.SetIO(os.Stdout, os.Stderr, os.Stdin)
	// 	err = s.SetConfig(provisionedHost.SSHAuthConfig)
	// 	if err != nil {
	// 		cliLogger.Errorf("Error applying configuration: %v", err)
	// 		os.Exit(1)
	// 	}
	// 	err = s.LaunchInteractiveShell()
	// 	if err != nil {
	// 		cliLogger.Errorf("interactive shell error: %v", err)
	// 		os.Exit(1)
	// 	}
	// }

	return nil
}
