package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Script defines a configurable type for an executable script object within the laforge configuration
type Script struct {
	ID           string            `hcl:",label" json:"id,omitempty"`
	Name         string            `hcl:"name,attr" json:"name,omitempty"`
	Language     string            `hcl:"language,attr" json:"language,omitempty"`
	Description  string            `hcl:"description,attr" json:"description,omitempty"`
	Maintainer   *User             `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Source       string            `hcl:"source,attr" json:"source,omitempty"`
	SourceType   string            `hcl:"source_type,attr" json:"source_type,omitempty"`
	Cooldown     int               `hcl:"cooldown,attr" json:"cooldown,omitempty"`
	IgnoreErrors bool              `hcl:"ignore_errors,attr" json:"ignore_errors,omitempty"`
	Args         []string          `hcl:"args,attr" json:"args,omitempty"`
	IO           IO                `hcl:"io,block" json:"io,omitempty"`
	Disabled     bool              `hcl:"disabled,attr" json:"disabled,omitempty"`
	Vars         map[string]string `hcl:"vars,attr" json:"vars,omitempty"`
	Tags         map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	OnConflict   OnConflict        `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller       Caller            `json:"-"`
}

// GetCaller implements the Mergeable interface
func (s *Script) GetCaller() Caller {
	return s.Caller
}

// GetID implements the Mergeable interface
func (s *Script) GetID() string {
	return s.ID
}

// GetOnConflict implements the Mergeable interface
func (s *Script) GetOnConflict() OnConflict {
	return s.OnConflict
}

// SetCaller implements the Mergeable interface
func (s *Script) SetCaller(c Caller) {
	s.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (s *Script) SetOnConflict(o OnConflict) {
	s.OnConflict = o
}

// Kind implements the Provisioner interface
func (s *Script) Kind() string {
	return "script"
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
	for _, x := range s.Args {
		ret = append(ret, x)
	}
	return strings.Join(ret, " ")
}

// Base is a template helper function to return the base filename of a source script
func (s *Script) Base() string {
	return filepath.Base(s.Source)
}

// ResolveSource attempts to locate the referenced source file with a laforge base configuration
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
	pr.Mapping[s.Source] = lfr
	return nil
}
