// Package core contains the entire implementation of the Laforge configuration language. It includes
// it's own loader and dependency resolution mechanisms and is meant to be the source of truth around
// declaration logic.
package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/xlab/treeprint"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

const (
	// TeamContext is a context level representing a full team.laforge, build.laforge, env.laforge, base.laforge, and global.laforge=
	TeamContext StateContext = iota

	// BuildContext is a context level representing being within an environment's build directory
	BuildContext

	// EnvContext is a context level representing a full env.laforge + base.laforge was found
	EnvContext

	// BaseContext is a context level representing just a base.laforge was located
	BaseContext

	// GlobalContext is a context level representing a valid global configuration was found
	GlobalContext

	// NoContext is a context level representing neither a env.laforge or base.laforge was found
	NoContext
)

var (
	// ErrContextViolation is an error thrown when an action is taken that attempts to exceed the current state's ScopeContext
	ErrContextViolation = errors.New("context scope violation")

	// ErrAbsPathDeclNotExist is thrown when an absolute path is listed in a file yet it cannot be resolved on the local system
	ErrAbsPathDeclNotExist = errors.New("absolute path was not found on local system")

	// ErrInvalidEnvName is thrown when an environment name does not meet specified regulations around environment naming conventions
	ErrInvalidEnvName = errors.New("environment names can only contain lowercase alphanumeric and dash characters (a valid subdomain)")

	// ErrDeprecatedID is thrown when an object using a deprecated ID format is loaded
	ErrDeprecatedID = errors.New("the object is using a deprecated ID format and should be updated")

	// ErrInvalidIDFormat is thrown when an object using an invalid ID format is loaded
	ErrInvalidIDFormat = errors.New("object ID does not meet valid ID requirements")

	// ValidIDRegexp is the format for the new object ID schema.
	ValidIDRegexp = regexp.MustCompile(`\A[a-z0-9][a-z0-9\-]{1,34}[a-z0-9]\z`)

	// ValidOldIDRegexp is the old ID schema that allowed additional characters (periods and underscores)
	// and did not validate minimum or maximum length. These formats often incurred incompatabilities with
	// existing systems and will be deprecated.
	ValidOldIDRegexp = regexp.MustCompile(`\A[a-zA-Z0-9][a-zA-Z0-9\-\_\.]{2,34}[a-zA-Z0-9]\z`)
)

