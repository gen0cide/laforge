package command

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/alecthomas/chroma/quick"
	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/deprecated/competition"
	input "github.com/tcnksm/go-input"
)

func CmdHost(c *cli.Context) {
	cli.ShowAppHelpAndExit(c, 0)
}

func CmdHostLs(c *cli.Context) {
	competition.ValidateEnv()
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME Config: " + err.Error())
	}
	env := comp.CurrentEnv()
	if env == nil {
		competition.LogFatal("Cannot load environment! (Check ~/.lf_env)")
	}
	hosts := env.ParseHosts()
	for name, host := range hosts {
		competition.Log("----------------------------------------------------------------------")
		competition.Log("HOST >>> " + name)
		competition.Log(filepath.Join(env.HostsDir(), fmt.Sprintf("%s.yml", name)))
		data, err := json.MarshalIndent(host, "", "  ")
		if err != nil {
			competition.LogError(fmt.Sprintf("Error printing host information: host=%s.yml error=%s", name, err.Error()))
		}
		err = quick.Highlight(os.Stdout, string(data), "json", "terminal", "vim")
		if err != nil {
			competition.LogError(fmt.Sprintf("Error printing host information: host=%s.yml error=%s", name, err.Error()))
		}
		fmt.Println()
	}
}

func CmdHostCreate(c *cli.Context) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "Enter hostname"
	newHostname, err := ui.Ask(query, &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
		ValidateFunc: func(s string) error {
			doesMatch, err := regexp.MatchString("^[a-z0-9]{1,16}$", s)
			if doesMatch || err == nil {
				return nil
			}

			return fmt.Errorf("hostname must be lowercase alphanumeric between 1 and 16 characters")
		},
	})
	if err != nil {
		competition.LogFatal("Fatal Error: " + err.Error())
	}

	query = "Enter Last Octet of IP"
	lastOctetString, err := ui.Ask(query, &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
		ValidateFunc: func(s string) error {
			doesMatch, err := regexp.MatchString("^[0-9]{1,3}$", s)
			if !doesMatch || err != nil {
				return fmt.Errorf("last octet must be a number (5-253)")
			}
			octet, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("not a valid integer")
			}

			if octet > 4 && octet < 254 {
				return nil
			}
			return fmt.Errorf("Last octet must be a number (5-253)")
		},
	})
	if err != nil {
		competition.LogFatal("Fatal Error: " + err.Error())
	}
	lastOctet, _ := strconv.Atoi(lastOctetString)

	query = "Enter instance size"
	instanceSize, err := ui.Ask(query, &input.Options{
		// Read the default val from env var
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		competition.LogFatal("Fatal Error: " + err.Error())
	}

	host := competition.Host{
		Hostname:     newHostname,
		LastOctet:    lastOctet,
		InstanceSize: instanceSize,
	}

	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME: " + err.Error())
	}

	comp.CurrentEnv().CreateHost(&host)
}
