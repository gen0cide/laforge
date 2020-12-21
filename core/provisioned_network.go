package core

import (
	"context"
	"fmt"
	"path"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/pkg/errors"
)

// ProvisionedNetwork is a build artifact type to denote a network inside a team's provisioend infrastructure.
//easyjson:json
type ProvisionedNetwork struct {
	ID               string                      `hcl:"id,label" json:"id,omitempty"`
	Name             string                      `hcl:"name,attr" json:"name,omitempty"`
	CIDR             string                      `hcl:"cidr,attr" json:"cidr,omitempty"`
	NetworkID        string                      `hcl:"network_id,attr" json:"network_id,omitempty"`
	ProvisionedHosts map[string]*ProvisionedHost `json:"provisioned_hosts"`
	Status           Status                      `hcl:"status,optional" json:"status"`
	Network          *Network                    `json:"-"`
	Team             *Team                       `json:"-"`
	Build            *Build                      `json:"-"`
	Environment      *Environment                `json:"-"`
	Competition      *Competition                `json:"-"`
	OnConflict       *OnConflict                 `json:"-"`
	Caller           Caller                      `json:"-"`
	Dir              string                      `json:"-"`
}

// Hash implements the Hasher interface
func (p *ProvisionedNetwork) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"name=%v cidr=%v net=%v team=%v status=%v",
			p.Name,
			p.CIDR,
			p.Network.Hash(),
			p.Team.Hash(),
			p.Status.Hash(),
		),
	)
}

// Path implements the Pather interface
func (p *ProvisionedNetwork) Path() string {
	return p.ID
}

// Base implements the Pather interface
func (p *ProvisionedNetwork) Base() string {
	return path.Base(p.ID)
}

// ValidatePath implements the Pather interface
func (p *ProvisionedNetwork) ValidatePath() error {
	if err := ValidateGenericPath(p.Path()); err != nil {
		return err
	}
	return nil
}

// GetCaller implements the Mergeable interface
func (p *ProvisionedNetwork) GetCaller() Caller {
	return p.Caller
}

// LaforgeID implements the Mergeable interface
func (p *ProvisionedNetwork) LaforgeID() string {
	return p.ID
}

// ParentLaforgeID returns the Team's parent build ID
func (p *ProvisionedNetwork) ParentLaforgeID() string {
	return path.Dir(path.Dir(p.LaforgeID()))
}

// GetOnConflict implements the Mergeable interface
func (p *ProvisionedNetwork) GetOnConflict() OnConflict {
	if p.OnConflict == nil {
		return OnConflict{
			Do:     "default",
			Append: true,
		}
	}
	return *p.OnConflict
}

// SetCaller implements the Mergeable interface
func (p *ProvisionedNetwork) SetCaller(ca Caller) {
	p.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (p *ProvisionedNetwork) SetOnConflict(o OnConflict) {
	p.OnConflict = &o
}

// Swap implements the Mergeable interface
func (p *ProvisionedNetwork) Swap(m Mergeable) error {
	rawVal, ok := m.(*ProvisionedNetwork)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", p, m)
	}
	*p = *rawVal
	return nil
}

// SetID increments the revision and sets the team ID if needed
func (p *ProvisionedNetwork) SetID() string {
	if p.ID == "" {
		p.ID = path.Join(p.Team.Path(), "networks", p.Network.Base())
	}
	if p.NetworkID == "" {
		p.NetworkID = p.Network.Path()
	}
	return p.ID
}

// CreateProvisionedHost creates the actual provisioned host object and assigns the parental objects accordingly.
func (p *ProvisionedNetwork) CreateProvisionedHost(host *Host) *ProvisionedHost {
	ph := &ProvisionedHost{
		Host:               host,
		SubnetIP:           host.CalcIP(p.CIDR),
		ProvisioningSteps:  map[string]*ProvisioningStep{},
		StepsByOffset:      []*ProvisioningStep{},
		ProvisionedNetwork: p,
		Team:               p.Team,
		Build:              p.Build,
		Environment:        p.Environment,
		Competition:        p.Competition,
	}
	p.ProvisionedHosts[ph.SetID()] = ph
	ph.Conn = ph.CreateConnection()
	ph.Conn.SetID()
	return ph
}

// CreateProvisionedHosts enumerates the parent environment's host by network and creates provisioned host objects in this tree.
func (p *ProvisionedNetwork) CreateProvisionedHosts() error {
	for _, h := range p.Team.Environment.HostByNetwork[p.Network.Path()] {
		ph := p.CreateProvisionedHost(h)
		err := ph.CreateProvisioningSteps()
		if err != nil {
			return err
		}
	}
	return nil
}

// Gather implements the Dependency interface
func (p *ProvisionedNetwork) Gather(g *Snapshot) error {
	// var err error
	// for _, h := range p.ProvisionedHosts {

	// 	// err = g.Relate(p, h)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	g.AddNode(h)
	// 	err = h.Gather(g)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// err = g.Relate(p, p.Network)
	// if err != nil {
	// 	return err
	// }
	// err = g.Relate(p.Network, p)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// CreateProvisionedNetworkEntry ...
func (p *ProvisionedNetwork) CreateProvisionedNetworkEntry(ctx context.Context,build *ent.Build, team *ent.Team, client *ent.Client) (*ent.ProvisionedNetwork, error) {
	status, err := p.Status.CreateStatusEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating provisioned network: %v", err)
		return nil, err
	}

	network, err := p.Network.CreateNetworkEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating provisioned network: %v", err)
		return nil, err
	}

	pn, err := client.ProvisionedNetwork.
		Create().
		SetName(p.Name).
		SetCidr(p.CIDR).
		AddStatus(status).
		AddNetwork(network).
		AddBuild(build).
		AddProvisionedNetworkToTeam(team).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating provisioned network: %v", err)
		return nil, err
	}

	for _, v := range p.ProvisionedHosts {
		_, err := v.CreateProvisionedHostEntry(ctx, pn, client)

		if err != nil {
			cli.Logger.Debugf("failed creating provisioned network: %v", err)
			return nil, err
		}
	}

	cli.Logger.Debugf("provisioned network was created: ", pn)
	return pn, nil
}
