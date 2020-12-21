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

// ProvisioningStep is a build artifact type to denote a specific step inside of a provisioned host
//easyjson:json
type ProvisioningStep struct {
	ID                 string              `hcl:"id,label" json:"id,omitempty"`
	ProvisionerID      string              `hcl:"provisioner_id,attr" json:"provisioner_id,omitempty"`
	ProvisionerType    string              `hcl:"provisioner_type,attr" json:"provisioner_type,omitempty"`
	StepNumber         int                 `hcl:"step_number,attr" json:"step_number,omitempty"`
	Status             string              `hcl:"status,optional" json:"status,omitempty"`
	ProvisionedHost    *ProvisionedHost    `json:"-"`
	ProvisionedNetwork *ProvisionedNetwork `json:"-"`
	Host               *Host               `json:"-"`
	Network            *Network            `json:"-"`
	Team               *Team               `json:"-"`
	Build              *Build              `json:"-"`
	Environment        *Environment        `json:"-"`
	Competition        *Competition        `json:"-"`
	Provisioner        Provisioner         `json:"-"`
	Script             *Script             `json:"-"`
	Command            *Command            `json:"-"`
	RemoteFile         *RemoteFile         `json:"-"`
	DNSRecord          *DNSRecord          `json:"-"`
	OnConflict         *OnConflict         `json:"-"`
	Caller             Caller              `json:"-"`
	Dir                string              `json:"-"`
}

// Hash implements the Hasher interface
func (p *ProvisioningStep) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"pid=%v ptype=%v phash=%v snum=%v",
			p.ProvisionerID,
			p.ProvisionerType,
			p.Provisioner.Hash(),
			p.StepNumber,
		),
	)
}

// Path implements the Pather interface
func (p *ProvisioningStep) Path() string {
	return p.ID
}

// Base implements the Pather interface
func (p *ProvisioningStep) Base() string {
	return path.Base(p.ID)
}

// ValidatePath implements the Pather interface
func (p *ProvisioningStep) ValidatePath() error {
	if err := ValidateGenericPath(p.Path()); err != nil {
		return err
	}
	return nil
}

// GetCaller implements the Mergeable interface
func (p *ProvisioningStep) GetCaller() Caller {
	return p.Caller
}

// LaforgeID implements the Mergeable interface
func (p *ProvisioningStep) LaforgeID() string {
	return p.ID
}

// ParentLaforgeID returns the Team's parent build ID
func (p *ProvisioningStep) ParentLaforgeID() string {
	return path.Dir(path.Dir(p.LaforgeID()))
}

// GetOnConflict implements the Mergeable interface
func (p *ProvisioningStep) GetOnConflict() OnConflict {
	if p.OnConflict == nil {
		return OnConflict{
			Do:     "default",
			Append: true,
		}
	}
	return *p.OnConflict
}

// SetCaller implements the Mergeable interface
func (p *ProvisioningStep) SetCaller(ca Caller) {
	p.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (p *ProvisioningStep) SetOnConflict(o OnConflict) {
	p.OnConflict = &o
}

// Swap implements the Mergeable interface
func (p *ProvisioningStep) Swap(m Mergeable) error {
	rawVal, ok := m.(*ProvisioningStep)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", p, m)
	}
	*p = *rawVal
	return nil
}

// SetID increments the revision and sets the team ID if needed
func (p *ProvisioningStep) SetID() string {
	if p.ID == "" {
		p.ID = path.Join(p.ProvisionedHost.Path(), "steps", fmt.Sprintf("%d-%s", p.StepNumber, p.Provisioner.Base()))
	}

	switch v := p.Provisioner.(type) {
	case *Command:
		p.Command = v
	case *DNSRecord:
		p.DNSRecord = v
	case *RemoteFile:
		p.RemoteFile = v
	case *Script:
		p.Script = v
	}

	return p.ID
}

