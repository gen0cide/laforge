package core

import "github.com/pkg/errors"

// DNS represents a configurable type for the creation of competition DNS infrastructure
type DNS struct {
	ID         string     `hcl:",label" json:"id,omitempty"`
	Type       string     `hcl:"type,attr" json:"type,omitempty"`
	RootDomain string     `hcl:"root_domain,attr" json:"root_domain,omitempty"`
	DNSServers []string   `hcl:"dns_servers,attr" json:"dns_servers,omitempty"`
	NTPServers []string   `hcl:"ntp_servers,attr" json:"ntp_servers,omitempty"`
	OnConflict OnConflict `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller     Caller     `json:"-"`
}

// DNSRecord is a configurable type for defining DNS entries related to this host in the core DNS infrastructure (if enabled)
type DNSRecord struct {
	ID         string     `hcl:",label" json:"id,omitempty"`
	Name       string     `hcl:"name,attr" json:"name,omitempty"`
	Value      string     `hcl:"value,attr" json:"value,omitempty"`
	Type       string     `hcl:"type,attr" json:"type,omitempty"`
	Disabled   bool       `hcl:"disabled,attr" json:"disabled,omitempty"`
	OnConflict OnConflict `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller     Caller     `json:"-"`
}

// GetCaller implements the Mergeable interface
func (r *DNSRecord) GetCaller() Caller {
	return r.Caller
}

// GetID implements the Mergeable interface
func (r *DNSRecord) GetID() string {
	return r.ID
}

// GetOnConflict implements the Mergeable interface
func (r *DNSRecord) GetOnConflict() OnConflict {
	return r.OnConflict
}

// SetCaller implements the Mergeable interface
func (r *DNSRecord) SetCaller(c Caller) {
	r.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (r *DNSRecord) SetOnConflict(o OnConflict) {
	r.OnConflict = o
}

// Swap implements the Mergeable interface
func (r *DNSRecord) Swap(m Mergeable) error {
	rawVal, ok := m.(*DNSRecord)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", r, m)
	}
	*r = *rawVal
	return nil
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
