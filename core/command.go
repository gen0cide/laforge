package core

import (
	"strings"

	"github.com/pkg/errors"
)

// Command represents an executable command that can be defined as part of a host configuration step
type Command struct {
	ID           string            `hcl:"id,label" json:"id,omitempty"`
	Name         string            `hcl:"name,attr" json:"name,omitempty"`
	Description  string            `hcl:"description,attr" json:"description,omitempty"`
	Program      string            `hcl:"program,attr" json:"program,omitempty"`
	Args         []string          `hcl:"args,attr" json:"args,omitempty"`
	IgnoreErrors bool              `hcl:"ignore_errors,attr" json:"ignore_errors,omitempty"`
	Cooldown     int               `hcl:"cooldown,attr" json:"cooldown,omitempty"`
	IO           *IO               `hcl:"io,block" json:"io,omitempty"`
	Disabled     bool              `hcl:"disabled,attr" json:"disabled,omitempty"`
	Vars         map[string]string `hcl:"vars,attr" json:"vars,omitempty"`
	Tags         map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	OnConflict   *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Maintainer   *User             `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Caller       Caller            `json:"-"`
}

// GetCaller implements the Mergeable interface
func (c *Command) GetCaller() Caller {
	return c.Caller
}

// GetID implements the Mergeable interface
func (c *Command) GetID() string {
	return c.ID
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
	return "command"
}

// CommandString is a template helper function to embed commands into the output
func (c *Command) CommandString() string {
	cmd := []string{c.Program}
	for _, x := range c.Args {
		cmd = append(cmd, x)
	}
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
