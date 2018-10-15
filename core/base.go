// Package core contains the entire implementation of the Laforge configuration language. It includes
// it's own loader and dependency resolution mechanisms and is meant to be the source of truth around
// declaration logic.
package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
)

// Laforge defines the type that holds the global namespace within the laforge configuration engine
type Laforge struct {
	Filename                string                      `json:"filename"`
	Includes                []string                    `json:"-"`
	DependencyGraph         treeprint.Tree              `json:"-"`
	BaseDir                 string                      `hcl:"base_dir,attr" json:"base_dir,omitempty"`
	CurrDir                 string                      `json:"-"`
	User                    User                        `hcl:"user,block" json:"user,omitempty"`
	DefinedCompetitions     []*Competition              `hcl:"competition,block" json:"competitions,omitempty"`
	DefinedEnvironments     []*Environment              `hcl:"environment,block" json:"environments,omitempty"`
	DefinedBuilds           []*Build                    `hcl:"build,block" json:"builds,omitempty"`
	DefinedTeams            []*Team                     `hcl:"team,block" json:"teams,omitempty"`
	DefinedHosts            []*Host                     `hcl:"host,block" json:"hosts,omitempty"`
	DefinedNetworks         []*Network                  `hcl:"network,block" json:"networks,omitempty"`
	DefinedIdentities       []*Identity                 `hcl:"identity,block" json:"identities,omitempty"`
	DefinedScripts          []*Script                   `hcl:"script,block" json:"scripts,omitempty"`
	DefinedCommands         []*Command                  `hcl:"command,block" json:"defined_commands,omitempty"`
	DefinedRemoteFiles      []*RemoteFile               `hcl:"remote_file,block" json:"defined_files,omitempty"`
	DefinedDNSRecords       []*DNSRecord                `hcl:"dns_record,block" json:"defined_dns_records,omitempty"`
	DefinedProvisionedHosts []*ProvisionedHost          `hcl:"provisioned_host,block" json:"provisioned_hosts,omitempty"`
	Hosts                   map[string]*Host            `json:"-"`
	Networks                map[string]*Network         `json:"-"`
	Identities              map[string]*Identity        `json:"-"`
	Scripts                 map[string]*Script          `json:"-"`
	Commands                map[string]*Command         `json:"-"`
	RemoteFiles             map[string]*RemoteFile      `json:"-"`
	DNSRecords              map[string]*DNSRecord       `json:"-"`
	Teams                   map[string]*Team            `json:"-"`
	Builds                  map[string]*Build           `json:"-"`
	ProvisionedHosts        map[string]*ProvisionedHost `json:"-"`
	Environments            map[string]*Environment     `json:"-"`
	Competitions            map[string]*Competition     `json:"-"`
	Caller                  Caller                      `json:"-"`
	ValidTeam               bool                        `json:"-"`
	ValidBuild              bool                        `json:"-"`
	ValidEnv                bool                        `json:"-"`
	ValidBase               bool                        `json:"-"`
	ValidGlobal             bool                        `json:"-"`
	ClearToBuild            bool                        `json:"-"`
	TeamRoot                string                      `json:"-"`
	BuildRoot               string                      `json:"-"`
	EnvRoot                 string                      `json:"-"`
	BaseRoot                string                      `json:"-"`
	GlobalRoot              string                      `json:"-"`
	TeamAbsPath             string                      `json:"-"`
	BuildAbsPath            string                      `json:"-"`
	EnvAbsPath              string                      `json:"-"`
	BaseAbsPath             string                      `json:"-"`
	TeamContextID           string                      `json:"-"`
	BuildContextID          string                      `json:"-"`
	EnvContextID            string                      `json:"-"`
	BaseContextID           string                      `json:"-"`
	GlobalContextID         string                      `json:"-"`
	CurrentEnv              *Environment                `json:"-"`
	CurrentBuild            *Build                      `json:"-"`
	CurrentTeam             *Team                       `json:"-"`
	CurrentCompetition      *Competition                `json:"-"`
	InitialContext          StateContext                `json:"-"`
	PathRegistry            *PathRegistry               `json:"-"`
}

