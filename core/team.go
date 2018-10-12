package core

import (
	"github.com/pkg/errors"
)

// Team represents a team specific object existing within an environment
type Team struct {
	ID               string                      `hcl:"id,label" json:"id,omitempty"`
	TeamNumber       int                         `hcl:"team_number,attr" json:"team_number,omitempty"`
	BuildID          string                      `hcl:"build_id,attr" json:"build_id,omitempty"`
	EnvironmentID    string                      `hcl:"environment_id,attr" json:"environment_id,omitempty"`
	Config           map[string]string           `hcl:"config,attr" json:"config,omitempty"`
	Tags             map[string]string           `hcl:"tags,attr" json:"tags,omitempty"`
	ProvisionedHosts []*ProvisionedHost          `hcl:"provisioned_host,block" json:"provisioned_hosts,omitempty"`
	ActiveHosts      map[string]*ProvisionedHost `json:"active_hosts,omitempty"`
	Maintainer       *User                       `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Build            *Build                      `json:"-"`
	Environment      *Environment                `json:"-"`
	RelBuildPath     string                      `json:"-"`
	OnConflict       OnConflict                  `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller           Caller                      `json:"-"`
}

// GetCaller implements the Mergeable interface
func (t *Team) GetCaller() Caller {
	return t.Caller
}

// GetID implements the Mergeable interface
func (t *Team) GetID() string {
	return t.ID
}

// GetOnConflict implements the Mergeable interface
func (t *Team) GetOnConflict() OnConflict {
	return t.OnConflict
}

// SetCaller implements the Mergeable interface
func (t *Team) SetCaller(ca Caller) {
	t.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (t *Team) SetOnConflict(o OnConflict) {
	t.OnConflict = o
}

// Swap implements the Mergeable interface
func (t *Team) Swap(m Mergeable) error {
	rawVal, ok := m.(*Team)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", t, m)
	}
	*t = *rawVal
	return nil
}
