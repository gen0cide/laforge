package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
)

func CmdSsh(c *cli.Context) {
	fmt.Println(competition.GetPublicIP())

}
