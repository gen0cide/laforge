package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/pkg/errors"
)

// Script defines a configurable type for an executable script object within the laforge configuration
//easyjson:json
//nolint:maligned
type Script struct {
	ID           string            `hcl:"id,label" json:"id,omitempty"`
	Name         string            `hcl:"name,attr" json:"name,omitempty"`
	Language     string            `hcl:"language,attr" json:"language,omitempty"`
	Description  string            `hcl:"description,optional" json:"description,omitempty"`
	Maintainer   *User             `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Source       string            `hcl:"source,attr" json:"source,omitempty"`
	SourceType   string            `hcl:"source_type,attr" json:"source_type,omitempty"`
	Cooldown     int               `hcl:"cooldown,optional" json:"cooldown,omitempty"`
	Timeout      int               `hcl:"timeout,optional" json:"timeout,omitempty"`
	IgnoreErrors bool              `hcl:"ignore_errors,optional" json:"ignore_errors,omitempty"`
	Args         []string          `hcl:"args,optional" json:"args,omitempty"`
	IO           *IO               `hcl:"io,block" json:"io,omitempty"`
	Disabled     bool              `hcl:"disabled,optional" json:"disabled,omitempty"`
	Vars         map[string]string `hcl:"vars,optional" json:"vars,omitempty"`
	Tags         map[string]string `hcl:"tags,optional" json:"tags,omitempty"`
	OnConflict   *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Findings     []*Finding        `hcl:"finding,block" json:"findings,omitempty"`
	AbsPath      string            `json:"-"`
	Caller       Caller            `json:"-"`
}

// Hash implements the Hasher interface
func (s *Script) Hash() uint64 {
	iostr := "n/a"
	if s.IO != nil {
		iostr = s.IO.Stderr + s.IO.Stdin + s.IO.Stdout
	}

	return xxhash.Sum64String(
		fmt.Sprintf(
			"language=%v sourcetype=%v cooldown=%v ignoreerrors=%v args=%v io=%v disabled=%v vars=%v source=%v",
			s.Language,
			s.SourceType,
			s.Cooldown,
			s.IgnoreErrors,
			strings.Join(s.Args, `,`),
			iostr,
			s.Disabled,
			s.Vars,
			s.ResourceHash(),
		),
	)
}

// Path implements the Pather interface
func (s *Script) Path() string {
	return s.ID
}

// Base implements the Pather interface
func (s *Script) Base() string {
	return path.Base(s.ID)
}

// ValidatePath implements the Pather interface
func (s *Script) ValidatePath() error {
	if err := ValidateGenericPath(s.Path()); err != nil {
		return err
	}
	if topdir := strings.Split(s.Path(), `/`); topdir[1] != "scripts" {
		return fmt.Errorf("path %s is not rooted in /%s", s.Path(), topdir[1])
	}
	return nil
}

// ResourceHash implements the ResourceHasher interface
func (s *Script) ResourceHash() uint64 {
	dep, err := ioutil.ReadFile(s.AbsPath)
	if err != nil {
		fmt.Printf("dependency error for %s: %s could not be read: %v", s.Path(), s.AbsPath, err)
		return 666
	}
	return xxhash.Sum64(dep)
}

// GetCaller implements the Mergeable interface
func (s *Script) GetCaller() Caller {
	return s.Caller
}

// LaforgeID implements the Mergeable interface
func (s *Script) LaforgeID() string {
	return s.ID
}

// ParentLaforgeID implements the Dependency interface
func (s *Script) ParentLaforgeID() string {
	return s.Path()
}

// Gather implements the Dependency interface
func (s *Script) Gather(g *Snapshot) error {
	return nil
}

// Fullpath implements the Pather interface
func (s *Script) Fullpath() string {
	return s.LaforgeID()
}

// GetOnConflict implements the Mergeable interface
func (s *Script) GetOnConflict() OnConflict {
	if s.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *s.OnConflict
}

// SetCaller implements the Mergeable interface
func (s *Script) SetCaller(c Caller) {
	s.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (s *Script) SetOnConflict(o OnConflict) {
	s.OnConflict = &o
}

// Kind implements the Provisioner interface
func (s *Script) Kind() string {
	return ObjectTypeScript.String()
}

// Swap implements the Mergeable interface
func (s *Script) Swap(m Mergeable) error {
	rawVal, ok := m.(*Script)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", s, m)
	}
	*s = *rawVal
	return nil
}

// ArgString is a template helper function to embed the arg string into the output
func (s *Script) ArgString() string {
	if len(s.Args) == 0 {
		return ""
	}
	ret := []string{" "}
	ret = append(ret, s.Args...)
	return strings.Join(ret, " ")
}

// SourceBase is a template helper function to return the base filename of a source script
func (s *Script) SourceBase() string {
	return filepath.Base(s.Source)
}

// ResolveSource attempts to locate the referenced source file with a laforge base configuration
//nolint:dupl
func (s *Script) ResolveSource(base *Laforge, pr *PathResolver, caller CallFile) error {
	if s.Source == "" {
		return nil
	}
	if s.SourceType != "" && s.SourceType != "local" {
		return nil
	}
	cwd, _ := os.Getwd()
	testSrc := s.Source
	if !filepath.IsAbs(s.Source) {
		testSrc = filepath.Join(caller.CallerDir, s.Source)
	}
	if !PathExists(testSrc) {
		pr.Unresolved[s.Source] = true
		return errors.Wrapf(ErrAbsPathDeclNotExist, "caller=%s path=%s", caller.CallerFile, s.Source)
	}
	rel, _ := filepath.Rel(cwd, testSrc)
	rel2, _ := filepath.Rel(caller.CallerDir, testSrc)
	lfr := &LocalFileRef{
		Base:          filepath.Base(testSrc),
		AbsPath:       testSrc,
		RelPath:       rel,
		Cwd:           cwd,
		DeclaredPath:  s.Source,
		RelToCallFile: rel2,
	}
	s.AbsPath = testSrc
	pr.Mapping[s.Source] = lfr
	return nil
}

// CreateScriptEntry ...
func (s *Script) CreateScriptEntry(ph *ent.ProvisionedHost, ctx context.Context, client *ent.Client) (*ent.Script, error) {
	tag, err := CreateTagEntry(s.ID, s.Tags, ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating script: %v", err)
		return nil, err
	}

	user, err := s.Maintainer.CreateUserEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating script: %v", err)
		return nil, err
	}

	script, err := client.Script.
		Create().
		SetName(s.Name).
		SetLanguage(s.Language).
		SetDescription(s.Description).
		SetSource(s.Source).
		SetSourceType(s.SourceType).
		SetCooldown(s.Cooldown).
		SetTimeout(s.Timeout).
		SetIgnoreErrors(s.IgnoreErrors).
		SetArgs(s.Args).
		SetDisabled(s.Disabled).
		SetVars(s.Vars).
		SetAbsPath(s.AbsPath).
		AddTag(tag).
		AddMaintainer(user).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating script: %v", err)
		return nil, err
	}

	for _, v := range s.Findings {
		_, err := v.CreateFindingEntry(ph, script, ctx, client)

		if err != nil {
			cli.Logger.Debugf("failed creating script: %v", err)
			return nil, err
		}
	}

	cli.Logger.Debugf("script was created: ", script)
	return script, nil
}
