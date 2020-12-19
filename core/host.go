package core

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"path"
	"sort"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/pkg/errors"
)

const (
	hostsDir = `hosts`
)

// Host defines a configurable type for customizing host parameters within the infrastructure.
//easyjson:json
type Host struct {
	ID               string                 `cty:"id" hcl:"id,label" json:"id,omitempty"`
	Hostname         string                 `cty:"hostname" hcl:"hostname,attr" json:"hostname,omitempty"`
	Description      string                 `cty:"description" hcl:"description,optional" json:"description,omitempty"`
	OS               string                 `cty:"os" hcl:"os,attr" json:"os,omitempty"`
	AMI              string                 `cty:"ami" hcl:"ami,optional" json:"ami,omitempty"`
	LastOctet        int                    `cty:"last_octet" hcl:"last_octet,attr" json:"last_octet,omitempty"`
	InstanceSize     string                 `cty:"instance_size" hcl:"instance_size,attr" json:"instance_size,omitempty"`
	AllowMACChanges  bool                   `cty:"allow_mac_changes" hcl:"allow_mac_changes,optional" json:"allow_mac_changes,omitempty"`
	Disk             Disk                   `cty:"disk" hcl:"disk,block" json:"disk,omitempty"`
	ProvisionSteps   []string               `cty:"provision_steps" hcl:"provision_steps,optional" json:"provision_steps,omitempty"`
	ExposedTCPPorts  []string               `cty:"exposed_tcp_ports" hcl:"exposed_tcp_ports,optional" json:"exposed_tcp_ports,omitempty"`
	ExposedUDPPorts  []string               `cty:"exposed_udp_ports" hcl:"exposed_udp_ports,optional" json:"exposed_udp_ports,omitempty"`
	OverridePassword string                 `cty:"override_password" hcl:"override_password,optional" json:"override_password,omitempty"`
	UserGroups       []string               `cty:"user_groups" hcl:"user_groups,optional" json:"user_groups,omitempty"`
	Dependencies     []*HostDependency      `cty:"depends_on" hcl:"depends_on,block" json:"depends_on,omitempty"`
	IO               *IO                    `cty:"io" hcl:"io,block" json:"io,omitempty"`
	Vars             map[string]string      `cty:"vars" hcl:"vars,optional" json:"vars,omitempty"`
	Tags             map[string]string      `cty:"tags" hcl:"tags,optional" json:"tags,omitempty"`
	Maintainer       *User                  `cty:"maintainer" hcl:"maintainer,block" json:"maintainer,omitempty"`
	OnConflict       *OnConflict            `cty:"on_conflict" hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Provisioners     []Provisioner          `json:"-"`
	Caller           Caller                 `json:"-"`
	Scripts          map[string]*Script     `json:"-"`
	Commands         map[string]*Command    `json:"-"`
	RemoteFiles      map[string]*RemoteFile `json:"-"`
	DNSRecords       map[string]*DNSRecord  `json:"-"`
}

// Disk is a configurable type for setting the root volume's disk size in GB
//easyjson:json
type Disk struct {
	Size int `cty:"size" hcl:"size,attr" json:"size,omitempty"`
}

// HostDependency is a configurable type for defining host or network dependencies to allow a dependency graph to be honored during deployment
//easyjson:json
type HostDependency struct {
	HostID     string      `cty:"host" hcl:"host,attr" json:"host,omitempty"`
	NetworkID  string      `cty:"network" hcl:"network,attr" json:"network,omitempty"`
	Step       string      `cty:"step" hcl:"step,optional" json:"step,omitempty"`
	StepID     int         `cty:"step_id" hcl:"step_id,optional" json:"step_id,omitempty"`
	OnConflict *OnConflict `cty:"on_conflict" hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Host       *Host       `json:"-"`
	Network    *Network    `json:"-"`
}

// Hash implements the Hasher interface
func (h *HostDependency) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"hid=%v nid=%v step=%v stepid=%v hh=%v nh=%v",
			h.HostID,
			h.NetworkID,
			h.Step,
			h.StepID,
			h.Host.Hash(),
			h.Network.Hash(),
		),
	)
}

// Hash implements the Hasher interface
func (h *Host) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"hn=%v os=%v ami=%v lo=%v isize=%v disk=%v ps=%v opass=%v ug=%v ph=%v vars=%v",
			h.Hostname,
			h.OS,
			h.AMI,
			h.LastOctet,
			h.InstanceSize,
			h.Disk,
			strings.Join(h.ProvisionSteps, `,`),
			h.OverridePassword,
			h.UserGroups,
			h.GetProvisionersHash(),
			HashConfigMap(h.Vars),
		),
	)
}

