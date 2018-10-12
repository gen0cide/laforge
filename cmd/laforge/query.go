package main

import (
	"fmt"
	"os"

	"github.com/gen0cide/laforge/core"
	"github.com/urfave/cli"
)

var (
	queryCommand = cli.Command{
		Name:      "query",
		Usage:     "gathers information about elements within the configuration state",
		UsageText: "laforge query QUERY",
		Action:    performquery,
	}
)

func performquery(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		cliLogger.Errorf("Error encountered during bootstrap: %v", err)
		os.Exit(1)
	}

	err = base.AssertMinContext(core.EnvContext)
	if err != nil {
		cliLogger.Errorf("Must be in a team context to use a shell connector: %v", err)
		os.Exit(1)
	}

	mappings := map[string][]string{}

	for _, host := range base.Environment.IncludedHosts {
		for _, x := range host.Provisioners {
			if x.Kind() != "script" {
				continue
			}
			script := x.(*core.Script)
			_, ok := mappings[script.ID]
			if !ok {
				mappings[script.ID] = []string{}
			}

			mappings[script.ID] = append(mappings[script.ID], host.Hostname)
		}
	}

	fmt.Println("script_id,hostname")
	for scriptID, hosts := range mappings {
		fmt.Printf("%s,\n", scriptID)
		for _, x := range hosts {
			fmt.Printf(",%s\n", x)
		}

	}

	// rs, err := core.GetState(base)
	// if err != nil {
	// 	panic(err)
	// }

	// core.SetLogLevel("info")
	// cliLogger.Warnf("Environment: %s", color.GreenString("%s", base.Environment.ID))
	// cliLogger.Warnf("Builder: %s", color.GreenString("%s", base.Environment.Builder))
	// cliLogger.Warnf("Team Number: %s", color.GreenString("%d", base.Team.TeamNumber))
	// cliLogger.Infof("Host Information Table")

	// table := tablewriter.NewWriter(os.Stdout)
	// table.SetHeader([]string{"Host", "Public IP", "Private IP"})

	// for _, v := range rs.Hosts {
	// 	table.Append(v.TableInfo())
	// }
	// table.Render() // Send output
	return nil
}
