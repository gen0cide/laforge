package laforge

import (
	"github.com/pkg/errors"
)

// Command represents an executable command that can be defined as part of a host configuration step
type Command struct {
	ID           string     `hcl:",label" json:"id,omitempty"`
	Name         string     `hcl:"name,attr" json:"name,omitempty"`
	Description  string     `hcl:"description,attr" json:"description,omitempty"`
	Program      string     `hcl:"program,attr" json:"program,omitempty"`
	Args         []string   `hcl:"args,attr" json:"args,omitempty"`
	IgnoreErrors bool       `hcl:"ignore_errors,attr" json:"ignore_errors,omitempty"`
	Cooldown     int        `hcl:"cooldown,attr" json:"cooldown,omitempty"`
	IO           IO         `hcl:"io,block" json:"io,omitempty"`
	Disabled     bool       `hcl:"disabled,attr" json:"disabled,omitempty"`
	OnConflict   OnConflict `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller       Caller     `json:"-"`
	Maintainer   *User      `hcl:"maintainer,block" json:"maintainer,omitempty"`
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
	return c.OnConflict
}

// SetCaller implements the Mergeable interface
func (c *Command) SetCaller(ca Caller) {
	c.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (c *Command) SetOnConflict(o OnConflict) {
	c.OnConflict = o
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
