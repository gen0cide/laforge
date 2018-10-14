package main

import (
	"errors"
	"io/ioutil"

	"github.com/gen0cide/laforge/core"
	"github.com/k0kubun/pp"
	"github.com/urfave/cli"
)

var (
	betaCommand = cli.Command{
		Name:      "beta",
		Usage:     "Developer subcommand for testing new Laforge functionality.",
		UsageText: "laforge beta",
		Action:    performbeta,
	}
)

func performbeta(c *cli.Context) error {
	arg0 := c.Args().Get(0)
	if arg0 == "" {
		return errors.New("must provide a file to show")
	}

	data, err := ioutil.ReadFile(arg0)
	if err != nil {
		return err
	}

	hostObj := &core.Host{}

	err = core.HCLBytesToObject(data, hostObj)

	if err != nil {
		return err
	}

	pp.Println(hostObj)
	return nil
}
