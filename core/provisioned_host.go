package core

import (
	"context"
	"fmt"
	"path"
	"path/filepath"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/pkg/errors"
)

const (
	// NullIP is represents an IP that is unknown to our systems
	NullIP string = `0.0.0.0`
)

// ProvisionedHost is a build artifact type to denote a host inside a team's provisioend infrastructure.
//easyjson:json
type ProvisionedHost struct {
	ID                 string                       `hcl:"id,label" json:"id,omitempty"`
	HostID             string                       `hcl:"host_id,attr" json:"host_id,omitempty"`
	SubnetIP           string                       `hcl:"subnet_ip,attr" json:"subnet_ip,omitempty"`
	Conn               *Connection                  `hcl:"connection,block" json:"connection"`
	Status             Status                       `hcl:"status,optional" json:"status"`
	ProvisioningSteps  map[string]*ProvisioningStep `json:"provisioning_steps"`
	StepsByOffset      []*ProvisioningStep          `json:"-"`
	ProvisionedNetwork *ProvisionedNetwork          `json:"-"`
	Team               *Team                        `json:"-"`
	Build              *Build                       `json:"-"`
	Environment        *Environment                 `json:"-"`
	Competition        *Competition                 `json:"-"`
	Network            *Network                     `json:"-"`
	Host               *Host                        `json:"-"`
	OnConflict         *OnConflict                  `json:"-"`
	Caller             Caller                       `json:"-"`
	Dir                string                       `json:"-"`
}

// Hash implements the Hasher interface
func (p *ProvisionedHost) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"hid=%v cidr=%v host=%v status=%v",
			p.HostID,
			p.SubnetIP,
			p.Host.Hash(),
			p.Status.Hash(),
		),
	)
}

// Path implements the Pather interface
func (p *ProvisionedHost) Path() string {
	return p.ID
}

// Base implements the Pather interface
func (p *ProvisionedHost) Base() string {
	return path.Base(p.ID)
}

// ValidatePath implements the Pather interface
func (p *ProvisionedHost) ValidatePath() error {
	if err := ValidateGenericPath(p.Path()); err != nil {
		return err
	}
	return nil
}

// GetCaller implements the Mergeable interface
func (p *ProvisionedHost) GetCaller() Caller {
	return p.Caller
}

// LaforgeID implements the Mergeable interface
// This will be: /envs/$env_base/$build_base/teams/$team_base/networks/$network_base/$host_base
func (p *ProvisionedHost) LaforgeID() string {
	return p.ID
}

// ParentLaforgeID returns the Team's parent build ID
func (p *ProvisionedHost) ParentLaforgeID() string {
	return path.Dir(path.Dir(p.LaforgeID()))
}

// GetOnConflict implements the Mergeable interface
func (p *ProvisionedHost) GetOnConflict() OnConflict {
	if p.OnConflict == nil {
		return OnConflict{
			Do:     "default",
			Append: true,
		}
	}
	return *p.OnConflict
}