// DependencyCount is a helper function used to aggregate the number of dependencies a host has recursively
func (h *Host) DependencyCount(e *Environment) int {
	ret := 0
	if h.Dependencies == nil || len(h.Dependencies) == 0 {
		return ret
	}
	for _, x := range h.Dependencies {
		found, ok := e.IncludedHosts[x.HostID]
		if !ok {
			continue
		}
		ret += found.DependencyCount(e)
	}
	return ret
}

// GetDependencyHash returns the host's dependency hash
func (h *Host) GetDependencyHash() string {
	p := []string{}
	for _, x := range h.Dependencies {
		p = append(p, fmt.Sprintf("%d", x.Hash()))
	}
	sort.Strings(p)
	return strings.Join(p, ",")
}

// GetProvisionersHash returns a concatinated string of the host's provisioners hashes
func (h *Host) GetProvisionersHash() uint64 {
	p := ChecksumList{}
	for _, x := range h.Scripts {
		p = append(p, x.Hash())
	}
	for _, x := range h.Commands {
		p = append(p, x.Hash())
	}
	for _, x := range h.DNSRecords {
		p = append(p, x.Hash())
	}
	for _, x := range h.RemoteFiles {
		p = append(p, x.Hash())
	}
	return p.Hash()
}

// Path implements the Pather interface
func (h *Host) Path() string {
	return h.ID
}

// Base implements the Pather interface
func (h *Host) Base() string {
	return path.Base(h.ID)
}

// ValidatePath implements the Pather interface
func (h *Host) ValidatePath() error {
	if err := ValidateGenericPath(h.Path()); err != nil {
		return err
	}
	if topdir := strings.Split(h.Path(), `/`); topdir[1] != hostsDir {
		return fmt.Errorf("path %s is not rooted in /%s", h.Path(), topdir[1])
	}
	return nil
}

// HasTag is a template helper function to return true/false if the host contains a tag of a specific key
func (h *Host) HasTag(tag string) bool {
	_, t := h.Tags[tag]
	return t
}

// TagEquals is a template helper function to return true/false if the host contains a tag key of a specific value
func (h *Host) TagEquals(tag, value string) bool {
	v, t := h.Tags[tag]
	if !t {
		return false
	}
	if v == value {
		return true
	}
	return false
}

// FinalStepID gets the final step ID for a host's offset
func (h *Host) FinalStepID() int {
	if len(h.Provisioners) == 0 {
		return -1
	}
	return len(h.Provisioners) - 1
}

// GetCaller implements the Mergeable interface
func (h *Host) GetCaller() Caller {
	return h.Caller
}

// LaforgeID implements the Mergeable interface
func (h *Host) LaforgeID() string {
	return h.ID
}

// ParentLaforgeID implements the Dependency interface
func (h *Host) ParentLaforgeID() string {
	return h.Path()
}

// Gather implements the Dependency interface
func (h *Host) Gather(g *Snapshot) error {
	return nil
}

// GetOnConflict implements the Mergeable interface
func (h *Host) GetOnConflict() OnConflict {
	if h.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *h.OnConflict
}

