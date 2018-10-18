package core

import (
	"fmt"
	"path"

	"github.com/cespare/xxhash"
	"github.com/pkg/errors"
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
		p.ID = path.Join(p.ProvisionedNetwork.Path(), "hosts", p.Host.Base())
	}
	if p.HostID == "" {
		p.HostID = p.Host.Path()
	}
	return p.ID
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
	var err error
	for _, s := range p.StepsByOffset {
		err = g.Relate(p, s)
		if err != nil {
			return err
		}
		if s.StepNumber != 0 {
			previousSteps := p.StepsByOffset[0:s.StepNumber]
			for _, x := range previousSteps {
				err = g.Relate(x, s)
				if err != nil {
					return err
				}
			}
		}
		err = s.Gather(g)
		if err != nil {
			return err
		}
	}
	for _, s := range p.Host.Dependencies {
		hd, err := p.ProvisionedNetwork.Team.LocateProvisionedHost(s.NetworkID, s.HostID)
		if err != nil {
			return err
		}
		finalStepOffset := hd.Host.FinalStepID()
		if finalStepOffset != -1 {
			located := false
			for _, pstep := range hd.ProvisioningSteps {
				if pstep.StepNumber == finalStepOffset {
					err = g.Relate(pstep, p)
					if err != nil {
						return err
					}
					located = true
					break
				}
			}
			if !located {
				return fmt.Errorf("there is no provisioning step with offset %d for host %s", finalStepOffset, hd.Path())
			}
		} else {
			err = g.Relate(hd, p)
			if err != nil {
				return err
			}
		}
	}
	err = g.Relate(p.Host, p)
	if err != nil {
		return err
	}
	return nil
}
