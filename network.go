package laforge

import (
	"fmt"

	"github.com/imdario/mergo"
)

// Network defines a network within a competition environment
type Network struct {
	Name       string            `hcl:"name,label" json:"name,omitempty"`
	CIDR       string            `hcl:"cidr,attr" json:"cidr,omitempty"`
	VDIVisible bool              `hcl:"vdi_visible,attr" json:"vdi_visible,omitempty"`
	Vars       map[string]string `hcl:"vars,attr" json:"vars,omitempty"`
	Tags       map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	OnConflict OnConflict        `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller     Caller            `json:"-"`
}

// IncludedNetwork is a configuration type used to parse included_hosts out of an environment config.
type IncludedNetwork struct {
	Name  string   `hcl:"name,label" json:"name,omitempty"`
	Hosts []string `hcl:"included_hosts,attr" json:"included_hosts:omitempty"`
}

// Update performs a patching operation on source (n) with diff (diff), using the diff's merge conflict settings as appropriate.
func (n *Network) Update(diff *Network) error {
	switch diff.OnConflict.Do {
	case "":
		return mergo.Merge(n, diff, mergo.WithOverride)
	case "overwrite":
		conflict := n.OnConflict
		*n = *diff
		n.OnConflict = conflict
		return nil
	case "inherit":
		callerCopy := diff.Caller
		conflict := n.OnConflict
		err := mergo.Merge(diff, n, mergo.WithOverride)
		*n = *diff
		n.Caller = callerCopy
		n.OnConflict = conflict
		return err
	case "panic":
		return NewMergeConflict(n, diff, n.Name, diff.Name, n.Caller.Current(), diff.Caller.Current())
	default:
		return fmt.Errorf("invalid conflict strategy %s in %s", diff.OnConflict.Do, diff.Caller.Current().CallerFile)
	}
}
