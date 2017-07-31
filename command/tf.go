package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
)

func CmdTf(c *cli.Context) {
	a := competition.TFVar{
		Name:        "mygod",
		Description: "what in the hell",
		Value:       "ohyea1234",
	}

	b, err := competition.TFRender(a)
	if err != nil {
		competition.LogFatal(err.Error())
	}
	fmt.Printf("%s\n", b)

}
