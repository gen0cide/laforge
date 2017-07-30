package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/command"
)

var GlobalFlags = []cli.Flag{}

// Subcommands: []cli.Command{
// 	{
// 	  Name:  "add",
// 	  Usage: "add a new template",
// 	  Action: func(c *cli.Context) error {
// 	    fmt.Println("new task template: ", c.Args().First())
// 	    return nil
// 	  },
// 	},
// 	{
// 	  Name:  "remove",
// 	  Usage: "remove an existing template",
// 	  Action: func(c *cli.Context) error {
// 	    fmt.Println("removed task template: ", c.Args().First())
// 	    return nil
// 	  },
// 	},
// },

var Commands = []cli.Command{
	{
		Name:   "init",
		Usage:  "Configure your current LaForge environment.",
		Action: command.CmdInit,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "doctor",
		Usage:  "Check to see if you have the required dependencies to use LaForge.",
		Action: command.CmdDoctor,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "cd",
		Usage:  "Changes your current directory to your comeptition environment's directory.",
		Action: command.CmdCd,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "ls",
		Usage:  "List the current competition environments located in the LF_HOME path.",
		Action: command.CmdLs,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "env",
		Usage:  "Show LaForge's current configuration and competition environment.",
		Action: command.CmdEnv,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "use",
		Usage:  "Switch the current LaForge context to the specified competition environment.",
		Action: command.CmdUse,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "create",
		Usage:  "Create a new skeleton competition environment in the LF_HOME path.",
		Action: command.CmdCreate,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "validate",
		Usage:  "Test the current competition's configuration for errors.",
		Action: command.CmdValidate,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "build",
		Usage:  "Build the current competition environment into a terraform configuration.",
		Action: command.CmdBuild,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "tf",
		Usage:  "Perform terraform functions on the current competition environment.",
		Action: command.CmdTf,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "ssh",
		Usage:  "Wrapper for SSH functionality in the current competition environment.",
		Action: command.CmdSsh,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "pkg",
		Usage:  "Custom application package management within the current environment.",
		Action: command.CmdPkg,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