// Opt defines a basic HCLv2 option label:
//	config "keyName" {
//		value = "valueData"
//	}
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
type OnConflict struct {
	Do     string `cty:"do" hcl:"do,attr" json:"do,omitempty"`
	Append bool   `cty:"append" hcl:"append,optional" json:"append,omitempty"`
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
	l.ProvisionedHosts = map[string]*ProvisionedHost{}
	l.Builds = map[string]*Build{}
	l.Competitions = map[string]*Competition{}
	l.Environments = map[string]*Environment{}
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
			Logger.Errorf("%T %s had a source location that was not found: %v", x, x.ID, err)
		}
		l.Identities[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedScripts {
		err := x.ResolveSource(l, currPathResolver, l.Caller.Current())
		if err != nil {
			Logger.Errorf("%T %s had a source location that was not found: %v", x, x.ID, err)
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
			Logger.Errorf("%T %s had a source location that was not found: %v", x, x.ID, err)
		}
		l.RemoteFiles[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedTeams {
		l.Teams[x.GetID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedDNSRecords {
		l.DNSRecords[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedProvisionedHosts {
		l.ProvisionedHosts[x.GetID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedBuilds {
		l.Builds[x.GetID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedEnvironments {
		l.Environments[x.GetID()] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedCompetitions {
		l.Competitions[x.GetID()] = x
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
	newUser := l.User
	err = mergo.Merge(&newUser, diff.User, mergo.WithOverride)
	if err != nil {
		return l, errors.WithStack(err)
	}
	l.User = newUser
	// if l.Competition == nil && diff.Competition != nil {
	// 	l.Competition = diff.Competition
	// } else if l.Competition != nil && diff.Competition != nil {
	// 	res, err := SmartMerge(l.Competition, diff.Competition, false)
	// 	if err != nil {
	// 		return l, errors.WithStack(err)
	// 	}
	// 	orig, ok := res.(*Competition)
	// 	if !ok {
	// 		return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
	// 	}
	// 	l.Competition = orig
	// }
	// if l.Environment == nil && diff.Environment != nil {
	// 	l.Environment = diff.Environment
	// } else if l.Environment != nil && diff.Environment != nil {
	// 	res, err := SmartMerge(l.Environment, diff.Environment, true)
	// 	if err != nil {
	// 		return l, errors.WithStack(err)
	// 	}
	// 	orig, ok := res.(*Environment)
	// 	if !ok {
	// 		return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
	// 	}
	// 	l.Environment = orig
	// 	if l.Environment.BaseDir == "" && l.EnvRoot != "" {
	// 		l.Environment.BaseDir = l.EnvRoot
	// 	}
	// }
	// if l.Team == nil && diff.Team != nil {
	// 	l.Team = diff.Team
	// } else if l.Team != nil && diff.Team != nil {
	// 	res, err := SmartMerge(l.Team, diff.Team, true)
	// 	if err != nil {
	// 		return l, errors.WithStack(err)
	// 	}
	// 	orig, ok := res.(*Team)
	// 	if !ok {
	// 		return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
	// 	}
	// 	l.Team = orig
	// }
	// if l.Build == nil && diff.Build != nil {
	// 	l.Build = diff.Build
	// } else if l.Build != nil && diff.Build != nil {
	// 	res, err := SmartMerge(l.Build, diff.Build, true)
	// 	if err != nil {
	// 		return l, errors.WithStack(err)
	// 	}
	// 	orig, ok := res.(*Build)
	// 	if !ok {
	// 		return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
	// 	}
	// 	l.Build = orig
	// }

	// if diff.Build != nil {
	// 	l.Build = diff.Build
	// }

	return l, nil
}

// Mask attempts to apply a differential update betweeen base and layer, returning a modified base and any errors it encountered.
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
		comp, ok := l.Competitions[e.GetParentID()]
		if ok {
			e.Competition = comp
		}
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
		env, ok := l.Environments[b.GetParentID()]
		if ok {
			b.Environment = env
			if b.Environment.Competition != nil {
				b.Competition = b.Environment.Competition
			}
		}
	}
	return nil
}

// IndexTeamDependencies enumerates all known teams and ensures they have proper associations
func (l *Laforge) IndexTeamDependencies() error {
	for _, t := range l.Teams {
		build, ok := l.Builds[t.GetParentID()]
		if ok {
			t.Build = build
			if t.Build.Environment != nil {
				t.Environment = t.Build.Environment
				if t.Environment.Competition != nil {
					t.Competition = t.Environment.Competition
				}
				if t.Environment.Teams == nil {
					t.Environment.Teams = map[int]*Team{}
				}
				t.Environment.Teams[t.TeamNumber] = t
				if t.Build.Teams == nil {
					t.Build.Teams = map[int]*Team{}
				}
				t.Build.Teams[t.TeamNumber] = t
			}
		}
	}
	return nil
}

// IndexProvisionedHostDependencies enumerates all known provisioned hosts and ensures they have proper associations
func (l *Laforge) IndexProvisionedHostDependencies() error {
	for _, p := range l.ProvisionedHosts {
		team, ok := l.Teams[p.GetParentID()]
		if ok {
			p.Team = team
			if p.Team.Build != nil {
				p.Build = p.Team.Build
				if p.Team.Hosts == nil {
					p.Team.Hosts = map[string]*ProvisionedHost{}
				}
				p.Team.Hosts[p.ID] = p
				if p.Build.Environment != nil {
					p.Environment = p.Build.Environment
					if p.Environment.Competition != nil {
						p.Competition = p.Environment.Competition
					}
				}
			}
		}
		host, ok := l.Hosts[p.HostID]
		if ok {
			p.Host = host
		}
		network, ok := l.Networks[p.NetworkID]
		if ok {
			p.Network = network
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
	clone.IndexEnvironmentDependencies()
	clone.IndexBuildDependencies()
	clone.IndexTeamDependencies()
	clone.IndexProvisionedHostDependencies()
	t := &Team{}
	tData, err := ioutil.ReadFile(teamconfig)
	if err != nil {
		return nil, err
	}
	err = HCLBytesToObject(tData, t)
	if err != nil {
		return nil, err
	}
	currTeam, ok := clone.Teams[t.GetID()]
	if !ok {
		return nil, fmt.Errorf("could not find team %s in the current environment", t.GetID())
	}
	clone.CurrentTeam = currTeam
	clone.CurrentBuild = currTeam.Build
	clone.CurrentEnv = currTeam.Environment
	clone.CurrentCompetition = currTeam.Competition
	clone.TeamContextID = clone.CurrentTeam.GetID()
	clone.BuildContextID = clone.CurrentBuild.GetID()
	clone.EnvContextID = clone.CurrentEnv.GetID()
	clone.BaseContextID = clone.CurrentCompetition.GetID()
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
	clone.IndexEnvironmentDependencies()
	clone.IndexBuildDependencies()
	clone.IndexTeamDependencies()
	clone.IndexProvisionedHostDependencies()
	b := &Build{}
	tData, err := ioutil.ReadFile(buildconfig)
	if err != nil {
		return nil, err
	}
	err = HCLBytesToObject(tData, b)
	if err != nil {
		return nil, err
	}
	currBuild, ok := clone.Builds[b.GetID()]
	if !ok {
		return nil, fmt.Errorf("could not find build %s in the current environment", b.GetID())
	}
	clone.CurrentBuild = currBuild
	clone.CurrentEnv = currBuild.Environment
	clone.CurrentCompetition = currBuild.Competition
	clone.BuildContextID = clone.CurrentBuild.GetID()
	clone.EnvContextID = clone.CurrentEnv.GetID()
	clone.BaseContextID = clone.CurrentCompetition.GetID()
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
	clone.IndexEnvironmentDependencies()
	e := &Environment{}
	tData, err := ioutil.ReadFile(envconfig)
	if err != nil {
		return nil, err
	}
	err = HCLBytesToObject(tData, e)
	if err != nil {
		return nil, err
	}
	currEnv, ok := clone.Environments[e.GetID()]
	if !ok {
		return nil, fmt.Errorf("could not find build %s in the current environment", e.GetID())
	}
	clone.CurrentEnv = currEnv
	clone.CurrentCompetition = currEnv.Competition
	clone.EnvContextID = clone.CurrentEnv.GetID()
	clone.BaseContextID = clone.CurrentCompetition.GetID()
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
	return filepath.Join(l.EnvRoot, "env.laforge")
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
		Logger.Errorf("No config found! Run `laforge configure` before continuing!")
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
	err = nil
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
		Logger.Infof("No base.laforge or env.laforge found in your current directory tree!")
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
		"commands",
		"hosts",
		"networks",
		"identities",
		"files",
		"envs",
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
			os.RemoveAll(dpath)
		}
		os.MkdirAll(dpath, 0755)
		keeper := filepath.Join(dpath, ".gitkeep")
		newFile, err := os.Create(keeper)
		if err != nil {
			return errors.WithMessage(err, fmt.Sprintf("cannot touch .gitkeep inside base directory subfolder %s", d))
		}
		newFile.Close()
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

	Logger.Debugf("currabs = %s", currabs)
	defcomp, err := RenderHCLv2Object(defaultCompetition(filepath.Base(currabs)))
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = newFile.Write(defcomp)
	if err != nil {
		return errors.WithStack(err)
	}
	newFile.Close()
	return nil
}
