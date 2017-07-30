package command

import (
	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
)

func CmdLs(c *cli.Context) {
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME Config: " + err.Error())
	}

	competition.LogEnvs(comp.GetEnvs())
}
