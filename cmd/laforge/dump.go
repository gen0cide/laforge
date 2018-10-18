package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gen0cide/laforge/core"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/k0kubun/pp"
	"github.com/urfave/cli"
)

var (
	dumpCommand = cli.Command{
		Name:      "dump",
		Usage:     "dumps the current configuration state in a pretty printed output",
		UsageText: "laforge dump",
		Action:    performdump,
	}
)

func performdump(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		if _, ok := err.(hcl.Diagnostics); ok {
			return errors.New("aborted due to parsing error")
		}
		return err
	}

	if c.Args().Get(0) == "" {
		pp.Println(base)
		return nil
	}

	switch strings.ToLower(c.Args().Get(0)) {
	case "build":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.Builds)
			os.Exit(0)
		}
		rec, found := base.Builds[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "build")
		}
		pp.Println(rec)
	case "competition":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.Competitions)
			os.Exit(0)
		}
		rec, found := base.Competitions[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "competition")
		}
		pp.Println(rec)
	case "environment":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.Environments)
			os.Exit(0)
		}
		rec, found := base.Environments[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "environment")
		}
		pp.Println(rec)
	case "dns_record":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.DNSRecords)
			os.Exit(0)
		}
		rec, found := base.DNSRecords[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "dns_record")
		}
		pp.Println(rec)
	case "command":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.Commands)
			os.Exit(0)
		}
		rec, found := base.DNSRecords[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "command")
		}
		pp.Println(rec)
	case "host":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.Hosts)
			os.Exit(0)
		}
		rec, found := base.Hosts[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "host")
		}
		pp.Println(rec)
	case "identity":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.Identities)
			os.Exit(0)
		}
		rec, found := base.Identities[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "identity")
		}
		pp.Println(rec)
	case "network":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.Networks)
			os.Exit(0)
		}
		rec, found := base.Networks[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "network")
		}
		pp.Println(rec)
	case "remote_file":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.RemoteFiles)
			os.Exit(0)
		}
		rec, found := base.RemoteFiles[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "remote_file")
		}
		pp.Println(rec)
	case "script":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.Scripts)
			os.Exit(0)
		}
		rec, found := base.Scripts[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "script")
		}
		pp.Println(rec)
	case "team":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.Teams)
			os.Exit(0)
		}
		rec, found := base.Teams[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "team")
		}
		pp.Println(rec)
	case "provisioned_host":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.ProvisionedHosts)
			os.Exit(0)
		}
		rec, found := base.ProvisionedHosts[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "provisioned_host")
		}
		pp.Println(rec)
	case "provisioned_network":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.ProvisionedNetworks)
			os.Exit(0)
		}
		rec, found := base.ProvisionedNetworks[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "provisioned_network")
		}
		pp.Println(rec)
	case "connection":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.Connections)
			os.Exit(0)
		}
		rec, found := base.Connections[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "connection")
		}
		pp.Println(rec)
	case "provisioning_step":
		param := c.Args().Get(1)
		if len(param) == 0 {
			pp.Println(base.ProvisioningSteps)
			os.Exit(0)
		}
		rec, found := base.ProvisioningSteps[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "provisioning_step")
		}
		pp.Println(rec)
	default:
		return errors.New("argument is not a known datatype")
	}

	return nil
}
