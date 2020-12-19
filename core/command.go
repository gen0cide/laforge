package core

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"

	"github.com/pkg/errors"
)

// Command represents an executable command that can be defined as part of a host configuration step
//easyjson:json
//nolint:maligned
type Command struct {
	ID           string            `hcl:"id,label" json:"id,omitempty"`
	Name         string            `hcl:"name,attr" json:"name,omitempty"`
	Description  string            `hcl:"description,attr" json:"description,omitempty"`
	Program      string            `hcl:"program,attr" json:"program,omitempty"`
	Args         []string          `hcl:"args,attr" json:"args,omitempty"`
	IgnoreErrors bool              `hcl:"ignore_errors,attr" json:"ignore_errors,omitempty"`
	Cooldown     int               `hcl:"cooldown,attr" json:"cooldown,omitempty"`
	Timeout      int               `hcl:"timeout,attr" json:"timeout,omitempty"`
	Disabled     bool              `hcl:"disabled,attr" json:"disabled,omitempty"`
	Vars         map[string]string `hcl:"vars,attr" json:"vars,omitempty"`
	Tags         map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	IO           *IO               `hcl:"io,block" json:"io,omitempty"`
	OnConflict   *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Maintainer   *User             `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Caller       Caller            `json:"-"`
}

// Hash implements the Hasher interface
func (c *Command) Hash() uint64 {
	iostr := "n/a"
	if c.IO != nil {
		iostr = c.IO.Stderr + c.IO.Stdin + c.IO.Stdout
	}

	return xxhash.Sum64String(
		fmt.Sprintf(
			"program=%v args=%v ignoreerrors=%v cooldown=%v io=%v disabled=%v vars=%v",
			c.Program,
			strings.Join(c.Args, ","),
			c.IgnoreErrors,
			c.Cooldown,
			iostr,
			c.Disabled,
			c.Vars,
		),
	)
}

// Path implements the Pather interface
func (c *Command) Path() string {
	return c.ID
}

// Base implements the Pather interface
func (c *Command) Base() string {
	return path.Base(c.ID)
}

// ValidatePath implements the Pather interface
func (c *Command) ValidatePath() error {
	if err := ValidateGenericPath(c.Path()); err != nil {
		return err
	}
	if topdir := strings.Split(c.Path(), `/`); topdir[1] != commandsDir {
		return fmt.Errorf("path %s is not rooted in /%s", c.Path(), topdir[1])
	}
	return nil
}

// GetCaller implements the Mergeable interface
func (c *Command) GetCaller() Caller {
	return c.Caller
}

// LaforgeID implements the Mergeable interface
func (c *Command) LaforgeID() string {
	return c.ID
}

// Fullpath implements the Pather interface
func (c *Command) Fullpath() string {
	return c.LaforgeID()
}

// ParentLaforgeID implements the Dependency interface
func (c *Command) ParentLaforgeID() string {
	return c.Path()
}

// Gather implements the Dependency interface
func (c *Command) Gather(g *Snapshot) error {
	return nil
}

// GetOnConflict implements the Mergeable interface
func (c *Command) GetOnConflict() OnConflict {
	if c.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *c.OnConflict
}

// SetCaller implements the Mergeable interface
func (c *Command) SetCaller(ca Caller) {
	c.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (c *Command) SetOnConflict(o OnConflict) {
	c.OnConflict = &o
}

// Kind implements the Provisioner interface
func (c *Command) Kind() string {
	return ObjectTypeCommand.String()
}

// CommandString is a template helper function to embed commands into the output
func (c *Command) CommandString() string {
	cmd := []string{c.Program}
	cmd = append(cmd, c.Args...)
	return strings.Join(cmd, " ")
}

// Swap implements the Mergeable interface
func (c *Command) Swap(m Mergeable) error {
	rawVal, ok := m.(*Command)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", c, m)
	}
	*c = *rawVal
	return nil
}

// CreateCommandEntry ...
func (c *Command) CreateCommandEntry(ctx context.Context, client *ent.Client) (*ent.Command, error) {
	tag, err := CreateTagEntry(c.ID, c.Tags, ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating command: %v", err)
		return nil, err
	}

	user, err := c.Maintainer.CreateUserEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating command: %v", err)
		return nil, err
	}

	command, err := client.Command.
		Create().
		SetName(c.Name).
		SetDescription(c.Description).
		SetProgram(c.Program).
		SetArgs(c.Args).
		SetIgnoreErrors(c.IgnoreErrors).
		SetDisabled(c.Disabled).
		SetCooldown(c.Cooldown).
		SetTimeout(c.Timeout).
		SetVars(c.Vars).
		AddUser(user).
		AddTag(tag).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating command: %v", err)
		return nil, err
	}

	cli.Logger.Debugf("command was created: ", command)
	return command, nil
}
