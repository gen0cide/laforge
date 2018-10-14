package core

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

// Team represents a team specific object existing within an environment
type Team struct {
	ID               string                      `hcl:"id,label" json:"id,omitempty"`
	TeamNumber       int                         `hcl:"team_number,attr" json:"team_number,omitempty"`
	BuildID          string                      `hcl:"build_id,attr" json:"build_id,omitempty"`
	EnvironmentID    string                      `hcl:"environment_id,attr" json:"environment_id,omitempty"`
	CompetitionID    string                      `hcl:"competition_id,attr" json:"competition_id,omitempty"`
	Config           map[string]string           `hcl:"config,attr" json:"config,omitempty"`
	Tags             map[string]string           `hcl:"tags,attr" json:"tags,omitempty"`
	ProvisionedHosts []*ProvisionedHost          `hcl:"provisioned_host,block" json:"provisioned_hosts,omitempty"`
	OnConflict       OnConflict                  `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Revision         int64                       `hcl:"revision,attr" json:"revision,omitempty"`
	Maintainer       *User                       `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Hosts            map[string]*ProvisionedHost `json:"-"`
	Build            *Build                      `json:"-"`
	Environment      *Environment                `json:"-"`
	Competition      *Competition                `json:"-"`
	RelBuildPath     string                      `json:"-"`
	TeamRoot         string                      `json:"-"`
	Caller           Caller                      `json:"-"`
}

// GetCaller implements the Mergeable interface
func (t *Team) GetCaller() Caller {
	return t.Caller
}

// GetID implements the Mergeable interface
func (t *Team) GetID() string {
	return filepath.Join(t.CompetitionID, t.EnvironmentID, t.BuildID, t.ID)
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

// SetID increments the revision and sets the team ID if needed
func (t *Team) SetID() string {
	t.Revision++
	if t.BuildID == "" && t.Build != nil {
		t.BuildID = t.Build.ID
	}
	if t.EnvironmentID == "" && t.Environment != nil {
		t.EnvironmentID = t.Environment.ID
	}
	if t.CompetitionID == "" && t.Competition != nil {
		t.CompetitionID = t.Competition.ID
	}
	t.ID = fmt.Sprintf("%d", t.TeamNumber)
	return t.ID
}
