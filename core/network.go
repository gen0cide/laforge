package core

import (
	"strings"

	"github.com/pkg/errors"
)

// Network defines a network within a competition environment
type Network struct {
	ID         string            `hcl:"id,label" json:"id,omitempty"`
	Name       string            `hcl:"name,attr" json:"name,omitempty"`
	CIDR       string            `hcl:"cidr,attr" json:"cidr,omitempty"`
	VDIVisible bool              `hcl:"vdi_visible,attr" json:"vdi_visible,omitempty"`
	Vars       map[string]string `hcl:"vars,attr" json:"vars,omitempty"`
	Tags       map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	OnConflict *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller     Caller            `json:"-"`
}

// IncludedNetwork is a configuration type used to parse included_hosts out of an environment config.
type IncludedNetwork struct {
	Name  string   `hcl:"name,label" json:"name,omitempty"`
	Hosts []string `hcl:"included_hosts,attr" json:"included_hosts,omitempty"`
}

// GetCaller implements the Mergeable interface
func (n *Network) GetCaller() Caller {
	return n.Caller
}

// GetID implements the Mergeable interface
func (n *Network) GetID() string {
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
