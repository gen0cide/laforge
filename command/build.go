package command

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/bradfitz/iter"
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

	var tpl bytes.Buffer

	tmpl := template.New("test")
	tmpl.Funcs(template.FuncMap{"N": iter.N})
	newTmpl, err := tmpl.Parse(string(competition.MustAsset("infra.tf")))
	if err != nil {
		panic(err)
	}

	if err := newTmpl.Execute(&tpl, tb); err != nil {
		panic(err)
	}

	finalTFTemplate, err := printer.Format(tpl.Bytes())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(finalTFTemplate))

}
