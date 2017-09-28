package command

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
)

func CmdSsh(c *cli.Context) {
	cli.ShowAppHelpAndExit(c, 0)
}

func CmdSshPassword(c *cli.Context) {
	_, env := InitConfig()
	TFCheck()
	podID := c.Args().Get(0)
	if len(podID) < 1 {
		competition.LogFatal("You did not provide a Pod ID to use.")
	}
	podVal, err := strconv.Atoi(podID)
	if err != nil {
		competition.LogFatal("You did not supply a valid team number.")
	}
	dp := env.PodPassword(podVal)
	competition.Log(fmt.Sprintf("Determined Password: %s", dp))
}

func CmdSshConfig(c *cli.Context) {
	_, env := InitConfig()
	env.GenerateSSHConfig()
	competition.Log("SSH Config successfully saved to: " + env.SSHConfigPath())
	return
}
