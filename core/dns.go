package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/pkg/errors"
)

// DNS represents a configurable type for the creation of competition DNS infrastructure
//easyjson:json
type DNS struct {
	ID         string            `hcl:"id,label" json:"id,omitempty"`
	Type       string            `hcl:"type,attr" json:"type,omitempty"`
	RootDomain string            `hcl:"root_domain,attr" json:"root_domain,omitempty"`
	DNSServers []string          `hcl:"dns_servers,attr" json:"dns_servers,omitempty"`
	NTPServers []string          `hcl:"ntp_servers,optional" json:"ntp_servers,omitempty"`
	Config     map[string]string `hcl:"config,optional" json:"config,omitempty"`
	OnConflict *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller     Caller            `json:"-"`
}

// Hash implements the Hasher interface
func (d *DNS) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"type=%v rd=%v ns=%v ntp=%v config=%v",
			d.Type,
			d.RootDomain,
			strings.Join(d.DNSServers, ","),
			strings.Join(d.NTPServers, ","),
			d.Config,
		),
	)
}

// GetCaller implements the Mergeable interface
func (d *DNS) GetCaller() Caller {
	return d.Caller
}

// LaforgeID implements the Mergeable interface
func (d *DNS) LaforgeID() string {
	return d.ID
}

// GetOnConflict implements the Mergeable interface
func (d *DNS) GetOnConflict() OnConflict {
	if d.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *d.OnConflict
}

// SetCaller implements the Mergeable interface
func (d *DNS) SetCaller(c Caller) {
	d.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (d *DNS) SetOnConflict(o OnConflict) {
	d.OnConflict = &o
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

// CreateDNSEntry ...
func (d *DNS) CreateDNSEntry(ctx context.Context, client *ent.Client) (*ent.DNS, error) {
	dns, err := client.DNS.
		Create().
		SetType(d.Type).
		SetRootDomain(d.RootDomain).
		SetDNSServers(d.DNSServers).
		SetNtpServers(d.NTPServers).
		SetConfig(d.Config).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating dns: %v", err)
		return nil, err
	}

	cli.Logger.Debugf("dns was created: ", dns)
	return dns, nil
}
