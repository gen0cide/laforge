package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

const (
	envFile = `env.laforge`
)

var (
	// ValidEnvNameRegexp is a regular expression that can be used to validate environment names
	ValidEnvNameRegexp = regexp.MustCompile(`\A[a-z0-9][a-z0-9\-]*?[a-z0-9]\z`)
)

// Environment represents the basic configurable type for a Laforge environment container
//easyjson:json
type Environment struct {
	ID               string              `hcl:"id,label" json:"id,omitempty"`
	CompetitionID    string              `hcl:"competition_id,attr" json:"competition_id,omitempty"`
	Name             string              `hcl:"name,attr" json:"name,omitempty"`
	Description      string              `hcl:"description,attr" json:"description,omitempty"`
	Builder          string              `hcl:"builder,attr" json:"builder,omitempty"`
	TeamCount        int                 `hcl:"team_count,attr" json:"team_count,omitempty"`
	AdminCIDRs       []string            `hcl:"admin_ranges,attr" json:"admin_ranges,omitempty"`
	ExposedVDIPorts  []string            `hcl:"vdi_allowed_tcp_ports" json:"vdi_allowed_tcp_ports,omitempty"`
	Config           map[string]string   `hcl:"config,optional" json:"config,omitempty"`
	Tags             map[string]string   `hcl:"tags,optional" json:"tags,omitempty"`
	Networks         []*IncludedNetwork  `hcl:"included_network,block" json:"included_networks,omitempty"`
	Maintainer       *User               `hcl:"maintainer,block" json:"maintainer,omitempty"`
	OnConflict       *OnConflict         `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	BaseDir          string              `hcl:"base_dir,optional" json:"base_dir,omitempty"`
	Revision         int64               `hcl:"revision,optional" json:"revision,omitempty"`
	Build            *Build              `json:"-"`
	IncludedNetworks map[string]*Network `json:"-"`
	IncludedHosts    map[string]*Host    `json:"-"`
	HostByNetwork    map[string][]*Host  `json:"-"`
	Teams            map[string]*Team    `json:"-"`
	Caller           Caller              `json:"-"`
	Competition      *Competition        `json:"-"`
}

// Hash implements the Hasher interface
func (e *Environment) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"name=%v builder=%v tc=%v acidrs=%v conf=%v",
			e.Name,
			e.Builder,
			e.TeamCount,
			strings.Join(e.AdminCIDRs, ","),
			HashConfigMap(e.Config),
		),
	)
}

// Path implements the Pather interface
func (e *Environment) Path() string {
	return e.ID
}

// Base implements the Pather interface
func (e *Environment) Base() string {
	return path.Base(e.ID)
}

// ValidatePath implements the Pather interface
func (e *Environment) ValidatePath() error {
	if err := ValidateGenericPath(e.Path()); err != nil {
		return err
	}
	if topdir := strings.Split(e.Path(), `/`); topdir[1] != envsDir {
		return fmt.Errorf("path %s is not rooted in /%s", e.Path(), topdir[1])
	}
	return nil
}

// GetCaller implements the Mergeable interface
func (e *Environment) GetCaller() Caller {
	return e.Caller
}

// LaforgeID implements the Mergeable interface
func (e *Environment) LaforgeID() string {
	return e.ID
}

// ParentLaforgeID returns the Team's parent build ID
func (e *Environment) ParentLaforgeID() string {
	return e.CompetitionID
}

// GetOnConflict implements the Mergeable interface
func (e *Environment) GetOnConflict() OnConflict {
	if e.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *e.OnConflict
}

// SetCaller implements the Mergeable interface
func (e *Environment) SetCaller(c Caller) {
	e.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (e *Environment) SetOnConflict(o OnConflict) {
	e.OnConflict = &o
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
		inet[n.Name] = ObjectTypeIncluded.String()
		e.HostByNetwork[n.Name] = []*Host{}
		for _, h := range n.Hosts {
			ihost[h] = ObjectTypeIncluded.String()
		}
	}
	for name, net := range base.Networks {
		status, found := inet[name]
		if !found {
			cli.Logger.Debugf("Skipping network %s", name)
			continue
		}
		if status == ObjectTypeIncluded.String() {
			e.IncludedNetworks[name] = net
			inet[name] = "resolved"
			cli.Logger.Debugf("Resolved network %s", name)
		}
	}
	for name, host := range base.Hosts {
		status, found := ihost[name]
		if !found {
			cli.Logger.Debugf("Skipping host %s", name)
			continue
		}
		if status == ObjectTypeIncluded.String() {
			e.IncludedHosts[name] = host
			ihost[name] = "resolved"
			cli.Logger.Debugf("Resolved host %s", name)
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
		if status == ObjectTypeIncluded.String() {
			return fmt.Errorf("no configuration for network %s", net)
		}
	}
	for host, status := range ihost {
		if status == ObjectTypeIncluded.String() {
			return fmt.Errorf("no configuration for host %s", host)
		}
	}
	return nil
}

// ValidEnvName is a helper function to determine if a supplied name is a valid environment name
func ValidEnvName(name string) bool {
	if len(name) > 16 {
		return false
	}
	return ValidEnvNameRegexp.MatchString(name)
}

// ValidID checks to determine if the string conforms to laforge ID schemas.
func ValidID(id string) bool {
	return ValidIDRegexp.MatchString(id)
}

// OutdatedID checks to determine if the ID is in an outdated format.
func OutdatedID(id string) bool {
	return ValidOldIDRegexp.MatchString(id)
}

// InitializeEnv attempts to initialize a new environment of a given name
func (l *Laforge) InitializeEnv(name string, overwrite bool) error {
	err := l.AssertExactContext(BaseContext)
	if err != nil && !overwrite {
		return errors.WithStack(err)
	}

	if !ValidEnvName(name) {
		return errors.WithStack(ErrInvalidEnvName)
	}

	envDir := filepath.Join(l.BaseRoot, envsDir, name)
	envDefPath := filepath.Join(envDir, envFile)

	_, dirErr := os.Stat(envDir)
	_, defErr := os.Stat(envDefPath)

	if defErr == nil || dirErr == nil {
		if !overwrite {
			return fmt.Errorf("Cannot initialize env directory - path is dirty: %s (--force/-f to overwrite)", envDir)
		}
		//nolint:errcheck,gosec
		os.RemoveAll(envDir)
	}

	//nolint:errcheck,gosec
	os.MkdirAll(envDir, 0755)
	keeper := filepath.Join(envDir, ".gitkeep")
	newFile, err := os.Create(keeper)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("cannot touch .gitkeep inside env directory subfolder %s", envDir))
	}

	//nolint:errcheck,gosec
	newFile.Close()

	envData, err := RenderHCLv2Object(baseEnvironment(name, l.User))
	if err != nil {
		return errors.WithStack(err)
	}

	return ioutil.WriteFile(envDefPath, envData, 0644)
}

// CleanBuildDirectory attempts to purge the current environment context's build directory (rm -r basically)
func (l *Laforge) CleanBuildDirectory(overwrite bool) error {
	pn := ""
	for cf := range l.PathRegistry.DB {
		if filepath.Base(cf.CallerFile) == envFile {
			pn = cf.CallerDir
			break
		}
	}

	if pn == "" {
		return fmt.Errorf("strange things have happened trying to find the env.laforge path")
	}

	buildDir := filepath.Join(pn, "build")
	_, e0 := os.Stat(buildDir)
	if e0 == nil {
		if !overwrite {
			return fmt.Errorf("Cannot clean build directory - path is dirty: %s (--force/-f to overwrite)", buildDir)
		}
	}

	return os.RemoveAll(buildDir)
}

// CreateBuild creates a fresh build object that is a child of the environment e.
func (e *Environment) CreateBuild() *Build {
	b := &Build{
		Dir:         filepath.Join(e.Caller.Current().CallerDir, e.Builder),
		TeamCount:   e.TeamCount,
		Config:      e.Config,
		Tags:        e.Tags,
		Teams:       map[string]*Team{},
		Environment: e,
		Competition: e.Competition,
		Maintainer:  e.Maintainer,
	}
	b.SetID()
	return b
}

// Gather implements the Dependency interface
func (e *Environment) Gather(g *Snapshot) error {
	// g.AddNode(e)
	// for nid, net := range e.IncludedNetworks {
	// 	g.AddNode(net)
	// 	for _, host := range e.HostByNetwork[nid] {
	// 		g.AddNode(host)
	// 		for _, x := range host.Scripts {
	// 			g.AddNode(x)
	// 		}
	// 		for _, x := range host.DNSRecords {
	// 			g.AddNode(x)
	// 		}
	// 		for _, x := range host.Commands {
	// 			g.AddNode(x)
	// 		}
	// 		for _, x := range host.RemoteFiles {
	// 			g.AddNode(x)
	// 		}
	// 	}

	// }
	// return e.Build.Gather(g)
	return nil
}

// GetAllEnvs recursively traverses the BaseRoot/envs/ folder looking for valid environments.
func (l *Laforge) GetAllEnvs() (map[string]*Laforge, error) {
	emap := map[string]*Laforge{}
	err := l.AssertMinContext(BaseContext)
	if err != nil {
		return emap, errors.WithStack(err)
	}

	basePath := filepath.Join(l.BaseRoot, envsDir)

	envConfigs := []string{}

	err = godirwalk.Walk(basePath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.Name() == envFile {
				envConfigs = append(envConfigs, osPathname)
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			cli.Logger.Debugf("Envwalker encountered a filesystem issue at %s: %v", osPathname, err)
			return godirwalk.SkipNode
		},
	})
	if err != nil {
		return emap, err
	}

	wg := new(sync.WaitGroup)
	resChan := make(chan *Laforge, len(envConfigs))
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	baseCfg := filepath.Join(l.BaseRoot, "base.laforge")

	for _, p := range envConfigs {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			lf, err := LoadFiles(baseCfg, s)
			if err != nil {
				errChan <- errors.Wrapf(err, "issue parsing environment for %s", s)
				return
			}
			resChan <- lf
		}(p)
	}

	go func() {
		wg.Wait()
		finChan <- true
	}()

	for {
		select {
		case res := <-resChan:
			_ = res
		case resErr := <-errChan:
			cli.Logger.Errorf("An error was found with an environment: %v", resErr)
		case <-finChan:
			return emap, nil
		}
	}
}

// CreateEnvironmentEntry ...
func (e *Environment) CreateEnvironmentEntry(ctx context.Context, client *ent.Client) (*ent.Environment, error) {
	tag, err := CreateTagEntry(e.ID, e.Tags, ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating environment: %v", err)
		return nil, err
	}

	user, err := e.Maintainer.CreateUserEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating environment: %v", err)
		return nil, err
	}

	build, err := e.Build.CreateBuildEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating environment: %v", err)
		return nil, err
	}

	competition, err := e.Competition.CreateCompetitionEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating environment: %v", err)
		return nil, err
	}

	environment, err := client.Environment.
		Create().
		SetCompetitionID(e.CompetitionID).
		SetName(e.Name).
		SetDescription(e.Description).
		SetBuilder(e.Builder).
		SetTeamCount(e.TeamCount).
		SetRevision(int(e.Revision)).
		SetAdminCidrs(e.AdminCIDRs).
		SetExposedVdiPorts(e.ExposedVDIPorts).
		SetConfig(e.Config).
		AddTag(tag).
		AddUser(user).
		AddBuild(build).
		AddCompetition(competition).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating environment: %v", err)
		return nil, err
	}

	for _, v := range e.Build.Teams {
		_, err := v.CreateTeamEntry(environment, build, ctx, client)

		if err != nil {
			cli.Logger.Debugf("failed creating environment: %v", err)
			return nil, err
		}
	}

	cli.Logger.Debugf("environment was created: ", environment)
	return environment, nil
}
