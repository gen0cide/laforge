// Package core contains the entire implementation of the Laforge configuration language. It includes
// it's own loader and dependency resolution mechanisms and is meant to be the source of truth around
// declaration logic.
package core

import (
	"fmt"
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

	//GlobalContext is a context level representing a valid global configuration was found
	GlobalContext

	// NoContext is a context level representing neither a env.laforge or base.laforge was found
	NoContext
)

var (
	// ErrContextViolation is an error thrown when an action is taken that attempts to exceed the current state's ScopeContext
	ErrContextViolation = errors.New("context scope violation")

	//ErrAbsPathDeclNotExist is thrown when an absolute path is listed in a file yet it cannot be resolved on the local system
	ErrAbsPathDeclNotExist = errors.New("absolute path was not found on local system")

	// ErrInvalidEnvName is thrown when an environment name does not meet specified regulations around environment naming conventions
	ErrInvalidEnvName = errors.New("environment names can only contain lowercase alphanumeric and dash characters (a valid subdomain)")
)

// Laforge defines the type that holds the global namespace within the laforge configuration engine
type Laforge struct {
	Filename          string                 `json:"filename"`
	Includes          []string               `json:"include,omitempty"`
	DependencyGraph   treeprint.Tree         `json:"-"`
	BaseDir           string                 `hcl:"base_dir,attr" json:"base_dir,omitempty"`
	CurrDir           string                 `json:"current_dir,omitempty"`
	User              User                   `hcl:"user,block" cty:"user" json:"user,omitempty"`
	Competition       *Competition           `hcl:"competition,block" json:"competition,omitempty"`
	Environment       *Environment           `hcl:"environment,block" json:"environment,omitempty"`
	Build             *Build                 `hcl:"build,block" json:"build,omitempty"`
	Team              *Team                  `hcl:"team,block" json:"team,omitempty"`
	DefinedHosts      []*Host                `hcl:"host,block" json:"defined_hosts,omitempty"`
	DefinedNetworks   []*Network             `hcl:"network,block" json:"defined_networks,omitempty"`
	DefinedIdentities []*Identity            `hcl:"identity,block" json:"defined_identities,omitempty"`
	DefinedScripts    []*Script              `hcl:"script,block" json:"defined_scripts,omitempty"`
	DefinedCommands   []*Command             `hcl:"command,block" json:"defined_commands,omitempty"`
	DefinedFiles      []*RemoteFile          `hcl:"remote_file,block" json:"defined_files,omitempty"`
	DefinedDNSRecords []*DNSRecord           `hcl:"dns_record,block" json:"defined_dns_records,omitempty"`
	Hosts             map[string]*Host       `json:"-"`
	Networks          map[string]*Network    `json:"-"`
	Identities        map[string]*Identity   `json:"-"`
	Scripts           map[string]*Script     `json:"-"`
	Commands          map[string]*Command    `json:"-"`
	Files             map[string]*RemoteFile `json:"-"`
	DNSRecords        map[string]*DNSRecord  `json:"-"`
	Caller            Caller                 `json:"-"`
	ValidTeam         bool                   `json:"-"`
	ValidBuild        bool                   `json:"-"`
	ValidEnv          bool                   `json:"-"`
	ValidBase         bool                   `json:"-"`
	ValidGlobal       bool                   `json:"-"`
	ClearToBuild      bool                   `json:"-"`
	TeamRoot          string                 `json:"-"`
	BuildRoot         string                 `json:"-"`
	EnvRoot           string                 `json:"-"`
	BaseRoot          string                 `json:"-"`
	GlobalRoot        string                 `json:"-"`
	InitialContext    StateContext           `json:"-"`
	PathRegistry      *PathRegistry          `json:"-"`
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
	Do     string `hcl:"do,attr" json:"do,omitempty"`
	Append bool   `hcl:"append,attr" json:"append,omitempty"`
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
	l.Files = map[string]*RemoteFile{}
	l.DNSRecords = map[string]*DNSRecord{}
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
	for _, x := range l.DefinedFiles {
		err := x.ResolveSource(l, currPathResolver, l.Caller.Current())
		if err != nil {
			Logger.Errorf("%T %s had a source location that was not found: %v", x, x.ID, err)
		}
		l.Files[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.DefinedDNSRecords {
		l.DNSRecords[x.ID] = x
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
	if l.Competition == nil && diff.Competition != nil {
		l.Competition = diff.Competition
	} else if l.Competition != nil && diff.Competition != nil {
		res, err := SmartMerge(l.Competition, diff.Competition, false)
		if err != nil {
			return l, errors.WithStack(err)
		}
		orig, ok := res.(*Competition)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
		l.Competition = orig
	}
	if l.Environment == nil && diff.Environment != nil {
		l.Environment = diff.Environment
	} else if l.Environment != nil && diff.Environment != nil {
		res, err := SmartMerge(l.Environment, diff.Environment, true)
		if err != nil {
			return l, errors.WithStack(err)
		}
		orig, ok := res.(*Environment)
		if !ok {
			return nil, errors.WithStack(errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", orig, res))
		}
		l.Environment = orig
		if l.Environment.BaseDir == "" && l.EnvRoot != "" {
			l.Environment.BaseDir = l.EnvRoot
		}
	}
	if diff.Team != nil {
		l.Team = diff.Team
	}

	if diff.Build != nil {
		l.Build = diff.Build
	}

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
	for name, obj := range layer.Files {
		orig, found := base.Files[name]
		if !found {
			base.Files[name] = obj
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

// LoadFromContext attempts to bootstrap the given state from it's assumed Context Level
func (l *Laforge) LoadFromContext() error {
	var clone *Laforge
	var err error
	switch l.GetContext() {
	case TeamContext:
		clone, err = LoadFiles(l.GlobalConfigFile(), l.TeamConfigFile())
		if err != nil {
			return err
		}

		// if clone.Team != nil {
		// 	if l.PathRegistry == nil {
		// 		l.PathRegistry = &PathRegistry{
		// 			DB: map[CallFile]*PathResolver{},
		// 		}
		// 	}
		// 	if l.PathRegistry.DB[l.Caller.Current()] == nil {
		// 		l.PathRegistry.DB[l.Caller.Current()] = &PathResolver{
		// 			Mapping:    map[string]*LocalFileRef{},
		// 			Unresolved: map[string]bool{},
		// 		}
		// 	}
		// 	currPathResolver := l.PathRegistry.DB[l.Caller.Current()]
		// 	err = clone.Build.LoadDBFile(l, currPathResolver, l.Caller.Current())
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	case BuildContext:
		clone, err = LoadFiles(l.GlobalConfigFile(), l.BuildConfigFile())
		if err != nil {
			return err
		}

		if clone.Build != nil {
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
			err = clone.Build.LoadDBFile(l, currPathResolver, l.Caller.Current())
			if err != nil {
				return err
			}
		}
	case EnvContext:
		clone, err = LoadFiles(l.GlobalConfigFile(), l.EnvConfigFile())
		if err != nil {
			return err
		}
		if clone != nil && clone.Environment != nil {
			err = clone.IndexHostDependencies()
			if err != nil {
				return err
			}
			err = clone.Environment.ResolveIncludedNetworks(clone)
			if err != nil {
				return err
			}
		}
	case BaseContext:
		clone, err = LoadFiles(l.GlobalConfigFile(), l.BaseConfigFile())
		if err != nil {
			return err
		}
		if clone != nil {
			err = clone.IndexHostDependencies()
			if err != nil {
				return err
			}
		}
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
	if base.Competition != nil && base.Competition.BaseDir == "" {
		base.Competition.BaseDir = base.BaseRoot
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
