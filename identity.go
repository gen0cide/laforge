package laforge

import (
	"fmt"

	"github.com/imdario/mergo"
)

// Identity defines a generic human identity primative that can be extended into Employee, Customer, Client, etc.
type Identity struct {
	ID          string     `hcl:",label" json:"id,omitempty"`
	Firstname   string     `hcl:"firstname,attr" json:"firstname,omitempty"`
	Lastname    string     `hcl:"lastname,attr" json:"lastname,omitempty"`
	Email       string     `hcl:"email,attr" json:"email,omitempty"`
	Password    string     `hcl:"password,attr" json:"password,omitempty"`
	Description string     `hcl:"description,attr" json:"description,omitempty"`
	AvatarFile  string     `hcl:"avatar_file,attr" json:"avatar_file,omitempty"`
	OnConflict  OnConflict `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller      Caller     `json:"-"`
}

// Update performs a patching operation on source (i) with diff (diff), using the diff's merge conflict settings as appropriate.
func (i *Identity) Update(diff *Identity) error {
	newCaller := i.Caller.Stack(diff.Caller)
	switch diff.OnConflict.Do {
	case "":
		err := mergo.Merge(i, diff, mergo.WithOverride)
		i.Caller = newCaller
		return err
	case "overwrite":
		conflict := i.OnConflict
		*i = *diff
		i.OnConflict = conflict
		return nil
	case "inherit":
		conflict := i.OnConflict
		err := mergo.Merge(diff, i, mergo.WithOverride)
		*i = *diff
		i.Caller = newCaller
		i.OnConflict = conflict
		return err
	case "panic":
		return NewMergeConflict(i, diff, i.ID, diff.ID, i.Caller.Current(), diff.Caller.Current())
	default:
		return fmt.Errorf("invalid conflict strategy %s in %s", diff.OnConflict.Do, diff.Caller.Current().CallerFile)
	}
}
