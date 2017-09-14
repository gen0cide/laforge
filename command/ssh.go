package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
)

func CmdSsh(c *cli.Context) {
	TFCheck()
	this := "alextest-a0"
	dp := competition.DeterminedPassword(this)
	competition.Log(fmt.Sprintf("Determined Password: %s", dp))
}
