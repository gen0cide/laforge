package main

import (
	"io/ioutil"

	"github.com/alecthomas/chroma/quick"
	"github.com/fatih/color"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/gen0cide/laforge"
	"github.com/urfave/cli"
)

// laforge global => prompts the user for specific information relating to the global laforge configuration

var (
	globalCommand = cli.Command{
		Name:      "global-config",
		Usage:     "configures global settings (similar to git config.global)",
		UsageText: "laforge global-config",
		Action:    performGlobal,
		Flags:     []cli.Flag{},
	}

	testUser = laforge.Laforge{
		User: laforge.User{
			Name:  "kimbo slice",
			ID:    "kimbo",
			Email: "kimbo@hardhitters.com",
		},
	}
)

func performGlobal(c *cli.Context) error {
	laforge.SetLogLevel("info")
	gcl, err := laforge.LocateGlobalConfig()
	if err != nil {
		if err != laforge.ErrNoGlobalConfig {
			return err
		}
		cliLogger.Infof("No config found!")
		return laforge.UserWizard()
	}
	data, err := ioutil.ReadFile(gcl)
	if err != nil {
		return err
	}
	cliLogger.Infof("Existing config found:")
	err = quick.Highlight(color.Output, string(data), "terraform", "terminal", "monokai")
	if err != nil {
		return err
	}
	name := false
	prompt := &survey.Confirm{
		Message: "Do you want to reconfigure?",
	}
	survey.AskOne(prompt, &name, nil)
	if name {
		return laforge.UserWizard()
	}
	return nil
}
