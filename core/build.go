package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Build represents the output of a laforge build
type Build struct {
	ID            string            `hcl:"id,label" json:"id,omitempty"`
	TeamCount     int               `hcl:"team_count,attr" json:"team_count,omitempty"`
	EnvironmentID string            `hcl:"environment_id,attr" json:"environment_id,omitempty"`
	CompetitionID string            `hcl:"competition_id,attr" json:"competition_id,omitempty"`
	Config        map[string]string `hcl:"config,attr" json:"config,omitempty"`
	Tags          map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	Maintainer    *User             `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Revision      int64             `hcl:"revision,attr" json:"revision,omitempty"`
	OnConflict    *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Environment   *Environment      `json:"-"`
	Competition   *Competition      `json:"-"`
	RelEnvPath    string            `json:"-"`
	Dir           string            `json:"-"`
	Caller        Caller            `json:"-"`
	LocalDBFile   *LocalFileRef     `json:"-"`
	Teams         map[int]*Team     `json:"-"`
}

// GetCaller implements the Mergeable interface
func (b *Build) GetCaller() Caller {
	return b.Caller
}

// GetID implements the Mergeable interface
func (b *Build) GetID() string {
	return filepath.Join(b.CompetitionID, b.EnvironmentID, b.ID)
}

// GetParentID returns the build's parent environment ID
func (b *Build) GetParentID() string {
	return filepath.Join(b.CompetitionID, b.EnvironmentID)
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
	b.Revision++
	if b.EnvironmentID == "" && b.Environment != nil {
		b.EnvironmentID = b.Environment.ID
	}
	if b.CompetitionID == "" && b.Competition != nil {
		b.CompetitionID = b.Competition.ID
	}
	return b.ID
}

// AssetForTeam is a template helper function that returns the location of team specific assets
func (b *Build) AssetForTeam(teamID int, assetName string) string {
	return filepath.Join(b.Dir, fmt.Sprintf("%d", teamID), "assets", assetName)
}

// RelAssetForTeam is a template helper function that returns the relative location of team specific assets
func (b *Build) RelAssetForTeam(teamID int, hostname, assetName string) string {
	return strings.Replace(filepath.Join(".", hostname, "assets", assetName), "\\", "/", -1)
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

	// if update {
	// 	clone, err := InitializeEnvContext(l.GlobalConfigFile(), l.EnvConfigFile())
	// 	if err != nil {
	// 		return err
	// 	}

	// 	l, err = Mask(l, clone)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	dirs := []string{buildDir, bdbDir, teamsDir}
	for _, d := range dirs {
		os.MkdirAll(d, 0755)
		err = TouchGitKeep(d)
		if err != nil {
			return err
		}
	}

	// builder := l.Environments.Builder
	builder := l.CurrentEnv.Builder
	if builder == "" {
		builder = "default"
	}
	bid := builder

	relEnvPath, err := filepath.Rel(buildDir, filepath.Join(l.EnvRoot, "env.laforge"))

	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not get relative path of build directory %s to env root %s", buildDir, l.EnvRoot)
	}

	b := &Build{
		ID:            bid,
		Dir:           buildDir,
		TeamCount:     l.CurrentEnv.TeamCount,
		EnvironmentID: l.CurrentEnv.ID,
		CompetitionID: l.CurrentEnv.CompetitionID,
		Config:        l.CurrentEnv.Config,
		Tags:          l.CurrentEnv.Tags,
		Teams:         map[int]*Team{},
		Environment:   l.CurrentEnv,
		Competition:   l.CurrentEnv.Competition,
		Maintainer:    &l.User,
		RelEnvPath:    relEnvPath,
	}

	bconfData, err := RenderHCLv2Object(b)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not generate build config for %s", bid)
	}

	err = ioutil.WriteFile(buildDefPath, bconfData, 0644)
	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not write build.laforge for build %s", bid)
	}

	l.CurrentBuild = b
	l.BuildContextID = b.GetID()
	l.ClearToBuild = true
	return nil
}
