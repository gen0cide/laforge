package command

import (
	"io/ioutil"
	"os"

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

	os.RemoveAll(env.TfScriptsDir())
	os.MkdirAll(env.TfScriptsDir(), 0755)

	finalTFTemplate, err := printer.Format(competition.RenderTB("infra.tf", &tb))
	if err != nil {
		competition.LogFatal("Terraform Configuration Syntax Error - Contact alex ASAP.")
	}

	ioutil.WriteFile(env.TfFile(), finalTFTemplate, 0644)

	competition.Log("Wrote Terraform configuration to: " + env.TfFile())
}
