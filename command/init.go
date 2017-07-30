package command

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	comp "github.com/gen0cide/laforge/competition"
	input "github.com/tcnksm/go-input"
)

func CmdInit(c *cli.Context) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "What directory would you like to set LF_HOME to?"
	newHome, err := ui.Ask(query, &input.Options{
		// Read the default val from env var
		Default:  os.Getenv("LF_HOME"),
		Required: true,
		Loop:     true,
	})
	if err != nil {
		comp.LogFatal("Fatal Error: " + err.Error())
	}
	os.Setenv("LF_HOME", newHome)
	if !comp.HomeExists() {
		query = "LF_HOME doesn't exist. Create it? [Y/n]"
		createDir, err := ui.Ask(query, &input.Options{
			Required: true,
			Loop:     true,
			ValidateFunc: func(s string) error {
				if s != "Y" && s != "n" {
					return fmt.Errorf("input must be Y or n")
				}
				return nil
			},
		})
		if err != nil {
			comp.LogFatal("Fatal Error: " + err.Error())
		}
		if createDir == "Y" {
			comp.Log("Creating LF_HOME...")
			comp.CreateHome()
			comp.Log("Created! Be sure to edit config/config.yml")
		}
	}
	comp.ValidateHome()
	comp.Log("LF_HOME is valid! To begin using laforge, issue the following command:")
	comp.LogPlain("export LF_HOME=" + comp.GetHome())
	comp.Log("Make sure to add this to your rc files:")
	comp.LogPlain("export LF_HOME=" + comp.GetHome())
}
