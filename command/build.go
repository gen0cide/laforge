package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
	"github.com/hashicorp/hcl/hcl/printer"
)

func CmdBuild(c *cli.Context) {
	competition.ValidateEnv()
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME Config: " + err.Error())
	}
	env := comp.CurrentEnv()
	if env == nil {
		competition.LogFatal("Cannot load environment! (Check ~/.lf_env)")
	}
	tb := competition.TemplateBuilder{
		Environment: env,
		Competition: comp,
	}

	finalTFTemplate, err := printer.Format(competition.RenderTB("infra.tf", &tb))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(finalTFTemplate))

}
