package command

import (
	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
)

func CmdDoctor(c *cli.Context) {
	competition.Log("Download Terraform binary into your path: https://www.terraform.io/downloads.html")
}
