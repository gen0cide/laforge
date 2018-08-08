package core

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// RemoteFile is a configurable type that defines a static file that will be placed on a configured target host.
type RemoteFile struct {
	ID          string     `hcl:",label" json:"id,omitempty"`
	SourceType  string     `hcl:"source_type,attr" json:"source_type,omitempty"`
	Source      string     `hcl:"source_path,attr" json:"source_path,omitempty"`
	Destination string     `hcl:"destination,label" json:"destination,omitempty"`
	Perms       string     `hcl:"perms,attr" json:"perms,omitempty"`
	Disabled    bool       `hcl:"disabled,attr" json:"disabled,omitempty"`
	OnConflict  OnConflict `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller      Caller     `json:"-"`
}

// GetCaller implements the Mergeable interface
func (r *RemoteFile) GetCaller() Caller {
	return r.Caller
}

// GetID implements the Mergeable interface
func (r *RemoteFile) GetID() string {
	return r.ID
}

// GetOnConflict implements the Mergeable interface
func (r *RemoteFile) GetOnConflict() OnConflict {
	return r.OnConflict
}

// SetCaller implements the Mergeable interface
func (r *RemoteFile) SetCaller(c Caller) {
	r.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (r *RemoteFile) SetOnConflict(o OnConflict) {
	r.OnConflict = o
}

// Swap implements the Mergeable interface
func (r *RemoteFile) Swap(m Mergeable) error {
	rawVal, ok := m.(*RemoteFile)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", r, m)
	}
	*r = *rawVal
	return nil
}

// ResolveSource attempts to locate the referenced source file with a laforge base configuration
func (r *RemoteFile) ResolveSource(base *Laforge, pr *PathResolver, caller CallFile) error {
	if r.Source == "" {
		return nil
	}
	if r.SourceType != "" && r.SourceType != "local" {
		return nil
	}
	cwd, _ := os.Getwd()
	testSrc := r.Source
	if !filepath.IsAbs(r.Source) {
		testSrc = filepath.Join(caller.CallerDir, r.Source)
	}
	if !PathExists(testSrc) {
		pr.Unresolved[r.Source] = true
		return errors.Wrapf(ErrAbsPathDeclNotExist, "caller=%s path=%s", caller.CallerFile, r.Source)
	}
	rel, _ := filepath.Rel(cwd, testSrc)
	rel2, _ := filepath.Rel(caller.CallerDir, testSrc)
	lfr := &LocalFileRef{
		Base:          filepath.Base(testSrc),
		AbsPath:       testSrc,
		RelPath:       rel,
		Cwd:           cwd,
		DeclaredPath:  r.Source,
		RelToCallFile: rel2,
	}
	pr.Mapping[r.Source] = lfr
	return nil
}
