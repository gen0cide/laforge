package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/pkg/errors"
)

const (
	envsDir     = `envs`
	commandsDir = `commands`
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
	Teams        map[string]*Team  `json:"-"`
}

// HashConfigMap is used to hash the configuration map in a deterministic order
func HashConfigMap(m map[string]string) []uint64 {
	data := []uint64{}
	for k, v := range m {
		data = append(data, xxhash.Sum64String(k))
		data = append(data, xxhash.Sum64String(v))
	}
	sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })
	return data
}

// Hash implements the Hasher interface
func (b *Build) Hash() uint64 {
	chash := HashConfigMap(b.Config)
	return xxhash.Sum64String(
		fmt.Sprintf(
			"teamcount=%v config=%v",
			b.TeamCount,
			chash,
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
	if topdir := strings.Split(b.Path(), `/`); topdir[1] != envsDir {
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
	return strings.Replace(filepath.Join(".", "networks", networkBase, hostsDir, hostBase, "assets", assetName), "\\", "/", -1)
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
			err = os.RemoveAll(buildDir)
			if err != nil {
				return err
			}
		}
	}

	dirs := []string{buildDir, bdbDir, teamsDir}
	for _, d := range dirs {
		//nolint:gosec
		err := os.MkdirAll(d, 0755)
		if err != nil {
			return err
		}
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

	relEnvPath, err := filepath.Rel(buildDir, filepath.Join(l.EnvRoot, envFile))

	if err != nil {
		return errors.Wrapf(errors.WithStack(err), "could not get relative path of build directory %s to env root %s", buildDir, l.EnvRoot)
	}

	b := &Build{
		Dir:         buildDir,
		TeamCount:   l.CurrentEnv.TeamCount,
		Config:      l.CurrentEnv.Config,
		Tags:        l.CurrentEnv.Tags,
		Teams:       map[string]*Team{},
		Environment: l.CurrentEnv,
		Competition: l.CurrentEnv.Competition,
		Maintainer:  l.User,
		RelEnvPath:  relEnvPath,
	}

	l.CurrentEnv.Build = b

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
	dbfile := filepath.Join(buildDir, "build.db")

	state := NewState()
	state.Base = l
	l.StateManager = state

	snap, err := NewSnapshotFromEnv(l.CurrentEnv, false)
	if err != nil {
		return err
	}

	state.SetCurrent(snap)

	err = state.LocateRevisions()
	if err != nil {
		return err
	}
	err = state.GenerateCurrentRevs()
	if err != nil {
		return err
	}

	envRev := state.NewRevs[l.CurrentEnv.Path()]
	buildRev := state.NewRevs[l.CurrentBuild.Path()]

	err = ioutil.WriteFile(filepath.Join(buildDir, ".env.lfrevision"), []byte(envRev.ToJSONString()), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(buildDir, ".build.lfrevision"), []byte(buildRev.ToJSONString()), 0644)
	if err != nil {
		return err
	}

	err = state.Open(dbfile)
	if err != nil {
		return err
	}

	return nil
}

// CreateTeams enumerates the build's team count and generates children team objects
func (b *Build) CreateTeams() error {
	if len(b.Teams) != 0 {
		return errors.New("build already is populated with teams")
	}
	for i := 0; i < b.TeamCount; i++ {
		t := b.CreateTeam(i)
		err := t.CreateProvisionResources()
		if err != nil {
			return err
		}
	}
	return nil
}

// Gather implements the Dependency interface
func (b *Build) Gather(g *Snapshot) error {
	err := g.Relate(b.Environment, b)
	if err != nil {
		return err
	}
	for _, t := range b.Teams {
		err = g.Relate(b, t)
		if err != nil {
			return err
		}
		err = t.Gather(g)
		if err != nil {
			return err
		}
	}
	return nil
}

// Associate walks the build and creates edges in the graph
func (b *Build) Associate(g *Snapshot) error {
	for _, t := range b.Teams {
		err := t.Associate(g)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateTeam creates a new team of a given team index for the build.
func (b *Build) CreateTeam(tid int) *Team {
	t := &Team{
		TeamNumber:          tid,
		Build:               b,
		Environment:         b.Environment,
		Competition:         b.Competition,
		ProvisionedNetworks: map[string]*ProvisionedNetwork{},
	}

	b.Teams[t.SetID()] = t
	return t
}

// CreateBuildEntry ...
func (b *Build) CreateBuildEntry(ctx context.Context, client *ent.Client) (*ent.Build, error) {
	user, err := b.Maintainer.CreateUserEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating build: %v", err)
		return nil, err
	}

	tag, err := CreateTagEntry(b.ID, b.Tags, ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating build: %v", err)
		return nil, err
	}

	build, err := client.Build.
		Create().
		SetRevision(int(b.Revision)+1). // Maybe not out of range now?
		SetConfig(b.Config).
		AddMaintainer(user).
		AddTag(tag).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating build: %v", err)
		return nil, err
	}

	cli.Logger.Debugf("build was created: ", build)
	return build, nil
}
