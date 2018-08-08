package core

import (
	"github.com/pkg/errors"
)

// Competition is a configurable type that holds competition wide settings
type Competition struct {
	ID           string            `hcl:",label" json:"id,omitempty"`
	BaseDir      string            `hcl:"base_dir,attr" json:"base_dir,omitempty"`
	RootPassword string            `hcl:"root_password,attr" json:"root_password,omitempty"`
	DNS          *DNS              `hcl:"dns,block" json:"dns,omitempty"`
	Remote       *Remote           `hcl:"remote,block" json:"remote,omitempty"`
	Config       map[string]string `hcl:"config,attr" json:"config,omitempty"`
	OnConflict   OnConflict        `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller       Caller            `json:"-"`
}

// GetCaller implements the Mergeable interface
func (c *Competition) GetCaller() Caller {
	return c.Caller
}

// GetID implements the Mergeable interface
func (c *Competition) GetID() string {
	return c.ID
}

// GetOnConflict implements the Mergeable interface
func (c *Competition) GetOnConflict() OnConflict {
	return c.OnConflict
}

// SetCaller implements the Mergeable interface
func (c *Competition) SetCaller(ca Caller) {
	c.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (c *Competition) SetOnConflict(oc OnConflict) {
	c.OnConflict = oc
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
