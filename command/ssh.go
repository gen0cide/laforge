package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
)

func CmdSsh(c *cli.Context) {
	TFCheck()
	hostName := c.Args().Get(0)
	if len(hostName) < 1 {
		competition.LogFatal("You did not provide an environment to use.")
	}
	dp := competition.DeterminedPassword(hostName)
	competition.Log(fmt.Sprintf("Determined Password: %s", dp))
}