// Gather implements the Dependency interface
func (p *ProvisioningStep) Gather(g *Snapshot) error {
	// switch v := p.Provisioner.(type) {
	// case *Command:
	// 	// err := g.Relate(p.Environment, v)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	g.AddNode(v)
	// 	// err := g.Relate(p.Host, v)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	// err = g.Relate(v, p)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// case *DNSRecord:
	// 	// err := g.Relate(p.Environment, v)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	g.AddNode(v)
	// 	// err := g.Relate(p.Host, v)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	// err = g.Relate(v, p)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// case *RemoteFile:
	// 	// err := g.Relate(p.Environment, v)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	g.AddNode(v)
	// 	// err := g.Relate(p.Host, v)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	// err = g.Relate(v, p)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// case *Script:
	// 	// err := g.Relate(p.Environment, v)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	g.AddNode(v)
	// 	// err := g.Relate(p.Host, v)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	// err = g.Relate(v, p)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// default:
	// 	return fmt.Errorf("invalid provisioner type for %s: %T", p.Path(), p.Provisioner)
	// }
	return nil
}

// CreateProvisioningStepEntry ...
func (p *ProvisioningStep) CreateProvisioningStepEntry(ctx context.Context, ph *ent.ProvisionedHost, client *ent.Client) (*ent.ProvisioningStep, error) {
	if p.Script != nil {
		script, err := p.Script.CreateScriptEntry(ph, ctx, client)

		if err != nil {
			cli.Logger.Debugf("failed creating provisioning step: %v", err)
			return nil, err
		}

		ps, err := client.ProvisioningStep.
			Create().
			SetProvisionerType(p.ProvisionerType).
			SetStepNumber(p.StepNumber).
			SetStatus(p.Status).
			AddProvisionedHost(ph).
			AddScript(script).
			Save(ctx)

		if err != nil {
			cli.Logger.Debugf("failed creating provisioning step: %v", err)
			return nil, err
		}

		cli.Logger.Debugf("provisioning step was created: ", ps)
		return ps, nil
	} else if p.Command != nil {
		command, err := p.Command.CreateCommandEntry(ctx, client)

		if err != nil {
			cli.Logger.Debugf("failed creating provisioning step: %v", err)
			return nil, err
		}

		ps, err := client.ProvisioningStep.
			Create().
			SetProvisionerType(p.ProvisionerType).
			SetStepNumber(p.StepNumber).
			SetStatus(p.Status).
			AddProvisionedHost(ph).
			AddCommand(command).
			Save(ctx)

		if err != nil {
			cli.Logger.Debugf("failed creating provisioning step: %v", err)
			return nil, err
		}

		cli.Logger.Debugf("provisioning step was created: ", ps)
		return ps, nil
	} else if p.DNSRecord != nil {
		dnsrecord, err := p.DNSRecord.CreateDNSRecordEntry(ctx, client)

		if err != nil {
			cli.Logger.Debugf("failed creating provisioning step: %v", err)
			return nil, err
		}

		ps, err := client.ProvisioningStep.
			Create().
			SetProvisionerType(p.ProvisionerType).
			SetStepNumber(p.StepNumber).
			SetStatus(p.Status).
			AddProvisionedHost(ph).
			AddDNSRecord(dnsrecord).
			Save(ctx)

		if err != nil {
			cli.Logger.Debugf("failed creating provisioning step: %v", err)
			return nil, err
		}

		cli.Logger.Debugf("provisioning step was created: ", ps)
		return ps, nil
	} else if p.RemoteFile != nil {
		remotefile, err := p.RemoteFile.CreateRemoteFileEntry(ctx, client)

		if err != nil {
			cli.Logger.Debugf("failed creating provisioning step: %v", err)
			return nil, err
		}

		ps, err := client.ProvisioningStep.
			Create().
			SetProvisionerType(p.ProvisionerType).
			SetStepNumber(p.StepNumber).
			SetStatus(p.Status).
			AddProvisionedHost(ph).
			AddRemoteFile(remotefile).
			Save(ctx)

		if err != nil {
			cli.Logger.Debugf("failed creating provisioning step: %v", err)
			return nil, err
		}

		cli.Logger.Debugf("provisioning step was created: ", ps)
		return ps, nil
	}

	return nil, nil
}
