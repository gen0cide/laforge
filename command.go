package laforge

import (
	"fmt"

	"github.com/imdario/mergo"
)

// Command represents an executable command that can be defined as part of a host configuration step
type Command struct {
	Name         string     `hcl:"name,label" json:"name,omitempty"`
	Program      string     `hcl:"program,attr" json:"program,omitempty"`
	Args         []string   `hcl:"args,attr" json:"args,omitempty"`
	IgnoreErrors bool       `hcl:"ignore_errors,attr" json:"ignore_errors,omitempty"`
	Cooldown     int        `hcl:"cooldown,attr" json:"cooldown,omitempty"`
	IO           IO         `hcl:"io,block" json:"io,omitempty"`
	Disabled     bool       `hcl:"disabled,attr" json:"disabled,omitempty"`
	OnConflict   OnConflict `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller       Caller     `json:"-"`
}

// Update performs a patching operation on source (c) with diff (diff), using the diff's merge conflict settings as appropriate.
func (c *Command) Update(diff *Command) error {
	switch diff.OnConflict.Do {
	case "":
		return mergo.Merge(c, diff, mergo.WithOverride)
	case "overwrite":
		conflict := c.OnConflict
		*c = *diff
		c.OnConflict = conflict
		return nil
	case "inherit":
		callerCopy := diff.Caller
		conflict := c.OnConflict
		err := mergo.Merge(diff, c, mergo.WithOverride)
		*c = *diff
		c.Caller = callerCopy
		c.OnConflict = conflict
		return err
	case "panic":
		return NewMergeConflict(c, diff, c.Name, diff.Name, c.Caller.Current(), diff.Caller.Current())
	default:
		return fmt.Errorf("invalid conflict strategy %s in %s", diff.OnConflict.Do, diff.Caller.Current().CallerFile)
	}
}
