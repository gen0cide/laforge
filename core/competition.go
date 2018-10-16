package core

import (
	"fmt"

	"github.com/cespare/xxhash"
	"github.com/pkg/errors"
)

// Competition is a configurable type that holds competition wide settings
//easyjson:json
type Competition struct {
	ID           string            `hcl:"id,label" json:"id,omitempty"`
	BaseDir      string            `hcl:"base_dir,optional" json:"base_dir,omitempty"`
	RootPassword string            `hcl:"root_password,attr" json:"root_password,omitempty"`
	DNS          *DNS              `hcl:"dns,block" json:"dns,omitempty"`
	Remote       *Remote           `hcl:"remote,block" json:"remote,omitempty"`
	Config       map[string]string `hcl:"config,optional" json:"config,omitempty"`
	OnConflict   *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller       Caller            `json:"-"`
}

// Hash implements the Hasher interface
func (c *Competition) Hash() uint64 {
	dh := uint64(666)
	rh := uint64(666)
	if c.DNS != nil {
		dh = c.DNS.Hash()
	}
	if c.Remote != nil {
		rh = c.Remote.Hash()
	}

	return xxhash.Sum64String(
		fmt.Sprintf(
			"rpw=%v dns=%v remote=%v config=%v",
			c.RootPassword,
			dh,
			rh,
			c.Config,
		),
	)
}

// Path implements the Pather interface
func (c *Competition) Path() string {
	return c.ID
}

// Base implements the Pather interface
func (c *Competition) Base() string {
	return c.ID
}

// ValidatePath implements the Pather interface
func (c *Competition) ValidatePath() error {
	return nil
}

// GetCaller implements the Mergeable interface
func (c *Competition) GetCaller() Caller {
	return c.Caller
}

// LaforgeID implements the Mergeable interface
func (c *Competition) LaforgeID() string {
	return c.ID
}

// GetOnConflict implements the Mergeable interface
func (c *Competition) GetOnConflict() OnConflict {
	if c.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *c.OnConflict
}

// SetCaller implements the Mergeable interface
func (c *Competition) SetCaller(ca Caller) {
	c.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (c *Competition) SetOnConflict(oc OnConflict) {
	c.OnConflict = &oc
}

// Swap implements the Mergeable interface
func (c *Competition) Swap(m Mergeable) error {
	rawVal, ok := m.(*Competition)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", c, m)
	}
	*c = *rawVal
	return nil
}

// PasswordForHost is a template helper function to allow a overridden password to be retrieved
func (c *Competition) PasswordForHost(h *Host) string {
	if h == nil {
		return c.RootPassword
	}

	if h.OverridePassword == "" {
		return c.RootPassword
	}

	return h.OverridePassword
}
