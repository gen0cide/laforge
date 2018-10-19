package core

import (
	"fmt"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

// Hasher is an interface to allow types to be checksumed for potentially build breaking changes
type Hasher interface {
	Hash() uint64
}

// Pather is an interface to define hosts which need to conform to valid pathing schemes
type Pather interface {
	Path() string
	Base() string
	ValidatePath() error
}

// Dependency is an interface to define a laforge object that can be represented on the graph
type Dependency interface {
	Pather
	Hasher
	ParentLaforgeID() string
	Gather(g *Snapshot) error
}

// ResourceHasher is an interface to define types who have file dependencies to checkum them
type ResourceHasher interface {
	ResourceHash() uint64
}

var (
	genericPathRegexp      = regexp.MustCompile(`^\/[a-z0-9\-\/]{3,}[a-z0-9]$`)
	consecutiveSlashRegexp = regexp.MustCompile(`\/\/`)

	// ErrPathEndsInSlash is thrown when a path ends in a trailing slash
	ErrPathEndsInSlash = errors.New("path ends in a trailing slash")

	// ErrPathContainsInvalidChars is thrown when a path contains invalid characters
	ErrPathContainsInvalidChars = errors.New("path contains invalid characters")

	// ErrPathContainsDuplicateSlash is thrown when a path contains two consecutive slashes
	ErrPathContainsDuplicateSlash = errors.New("path contains consecutive slash characters")
)

// ValidateGenericPath covers basic rules validating a path generically for invalid schema
func ValidateGenericPath(p string) error {
	if !genericPathRegexp.MatchString(p) {
		return ErrPathContainsInvalidChars
	}
	if consecutiveSlashRegexp.MatchString(p) {
		return ErrPathContainsDuplicateSlash
	}
	return nil
}

// Metadata stores metadata about different structs within the environment
//easyjson:json
type Metadata struct {
	ID         string         `json:"id,omitempty"`
	GID        int            `json:"gid"`
	GCost      int64          `json:"gcost"`
	ObjectType string         `json:"object_type,omitempty"`
	Dependency Dependency     `json:"-"`
	Checksum   uint64         `json:"checksum,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	ModifiedAt time.Time      `json:"modified_at,omitempty"`
	ParentIDs  []string       `json:"parent_ids,omitempty"`
	ParentDeps []DotNode      `json:"-"`
	ParentGIDs []int          `json:"parent_gids"`
	ChildDeps  []DotNode      `json:"-"`
	ChildGIDs  []int          `json:"child_gids"`
	ChildIDs   []string       `json:"dependency_ids,omitempty"`
	Resources  []MetaResource `json:"resources,omitempty"`
}

// GetID implements the DotNode interface
func (m *Metadata) GetID() string {
	return m.ID
}

// GetGID implements the DotNode interface
func (m *Metadata) GetGID() int {
	return m.GID
}

// GetGCost implements the DotNode interface
func (m *Metadata) GetGCost() int64 {
	return m.GCost
}

// GetChecksum implements the DotNode interface
func (m *Metadata) GetChecksum() uint64 {
	return m.Checksum
}

// Name implements the DotNode interface
func (m *Metadata) Name() string {
	return m.ID
}

// Label implements the DotNode interface
func (m *Metadata) Label() string {
	return fmt.Sprintf("%s|type = %s|primary_parent = %s|checksum = %x|num_parents = %d|num_children = %d",
		m.ID,
		m.ObjectType,
		m.Dependency.ParentLaforgeID(),
		m.Checksum,
		len(m.ParentDeps),
		len(m.ChildDeps),
	)
}

// Shape implements the DotNode interface
func (m *Metadata) Shape() string {
	return "record"
}

// Style implements the DotNode interface
func (m *Metadata) Style() string {
	return "solid"
}

// Children implements the DotNode interface
func (m *Metadata) Children() []DotNode {
	return m.ChildDeps
}

// Parents implements the DotNode interface
func (m *Metadata) Parents() []DotNode {
	return m.ParentDeps
}

// CalculateChecksum assigns the metadata object's checksum field with the dependency's hash
func (m *Metadata) CalculateChecksum() {
	m.Checksum = m.Dependency.Hash()
}

// MetaResource stores information about a local file dependency. This can be a directory.
// If the resource is a directory, it will be recursively gzip'd and that will be checksum'd.
// If the resource is a directory, size will be the size of the final gzip file.
// Note creation and modification date refer to meta resource validation, not the actual file.
//easyjson:json
type MetaResource struct {
	ID           string    `json:"id,omitempty"`
	PathFromBase string    `json:"path_from_base,omitempty"`
	Basename     string    `json:"basename,omitempty"`
	ParentIDs    []string  `json:"parent_ids,omitempty"`
	IsDir        bool      `json:"is_dir,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	ModifiedAt   time.Time `json:"modified_at,omitempty"`
	Checksum     uint64    `json:"checksum,omitempty"`
	Size         int       `json:"size,omitempty"`
}
