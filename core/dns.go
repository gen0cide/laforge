package core

import "github.com/pkg/errors"

// DNS represents a configurable type for the creation of competition DNS infrastructure
type DNS struct {
	ID         string            `hcl:"id,label" json:"id,omitempty"`
	Type       string            `hcl:"type,attr" json:"type,omitempty"`
	RootDomain string            `hcl:"root_domain,attr" json:"root_domain,omitempty"`
	DNSServers []string          `hcl:"dns_servers,attr" json:"dns_servers,omitempty"`
	NTPServers []string          `hcl:"ntp_servers,attr" json:"ntp_servers,omitempty"`
	Config     map[string]string `hcl:"config,attr" json:"config,omitempty"`
	OnConflict OnConflict        `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller     Caller            `json:"-"`
}

// GetCaller implements the Mergeable interface
func (d *DNS) GetCaller() Caller {
	return d.Caller
}

// GetID implements the Mergeable interface
func (d *DNS) GetID() string {
	return d.ID
}

// GetOnConflict implements the Mergeable interface
func (d *DNS) GetOnConflict() OnConflict {
	return d.OnConflict
}

// SetCaller implements the Mergeable interface
func (d *DNS) SetCaller(c Caller) {
	d.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (d *DNS) SetOnConflict(o OnConflict) {
	d.OnConflict = o
}

// Swap implements the Mergeable interface
func (d *DNS) Swap(m Mergeable) error {
	rawVal, ok := m.(*DNS)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", d, m)
	}
	*d = *rawVal
	return nil
}
