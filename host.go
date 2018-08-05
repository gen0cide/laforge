package laforge

import (
	"fmt"

	"github.com/pkg/errors"
)

// Host defines a configurable type for customizing host parameters within the infrastructure.
type Host struct {
	ID                   string                 `hcl:",label" json:"id,omitempty"`
	Hostname             string                 `hcl:"hostname,attr" cty:"hostname" json:"hostname,omitempty"`
	OS                   string                 `hcl:"os,attr" json:"os,omitempty"`
	AMI                  string                 `hcl:"ami,attr" json:"ami,omitempty"`
	LastOctet            int                    `hcl:"last_octet,attr" json:"last_octet,omitempty"`
	InstanceSize         string                 `hcl:"instance_size,attr" json:"instance_size,omitempty"`
	Disk                 Disk                   `hcl:"disk,block" json:"disk,omitempty"`
	IncludedDNSRecords   []string               `hcl:"dns_records,attr" json:"included_dns_records,omitempty"`
	ExposedTCPPorts      []string               `hcl:"exposed_tcp_ports,attr" json:"exposed_tcp_ports,omitempty"`
	ExposedUDPPorts      []string               `hcl:"exposed_udp_ports,attr" json:"exposed_udp_ports,omitempty"`
	IncludedScripts      []string               `hcl:"scripts,attr" json:"included_scripts,omitempty"`
	IncludedCommands     []string               `hcl:"commands,attr" json:"included_commands,omitempty"`
	IncludedDependencies []string               `hcl:"dependencies,attr" json:"included_dependencies,omitempty"`
	IncludedFiles        []string               `hcl:"files,attr" json:"included_files,omitempty"`
	Scripts              map[string]*Script     `json:"scripts,omitempty"`
	Commands             map[string]*Command    `json:"commands,omitempty"`
	Files                map[string]*RemoteFile `json:"files,omitempty"`
	DNSRecords           map[string]*DNSRecord  `json:"dns_records,omitempty"`
	Dependencies         map[string]interface{} `json:"dependencies,omitempty"`
	OverridePassword     string                 `hcl:"override_password,attr" json:"override_password,omitempty"`
	UserGroups           []string               `hcl:"user_groups,attr" json:"user_groups,omitempty"`
	IO                   IO                     `hcl:"io,block" json:"io,omitempty"`
	Vars                 map[string]string      `hcl:"vars,attr" json:"vars,omitempty"`
	Tags                 map[string]string      `hcl:"tags,attr" json:"tags,omitempty"`
	Maintainer           User                   `hcl:"maintainer,block" json:"maintainer,omitempty"`
	OnConflict           OnConflict             `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller               Caller                 `json:"-"`
}

// Disk is a configurable type for setting the root volume's disk size in GB
type Disk struct {
	Size int `hcl:"size,attr" json:"size,omitempty"`
}

// Dependency is a configurable type for defining host or network dependencies to allow a dependency graph to be honored during deployment
type Dependency struct {
	Host       string     `hcl:"host,attr" json:"host,omitempty"`
	Network    string     `hcl:"network,attr" json:"network,omitempty"`
	OnConflict OnConflict `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
}

// GetCaller implements the Mergeable interface
func (h *Host) GetCaller() Caller {
	return h.Caller
}

// GetID implements the Mergeable interface
func (h *Host) GetID() string {
	return h.ID
}

// GetOnConflict implements the Mergeable interface
func (h *Host) GetOnConflict() OnConflict {
	return h.OnConflict
}

// SetCaller implements the Mergeable interface
func (h *Host) SetCaller(c Caller) {
	h.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (h *Host) SetOnConflict(o OnConflict) {
	h.OnConflict = o
}

// Swap implements the Mergeable interface
func (h *Host) Swap(m Mergeable) error {
	rawVal, ok := m.(*Host)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", h, m)
	}
	*h = *rawVal
	return nil
}

// Index attempts to index all children dependencies of this type
func (h *Host) Index(base *Laforge) error {
	h.Scripts = map[string]*Script{}
	h.Commands = map[string]*Command{}
	h.Files = map[string]*RemoteFile{}
	h.DNSRecords = map[string]*DNSRecord{}
	iscripts := map[string]string{}
	icommands := map[string]string{}
	ifiles := map[string]string{}
	irecords := map[string]string{}
	for _, s := range h.IncludedScripts {
		Logger.Debugf("indexing script %s for host %s", s, h.ID)
		iscripts[s] = "included"
	}
	for _, c := range h.IncludedCommands {
		icommands[c] = "included"
	}
	for _, c := range h.IncludedFiles {
		ifiles[c] = "included"
	}
	for _, c := range h.IncludedDNSRecords {
		irecords[c] = "included"
	}
	for name, script := range base.Scripts {
		status, found := iscripts[name]
		if !found {
			continue
		}
		if status == "included" {
			h.Scripts[name] = script
			iscripts[name] = "resolved"
			Logger.Debugf("Resolved %T dependency %s for %s", script, script.ID, h.ID)
		}
	}
	for name, command := range base.Commands {
		status, found := icommands[name]
		if !found {
			continue
		}
		if status == "included" {
			h.Commands[name] = command
			icommands[name] = "resolved"
			Logger.Debugf("Resolved %T dependency %s for %s", command, command.ID, h.ID)
		}
	}
	for name, file := range base.Files {
		status, found := ifiles[name]
		if !found {
			continue
		}
		if status == "included" {
			h.Files[name] = file
			ifiles[name] = "resolved"
			Logger.Debugf("Resolved %T dependency %s for %s", file, file.ID, h.ID)
		}
	}
	for name, record := range base.DNSRecords {
		status, found := irecords[name]
		if !found {
			continue
		}
		if status == "included" {
			h.DNSRecords[name] = record
			irecords[name] = "resolved"
			Logger.Debugf("Resolved %T dependency %s for %s", record, record.ID, h.ID)
		}
	}
	for x, status := range iscripts {
		if status == "included" {
			return fmt.Errorf("unmet %s dependency %s for host %s\n%s", "script", x, h.ID, h.Caller.Error())
		}
	}
	for x, status := range icommands {
		if status == "included" {
			return fmt.Errorf("unmet %s dependency %s for host %s\n%s", "command", x, h.ID, h.Caller.Error())
		}
	}
	for x, status := range ifiles {
		if status == "included" {
			return fmt.Errorf("unmet %s dependency %s for host %s\n%s", "remote_file", x, h.ID, h.Caller.Error())
		}
	}
	for x, status := range irecords {
		if status == "included" {
			return fmt.Errorf("unmet %s dependency %s for host %s\n%s", "dns_record", x, h.ID, h.Caller.Error())
		}
	}
	return nil
}
