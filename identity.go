package laforge

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Identity defines a generic human identity primative that can be extended into Employee, Customer, Client, etc.
type Identity struct {
	ID          string            `hcl:",label" json:"id,omitempty"`
	Firstname   string            `hcl:"firstname,attr" json:"firstname,omitempty"`
	Lastname    string            `hcl:"lastname,attr" json:"lastname,omitempty"`
	Email       string            `hcl:"email,attr" json:"email,omitempty"`
	Password    string            `hcl:"password,attr" json:"password,omitempty"`
	Description string            `hcl:"description,attr" json:"description,omitempty"`
	AvatarFile  string            `hcl:"avatar_file,attr" json:"avatar_file,omitempty"`
	Vars        map[string]string `hcl:"vars,attr" json:"vars,omitempty"`
	OnConflict  OnConflict        `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller      Caller            `json:"-"`
}

// GetCaller implements the Mergeable interface
func (i *Identity) GetCaller() Caller {
	return i.Caller
}

// GetID implements the Mergeable interface
func (i *Identity) GetID() string {
	return i.ID
}

// GetOnConflict implements the Mergeable interface
func (i *Identity) GetOnConflict() OnConflict {
	return i.OnConflict
}

// SetCaller implements the Mergeable interface
func (i *Identity) SetCaller(c Caller) {
	i.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (i *Identity) SetOnConflict(o OnConflict) {
	i.OnConflict = o
}

// Swap implements the Mergeable interface
func (i *Identity) Swap(m Mergeable) error {
	rawVal, ok := m.(*Identity)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", i, m)
	}
	*i = *rawVal
	return nil
}

// ResolveSource attempts to locate the referenced source file with a laforge base configuration
func (i *Identity) ResolveSource(base *Laforge, pr *PathResolver, caller CallFile) error {
	if i.AvatarFile == "" {
		return nil
	}
	cwd, _ := os.Getwd()
	testSrc := i.AvatarFile
	if !filepath.IsAbs(i.AvatarFile) {
		testSrc = filepath.Join(caller.CallerDir, i.AvatarFile)
	}
	if !PathExists(testSrc) {
		pr.Unresolved[i.AvatarFile] = true
		return errors.Wrapf(ErrAbsPathDeclNotExist, "caller=%s path=%s", caller.CallerFile, i.AvatarFile)
	}
	rel, _ := filepath.Rel(cwd, testSrc)
	rel2, _ := filepath.Rel(caller.CallerDir, testSrc)
	lfr := &LocalFileRef{
		Base:          filepath.Base(testSrc),
		AbsPath:       testSrc,
		RelPath:       rel,
		Cwd:           cwd,
		DeclaredPath:  i.AvatarFile,
		RelToCallFile: rel2,
	}
	pr.Mapping[i.AvatarFile] = lfr
	return nil
}
