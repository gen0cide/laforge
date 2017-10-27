package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/command"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:   "init",
		Usage:  "Configure your current LaForge environment.",
		Action: command.CmdInit,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "cdn",
		Usage:  "Manipulate files in the CDN.",
		Action: command.CmdCDN,
		Flags:  []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "ls",
				Usage:  "List files in the CDN.",
				Action: command.CmdCDNLs,
			},
			{
				Name:   "link",
				Usage:  "Get the public URL for a file in the CDN.",
				Action: command.CmdCDNLink,
			},
			{
				Name:   "upload",
				Usage:  "Upload a file to the CDN.",
				Action: command.CmdCDNUpload,
			},
			{
				Name:   "rm",
				Usage:  "Delete a file from the CDN.",
				Action: command.CmdCDNRm,
			},
		},
	},
	{
		Name:   "doctor",
		Usage:  "Check to see if you have the required dependencies to use LaForge.",
		Action: command.CmdDoctor,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "network",
		Usage:  "Manage Networks within your current LaForge environment.",
		Action: command.CmdNetwork,
		Flags:  []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "ls",
				Usage:  "List the current environment's networks.",
				Action: command.CmdNetworkLs,
			},
			{
				Name:   "create",
				Usage:  "Create a new network.",
				Action: command.CmdNetworkCreate,
			},
		},
	},
	{
		Name:   "host",
		Usage:  "Manage Hosts within your networks.",
		Action: command.CmdHost,
		Flags:  []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "ls",
				Usage:  "List the current environment's hosts.",
				Action: command.CmdHostLs,
			},
			{
				Name:   "create",
				Usage:  "Create a new host.",
				Action: command.CmdHostCreate,
			},
		},
	},
	{
		Name:   "env",
		Usage:  "Manage LaForge competition environment.",
		Action: command.CmdEnv,
		Flags:  []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "ls",
				Usage:  "List the current competition environments located in the LF_HOME path.",
				Action: command.CmdEnvLs,
			},
			{
				Name:   "use",
				Usage:  "Switch the current LaForge context to the specified competition environment.",
				Action: command.CmdEnvUse,
			},
			{
				Name:   "create",
				Usage:  "List the current competition environments located in the LF_HOME path.",
				Action: command.CmdEnvCreate,
			},
			{
				Name:   "bashconfig",
				Usage:  "Generate a bash env config for some productive aliases.",
				Action: command.CmdEnvBashConfig,
			},
			{
				Name:   "team-password",
				Usage:  "Deterministically generate the password for a given Pod ID.",
				Action: command.CmdEnvPassword,
			},
			{
				Name:   "sshconfig",
				Usage:  "Write an SSH configuration to TF_HOME/environments/TF_ENV/ssh.conf",
				Action: command.CmdEnvSshConfig,
			},
		},
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
		Subcommands: []cli.Command{
			{
				Name:   "init",
				Usage:  "Initialize the terraform directory.",
				Action: command.CmdTfInit,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "team, t",
						Value: 0,
						Usage: "The team you wish to perform the terraform actions on. (Default = 0)",
					},
					cli.IntFlag{
						Name:  "parallelism, p",
						Value: 10,
						Usage: "The number of parallel workers you want terraform to use. (Default = 10)",
					},
				},
			},
			{
				Name:   "plan",
				Usage:  "Plan the terraform changes with the state delta.",
				Action: command.CmdTfPlan,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "team, t",
						Value: 0,
						Usage: "The team you wish to perform the terraform actions on. (Default = 0)",
					},
					cli.IntFlag{
						Name:  "parallelism, p",
						Value: 10,
						Usage: "The number of parallel workers you want terraform to use. (Default = 10)",
					},
				},
			},
			{
				Name:   "apply",
				Usage:  "Apply the current terraform plan.",
				Action: command.CmdTfApply,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "team, t",
						Value: 0,
						Usage: "The team you wish to perform the terraform actions on. (Default = 0)",
					},
					cli.IntFlag{
						Name:  "parallelism, p",
						Value: 10,
						Usage: "The number of parallel workers you want terraform to use. (Default = 10)",
					},
				},
			},
			{
				Name:   "output",
				Usage:  "Show the terraform outputs.",
				Action: command.CmdTfOutput,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "team, t",
						Value: 0,
						Usage: "The team you wish to perform the terraform actions on. (Default = 0)",
					},
				},
			},
			{
				Name:   "refresh",
				Usage:  "Refresh the terraform state.",
				Action: command.CmdTfRefresh,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "team, t",
						Value: 0,
						Usage: "The team you wish to perform the terraform actions on. (Default = 0)",
					},
					cli.IntFlag{
						Name:  "parallelism, p",
						Value: 10,
						Usage: "The number of parallel workers you want terraform to use. (Default = 10)",
					},
				},
			},
			{
				Name:   "taint",
				Usage:  "Taint an object in the dependency graph.",
				Action: command.CmdTfTaint,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "team, t",
						Value: 0,
						Usage: "The team you wish to perform the terraform actions on. (Default = 0)",
					},
				},
			},
			{
				Name:   "state",
				Usage:  "Show the current dependency tree.",
				Action: command.CmdTfState,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "team, t",
						Value: 0,
						Usage: "The team you wish to perform the terraform actions on. (Default = 0)",
					},
					cli.IntFlag{
						Name:  "parallelism, p",
						Value: 10,
						Usage: "The number of parallel workers you want terraform to use. (Default = 10)",
					},
				},
			},
			{
				Name:   "destroy",
				Usage:  "Destroy the current environment.",
				Action: command.CmdTfDestroy,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "team, t",
						Value: 0,
						Usage: "The team you wish to perform the terraform actions on. (Default = 0)",
					},
					cli.IntFlag{
						Name:  "parallelism, p",
						Value: 10,
						Usage: "The number of parallel workers you want terraform to use. (Default = 10)",
					},
				},
			},
		},
	},
	{
		Name:   "ssh",
		Usage:  "Allows for administrative SSH access to hosts in the current competition environment.",
		Action: command.CmdSsh,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "powershell",
		Usage:  "Creates an interactive powershell session for a host in the current competition environment.",
		Action: command.CmdPowershell,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "app",
		Usage:  "Custom application package management within the current environment.",
		Action: command.CmdApp,
		Flags:  []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "ls",
				Usage:  "List this competition's current applications.",
				Action: command.CmdAppLs,
			},
			{
				Name:   "create",
				Usage:  "Create a new skeleton application.",
				Action: command.CmdAppCreate,
			},
			{
				Name:   "pkg",
				Usage:  "Package and upload the current application to S3.",
				Action: command.CmdAppPkg,
			},
		},
	},
	{
		Name:   "update",
		Usage:  "Updates laforge to your the latest release.",
		Action: command.CmdUpdate,
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
