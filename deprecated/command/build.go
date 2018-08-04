package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/deprecated/competition"
	"github.com/hashicorp/hcl/hcl/printer"
)

func CmdBuild(c *cli.Context) {
	comp, env := InitConfig()

	var wg sync.WaitGroup

	teams := env.TeamIDs()

	wg.Add(len(teams))

	competition.Log("Creating Team Directories...")
	for _, t := range teams {
		os.RemoveAll(env.TfScriptsDir(t))
		os.MkdirAll(env.TfScriptsDir(t), 0755)
	}

	for _, t := range teams {
		go func(c *competition.Competition, e *competition.Environment, pid int) {
			defer wg.Done()
			tb := competition.TemplateBuilder{
				Environment: e,
				Competition: c,
				PodID:       pid,
			}
			raw := competition.RenderTBV2("infra-v2.tf", &tb)

			finalTFTemplate, err := printer.Format(raw)
			if err != nil {
				competition.LogError("Terraform Configuration Syntax Error: " + err.Error())
				competition.LogPlain(string(raw))
				competition.LogFatal(" - Contact alex ASAP.")
			}

			ioutil.WriteFile(env.TfFile(pid), finalTFTemplate, 0644)

			competition.Log(fmt.Sprintf("Config Generated For Team %d at %s", pid, env.TfFile(pid)))
		}(comp, env, t)
	}

	wg.Wait()
}
