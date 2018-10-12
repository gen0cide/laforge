package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"

	"github.com/pkg/errors"
)

// Build represents the output of a laforge build
type Build struct {
	ID               string             `hcl:"id,label" json:"id,omitempty"`
	DBFile           string             `hcl:"db_file,attr" json:"-"`
	Builder          string             `hcl:"builder,attr" json:"builder,omitempty"`
	TeamCount        int                `hcl:"team_count,attr" json:"team_count,omitempty"`
	EnvironmentID    string             `hcl:"environment_id,attr" json:"environment_id,omitempty"`
	Config           map[string]string  `hcl:"config,attr" json:"config,omitempty"`
	Tags             map[string]string  `hcl:"tags,attr" json:"tags,omitempty"`
	Maintainer       *User              `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Networks         []*IncludedNetwork `hcl:"included_network,block" json:"included_networks,omitempty"`
	EnvironmentCache *Environment       `json:"environment_cache,omitempty"`
	RelEnvPath       string             `json:"rel_env_path"`
	Dir              string             `json:"-"`
	Caller           Caller             `json:"-"`
	LocalDBFile      *LocalFileRef      `json:"-"`
	Teams            map[int]*Team      `json:"teams,omitempty"`
	OnConflict       OnConflict         `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
}

// AssetForTeam is a template helper function that returns the location of team specific assets
func (b *Build) AssetForTeam(teamID int, assetName string) string {
	return filepath.Join(b.Dir, fmt.Sprintf("%d", teamID), "assets", assetName)
}

// RelAssetForTeam is a template helper function that returns the relative location of team specific assets
func (b *Build) RelAssetForTeam(teamID int, hostname, assetName string) string {
	return strings.Replace(filepath.Join(".", hostname, "assets", assetName), "\\", "/", -1)
}

// MergeFromDB loads and merges the build's DB file into the current build object
func (b *Build) MergeFromDB() error {
	if b.LocalDBFile == nil {
		return errors.Errorf("could not resolve db_file for build %s", b.ID)
	}
	data, err := ioutil.ReadFile(b.LocalDBFile.AbsPath)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "error reading db_file for build %s", b.ID)
	}

	var newBuild *Build

	err = json.Unmarshal(data, newBuild)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "error parsing db_file JSON for build %s", b.ID)
	}

	err = mergo.Merge(b, newBuild, mergo.WithOverride)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "errors merging db_file into the state for build %s", b.ID)
	}
	return nil
}

// LoadDBFile attempts to locate the referenced source file with a laforge base configuration
func (b *Build) LoadDBFile(base *Laforge, pr *PathResolver, caller CallFile) error {
	if b.DBFile == "" {
		return errors.New("no build db_file defined")
	}
	cwd, _ := os.Getwd()
	testSrc := b.DBFile
	if !filepath.IsAbs(b.DBFile) {
		testSrc = filepath.Join(caller.CallerDir, b.DBFile)
	}
	if !PathExists(testSrc) {
		pr.Unresolved[b.DBFile] = true
		return errors.Wrapf(ErrAbsPathDeclNotExist, "caller=%s path=%s", caller.CallerFile, b.DBFile)
	}
	rel, _ := filepath.Rel(cwd, testSrc)
	rel2, _ := filepath.Rel(caller.CallerDir, testSrc)
	lfr := &LocalFileRef{
		Base:          filepath.Base(testSrc),
		AbsPath:       testSrc,
		RelPath:       rel,
		Cwd:           cwd,
		DeclaredPath:  b.DBFile,
		RelToCallFile: rel2,
	}
	b.LocalDBFile = lfr
	pr.Mapping[b.DBFile] = lfr
	return b.MergeFromDB()
}

// InitializeBuildDirectory creates a build directory structure and writes the build.db as a precursor to builder's taking over.
func InitializeBuildDirectory(l *Laforge, overwrite, update bool) error {
	err := l.AssertExactContext(EnvContext)
	if err != nil && !overwrite && !update {
		return errors.WithStack(err)
	}
	err = l.AssertMinContext(EnvContext)
	if err != nil {
		return errors.WithStack(err)
	}

	buildDir := filepath.Join(l.EnvRoot, "build")
	buildDefPath := filepath.Join(buildDir, "build.laforge")
	bdbDir := filepath.Join(buildDir, "data")
	bdbDefPath := filepath.Join(bdbDir, "build.db")
	teamsDir := filepath.Join(buildDir, "teams")

	_, e0 := os.Stat(buildDir)
	_, e1 := os.Stat(buildDefPath)
	_, e2 := os.Stat(bdbDir)
	_, e3 := os.Stat(bdbDefPath)

	if e0 == nil || e1 == nil || e2 == nil || e3 == nil {
		if !overwrite && !update {
			return fmt.Errorf("Cannot initialize build directory - path is dirty: %s (--force/-f to overwrite)", buildDir)
		}
		if !update {
			os.RemoveAll(buildDir)
		}
	}

	if update {
		clone, err := LoadFiles(l.GlobalConfigFile(), l.EnvConfigFile())
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
		l, err = Mask(l, clone)
		if err != nil {
			return err
		}
	}

	dirs := []string{buildDir, bdbDir, teamsDir}
	for _, d := range dirs {
		os.MkdirAll(d, 0755)
		err = TouchGitKeep(d)
		if err != nil {
			return err
		}
	}

	builder := l.Environment.Builder
	if builder == "" {
		builder = "default"
	}
	bid := fmt.Sprintf("%s_%s", l.Environment.ID, builder)

	relEnvPath, err := filepath.Rel(buildDir, filepath.Join(l.EnvRoot, "env.laforge"))

	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not get relative path of build directory %s to env root %s", buildDir, l.EnvRoot)
	}

	b := &Build{
		ID:               bid,
		Dir:              buildDir,
		DBFile:           "./data/build.db",
		Builder:          builder,
		TeamCount:        l.Environment.TeamCount,
		EnvironmentID:    l.Environment.ID,
		Config:           l.Environment.Config,
		Tags:             l.Environment.Tags,
		Maintainer:       &l.User,
		Networks:         l.Environment.Networks,
		EnvironmentCache: l.Environment,
		RelEnvPath:       relEnvPath,
		Teams:            map[int]*Team{},
	}

	jsonData, err := json.Marshal(b)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not generate build.db for build %s", bid)
	}

	err = ioutil.WriteFile(bdbDefPath, jsonData, 0644)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not write build.db for build %s", bid)
	}

	bconfData, err := RenderHCLv2Object(b)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not generate build config for %s", bid)
	}

	err = ioutil.WriteFile(buildDefPath, bconfData, 0644)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not write build.laforge for build %s", bid)
	}

	l.Build = b
	l.ClearToBuild = true
	return nil
}

// GetCaller implements the Mergeable interface
func (b *Build) GetCaller() Caller {
	return b.Caller
}

// GetID implements the Mergeable interface
func (b *Build) GetID() string {
	return b.ID
}

// GetOnConflict implements the Mergeable interface
func (b *Build) GetOnConflict() OnConflict {
	return b.OnConflict
}

// SetCaller implements the Mergeable interface
func (b *Build) SetCaller(ca Caller) {
	b.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (b *Build) SetOnConflict(o OnConflict) {
	b.OnConflict = o
}

// Swap implements the Mergeable interface
func (b *Build) Swap(m Mergeable) error {
	rawVal, ok := m.(*Build)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", b, m)
	}
	*b = *rawVal
	return nil
}
