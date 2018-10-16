package core

import (
	"fmt"
	"path"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/pkg/errors"
)

// Network defines a network within a competition environment
//easyjson:json
type Network struct {
	ID         string            `hcl:"id,label" json:"id,omitempty"`
	Name       string            `hcl:"name,attr" json:"name,omitempty"`
	CIDR       string            `hcl:"cidr,attr" json:"cidr,omitempty"`
	VDIVisible bool              `hcl:"vdi_visible,optional" json:"vdi_visible,omitempty"`
	Vars       map[string]string `hcl:"vars,optional" json:"vars,omitempty"`
	Tags       map[string]string `hcl:"tags,optional" json:"tags,omitempty"`
	OnConflict *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller     Caller            `json:"-"`
}

// IncludedNetwork is a configuration type used to parse included_hosts out of an environment config.
//easyjson:json
type IncludedNetwork struct {
	Name  string   `hcl:"name,label" json:"name,omitempty"`
	Hosts []string `hcl:"included_hosts,attr" json:"included_hosts,omitempty"`
}

func (i *IncludedNetwork) String() string {
	return fmt.Sprintf("[network name=%s hosts=%s]", i.Name, strings.Join(i.Hosts, `,`))
}

// Hash implements the Hasher interface
func (n *Network) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"name=%v cidr=%v vdivisible=%v vars=%v",
			n.Name,
			n.CIDR,
			n.VDIVisible,
			n.Vars,
		),
	)
}

// Path implements the Pather interface
func (n *Network) Path() string {
	return n.ID
}

// Base implements the Pather interface
func (n *Network) Base() string {
	return path.Base(n.ID)
}

// ValidatePath implements the Pather interface
func (n *Network) ValidatePath() error {
	if err := ValidateGenericPath(n.Path()); err != nil {
		return err
	}
	if topdir := strings.Split(n.Path(), `/`); topdir[1] != "networks" {
		return fmt.Errorf("path %s is not rooted in /%s", n.Path(), topdir[1])
	}
	return nil
}

// GetCaller implements the Mergeable interface
func (n *Network) GetCaller() Caller {
	return n.Caller
}

// LaforgeID implements the Mergeable interface
func (n *Network) LaforgeID() string {
	return n.ID
}

// GetOnConflict implements the Mergeable interface
func (n *Network) GetOnConflict() OnConflict {
	if n.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *n.OnConflict
}

// SetCaller implements the Mergeable interface
func (n *Network) SetCaller(c Caller) {
	n.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (n *Network) SetOnConflict(o OnConflict) {
	n.OnConflict = &o
}

// Swap implements the Mergeable interface
func (n *Network) Swap(m Mergeable) error {
	rawVal, ok := m.(*Network)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", n, m)
	}
	*n = *rawVal
	return nil
}

// Octet is a template helper function to get a network's octet at a specified offset
func (n *Network) Octet() string {
	if n.CIDR == "" {
		return "NO_CIDR"
	}
	octets := strings.Split(n.CIDR, ".")
	if len(octets) <= 3 {
		return "INVALID_CIDR"
	}

	return octets[2]
}
