package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

var (
	// ValidEnvNameRegexp is a regular expression that can be used to validate environment names
	ValidEnvNameRegexp = regexp.MustCompile(`\A[a-z0-9][a-z0-9\-]*?[a-z0-9]\z`)
)

// Environment represents the basic configurable type for a Laforge environment container
type Environment struct {
	ID               string              `hcl:"id,label" json:"id,omitempty"`
	CompetitionID    string              `hcl:"competition_id,attr" json:"competition_id,omitempty"`
	Name             string              `hcl:"name,attr" json:"name,omitempty"`
	Description      string              `hcl:"description,attr" json:"description,omitempty"`
	Builder          string              `hcl:"builder,attr" json:"builder,omitempty"`
	TeamCount        int                 `hcl:"team_count,attr" json:"team_count,omitempty"`
	AdminCIDRs       []string            `hcl:"admin_ranges,attr" json:"admin_ranges,omitempty"`
	Config           map[string]string   `hcl:"config,optional" json:"config,omitempty"`
	Tags             map[string]string   `hcl:"tags,optional" json:"tags,omitempty"`
	Networks         []*IncludedNetwork  `hcl:"included_network,block" json:"included_networks,omitempty"`
	Maintainer       *User               `hcl:"maintainer,block" json:"maintainer,omitempty"`
	OnConflict       *OnConflict         `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	BaseDir          string              `hcl:"base_dir,optional" json:"base_dir,omitempty"`
	Revision         int64               `hcl:"revision,optional" json:"revision,omitempty"`
	IncludedNetworks map[string]*Network `json:"-"`
	IncludedHosts    map[string]*Host    `json:"-"`
	HostByNetwork    map[string][]*Host  `json:"-"`
	Teams            map[int]*Team       `json:"-"`
	Caller           Caller              `json:"-"`
	Competition      *Competition        `json:"-"`
}

// GetCaller implements the Mergeable interface
func (e *Environment) GetCaller() Caller {
	return e.Caller
}

// GetID implements the Mergeable interface
func (e *Environment) GetID() string {
	return filepath.Join(e.CompetitionID, e.ID)
}

// GetParentID returns the Team's parent build ID
func (e *Environment) GetParentID() string {
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

// ValidEnvName is a helper function to determine if a supplied name is a valid environment name
func ValidEnvName(name string) bool {
	if len(name) > 16 {
		return false
	}
	return ValidEnvNameRegexp.MatchString(name)
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

	envDir := filepath.Join(l.BaseRoot, "envs", name)
	envDefPath := filepath.Join(envDir, "env.laforge")

	_, dirErr := os.Stat(envDir)
	_, defErr := os.Stat(envDefPath)

	if defErr == nil || dirErr == nil {
		if !overwrite {
			return fmt.Errorf("Cannot initialize env directory - path is dirty: %s (--force/-f to overwrite)", envDir)
		}
		os.RemoveAll(envDir)
	}

	os.MkdirAll(envDir, 0755)
	keeper := filepath.Join(envDir, ".gitkeep")
	newFile, err := os.Create(keeper)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("cannot touch .gitkeep inside env directory subfolder %s", envDir))
	}
	newFile.Close()

	envData, err := RenderHCLv2Object(baseEnvironment(name, &l.User))
	if err != nil {
		return errors.WithStack(err)
	}

	return ioutil.WriteFile(envDefPath, envData, 0644)
}

// CleanBuildDirectory attempts to purge the current environment context's build directory (rm -r basically)
func (l *Laforge) CleanBuildDirectory(overwrite bool) error {
	pn := ""
	for cf := range l.PathRegistry.DB {
		if filepath.Base(cf.CallerFile) == "env.laforge" {
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

// GetAllEnvs recursively traverses the BaseRoot/envs/ folder looking for valid environments.
func (l *Laforge) GetAllEnvs() (map[string]*Laforge, error) {
	emap := map[string]*Laforge{}
	err := l.AssertMinContext(BaseContext)
	if err != nil {
		return emap, errors.WithStack(err)
	}

	basePath := filepath.Join(l.BaseRoot, "envs")

	envConfigs := []string{}

	err = godirwalk.Walk(basePath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.Name() == "env.laforge" {
				envConfigs = append(envConfigs, osPathname)
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			Logger.Debugf("Envwalker encountered a filesystem issue at %s: %v", osPathname, err)
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
			return
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
			// if res.Environment == nil {
			// 	Logger.Errorf("Nil environment found during directory walk... (this should of errored, but didn't)")
			// 	continue
			// }
			// if emap[res.Environment.ID] != nil {
			// 	origPath := ""
			// 	badPath := ""
			// 	for cf := range emap[res.Environment.ID].PathRegistry.DB {
			// 		if filepath.Base(cf.CallerFile) == "env.laforge" {
			// 			origPath = cf.CallerDir
			// 			break
			// 		}
			// 	}
			// 	for cf := range res.PathRegistry.DB {
			// 		if filepath.Base(cf.CallerFile) == "env.laforge" {
			// 			badPath = cf.CallerDir
			// 			break
			// 		}
			// 	}
			// 	Logger.Errorf("Name collision between two environments! Check env.laforge environment IDs in these directories:\n  %s\n  %s", origPath, badPath)
			// 	continue
			// }
			// emap[res.Environment.ID] = res
		case resErr := <-errChan:
			Logger.Errorf("An error was found with an environment: %v", resErr)
		case <-finChan:
			return emap, nil
		}
	}
}