// Laforge defines the type that holds the global namespace within the laforge configuration engine
//easyjson:json
type Laforge struct {
	Filename                   string                         `json:"filename"`
	Includes                   []string                       `json:"-"`
	DependencyGraph            treeprint.Tree                 `json:"-"`
	CurrDir                    string                         `json:"-"`
	BaseDir                    string                         `hcl:"base_dir,optional" json:"base_dir,omitempty"`
	User                       *User                          `hcl:"user,block" json:"user,omitempty"`
	IncludePaths               []*Include                     `hcl:"include,block" json:"include_paths,omitempty"`
	DefinedCompetitions        []*Competition                 `hcl:"competition,block" json:"competitions,omitempty"`
	DefinedHosts               []*Host                        `hcl:"host,block" json:"hosts,omitempty"`
	DefinedNetworks            []*Network                     `hcl:"network,block" json:"networks,omitempty"`
	DefinedIdentities          []*Identity                    `hcl:"identity,block" json:"identities,omitempty"`
	DefinedScripts             []*Script                      `hcl:"script,block" json:"scripts,omitempty"`
	DefinedCommands            []*Command                     `hcl:"command,block" json:"defined_commands,omitempty"`
	DefinedRemoteFiles         []*RemoteFile                  `hcl:"remote_file,block" json:"defined_files,omitempty"`
	DefinedDNSRecords          []*DNSRecord                   `hcl:"dns_record,block" json:"defined_dns_records,omitempty"`
	DefinedEnvironments        []*Environment                 `hcl:"environment,block" json:"environments,omitempty"`
	DefinedBuilds              []*Build                       `hcl:"build,block" json:"builds,omitempty"`
	DefinedTeams               []*Team                        `hcl:"team,block" json:"teams,omitempty"`
	DefinedProvisionedNetworks []*ProvisionedNetwork          `hcl:"provisioned_network,block" json:"provisioned_networks,omitempty"`
	DefinedProvisionedHosts    []*ProvisionedHost             `hcl:"provisioned_host,block" json:"provisioned_hosts,omitempty"`
	DefinedProvisioningSteps   []*ProvisioningStep            `hcl:"provisioning_step,block" json:"provisioning_steps,omitempty"`
	DefinedConnections         []*Connection                  `hcl:"connection,block" json:"connections,omitempty"`
	Hosts                      map[string]*Host               `json:"-"`
	Networks                   map[string]*Network            `json:"-"`
	Identities                 map[string]*Identity           `json:"-"`
	Scripts                    map[string]*Script             `json:"-"`
	Commands                   map[string]*Command            `json:"-"`
	RemoteFiles                map[string]*RemoteFile         `json:"-"`
	DNSRecords                 map[string]*DNSRecord          `json:"-"`
	Competitions               map[string]*Competition        `json:"-"`
	Environments               map[string]*Environment        `json:"-"`
	Builds                     map[string]*Build              `json:"-"`
	Teams                      map[string]*Team               `json:"-"`
	ProvisionedNetworks        map[string]*ProvisionedNetwork `json:"-"`
	ProvisionedHosts           map[string]*ProvisionedHost    `json:"-"`
	ProvisioningSteps          map[string]*ProvisioningStep   `json:"-"`
	Connections                map[string]*Connection         `json:"-"`
	Caller                     Caller                         `json:"-"`
	ValidTeam                  bool                           `json:"-"`
	ValidBuild                 bool                           `json:"-"`
	ValidEnv                   bool                           `json:"-"`
	ValidBase                  bool                           `json:"-"`
	ValidGlobal                bool                           `json:"-"`
	ClearToBuild               bool                           `json:"-"`
	TeamRoot                   string                         `json:"-"`
	BuildRoot                  string                         `json:"-"`
	EnvRoot                    string                         `json:"-"`
	BaseRoot                   string                         `json:"-"`
	GlobalRoot                 string                         `json:"-"`
	TeamAbsPath                string                         `json:"-"`
	BuildAbsPath               string                         `json:"-"`
	EnvAbsPath                 string                         `json:"-"`
	BaseAbsPath                string                         `json:"-"`
	TeamContextID              string                         `json:"-"`
	BuildContextID             string                         `json:"-"`
	EnvContextID               string                         `json:"-"`
	BaseContextID              string                         `json:"-"`
	GlobalContextID            string                         `json:"-"`
	CurrentEnv                 *Environment                   `json:"-"`
	CurrentBuild               *Build                         `json:"-"`
	CurrentTeam                *Team                          `json:"-"`
	CurrentCompetition         *Competition                   `json:"-"`
	StateManager               *State                         `json:"-"`
	InitialContext             StateContext                   `json:"-"`
	PathRegistry               *PathRegistry                  `json:"-"`
}

// Include defines a named include type
//easyjson:json
type Include struct {
	Path string `hcl:"path,attr" json:"path,omitempty"`
}

// Opt defines a basic HCLv2 option label:
//	config "keyName" {
//		value = "valueData"
//	}
//easyjson:json
type Opt struct {
	Key   string `hcl:",label" json:"key,omitempty"`
	Value string `hcl:"value,attr" json:"value,omitempty"`
}

// OnConflict defines a configuration override for how to handle conflicting objects
// in the differential enumeration of objects. For example, if you have a file that has loaded
// a host definition (host "foo" {...}) and you have included an additional configuration
// in your env.laforge, you can specify an on_conflict { do = "" } setting within
// your block to control how the state of host "foo" has state changes applied to it.
//	host "foo" {
//		on_conflict {
//			do = "overwrite"
//		}
//	}
// The default behavior (no on_conflict block) is to perform an overriding MERGE on the original object
// with your changes. Think of this as updating only fields you have defined and are non-zero.
// If you wish to specify a different strategy, the following ones are valid:
//	- "overwrite" will replace the entirety of original "foo" with your definition, discarding any previous state.
//	- "inherit" will apply a merge in reverse - merging the original "foo" into your definition, overwriting any fields.
//	- "panic" will raise a runtime error and prevent further execution. This can be a very helpful way to avoid state on "root" definitions.
//easyjson:json
type OnConflict struct {
	Do     string `cty:"do" hcl:"do,attr" json:"do,omitempty"`
	Append bool   `cty:"append" hcl:"append,optional" json:"append,omitempty"`
}

