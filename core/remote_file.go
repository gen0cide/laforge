package core

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// RemoteFile is a configurable type that defines a static file that will be placed on a configured target host.
type RemoteFile struct {
	ID          string            `hcl:"id,label" json:"id,omitempty"`
	SourceType  string            `hcl:"source_type,attr" json:"source_type,omitempty"`
	Source      string            `hcl:"source,attr" json:"source,omitempty"`
	Destination string            `hcl:"destination,attr" json:"destination,omitempty"`
	Vars        map[string]string `hcl:"vars,attr" json:"vars,omitempty"`
	Tags        map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	Template    bool              `hcl:"template,attr" json:"template,omitempty"`
	Perms       string            `hcl:"perms,attr" json:"perms,omitempty"`
	Disabled    bool              `hcl:"disabled,attr" json:"disabled,omitempty"`
	OnConflict  *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Checksum    string            `hcl:"md5,attr" json:"md5,omitempty"`
	Caller      Caller            `json:"-"`
	AbsPath     string            `json:"-"`
	Ext         string            `json:"-"`
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
	if r.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *r.OnConflict
}

// SetCaller implements the Mergeable interface
func (r *RemoteFile) SetCaller(c Caller) {
	r.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (r *RemoteFile) SetOnConflict(o OnConflict) {
	r.OnConflict = &o
}

// Kind implements the Provisioner interface
func (r *RemoteFile) Kind() string {
	return "remote_file"
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
	r.AbsPath = testSrc
	pr.Mapping[r.Source] = lfr
	return nil
}

// MD5Sum returns the MD5 checksum of a local file
func (r *RemoteFile) MD5Sum() (string, error) {
	f, err := os.Open(r.AbsPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// CopyTo copies the local file to another location on the local machine
func (r *RemoteFile) CopyTo(dst string) error {
	in, err := os.Open(r.AbsPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// AssetName returns the asset's name calculated as intended
func (r *RemoteFile) AssetName() (string, error) {
	if r.AbsPath == "" {
		return "", errors.New("no absolute path determined")
	}

	if r.Checksum == "" {
		cs, err := r.MD5Sum()
		if err != nil {
			return "", err
		}

		r.Checksum = cs
	}

	if r.Ext == "" {
		r.Ext = filepath.Ext(r.AbsPath)
	}

	return fmt.Sprintf("%s%s", r.Checksum, r.Ext), nil

}