// SetCaller implements the Mergeable interface
func (h *Host) SetCaller(c Caller) {
	h.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (h *Host) SetOnConflict(o OnConflict) {
	h.OnConflict = &o
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

// CalcIP is used to calculate the IP of a host within a given subnet
func (h *Host) CalcIP(subnet string) string {
	ip, _, err := net.ParseCIDR(subnet)
	if err != nil {
		return fmt.Sprintf("ERR_INVALID_SUBNET_%s_FOR_HOST_%s", subnet, h.ID)
	}
	offset32 := uint32(h.LastOctet)
	ip32 := IPv42Int(ip)
	newIP := Int2IPv4(ip32 + offset32)
	return newIP.To4().String()
}

// IsWindows is a template helper function to determine if the underlying operating system is windows
func (h *Host) IsWindows() bool {
	switch strings.ToLower(h.OS) {
	case "w2k12":
		return true
	case "w2k16":
		return true
	case "w2k19":
		return true
	case "w2k12-sql":
		return true
	case "w2k16-sql":
		return true
	case "w2k19-sql":
		return true
	//nolint:goconst
	case "windows":
		return true
	default:
		return false
	}
}

// Index attempts to index all children dependencies of this type
func (h *Host) Index(base *Laforge) error {
	h.Scripts = map[string]*Script{}
	h.Commands = map[string]*Command{}
	h.RemoteFiles = map[string]*RemoteFile{}
	h.DNSRecords = map[string]*DNSRecord{}
	iprov := map[string]string{}
	h.Provisioners = []Provisioner{}

	for _, s := range h.ProvisionSteps {
		cli.Logger.Debugf("indexing provision step %s for host %s", s, h.ID)
		iprov[s] = ObjectTypeIncluded.String()
	}
	for name, script := range base.Scripts {
		status, found := iprov[name]
		if !found {
			continue
		}
		if status == ObjectTypeIncluded.String() {
			h.Scripts[name] = script
			iprov[name] = ObjectTypeScript.String()
			cli.Logger.Debugf("Resolved %T dependency %s for %s", script, script.ID, h.ID)
		}
	}
	for name, command := range base.Commands {
		status, found := iprov[name]
		if !found {
			continue
		}
		if status == ObjectTypeIncluded.String() {
			h.Commands[name] = command
			iprov[name] = ObjectTypeCommand.String()
			cli.Logger.Debugf("Resolved %T dependency %s for %s", command, command.ID, h.ID)
		}
	}
	for name, file := range base.RemoteFiles {
		status, found := iprov[name]
		if !found {
			continue
		}
		if status == ObjectTypeIncluded.String() {
			h.RemoteFiles[name] = file
			iprov[name] = ObjectTypeRemoteFile.String()
			cli.Logger.Debugf("Resolved %T dependency %s for %s", file, file.ID, h.ID)
		}
	}
	for name, record := range base.DNSRecords {
		status, found := iprov[name]
		if !found {
			continue
		}
		if status == ObjectTypeIncluded.String() {
			h.DNSRecords[name] = record
			iprov[name] = ObjectTypeDNSRecord.String()
			cli.Logger.Debugf("Resolved %T dependency %s for %s", record, record.ID, h.ID)
		}
	}
	for x, status := range iprov {
		if status == ObjectTypeIncluded.String() {
			return fmt.Errorf("unmet provision_step dependency %s for host %s\n%s", x, h.ID, h.Caller.Error())
		}
	}
	for _, s := range h.ProvisionSteps {
		switch iprov[s] {
		case ObjectTypeScript.String():
			h.Provisioners = append(h.Provisioners, h.Scripts[s])
		case ObjectTypeCommand.String():
			h.Provisioners = append(h.Provisioners, h.Commands[s])
		case ObjectTypeRemoteFile.String():
			h.Provisioners = append(h.Provisioners, h.RemoteFiles[s])
		case ObjectTypeDNSRecord.String():
			h.Provisioners = append(h.Provisioners, h.DNSRecords[s])
		default:
			return fmt.Errorf("unmet provision_step dependency %s for host %s\n%s", s, h.ID, h.Caller.Error())
		}
	}
	return nil
}

// CreateDiskEntry ...
func (d *Disk) CreateDiskEntry(ctx context.Context, client *ent.Client) (*ent.Disk, error) {
	disk, err := client.Disk.
		Create().
		SetSize(d.Size).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating disk: %v", err)
		return nil, err
	}

	cli.Logger.Debugf("disk was created: ", disk)
	return disk, nil
}

// CreateHostEntry ...
func (h *Host) CreateHostEntry(ctx context.Context, client *ent.Client) (*ent.Host, error) {
	disk, err := h.Disk.CreateDiskEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating host: %v", err)
		return nil, err
	}

	user, err := h.Maintainer.CreateUserEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating host: %v", err)
		return nil, err
	}

	tag, err := CreateTagEntry(h.ID, h.Tags, ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating host: %v", err)
		return nil, err
	}

	host, err := client.Host.
		Create().
		SetHostname(h.Hostname).
		SetDescription(h.Description).
		SetOS(h.OS).
		SetLastOctet(h.LastOctet).
		SetAllowMACChanges(h.AllowMACChanges).
		SetExposedTCPPorts(h.ExposedTCPPorts).
		SetExposedUDPPorts(h.ExposedUDPPorts).
		SetOverridePassword(h.OverridePassword).
		SetVars(h.Vars).
		SetUserGroups(h.UserGroups).
		// SetDependsOn(h.DependsOn). // Not in the Laforge Host Struct
		// SetScripts(h.Scripts).
		// SetCommands(h.Commands).
		// SetRemoteFiles(h.RemoteFiles).
		// SetDnsRecords(h.DNSRecords).
		AddDisk(disk).
		AddMaintainer(user).
		AddTag(tag).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating host: %v", err)
		return nil, err
	}

	cli.Logger.Debugf("host was created: ", host)
	return host, nil
}

// Fullpath implements the Pather interface
func (h *Host) Fullpath() string {
	return h.LaforgeID()
}

// IPv42Int converts net.IP address objects to their uint32 representation
func IPv42Int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

// Int2IPv4 converts uint32s to their net.IP object
func Int2IPv4(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