// CurrentStateManager attempts to return the current state manager, throwing an error if it doesn't exist
func (l *Laforge) CurrentStateManager() (*State, error) {
	if l.StateManager == nil {
		return nil, errors.New("cannot access a nil state manager")
	}

	return l.StateManager, nil
}

// StateContext is a type alias to the level of context we are currently executing in
type StateContext int

func (s StateContext) String() string {
	switch s {
	case TeamContext:
		return "TeamContext"
	case BuildContext:
		return "BuildContext"
	case EnvContext:
		return "EnvContext"
	case BaseContext:
		return "BaseContext"
	case GlobalContext:
		return "GlobalContext"
	default:
		return "NoContext"
	}
}

// AssertMinContext allows a program to assert a required context and fail gracefully
func (l *Laforge) AssertMinContext(s StateContext) error {
	if l.GetContext() <= s {
		return nil
	}
	return ErrContextViolation
}

// AssertExactContext allows a program to assert a required context and fail gracefully
func (l *Laforge) AssertExactContext(s StateContext) error {
	if l.GetContext() == s {
		return nil
	}
	return ErrContextViolation
}

// GetContext returns the current state's context
func (l *Laforge) GetContext() StateContext {
	if l.ValidTeam {
		return TeamContext
	}
	if l.ValidBuild {
		return BuildContext
	}
	if l.ValidEnv {
		return EnvContext
	}
	if l.ValidBase {
		return BaseContext
	}
	if l.ValidGlobal {
		return GlobalContext
	}
	return NoContext
}

