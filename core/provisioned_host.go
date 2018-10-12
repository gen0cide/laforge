package core

import "github.com/pkg/errors"

// ProvisionedHost defines a provisioned host within a team's environment (network neutral)
type ProvisionedHost struct {
	ID              string           `hcl:"id,label" json:"id,omitempty"`
	TeamID          string           `hcl:"team_id,attr" json:"team_id,omitempty"`
	HostID          string           `hcl:"host_id,attr" json:"host_id,omitempty"`
	Active          bool             `hcl:"active,attr" json:"active,omitempty"`
	LocalAddr       string           `hcl:"local_addr,attr" json:"local_addr,omitempty"`
	RemoteAddr      string           `hcl:"remote_addr,attr" json:"remote_addr,omitempty"`
	SSHAuthConfig   *SSHAuthConfig   `hcl:"ssh,block" json:"ssh_config,omitempty"`
	WinRMAuthConfig *WinRMAuthConfig `hcl:"winrm,block" json:"winrm_config,omitempty"`
	Team            *Team            `json:"-"`
	Host            *Host            `json:"-"`
	OnConflict      OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
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
	return p.ID
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
