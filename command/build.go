package command

import (
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
	"github.com/hashicorp/hcl/hcl/printer"
)

func CmdBuild(c *cli.Context) {
	comp, env := InitConfig()

	tb := competition.TemplateBuilder{
		Environment: env,
		Competition: comp,
	}

	finalTFTemplate, err := printer.Format(competition.RenderTB("infra.tf", &tb))
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(env.TfFile(), finalTFTemplate, 0644)

	competition.Log("Wrote Terraform configuration to: " + env.TfFile())

}