// SetCaller implements the Mergeable interface
func (p *ProvisionedHost) SetCaller(ca Caller) {
	p.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (p *ProvisionedHost) SetOnConflict(o OnConflict) {
	p.OnConflict = &o
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

// SetID increments the revision and sets the team ID if needed
func (p *ProvisionedHost) SetID() string {
	if p.ID == "" {
		p.ID = path.Join(p.ProvisionedNetwork.Path(), hostsDir, p.Host.Base())
	}
	if p.HostID == "" {
		p.HostID = p.Host.Path()
	}
	return p.ID
}

// ActualPassword attempts to get everything just right interms of what the actual password of this machine is
func (p *ProvisionedHost) ActualPassword() string {
	pass := p.Competition.RootPassword
	if p.Host.OverridePassword != "" {
		pass = p.Host.OverridePassword
	}
	return pass
}

// CreateConnection creates this host's skeleton connection file to be used
func (p *ProvisionedHost) CreateConnection() *Connection {
	if p.Conn != nil {
		return p.Conn
	}
	c := &Connection{
		ID:                 path.Join(p.Path(), "conn"),
		Competition:        p.Competition,
		Active:             false,
		LocalAddr:          p.SubnetIP,
		RemoteAddr:         NullIP,
		ResourceName:       path.Join(p.Path(), "conn"),
		Environment:        p.Environment,
		Build:              p.Build,
		Team:               p.Team,
		Network:            p.Network,
		Host:               p.Host,
		ProvisionedHost:    p,
		ProvisionedNetwork: p.ProvisionedNetwork,
	}
	if p.Host.IsWindows() {
		c.WinRMAuthConfig = &WinRMAuthConfig{
			RemoteAddr: NullIP,
			Port:       5985,
			HTTPS:      false,
			SkipVerify: true,
			User:       "Administrator",
			Password:   p.ActualPassword(),
		}
	} else {
		keyfile := path.Join(p.Build.Path(), "data", "ssh.pem")
		relp, err := filepath.Rel(p.Team.Path(), keyfile)
		if err != nil {
			panic("error attempting to construct relative path")
		}
		c.SSHAuthConfig = &SSHAuthConfig{
			RemoteAddr:   NullIP,
			Port:         22,
			User:         "root",
			IdentityFile: relp,
			Password:     p.ActualPassword(),
		}
	}

	p.Conn = c
	return c
}

// CreateProvisioningStep creates a new provisioning step object for the provisioned host, mapping parent objects.
func (p *ProvisionedHost) CreateProvisioningStep(pr Provisioner, offset int) *ProvisioningStep {
	ps := &ProvisioningStep{
		Provisioner:        pr,
		ProvisionerID:      pr.Path(),
		ProvisionerType:    pr.Kind(),
		StepNumber:         offset,
		ProvisionedHost:    p,
		ProvisionedNetwork: p.ProvisionedNetwork,
		Network:            p.Network,
		Host:               p.Host,
		Team:               p.Team,
		Build:              p.Build,
		Environment:        p.Environment,
		Competition:        p.Competition,
	}

	p.ProvisioningSteps[ps.SetID()] = ps
	p.StepsByOffset = append(p.StepsByOffset, ps)
	return ps
}

// CreateProvisioningSteps enumerates all the parent Host object's provisioning steps, mapping a custom ProvisioningStep object for this provisioned host.
func (p *ProvisionedHost) CreateProvisioningSteps() error {
	for sid, step := range p.Host.Provisioners {
		p.CreateProvisioningStep(step, sid)
	}
	return nil
}

// Gather implements the Dependency interface
func (p *ProvisionedHost) Gather(g *Snapshot) error {
	// var err error
	// g.AddNode(p.Conn)
	// for _, s := range p.StepsByOffset {
	// 	g.AddNode(s)
	// 	err = s.Gather(g)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// for _, s := range p.StepsByOffset {
	// 	g.AddNode(s)
	// 	err = s.Gather(g)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// CreateProvisionedHostEntry ...
func (p *ProvisionedHost) CreateProvisionedHostEntry(ctx context.Context, pn *ent.ProvisionedNetwork, client *ent.Client) (*ent.ProvisionedHost, error) {
	status, err := p.Status.CreateStatusEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating provisioned host: %v", err)
		return nil, err
	}

	host, err := p.Host.CreateHostEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating provisioned host: %v", err)
		return nil, err
	}

	ph, err := client.ProvisionedHost.
		Create().
		SetSubnetIP(p.SubnetIP).
		AddStatus(status).
		AddProvisionedNetwork(pn).
		AddHost(host).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating provisioned host: %v", err)
		return nil, err
	}

	for _, v := range p.ProvisioningSteps {
		_, err := v.CreateProvisioningStepEntry(ctx, ph, client)

		if err != nil {
			cli.Logger.Debugf("failed creating provisioned host: %v", err)
			return nil, err
		}
	}

	cli.Logger.Debugf("provisioned host was created: ", ph)
	return ph, nil
}
