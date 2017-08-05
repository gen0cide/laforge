package command

import (
	"github.com/codegangsta/cli"
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
	return
}

func CmdTfApply(c *cli.Context) {
	return
}

func CmdTfDestroy(c *cli.Context) {
	return
}

func CmdTfNuke(c *cli.Context) {
	return
}
