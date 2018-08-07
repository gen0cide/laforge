package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/gen0cide/laforge"
	"github.com/urfave/cli"
)

var (
	envCommand = cli.Command{
		Name:      "env",
		Usage:     "allows listing and the creation of new environments within a base competition config",
		UsageText: "options for environments",
		Subcommands: []cli.Command{
			{
				Name:   "list",
				Usage:  "list currently known environments",
				Action: listenv,
			},
			{
				Name:   "create",
				Usage:  "create a new environment",
				Action: createenv,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:        "force, f",
						Usage:       "force removes and deletes any conflicting directories (dangerous)",
						Destination: &overwrite,
					},
				},
			},
		},
	}
)

func listenv(c *cli.Context) error {
	base, err := laforge.Bootstrap()
	if err != nil {
		return err
	}

	envs, err := base.GetAllEnvs()
	if err != nil {
		return err
	}

	checkEnv := false
	if base.EnvRoot != "" {
		checkEnv = true
	}

	laforge.SetLogLevel("info")
	envList := []string{}
	for name, elf := range envs {
		label := ""
		if checkEnv && base.Environment.ID == elf.Environment.ID {
			label = fmt.Sprintf(" %s %s %s - %s", color.HiGreenString("*"), color.HiWhiteString("(current)"), color.HiGreenString(name), base.EnvRoot)
		} else {
			pn := ""
			for cf := range elf.PathRegistry.DB {
				if filepath.Base(cf.CallerFile) == "env.laforge" {
					pn = cf.CallerDir
					break
				}
			}
			label = fmt.Sprintf(" %s %s - %s", color.CyanString("*"), color.CyanString(name), pn)
		}
		envList = append(envList, label)
	}

	envListing := strings.Join(envList, "\n")

	cliLogger.Infof("Known Environments:\n%s", envListing)
	return nil
}

func createenv(c *cli.Context) error {
	base, err := laforge.Bootstrap()
	if err != nil {
		return err
	}

	name := c.Args().Get(0)
	if name == "" {
		return fmt.Errorf("must supply a name for the new environment! (laforge env create FOO)")
	}

	err = base.InitializeEnv(name, overwrite)
	if err != nil {
		return err
	}

	newPath := filepath.Join(base.BaseRoot, "envs", name)

	laforge.SetLogLevel("info")
	cliLogger.Infof("Successfully created new environment %s in directory %s", name, newPath)

	return nil
}
