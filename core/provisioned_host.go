package core

import (
	"path/filepath"

	"github.com/pkg/errors"
)

// ProvisionedHost defines a provisioned host within a team's environment (network neutral)
type ProvisionedHost struct {
	ID              string           `hcl:"id,label" json:"id,omitempty"`
	CompetitionID   string           `hcl:"competition_id,attr" json:"competition_id"`
	EnvironmentID   string           `hcl:"environment_id,attr" json:"environment_id,omitempty"`
	BuildID         string           `hcl:"build_id,attr" json:"build_id,omitempty"`
	TeamID          string           `hcl:"team_id,attr" json:"team_id,omitempty"`
	HostID          string           `hcl:"host_id,attr" json:"host_id,omitempty"`
	Active          bool             `hcl:"active,attr" json:"active,omitempty"`
	LocalAddr       string           `hcl:"local_addr,attr" json:"local_addr,omitempty"`
	RemoteAddr      string           `hcl:"remote_addr,attr" json:"remote_addr,omitempty"`
	SSHAuthConfig   *SSHAuthConfig   `hcl:"ssh,block" json:"ssh,omitempty"`
	WinRMAuthConfig *WinRMAuthConfig `hcl:"winrm,block" json:"winrm,omitempty"`
	Revision        int64            `hcl:"revision,attr" json:"revision,omitempty"`
	OnConflict      OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Competition     *Competition     `json:"-"`
	Environment     *Environment     `json:"-"`
	Build           *Build           `json:"-"`
	Team            *Team            `json:"-"`
	Host            *Host            `json:"-"`
	Caller          Caller           `json:"-"`
}

// IsSSH is a convenience method for checking if the provisioned host is setup for remote SSH
func (p *ProvisionedHost) IsSSH() bool {
	return p.SSHAuthConfig != nil
}

// IsWinRM is a convenience method for checking if the provisioned host is setup for remote WinRM
func (p *ProvisionedHost) IsWinRM() bool {
	return p.WinRMAuthConfig != nil
}

// GetCaller implements the Mergeable interface
func (p *ProvisionedHost) GetCaller() Caller {
	return p.Caller
}

// GetID implements the Mergeable interface
func (p *ProvisionedHost) GetID() string {
	return filepath.Join(p.CompetitionID, p.EnvironmentID, p.BuildID, p.TeamID, p.ID)
}

// GetOnConflict implements the Mergeable interface
func (p *ProvisionedHost) GetOnConflict() OnConflict {
	return p.OnConflict
}

// SetCaller implements the Mergeable interface
func (p *ProvisionedHost) SetCaller(ca Caller) {
	p.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (p *ProvisionedHost) SetOnConflict(o OnConflict) {
	p.OnConflict = o
}

// Swap implements the Mergeable interface
func (p *ProvisionedHost) Swap(m Mergeable) error {
	rawVal, ok := m.(*ProvisionedHost)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", p, m)
	}
	*p = *rawVal
	return nil
}

// SetID increments the revision and sets the ID if needed
func (p *ProvisionedHost) SetID() string {
	p.Revision++
	if p.TeamID == "" && p.Team != nil {
		p.TeamID = p.Team.ID
	}
	if p.BuildID == "" && p.Build != nil {
		p.BuildID = p.Build.ID
	}
	if p.EnvironmentID == "" && p.Environment != nil {
		p.EnvironmentID = p.Environment.ID
	}
	if p.CompetitionID == "" && p.Competition != nil {
		p.CompetitionID = p.Competition.ID
	}
	if p.ID == "" && p.Host != nil {
		p.ID = p.Host.ID
	}
	if p.ID == "" {
		p.ID = p.HostID
	}
	return p.ID
}
