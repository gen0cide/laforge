package laforge

import (
	"fmt"

	"github.com/imdario/mergo"
)

// Script defines a configurable type for an executable script object within the laforge configuration
type Script struct {
	Name         string            `hcl:"name,label" json:"name,omitempty"`
	Language     string            `hcl:"language,attr" json:"language,omitempty"`
	Description  string            `hcl:"description,attr" json:"description,omitempty"`
	Maintainer   User              `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Source       string            `hcl:"source,attr" json:"source,omitempty"`
	Cooldown     int               `hcl:"cooldown,attr" json:"cooldown,omitempty"`
	IgnoreErrors bool              `hcl:"ignore_errors,attr" json:"ignore_errors,omitempty"`
	Args         []string          `hcl:"args,attr" json:"args,omitempty"`
	IO           IO                `hcl:"io,block" json:"io,omitempty"`
	Disabled     bool              `hcl:"disabled,attr" json:"disabled,omitempty"`
	Vars         map[string]string `hcl:"vars,attr" json:"vars,omitempty"`
	Tags         map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	OnConflict   OnConflict        `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller       Caller            `json:"-"`
}

// Update performs a patching operation on source (s) with diff (diff), using the diff's merge conflict settings as appropriate.
func (s *Script) Update(diff *Script) error {
	newCaller := s.Caller.Stack(diff.Caller)
	switch diff.OnConflict.Do {
	case "":
		err := mergo.Merge(s, diff, mergo.WithOverride)
		s.Caller = newCaller
		return err
	case "overwrite":
		conflict := s.OnConflict
		*s = *diff
		s.OnConflict = conflict
		return nil
	case "inherit":
		conflict := s.OnConflict
		err := mergo.Merge(diff, s, mergo.WithOverride)
		*s = *diff
		s.Caller = newCaller
		s.OnConflict = conflict
		return err
	case "panic":
		return NewMergeConflict(s, diff, s.Name, diff.Name, s.Caller.Current(), diff.Caller.Current())
	default:
		return fmt.Errorf("invalid conflict strategy %s in %s", diff.OnConflict.Do, diff.Caller.Current().CallerFile)
	}
}

// ResolveSource attempts to locate the referenced source file with a laforge base configuration
func (s *Script) ResolveSource(base *Laforge) error {
	return nil
}
