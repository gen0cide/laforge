package command

import (
	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
)

func CmdUse(c *cli.Context) {
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
