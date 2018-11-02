package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/hashicorp/hcl2/hclwrite"

	"github.com/urfave/cli"
)

var (
	fmtoverwrite = false
	fmtCommand   = cli.Command{
		Name:      "fmt",
		Usage:     "formats a laforge configuration file and prints it to stdout",
		UsageText: "laforge fmt",
		Action:    performfmt,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "write, w",
				Usage:       "Overwrites the file in place instead of printing to STDOUT.",
				Destination: &fmtoverwrite,
			},
		},
	}
)

func performfmt(c *cli.Context) error {
	errored := false
	for _, f := range c.Args() {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			cliLogger.Errorf("Could not read file %s: %v", f, err)
			errored = true
			continue
		}

		fmtData := hclwrite.Format(data)
		_ = fmtData
		buf := new(bytes.Buffer)
		err = quick.Highlight(buf, string(data), "terraform", "tokens", "api")
		if err != nil {
			cliLogger.Errorf("Could not format file %s: %v", f, err)
			errored = true
			continue
		}

		if fmtoverwrite {
			fi, err := os.Stat(f)
			if err != nil {
				cliLogger.Errorf("could not stat file %s: %v", f, err)
				errored = true
				continue
			}

			err = ioutil.WriteFile(f, fmtData, fi.Mode())
			if err != nil {
				cliLogger.Errorf("could not write file %s: %v", f, err)
				errored = true
				continue
			}

			continue
		}

		fmt.Printf("%s\n", string(fmtData))
	}

	if errored {
		return errors.New("failure formatting files")
	}

	return nil
}
