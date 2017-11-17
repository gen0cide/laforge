package command

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alecthomas/chroma/quick"
	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
	input "github.com/tcnksm/go-input"
)

func CmdNetwork(c *cli.Context) {
	cli.ShowAppHelpAndExit(c, 0)
}

func CmdNetworkLs(c *cli.Context) {
	competition.ValidateEnv()
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME Config: " + err.Error())
	}
	env := comp.CurrentEnv()
	if env == nil {
		competition.LogFatal("Cannot load environment! (Check ~/.lf_env)")
	}
	networks := env.ParseNetworks()
	for name, network := range networks {
		competition.Log("----------------------------------------------------------------------")
		competition.Log("NETWORK >>> " + name)
		competition.Log(filepath.Join(env.NetworksDir(), fmt.Sprintf("%s.yml", name)))
		data, err := json.MarshalIndent(network, "", "  ")
		if err != nil {
			competition.LogError(fmt.Sprintf("Error printing network information: network=%s.yml error=%s", name, err.Error()))
		}
		err = quick.Highlight(os.Stdout, string(data), "json", "terminal", "vim")
		if err != nil {
			competition.LogError(fmt.Sprintf("Error printing network information: network=%s.yml error=%s", name, err.Error()))
		}
		fmt.Println()
	}
}

func CmdNetworkCreate(c *cli.Context) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "Enter network name"
	networkName, err := ui.Ask(query, &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
		ValidateFunc: func(s string) error {
			doesMatch, err := regexp.MatchString("^[a-z0-9]{1,16}$", s)
			if doesMatch || err == nil {
				return nil
			}

			return fmt.Errorf("network name must be lowercase alphanumeric between 1 and 16 characters")
		},
	})
	if err != nil {
		competition.LogFatal("Fatal Error: " + err.Error())
	}

	query = "Enter subdomain"
	subdomain, err := ui.Ask(query, &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
		ValidateFunc: func(s string) error {
			doesMatch, err := regexp.MatchString("^[a-z0-9]{1,8}$", s)
			if doesMatch || err == nil {
				return nil
			}

			return fmt.Errorf("subdomain must be lowercase alphanumeric between 1 and 8 characters")
		},
	})
	if err != nil {
		competition.LogFatal("Fatal Error: " + err.Error())
	}

	query = "Enter network CIDR"
	cidr, err := ui.Ask(query, &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
		ValidateFunc: func(s string) error {
			_, _, err := net.ParseCIDR(s)
			if err != nil {
				return fmt.Errorf("not a valid CIDR")
			}

			return nil
		},
	})
	if err != nil {
		competition.LogFatal("Fatal Error: " + err.Error())
	}

	network := competition.Network{
		Name:      networkName,
		Subdomain: subdomain,
		CIDR:      cidr,
	}

	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME: " + err.Error())
	}

	comp.CurrentEnv().CreateNetwork(&network)
}
