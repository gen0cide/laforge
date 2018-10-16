package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/pkg/errors"
)

// Build represents the output of a laforge build
//easyjson:json
type Build struct {
	ID           string            `hcl:"id,label" json:"id,omitempty"`
	TeamCount    int               `hcl:"team_count,attr" json:"team_count,omitempty"`
	Config       map[string]string `hcl:"config,attr" json:"config,omitempty"`
	Tags         map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	Maintainer   *User             `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Revision     int64             `hcl:"revision,attr" json:"revision,omitempty"`
	OnConflict   *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	PathFromBase string            `hcl:"path_from_base,optional" json:"path_from_base,omitempty"`
	Environment  *Environment      `json:"-"`
	Competition  *Competition      `json:"-"`
	RelEnvPath   string            `json:"-"`
	Dir          string            `json:"-"`
	Caller       Caller            `json:"-"`
	LocalDBFile  *LocalFileRef     `json:"-"`
	Teams        map[int]*Team     `json:"-"`
}

// Hash implements the Hasher interface
func (b *Build) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"teamcount=%v config=%v",
			b.TeamCount,
			b.Config,
		),
	)
}

// Path implements the Pather interface
func (b *Build) Path() string {
	return b.ID
}

// Base implements the Pather interface
func (b *Build) Base() string {
	return path.Base(b.ID)
}

// ValidatePath implements the Pather interface
func (b *Build) ValidatePath() error {
	if err := ValidateGenericPath(b.Path()); err != nil {
		return err
	}
	if topdir := strings.Split(b.Path(), `/`); topdir[1] != "envs" {
		return fmt.Errorf("path %s is not rooted in /%s", b.Path(), topdir)
	}
	return nil
}

// GetCaller implements the Mergeable interface
func (b *Build) GetCaller() Caller {
	return b.Caller
}

// LaforgeID implements the Mergeable interface
func (b *Build) LaforgeID() string {
	return b.ID
}

// ParentLaforgeID returns the build's parent environment ID
func (b *Build) ParentLaforgeID() string {
	return path.Dir(b.LaforgeID())
}

// GetOnConflict implements the Mergeable interface
func (b *Build) GetOnConflict() OnConflict {
	if b.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *b.OnConflict
}

// SetCaller implements the Mergeable interface
func (b *Build) SetCaller(ca Caller) {
	b.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (b *Build) SetOnConflict(o OnConflict) {
	b.OnConflict = &o
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

// SetID generates a unique build ID for this build
func (b *Build) SetID() string {
	if b.ID == "" {
		b.ID = path.Join(b.Environment.ID, b.Environment.Builder)
	}
	return b.ID
}

// AssetForTeam is a template helper function that returns the location of team specific assets
func (b *Build) AssetForTeam(teamID int, assetName string) string {
	return filepath.Join(b.Dir, fmt.Sprintf("%d", teamID), "assets", assetName)
}

// RelAssetForTeam is a template helper function that returns the relative location of team specific assets
func (b *Build) RelAssetForTeam(networkBase, hostBase, assetName string) string {
	return strings.Replace(filepath.Join(".", networkBase, hostBase, "assets", assetName), "\\", "/", -1)
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

	buildDir := filepath.Join(l.EnvRoot, l.CurrentEnv.Builder)
	buildDefPath := filepath.Join(buildDir, "build.laforge")
	bdbDir := filepath.Join(buildDir, "data")
	teamsDir := filepath.Join(buildDir, "teams")

	_, e0 := os.Stat(buildDir)
	_, e1 := os.Stat(buildDefPath)
	_, e2 := os.Stat(bdbDir)

	if e0 == nil || e1 == nil || e2 == nil {
		if !overwrite && !update {
			return fmt.Errorf("Cannot initialize build directory - path is dirty: %s (--force/-f to overwrite)", buildDir)
		}
		if !update {
			os.RemoveAll(buildDir)
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

	builder := l.CurrentEnv.Builder
	if builder == "" {
		builder = "null"
	}
	bid := builder

	relEnvPath, err := filepath.Rel(buildDir, filepath.Join(l.EnvRoot, "env.laforge"))

	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not get relative path of build directory %s to env root %s", buildDir, l.EnvRoot)
	}

	b := &Build{
		Dir:         buildDir,
		TeamCount:   l.CurrentEnv.TeamCount,
		Config:      l.CurrentEnv.Config,
		Tags:        l.CurrentEnv.Tags,
		Teams:       map[int]*Team{},
		Environment: l.CurrentEnv,
		Competition: l.CurrentEnv.Competition,
		Maintainer:  l.User,
		RelEnvPath:  relEnvPath,
	}

	b.SetID()

	bconfData, err := RenderHCLv2Object(b)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not generate build config for %s", bid)
	}

	err = ioutil.WriteFile(buildDefPath, bconfData, 0644)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not write build.laforge for build %s", bid)
	}

	l.CurrentBuild = b
	l.BuildContextID = b.LaforgeID()
	l.ClearToBuild = true
	return nil
}
