package command

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
	"github.com/masterzen/winrm"
)

func CmdPowershell(c *cli.Context) {
	hostname := c.Args().Get(0)
	if len(hostname) < 1 {
		competition.LogFatal("You did not provide a hostname to use.")
	}

	comp, env := InitConfig()
	sshHosts := env.NewSSHConfig()

	var (
		ip   string
		port = 5985
		user = "Administrator"
	)

	if val, ok := sshHosts.Hosts[hostname]; ok {
		ip = val
	} else {
		competition.LogFatal("Unknown host: " + hostname)
	}

	endpoint := winrm.NewEndpoint(ip, port, false, false, nil)
	client, err := winrm.NewClient(endpoint, user, comp.RootPassword)
	if err != nil {
		panic(err)
	}

	_, err = client.RunWithInput("powershell", os.Stdout, os.Stderr, os.Stdin)
	if err != nil {
		panic(err)
	}
}
