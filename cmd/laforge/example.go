package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gen0cide/laforge/core"

	"github.com/alecthomas/chroma/quick"
	"github.com/urfave/cli"
)

var (
	noColor        = false
	exampleCommand = cli.Command{
		Name:      "example",
		Usage:     "generates an example laforge configuration object for reference",
		UsageText: "laforge example [OPTIONS] TYPE",
		Action:    performexample,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "no-color, c",
				Usage:       "does not perform syntax highlighting on standard out - good for writing to files",
				Destination: &noColor,
			},
		},
	}
)

func performexample(c *cli.Context) error {
	validTypes := []string{}
	for k := range core.ExampleObjects {
		validTypes = append(validTypes, k)
	}
	requestedType := c.Args().Get(0)
	if requestedType == "" {
		return fmt.Errorf("example command must be passed a known type: %v", validTypes)
	}

	obj, err := core.ExampleObjectByName(requestedType)
	if err != nil {
		return err
	}

	if noColor {
		fmt.Fprintf(color.Output, "%s", string(obj))
	} else {
		err := quick.Highlight(color.Output, string(obj), "terraform", "terminal", "monokai")
		if err != nil {
			return err
		}
	}
	return nil
}
