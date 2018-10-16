package core

import (
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
	ID            string
	ObjectType    string
	Checksum      uint64
	CreatedAt     time.Time
	ModifiedAt    time.Time
	ParentIDs     []string
	DependencyIDs []string
	Resources     []MetaResource
}

// MetaResource stores information about a local file dependency. This can be a directory.
// If the resource is a directory, it will be recursively gzip'd and that will be checksum'd.
// If the resource is a directory, size will be the size of the final gzip file.
// Note creation and modification date refer to meta resource validation, not the actual file.
//easyjson:json
type MetaResource struct {
	ID           string
	PathFromBase string
	Basename     string
	ParentIDs    []string
	IsDir        string
	CreatedAt    time.Time
	ModifiedAt   time.Time
	Checksum     uint64
	Size         int
}
