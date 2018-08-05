package laforge

import (
	"fmt"

	"github.com/pkg/errors"
)

// Environment represents the basic configurable type for a Laforge environment container
type Environment struct {
	ID               string              `hcl:",label" json:"id,omitempty"`
	BaseDir          string              `json:"base_dir,omitempty"`
	Name             string              `hcl:"name,attr" json:"name,omitempty"`
	Type             string              `hcl:"type,attr" json:"type,omitempty"`
	Config           map[string]string   `hcl:"config,attr" json:"config,omitempty"`
	Vars             map[string]string   `hcl:"vars,attr" json:"vars,omitempty"`
	Tags             map[string]string   `hcl:"tags,attr" json:"tags,omitempty"`
	Networks         []*IncludedNetwork  `hcl:"network,block" json:"included_networks,omitempty"`
	IncludedNetworks map[string]*Network `json:"-"`
	IncludedHosts    map[string]*Host    `json:"-"`
	HostByNetwork    map[string][]*Host  `json:"-"`
	OnConflict       OnConflict          `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller           Caller              `json:"-"`
}

// GetCaller implements the Mergeable interface
func (e *Environment) GetCaller() Caller {
	return e.Caller
}

// GetID implements the Mergeable interface
func (e *Environment) GetID() string {
	return e.ID
}

// GetOnConflict implements the Mergeable interface
func (e *Environment) GetOnConflict() OnConflict {
	return e.OnConflict
}

// SetCaller implements the Mergeable interface
func (e *Environment) SetCaller(c Caller) {
	e.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (e *Environment) SetOnConflict(o OnConflict) {
	e.OnConflict = o
}

// Swap implements the Mergeable interface
func (e *Environment) Swap(m Mergeable) error {
	rawVal, ok := m.(*Environment)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", e, m)
	}
	*e = *rawVal
	return nil
}

// ResolveIncludedNetworks walks the included_networks and included_hosts within the environment configuration
// ensuring that they can be located in the base laforge namespace.
func (e *Environment) ResolveIncludedNetworks(base *Laforge) error {
	e.IncludedNetworks = map[string]*Network{}
	e.HostByNetwork = map[string][]*Host{}
	e.IncludedHosts = map[string]*Host{}
	inet := map[string]string{}
	ihost := map[string]string{}
	for _, n := range e.Networks {
		inet[n.Name] = "included"
		e.HostByNetwork[n.Name] = []*Host{}
		for _, h := range n.Hosts {
			ihost[h] = "included"
		}
	}
	for name, net := range base.Networks {
		status, found := inet[name]
		if !found {
			Logger.Debugf("Skipping network %s", name)
			continue
		}
		if status == "included" {
			e.IncludedNetworks[name] = net
			inet[name] = "resolved"
			Logger.Infof("Resolved network %s", name)
		}
	}
	for name, host := range base.Hosts {
		status, found := ihost[name]
		if !found {
			Logger.Debugf("Skipping host %s", name)
			continue
		}
		if status == "included" {
			e.IncludedHosts[name] = host
			ihost[name] = "resolved"
			Logger.Infof("Resolved host %s", name)
		}
	}
	for _, n := range e.Networks {
		for _, h := range n.Hosts {
			host, found := e.IncludedHosts[h]
			if !found {
				return fmt.Errorf("unknown host included: %s", h)
			}
			err := host.Index(base)
			if err != nil {
				return err
			}
			e.HostByNetwork[n.Name] = append(e.HostByNetwork[n.Name], host)
		}
	}
	for net, status := range inet {
		if status == "included" {
			return fmt.Errorf("no configuration for network %s", net)
		}
	}
	for host, status := range ihost {
		if status == "included" {
			return fmt.Errorf("no configuration for host %s", host)
		}
	}
	return nil
}
