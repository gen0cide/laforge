package core

import "github.com/pkg/errors"

// DNSRecord is a configurable type for defining DNS entries related to this host in the core DNS infrastructure (if enabled)
type DNSRecord struct {
	ID         string            `hcl:"id,label" json:"id,omitempty"`
	Name       string            `hcl:"name,attr" json:"name,omitempty"`
	Values     []string          `hcl:"values,attr" json:"values,omitempty"`
	Type       string            `hcl:"type,attr" json:"type,omitempty"`
	Zone       string            `hcl:"zone,attr" json:"zone,omitempty"`
	Vars       map[string]string `hcl:"vars,attr" json:"vars,omitempty"`
	Tags       map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	Disabled   bool              `hcl:"disabled,attr" json:"disabled,omitempty"`
	OnConflict OnConflict        `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller     Caller            `json:"-"`
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

// Kind implements the Provisioner interface
func (r *DNSRecord) Kind() string {
	return "dns_record"
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

// Inherited is a boolean condition that is triggered when a DNS record is not statically defined
func (r *DNSRecord) Inherited() bool {
	return len(r.Values) == 0
}

// SetValue is an override which allows you to set the value of a DNS record during a template run
func (r *DNSRecord) SetValue(val string) {
	r.Values = append(r.Values, val)
}
