package main

import (
	"github.com/gen0cide/laforge"
	"github.com/gen0cide/laforge/statusui"
	"github.com/urfave/cli"
)

// laforge status => show the current laforge configuration status

var (
	statusCommand = cli.Command{
		Name:      "status",
		Usage:     "displays the current laforge environment status",
		UsageText: "laforge status",
		Action:    performStatus,
		Flags:     []cli.Flag{},
	}
)

func performStatus(c *cli.Context) error {
	laforge.SetLogLevel("info")
	loadables := []string{}
	var base *laforge.Laforge
	gcl, err := laforge.LocateGlobalConfig()
	if err != nil {
		if err != laforge.ErrNoGlobalConfig {
			return err
		}
		cliLogger.Errorf("No config found! Run laforge global-config to fix!")
		return nil
	}
	loadables = append(loadables, gcl)
	err = nil
	ecl, err := laforge.LocateEnvConfig()
	if err != nil {
		if err != laforge.ErrNoConfigRootReached {
			return err
		}
	}
	err = nil
	bcl, err := laforge.LocateBaseConfig()
	if err != nil {
		if err != laforge.ErrNoConfigRootReached {
			return err
		}
	}
	err = nil
	if bcl == "" && ecl == "" {
		cliLogger.Errorf("No base.laforge or env.laforge found in your current directory tree!")
		return nil
	}
	if ecl != "" {
		loadables = append(loadables, ecl)
		base, err = laforge.LoadFiles(loadables...)
		if err != nil {
			return err
		}
	} else {
		loadables = append(loadables, bcl)
		base, err = laforge.LoadFiles(loadables...)
		if err != nil {
			return err
		}
	}

	return statusui.RenderLaforgeStatusUI(base)
}
