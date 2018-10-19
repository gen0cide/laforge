package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/gen0cide/laforge/core"
	"github.com/gobwas/glob"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/urfave/cli"
)

var (
	debugCommand = cli.Command{
		Name:      "debug",
		Usage:     "Allows for expirimental debugging features that only laforge developers would find useful",
		UsageText: "laforge debug ...",
		Action:    performdebug,
	}
)

func performdebug(c *cli.Context) error {
	base, err := core.Bootstrap()
	if err != nil {
		if _, ok := err.(hcl.Diagnostics); ok {
			return errors.New("aborted due to parsing error")
		}
		cliLogger.Errorf("Error encountered during bootstrap: %v", err)
		os.Exit(1)
	}

	err = base.AssertMinContext(core.BuildContext)
	if err != nil {
		cliLogger.Errorf("Must be in a team context to use this command: %v", err)
		os.Exit(1)
	}

	args := c.Args().Get(0)

	snap, err := core.NewSnapshotFromEnv(base.CurrentEnv)
	if err != nil {
		return err
	}

	build, ok := snap.Objects[path.Join(base.CurrentEnv.Path(), base.CurrentEnv.Builder)].(*core.Build)
	if !ok {
		return errors.New("builder was not able to resolve object of type Build")
	}
	base.CurrentBuild = build

	err = snap.Sort()
	if err != nil {
		panic(err)
	}

	if args == "" {
		data, err := json.MarshalIndent(snap, "", "  ")
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", string(data))
		return nil
	}

	if args != "" {
		pat := c.Args().Get(1)
		if pat == "" {
			return errors.New("must provide a comparison pattern")
		}
		g, err := glob.Compile(pat, '/')
		if err != nil {
			return err
		}

		core.SetLogLevel("info")

		for key, meta := range snap.Metadata {
			if !g.Match(key) {
				continue
			}
			cliLogger.Infof("*** INLINE OBJECT *** Parents Of %s (gid=%d) (checksum=%x):", key, meta.GID, meta.Checksum)
			for _, x := range meta.ParentDeps {
				fmt.Printf("  <- (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.GetChecksum()))
			}
			cliLogger.Infof("*** INLINE OBJECT *** Children Of %s (gid=%d) (checksum=%x):", key, meta.GID, meta.Checksum)
			for _, x := range meta.ChildDeps {
				fmt.Printf("  -> (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.GetChecksum()))
			}
		}

		statedata, err := ioutil.ReadFile(args)
		if err != nil {
			return err
		}

		savedSnap := core.NewEmptySnapshot()
		err = json.Unmarshal(statedata, savedSnap)
		if err != nil {
			return err
		}

		err = savedSnap.RebuildGraph()
		if err != nil {
			return err
		}

		for key, meta := range savedSnap.Metadata {
			if !g.Match(key) {
				continue
			}
			cliLogger.Infof("~~~ PERSISTED OBJECT ~~~ Parents Of %s (gid=%d) (checksum=%x):", key, meta.GID, meta.Checksum)
			for _, x := range meta.ParentDeps {
				fmt.Printf("  <- (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.GetChecksum()))
			}
			cliLogger.Infof("~~~ PERSISTED OBJECT ~~~ Children Of %s (gid=%d) (checksum=%x):", key, meta.GID, meta.Checksum)
			for _, x := range meta.ChildDeps {
				fmt.Printf("  -> (gid=%d) %s (checksum=%s)\n", x.GetGID(), color.YellowString(x.GetID()), color.CyanString("%x", x.GetChecksum()))
			}
		}

	}

	return nil
}