// CreateIndex maps out all of the known networks, identities, scripts, commands, and hosts within a laforge configuration snapshot.
func (l *Laforge) CreateIndex() {
	if l.PathRegistry == nil {
		l.PathRegistry = &PathRegistry{
			DB: map[CallFile]*PathResolver{},
		}
	}
	if l.PathRegistry.DB[l.Caller.Current()] == nil {
		l.PathRegistry.DB[l.Caller.Current()] = &PathResolver{
			Mapping:    map[string]*LocalFileRef{},
			Unresolved: map[string]bool{},
		}
	}
	currPathResolver := l.PathRegistry.DB[l.Caller.Current()]
	l.Hosts = map[string]*Host{}
	l.Networks = map[string]*Network{}
	l.Identities = map[string]*Identity{}
	l.Scripts = map[string]*Script{}
	l.Commands = map[string]*Command{}
	l.RemoteFiles = map[string]*RemoteFile{}
	l.DNSRecords = map[string]*DNSRecord{}
	l.Teams = map[string]*Team{}
	l.Builds = map[string]*Build{}
	l.Competitions = map[string]*Competition{}
	l.Environments = map[string]*Environment{}
	l.ProvisionedNetworks = map[string]*ProvisionedNetwork{}
	l.ProvisionedHosts = map[string]*ProvisionedHost{}
	l.ProvisioningSteps = map[string]*ProvisioningStep{}
	l.Connections = map[string]*Connection{}
	for _, x := range l.DefinedHosts {
		l.Hosts[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedNetworks {
		l.Networks[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedIdentities {
		err := x.ResolveSource(l, currPathResolver, l.Caller.Current())
		if err != nil {
			cli.Logger.Errorf("%T %s had a source location that was not found: %v", x, x.ID, err)
		}
		l.Identities[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedScripts {
		err := x.ResolveSource(l, currPathResolver, l.Caller.Current())
		if err != nil {
			cli.Logger.Errorf("%T %s had a source location that was not found: %v", x, x.ID, err)
		}
		l.Scripts[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedCommands {
		l.Commands[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedRemoteFiles {
		err := x.ResolveSource(l, currPathResolver, l.Caller.Current())
		if err != nil {
			cli.Logger.Errorf("%T %s had a source location that was not found: %v", x, x.ID, err)
		}
		l.RemoteFiles[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedConnections {
		l.Connections[x.LaforgeID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedProvisioningSteps {
		l.ProvisioningSteps[x.LaforgeID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedProvisionedHosts {
		l.ProvisionedHosts[x.LaforgeID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedProvisionedNetworks {
		l.ProvisionedNetworks[x.LaforgeID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedTeams {
		l.Teams[x.LaforgeID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedDNSRecords {
		l.DNSRecords[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedBuilds {
		l.Builds[x.LaforgeID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedEnvironments {
		l.Environments[x.LaforgeID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedCompetitions {
		l.Competitions[x.LaforgeID()] = x
		x.Caller = l.Caller
	}
}

// Update performs a patching operation on source (l) with diff (diff), using the diff's merge conflict settings as appropriate.
func (l *Laforge) Update(diff *Laforge) (*Laforge, error) {
	err := mergo.Merge(l.PathRegistry, diff.PathRegistry, mergo.WithOverride)
	if err != nil {
		panic(err)
	}
	l.Caller = diff.Caller
	if l.Filename != diff.Filename && diff.Filename != "" {
		l.Filename = diff.Filename
	}
	if l.BaseDir != diff.BaseDir && diff.BaseDir != "" {
		l.BaseDir = diff.BaseDir
	}
	baseUser := l.User
	if l.User == nil {
		l.User = &User{}
		baseUser = l.User
	}
	if diff.User == nil {
		diff.User = &User{}
	}
	err = mergo.Merge(baseUser, diff.User, mergo.WithOverride)
	if err != nil {
		return l, errors.WithStack(err)
	}
	l.User = baseUser
	return l, nil
}

// Mask attempts to apply a differential update betweeen base and layer, returning a modified base and any errors it encountered.
//nolint:gocyclo
func Mask(base, layer *Laforge) (*Laforge, error) {
	layer.CreateIndex()
	for name, obj := range layer.Hosts {
		orig, found := base.Hosts[name]
		if !found {
			base.Hosts[name] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*Host)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for name, obj := range layer.Networks {
		orig, found := base.Networks[name]
		if !found {
			base.Networks[name] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*Network)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for name, obj := range layer.Identities {
		orig, found := base.Identities[name]
		if !found {
			base.Identities[name] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*Identity)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for name, obj := range layer.Scripts {
		orig, found := base.Scripts[name]
		if !found {
			base.Scripts[name] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*Script)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for name, obj := range layer.Commands {
		orig, found := base.Commands[name]
		if !found {
			base.Commands[name] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*Command)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for name, obj := range layer.RemoteFiles {
		orig, found := base.RemoteFiles[name]
		if !found {
			base.RemoteFiles[name] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*RemoteFile)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for name, obj := range layer.DNSRecords {
		orig, found := base.DNSRecords[name]
		if !found {
			base.DNSRecords[name] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*DNSRecord)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}

	for id, obj := range layer.Competitions {
		orig, found := base.Competitions[id]
		if !found {
			base.Competitions[id] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*Competition)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for id, obj := range layer.Environments {
		orig, found := base.Environments[id]
		if !found {
			base.Environments[id] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*Environment)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for id, obj := range layer.Builds {
		orig, found := base.Builds[id]
		if !found {
			base.Builds[id] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*Build)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for id, obj := range layer.Teams {
		orig, found := base.Teams[id]
		if !found {
			base.Teams[id] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*Team)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for id, obj := range layer.ProvisionedNetworks {
		orig, found := base.ProvisionedNetworks[id]
		if !found {
			base.ProvisionedNetworks[id] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*ProvisionedNetwork)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for id, obj := range layer.ProvisionedHosts {
		orig, found := base.ProvisionedHosts[id]
		if !found {
			base.ProvisionedHosts[id] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*ProvisionedHost)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for id, obj := range layer.ProvisioningSteps {
		orig, found := base.ProvisioningSteps[id]
		if !found {
			base.ProvisioningSteps[id] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*ProvisioningStep)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	for id, obj := range layer.Connections {
		orig, found := base.Connections[id]
		if !found {
			base.Connections[id] = obj
			continue
		}
		res, err := SmartMerge(orig, obj, false)
		if err != nil {
			return nil, err
		}
		orig, ok := res.(*Connection)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
	}
	return base.Update(layer)
}

// IndexHostDependencies enumerates all Host objects in the state and indexes their dependencies, reporting errors
func (l *Laforge) IndexHostDependencies() error {
	for _, h := range l.Hosts {
		err := h.Index(l)
		if err != nil {
			return err
		}
	}
	return nil
}

// IndexEnvironmentDependencies enumerates all known environments and makes sure they have valid network inclusions
func (l *Laforge) IndexEnvironmentDependencies() error {
	for _, e := range l.Environments {
		comp, found := l.Competitions[e.ParentLaforgeID()]
		if !found {
			return fmt.Errorf("competition %s for env %s could not be located", e.ParentLaforgeID(), e.Path())
		}
		e.Competition = comp
		err := e.ResolveIncludedNetworks(l)
		if err != nil {
			return err
		}
	}
	return nil
}

// IndexBuildDependencies enumerates all known builds and ensures it has competition and environment associations
func (l *Laforge) IndexBuildDependencies() error {
	for _, b := range l.Builds {
		env, found := l.Environments[b.ParentLaforgeID()]
		if !found {
			return fmt.Errorf("environment %s for build %s could not be located", b.ParentLaforgeID(), b.Path())
		}
		b.Environment = env
		b.Competition = env.Competition
		if b.Teams == nil {
			b.Teams = map[string]*Team{}
		}
		b.Dir = b.Caller.Current().CallerDir
	}
	return nil
}

// IndexTeamDependencies enumerates all known teams and ensures they have proper associations
func (l *Laforge) IndexTeamDependencies() error {
	for _, team := range l.Teams {
		build, found := l.Builds[team.ParentLaforgeID()]
		if !found {
			return fmt.Errorf("build %s for team %s could not be located", team.ParentLaforgeID(), team.Path())
		}
		team.Build = build
		team.Environment = build.Environment
		team.Competition = build.Competition
		team.Dir = team.Caller.Current().CallerDir
		if team.ProvisionedNetworks == nil {
			team.ProvisionedNetworks = map[string]*ProvisionedNetwork{}
		}
	}
	return nil
}

// IndexProvisionedNetworkDependencies enumerates all known provisioned networks and ensures they have proper associations
func (l *Laforge) IndexProvisionedNetworkDependencies() error {
	for _, net := range l.ProvisionedNetworks {
		team, found := l.Teams[net.ParentLaforgeID()]
		if !found {
			return fmt.Errorf("team %s for provisioned network %s could not be located", net.ParentLaforgeID(), net.Path())
		}
		team.ProvisionedNetworks[net.Path()] = net
		net.Team = team
		net.Build = team.Build
		net.Environment = team.Environment
		net.Competition = team.Competition
		net.Dir = net.Caller.Current().CallerDir
		parent, found := l.Networks[net.NetworkID]
		if !found {
			return fmt.Errorf("network %s for provisioned network %s could not be located", net.NetworkID, net.Path())
		}
		net.Network = parent
		if net.ProvisionedHosts == nil {
			net.ProvisionedHosts = map[string]*ProvisionedHost{}
		}
	}
	return nil
}

// IndexProvisionedHostDependencies enumerates all known provisioned hosts and ensures they have proper associations
func (l *Laforge) IndexProvisionedHostDependencies() error {
	for _, host := range l.ProvisionedHosts {
		net, found := l.ProvisionedNetworks[host.ParentLaforgeID()]
		if !found {
			return fmt.Errorf("provisioned network %s for provisioned host %s could not be located", host.ParentLaforgeID(), host.Path())
		}
		net.ProvisionedHosts[host.Path()] = host
		host.ProvisionedNetwork = net
		host.Team = net.Team
		host.Build = net.Build
		host.Environment = net.Environment
		host.Competition = net.Competition
		host.Network = net.Network
		parent, found := l.Hosts[host.HostID]
		if !found {
			return fmt.Errorf("host %s for provisioned host %s could not be located", host.ParentLaforgeID(), host.Path())
		}
		host.Host = parent
		if host.ProvisioningSteps == nil {
			host.ProvisioningSteps = map[string]*ProvisioningStep{}
		}
		host.Dir = host.Caller.Current().CallerDir
	}
	return nil
}

// IndexConnectionDependencies enumerates all known connections and ensures they have proper associations
func (l *Laforge) IndexConnectionDependencies() error {
	for _, conn := range l.Connections {
		ph, found := l.ProvisionedHosts[conn.ParentLaforgeID()]
		if !found {
			return fmt.Errorf("provisioned host %s for conn %s could not be located", conn.ParentLaforgeID(), conn.Path())
		}
		conn.ProvisionedHost = ph
		ph.Conn = conn
		conn.ProvisionedNetwork = ph.ProvisionedNetwork
		conn.Host = ph.Host
		conn.Network = ph.Network
		conn.Team = ph.Team
		conn.Build = ph.Build
		conn.Environment = ph.Environment
		conn.Competition = ph.Competition
	}
	return nil
}

// IndexProvisioningStepDependencies enumerates all known provisioning steps and ensures they have proper associations
func (l *Laforge) IndexProvisioningStepDependencies() error {
	for _, ps := range l.ProvisioningSteps {
		ph, found := l.ProvisionedHosts[ps.ParentLaforgeID()]
		if !found {
			return fmt.Errorf("provisioned host %s for provisioning step %s could not be located", ps.ParentLaforgeID(), ps.Path())
		}
		ps.ProvisionedHost = ph
		ps.ProvisionedNetwork = ph.ProvisionedNetwork
		ps.Team = ph.Team
		ps.Build = ph.Build
		ps.Environment = ph.Environment
		ps.Competition = ph.Competition
		ps.Host = ph.Host
		ps.Network = ph.Network
		ph.ProvisioningSteps[ps.Path()] = ps
		ps.Dir = ph.Dir
		switch ps.ProvisionerType {
		case ObjectTypeCommand.String():
			prov, found := l.Commands[ps.ProvisionerID]
			if !found {
				return fmt.Errorf("command %s for provisioning step %s could not be located", ps.ProvisionerID, ps.Path())
			}
			ps.Command = prov
			ps.Provisioner = prov
		case ObjectTypeDNSRecord.String():
			prov, found := l.DNSRecords[ps.ProvisionerID]
			if !found {
				return fmt.Errorf("dns record %s for provisioning step %s could not be located", ps.ProvisionerID, ps.Path())
			}
			ps.DNSRecord = prov
			ps.Provisioner = prov
		case ObjectTypeRemoteFile.String():
			prov, found := l.RemoteFiles[ps.ProvisionerID]
			if !found {
				return fmt.Errorf("remote file %s for provisioning step %s could not be located", ps.ProvisionerID, ps.Path())
			}
			ps.RemoteFile = prov
			ps.Provisioner = prov
		case ObjectTypeScript.String():
			prov, found := l.Scripts[ps.ProvisionerID]
			if !found {
				return fmt.Errorf("script %s for provisioning step %s could not be located", ps.ProvisionerID, ps.Path())
			}
			ps.Script = prov
			ps.Provisioner = prov
		default:
			return fmt.Errorf("unknown provisioner type %s for provisioning step %s", ps.ProvisionerType, ps.Path())
		}
	}
	return nil
}

// InitializeTeamContext returns a base context preset with a team context configuration
func InitializeTeamContext(globalconfig, buildconfig, teamconfig string) (*Laforge, error) {
	clone, err := LoadFiles(globalconfig, buildconfig)

	if err != nil {
		return nil, err
	}

	err = clone.IndexHostDependencies()

	if err != nil {
		return nil, err
	}

	err = clone.IndexEnvironmentDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexBuildDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexTeamDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexProvisionedNetworkDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexProvisionedHostDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexConnectionDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexProvisioningStepDependencies()
	if err != nil {
		return nil, err
	}
	t := &Team{}
	//nolint:gosec
	tData, err := ioutil.ReadFile(teamconfig)
	if err != nil {
		return nil, err
	}
	err = HCLBytesToObject(tData, t)
	if err != nil {
		return nil, err
	}
	currTeam, ok := clone.Teams[t.LaforgeID()]
	if !ok {
		return nil, fmt.Errorf("could not find team %s in the current environment", t.LaforgeID())
	}
	clone.CurrentTeam = currTeam
	clone.CurrentBuild = currTeam.Build
	clone.CurrentEnv = currTeam.Environment
	clone.CurrentCompetition = currTeam.Competition
	clone.TeamContextID = clone.CurrentTeam.LaforgeID()
	clone.BuildContextID = clone.CurrentBuild.LaforgeID()
	clone.EnvContextID = clone.CurrentEnv.LaforgeID()
	clone.BaseContextID = clone.CurrentCompetition.LaforgeID()
	return clone, nil
}

// InitializeBuildContext returns a base context preset with a build context configuration
func InitializeBuildContext(globalconfig, buildconfig string) (*Laforge, error) {
	clone, err := LoadFiles(globalconfig, buildconfig)

	if err != nil {
		return nil, err
	}

	err = clone.IndexHostDependencies()

	if err != nil {
		return nil, err
	}

	err = clone.IndexEnvironmentDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexBuildDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexTeamDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexProvisionedNetworkDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexProvisionedHostDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexConnectionDependencies()
	if err != nil {
		return nil, err
	}
	err = clone.IndexProvisioningStepDependencies()
	if err != nil {
		return nil, err
	}
	b := &Build{}
	//nolint:gosec
	tData, err := ioutil.ReadFile(buildconfig)
	if err != nil {
		return nil, err
	}
	err = HCLBytesToObject(tData, b)
	if err != nil {
		return nil, err
	}
	currBuild, ok := clone.Builds[b.LaforgeID()]
	if !ok {
		return nil, fmt.Errorf("could not find build %s in the current environment", b.LaforgeID())
	}
	clone.CurrentBuild = currBuild
	clone.CurrentEnv = currBuild.Environment
	clone.CurrentCompetition = currBuild.Competition
	clone.BuildContextID = clone.CurrentBuild.LaforgeID()
	clone.EnvContextID = clone.CurrentEnv.LaforgeID()
	clone.BaseContextID = clone.CurrentCompetition.LaforgeID()
	return clone, nil
}

// InitializeEnvContext returns a base context preset with a env focused configuration
func InitializeEnvContext(globalconfig, envconfig string) (*Laforge, error) {
	clone, err := LoadFiles(globalconfig, envconfig)

	if err != nil {
		return nil, err
	}

	err = clone.IndexHostDependencies()

	if err != nil {
		return nil, err
	}

	err = clone.IndexEnvironmentDependencies()
	if err != nil {
		return nil, err
	}
	e := &Environment{}
	//nolint:gosec
	tData, err := ioutil.ReadFile(envconfig)
	if err != nil {
		return nil, err
	}
	err = HCLBytesToObject(tData, e)
	if err != nil {
		return nil, err
	}
	currEnv, ok := clone.Environments[e.LaforgeID()]
	if !ok {
		return nil, fmt.Errorf("could not find build %s in the current environment", e.LaforgeID())
	}
	clone.CurrentEnv = currEnv
	clone.CurrentCompetition = currEnv.Competition
	clone.EnvContextID = clone.CurrentEnv.LaforgeID()
	clone.BaseContextID = clone.CurrentCompetition.LaforgeID()
	return clone, nil
}

// InitializeBaseContext returns a base context preset with a base focused configuration
func InitializeBaseContext(globalconfig, baseconfig string) (*Laforge, error) {
	clone, err := LoadFiles(globalconfig, baseconfig)
	if err != nil {
		return nil, err
	}
	if clone != nil {
		err = clone.IndexHostDependencies()

		if err != nil {
			return nil, err
		}
	}
	return clone, nil
}

// LoadFromContext attempts to bootstrap the given state from it's assumed Context Level
func (l *Laforge) LoadFromContext() error {
	var clone *Laforge
	var err error
	switch l.GetContext() {
	case TeamContext:
		clone, err = InitializeTeamContext(l.GlobalConfigFile(), l.BuildConfigFile(), l.TeamConfigFile())
	case BuildContext:
		clone, err = InitializeBuildContext(l.GlobalConfigFile(), l.BuildConfigFile())
	case EnvContext:
		clone, err = InitializeEnvContext(l.GlobalConfigFile(), l.EnvConfigFile())
	case BaseContext:
		clone, err = InitializeBaseContext(l.GlobalConfigFile(), l.BaseConfigFile())
	default:
		return ErrContextViolation
	}
	if err != nil {
		return err
	}
	if clone != nil {
		err = mergo.Merge(l, clone, mergo.WithOverride, mergo.WithAppendSlice)
		if err != nil {
			return err
		}
	}
	return nil
}

// GlobalConfigFile is a helper method for creating an absolute path to the global configuration file
func (l *Laforge) GlobalConfigFile() string {
	return filepath.Join(l.GlobalRoot, "global.laforge")
}

// EnvConfigFile is a helper method for creating an absolute path to the environment configuration file
func (l *Laforge) EnvConfigFile() string {
	return filepath.Join(l.EnvRoot, envFile)
}

// BaseConfigFile is a helper method for creating an absolute path to the base configuration file
func (l *Laforge) BaseConfigFile() string {
	return filepath.Join(l.BaseRoot, "base.laforge")
}

// BuildConfigFile is a helper method for creating an absolute path to the build results file
func (l *Laforge) BuildConfigFile() string {
	return filepath.Join(l.BuildRoot, "build.laforge")
}

// TeamConfigFile is a helper method for creating an absolute path to the build results file
func (l *Laforge) TeamConfigFile() string {
	return filepath.Join(l.TeamRoot, "team.laforge")
}

// InitializeContext attempts to resolve the current state context from path traversal
func (l *Laforge) InitializeContext() error {
	gcl, err := LocateGlobalConfig()
	if err != nil {
		if err != ErrNoGlobalConfig {
			return err
		}
		cli.Logger.Errorf("No config found! Run `laforge configure` before continuing!")
		return errors.New("no global configuration found")
	}
	l.ValidGlobal = true
	l.GlobalRoot = filepath.Dir(gcl)
	tcl, err := LocateTeamConfig()
	if err != nil {
		if err != ErrNoConfigRootReached {
			return err
		}
	} else {
		tclabs, err := filepath.Abs(tcl)
		if err != nil {
			return err
		}
		l.ValidTeam = true
		l.TeamRoot = filepath.Dir(tclabs)
	}
	bucl, err := LocateBuildConfig()
	if err != nil {
		if err != ErrNoConfigRootReached {
			return err
		}
	} else {
		buclabs, err := filepath.Abs(bucl)
		if err != nil {
			return err
		}
		l.ValidBuild = true
		l.BuildRoot = filepath.Dir(buclabs)
	}
	ecl, err := LocateEnvConfig()
	if err != nil {
		if err != ErrNoConfigRootReached {
			return err
		}
	} else {
		eclabs, err := filepath.Abs(ecl)
		if err != nil {
			return err
		}
		l.ValidEnv = true
		l.EnvRoot = filepath.Dir(eclabs)
	}
	bcl, err := LocateBaseConfig()
	if err != nil {
		if err != ErrNoConfigRootReached {
			return err
		}
	} else {
		bclabs, err := filepath.Abs(bcl)
		if err != nil {
			return err
		}
		l.ValidBase = true
		l.BaseRoot = filepath.Dir(bclabs)
	}
	return nil
}

// Bootstrap performs a full lifecycle bootstrap of the laforge state in a context aware way.
func Bootstrap() (*Laforge, error) {
	base := &Laforge{}
	err := base.InitializeContext()
	if err != nil {
		return base, err
	}
	err = base.AssertMinContext(BaseContext)
	if err != nil {
		cli.Logger.Infof("No base.laforge or env.laforge found in your current directory tree!")
		return base, errors.Wrapf(ErrContextViolation, "no base.laforge or env.laforge found")
	}
	err = base.LoadFromContext()
	if err != nil {
		return base, err
	}

	if base.BaseDir == "" {
		base.BaseDir = base.BaseRoot
	}

	base.InitialContext = base.GetContext()
	return base, err
}

var (
	baseSubDirs = []string{
		"config",
		"scripts",
		commandsDir,
		hostsDir,
		"networks",
		"identities",
		"files",
		envsDir,
	}
)

// InitializeBaseDirectory configures a skeleton competition repository in the caller's current directory
func (l *Laforge) InitializeBaseDirectory(overwrite bool) error {
	err := l.AssertExactContext(GlobalContext)
	if err != nil && !overwrite {
		return errors.WithStack(err)
	}

	for _, d := range baseSubDirs {
		dpath := filepath.Join(l.CurrDir, d)
		if _, err := os.Stat(dpath); err == nil {
			if !overwrite {
				return fmt.Errorf("Cannot initialize base directory - %s folder already exists! (--force/-f to overwrite)", d)
			}
			err := os.RemoveAll(dpath)
			if err != nil {
				return err
			}
		}
		//nolint:gosec
		err = os.MkdirAll(dpath, 0755)
		if err != nil {
			return err
		}
		keeper := filepath.Join(dpath, ".gitkeep")
		newFile, err := os.Create(keeper)
		if err != nil {
			return errors.WithMessage(err, fmt.Sprintf("cannot touch .gitkeep inside base directory subfolder %s", d))
		}
		err = newFile.Close()
		if err != nil {
			return err
		}
	}

	basefile := filepath.Join(l.CurrDir, "base.laforge")
	newFile, err := os.Create(basefile)
	if err != nil {
		return errors.WithMessage(err, "error creating base.laforge file")
	}
	currabs, err := filepath.Abs(l.CurrDir)
	if err != nil {
		return errors.WithStack(err)
	}

	cli.Logger.Debugf("currabs = %s", currabs)
	defcomp, err := RenderHCLv2Object(defaultCompetition(filepath.Base(currabs)))
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = newFile.Write(defcomp)
	if err != nil {
		return errors.WithStack(err)
	}
	err = newFile.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
