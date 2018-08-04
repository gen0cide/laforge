package laforge

import (
	"fmt"

	"github.com/imdario/mergo"
)

// Host defines a configurable type for customizing host parameters within the infrastructure.
type Host struct {
	Hostname         string                 `hcl:"hostname,label" cty:"hostname" json:"hostname,omitempty"`
	OS               string                 `hcl:"os,attr" json:"os,omitempty"`
	AMI              string                 `hcl:"ami,attr" json:"ami,omitempty"`
	LastOctet        int                    `hcl:"last_octet,attr" json:"last_octet,omitempty"`
	InstanceSize     string                 `hcl:"instance_size,attr" json:"instance_size,omitempty"`
	Disk             Disk                   `hcl:"disk,block" json:"disk,omitempty"`
	DNS              []DNSEntry             `hcl:"dns,block" json:"dns,omitempty"`
	ExposedTCPPorts  []string               `hcl:"exposed_tcp_ports,attr" json:"exposed_tcp_ports,omitempty"`
	ExposedUDPPorts  []string               `hcl:"exposed_udp_ports,attr" json:"exposed_udp_ports,omitempty"`
	Scripts          []Script               `hcl:"script,block" json:"scripts,omitempty"`
	Commands         []Command              `hcl:"command,block" json:"commands,omitempty"`
	ScriptMap        map[string]*Script     `json:"-"`
	CommandMap       map[string]*Command    `json:"-"`
	RemoteFileMap    map[string]*RemoteFile `json:"-"`
	OverridePassword string                 `hcl:"override_password,attr" json:"override_password,omitempty"`
	UserGroups       []string               `hcl:"user_groups,attr" json:"user_groups,omitempty"`
	Dependencies     []Dependency           `hcl:"dependency,block" json:"dependencies,omitempty"`
	RemoteFiles      []RemoteFile           `hcl:"remote_file,block" json:"remote_files,omitempty"`
	IO               IO                     `hcl:"io,block" json:"io,omitempty"`
	Vars             map[string]string      `hcl:"vars,attr" json:"vars,omitempty"`
	Tags             map[string]string      `hcl:"tags,attr" json:"tags,omitempty"`
	Maintainer       User                   `hcl:"maintainer,block" json:"maintainer,omitempty"`
	OnConflict       OnConflict             `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller           Caller                 `json:"-"`
}

// Disk is a configurable type for setting the root volume's disk size in GB
type Disk struct {
	Size int `hcl:"size,attr" json:"size,omitempty"`
}

// DNSEntry is a configurable type for defining DNS entries related to this host in the core DNS infrastructure (if enabled)
type DNSEntry struct {
	Label      string     `hcl:"id,label" json:"id,omitempty"`
	Name       string     `hcl:"name,attr" json:"name,omitempty"`
	Value      string     `hcl:"value,attr" json:"value,omitempty"`
	Type       string     `hcl:"type,attr" json:"type,omitempty"`
	Disabled   bool       `hcl:"disabled,attr" json:"disabled,omitempty"`
	OnConflict OnConflict `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
}

// Dependency is a configurable type for defining host or network dependencies to allow a dependency graph to be honored during deployment
type Dependency struct {
	Host       string     `hcl:"host,attr" json:"host,omitempty"`
	Network    string     `hcl:"network,attr" json:"network,omitempty"`
	OnConflict OnConflict `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
}

// Update performs a patching operation on source (h) with diff (diff), using the diff's merge conflict settings as appropriate.
func (h *Host) Update(diff *Host) error {
	switch diff.OnConflict.Do {
	case "":
		return mergo.Merge(h, diff, mergo.WithOverride)
	case "overwrite":
		conflict := h.OnConflict
		*h = *diff
		h.OnConflict = conflict
		return nil
	case "inherit":
		callerCopy := diff.Caller
		conflict := h.OnConflict
		err := mergo.Merge(diff, h, mergo.WithOverride)
		*h = *diff
		h.Caller = callerCopy
		h.OnConflict = conflict
		return err
	case "panic":
		return NewMergeConflict(h, diff, h.Hostname, diff.Hostname, h.Caller.Current(), diff.Caller.Current())
	default:
		return fmt.Errorf("invalid conflict strategy %s in %s", diff.OnConflict.Do, diff.Caller.Current().CallerFile)
	}
}

// ResolveScripts attempts to match script declarations referenced within this host object and script globals listed in the base config
func (h *Host) ResolveScripts(base *Laforge) error {
	return nil
}

// ResolveCommands attempts to match command declarations referenced within this host object and command globals listed in the base config
func (h *Host) ResolveCommands(base *Laforge) error {
	return nil
}

// ResolveRemoteFiles attempts to match remote file declarations referenced within this host object and remote file globals listed in the base config
func (h *Host) ResolveRemoteFiles(base *Laforge) error {
	return nil
}
