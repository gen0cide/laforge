package command

import (
	"fmt"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
)

func CmdTf(c *cli.Context) {
	// a := competition.TFVar{
	// 	Name:        "mygod",
	// 	Description: "what in the hell",
	// 	Value:       "ohyea1234",
	// }

	// b, err := competition.TFRender(a)
	// if err != nil {
	// 	competition.LogFatal(err.Error())
	// }
	// fmt.Printf("%s\n", b)

	cli.ShowAppHelpAndExit(c, 0)

}

func CmdTfPlan(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	competition.Log("** Run the following command **")
	competition.LogPlain(fmt.Sprintf("cd %s && terraform plan", env.TfDir()))
}

func CmdTfApply(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	competition.Log("** Run the following command **")
	competition.LogPlain(fmt.Sprintf("cd %s && terraform apply", env.TfDir()))
}

func CmdTfDestroy(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	competition.Log("** Run the following command **")
	competition.LogPlain(fmt.Sprintf("cd %s && terraform destroy", env.TfDir()))
}

func CmdTfNuke(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	competition.Log("** Run the following command **")
	competition.LogPlain(fmt.Sprintf("cd %s && terraform destroy -force -parallelism=50", env.TfDir()))
}

func TFCheck() {
	_, err := exec.LookPath("terraform")
	if err != nil {
		competition.LogFatal("The terraform executable could not be found in your $PATH!\n\t* Download it at https://www.terraform.io/downloads.html")
	}
}
