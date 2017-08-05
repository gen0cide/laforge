package command

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
	input "github.com/tcnksm/go-input"
)

func CmdEnv(c *cli.Context) {
	cli.ShowAppHelpAndExit(c, 0)
}

func CmdEnvLs(c *cli.Context) {
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME Config: " + err.Error())
	}

	competition.LogEnvs(comp.GetEnvs())
}

func CmdEnvUse(c *cli.Context) {
	envName := c.Args().Get(0)
	if len(envName) < 1 {
		competition.LogFatal("You did not provide an environment to use.")
	}
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME Config: " + err.Error())
	}
	comp.ChangeEnv(envName)

}

func CmdEnvCreate(c *cli.Context) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "Enter environment name"
	newEnvName, err := ui.Ask(query, &input.Options{
		// Read the default val from env var
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		competition.LogFatal("Fatal Error: " + err.Error())
	}
	query = "Enter environment prefix"
	newEnvPrefix, err := ui.Ask(query, &input.Options{
		// Read the default val from env var
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		competition.LogFatal("Fatal Error: " + err.Error())
	}
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME: " + err.Error())
	}
	comp.CreateEnv(newEnvName, newEnvPrefix)
}
