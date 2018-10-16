package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/pkg/errors"
)

// Identity defines a generic human identity primative that can be extended into Employee, Customer, Client, etc.
//easyjson:json
type Identity struct {
	ID          string            `hcl:"id,label" json:"id,omitempty"`
	Firstname   string            `hcl:"firstname,attr" json:"firstname,omitempty"`
	Lastname    string            `hcl:"lastname,attr" json:"lastname,omitempty"`
	Email       string            `hcl:"email,attr" json:"email,omitempty"`
	Password    string            `hcl:"password,attr" json:"password,omitempty"`
	Description string            `hcl:"description,optional" json:"description,omitempty"`
	AvatarFile  string            `hcl:"avatar_file,optional" json:"avatar_file,omitempty"`
	Vars        map[string]string `hcl:"vars,optional" json:"vars,omitempty"`
	Tags        map[string]string `hcl:"tags,optional" json:"tags,omitempty"`
	OnConflict  *OnConflict       `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Caller      Caller            `json:"-"`
}

// Hash implements the Hasher interface
func (i *Identity) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"first=%v last=%v email=%v pass=%v avatarfile=%v vars=%v",
			i.Firstname,
			i.Lastname,
			i.Email,
			i.Password,
			i.AvatarFile,
			i.Vars,
		),
	)
}

// Path implements the Pather interface
func (i *Identity) Path() string {
	return i.ID
}

// Base implements the Pather interface
func (i *Identity) Base() string {
	return path.Base(i.ID)
}

// ValidatePath implements the Pather interface
func (i *Identity) ValidatePath() error {
	if err := ValidateGenericPath(i.Path()); err != nil {
		return err
	}
	if topdir := strings.Split(i.Path(), `/`); topdir[1] != "identities" {
		return fmt.Errorf("path %s is not rooted in /%s", i.Path(), topdir[1])
	}
	return nil
}

// GetCaller implements the Mergeable interface
func (i *Identity) GetCaller() Caller {
	return i.Caller
}

// LaforgeID implements the Mergeable interface
func (i *Identity) LaforgeID() string {
	return i.ID
}

// GetOnConflict implements the Mergeable interface
func (i *Identity) GetOnConflict() OnConflict {
	if i.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *i.OnConflict
}

// SetCaller implements the Mergeable interface
func (i *Identity) SetCaller(c Caller) {
	i.Caller = c
}

// SetOnConflict implements the Mergeable interface
func (i *Identity) SetOnConflict(o OnConflict) {
	i.OnConflict = &o
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
