package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gen0cide/laforge/core"
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
		return err
	}

	if c.Args().Get(0) == "" {
		pp.Println(base)
		return nil
	}

	switch strings.ToLower(c.Args().Get(0)) {
	case "build":
		pp.Println(base.Build)
	case "competition":
		pp.Println(base.Competition)
	case "dns":
		pp.Println(base.Competition.DNS)
	case "dns_record":
		param := c.Args().Get(1)
		if len(param) == 0 {
			return errors.New("second argument must be supplied with this type")
		}
		rec, found := base.DNSRecords[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "dns_record")
		}
		pp.Println(rec)
	case "command":
		param := c.Args().Get(1)
		if len(param) == 0 {
			return errors.New("second argument must be supplied with this type")
		}
		rec, found := base.DNSRecords[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "dns_record")
		}
		pp.Println(rec)
	case "environment":
		pp.Println(base.Environment)
	case "host":
		param := c.Args().Get(1)
		if len(param) == 0 {
			return errors.New("second argument must be supplied with this type")
		}
		rec, found := base.Hosts[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "host")
		}
		pp.Println(rec)
	case "identity":
		param := c.Args().Get(1)
		if len(param) == 0 {
			return errors.New("second argument must be supplied with this type")
		}
		rec, found := base.Identities[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "identity")
		}
		pp.Println(rec)
	case "network":
		param := c.Args().Get(1)
		if len(param) == 0 {
			return errors.New("second argument must be supplied with this type")
		}
		rec, found := base.Networks[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "network")
		}
		pp.Println(rec)
	case "remote_file":
		param := c.Args().Get(1)
		if len(param) == 0 {
			return errors.New("second argument must be supplied with this type")
		}
		rec, found := base.Files[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "remote_file")
		}
		pp.Println(rec)
	case "script":
		param := c.Args().Get(1)
		if len(param) == 0 {
			return errors.New("second argument must be supplied with this type")
		}
		rec, found := base.Scripts[param]
		if !found {
			return fmt.Errorf("object with id %s and type %s could not be found in tree", param, "script")
		}
		pp.Println(rec)
	case "team":
		pp.Println(base.Team)
	default:
		return errors.New("argument is not a known datatype")
	}

	return nil
}
