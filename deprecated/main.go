package main

//go:generate go-bindata -nometadata -pkg competition -prefix templates/ -o competition/bindata.go templates templates/terraform

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "Alex Levinson"
	app.Email = "gen0cide.threats@gmail.com"
	app.Usage = "Competition infrastructure management for the cloud."

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
