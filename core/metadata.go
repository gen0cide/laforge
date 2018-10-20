package core

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/dag"

	"github.com/gen0cide/laforge/core/graph"
	"github.com/pkg/errors"
)

// Pather is an interface to define hosts which need to conform to valid pathing schemes
type Pather interface {
	Path() string
	Base() string
	ValidatePath() error
}

// Dependency is an interface to define a laforge object that can be represented on the graph
type Dependency interface {
	Pather
	graph.Hasher
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
	Name         string               `json:"name"`
	ID           string               `json:"id,omitempty"`
	GID          int                  `json:"gid"`
	GCost        int64                `json:"gcost"`
	ObjectType   string               `json:"object_type,omitempty"`
	Dependency   Dependency           `json:"-"`
	Tainted      bool                 `json:"tainted,omitempty"`
	Addition     bool                 `json:"addition,omitempty"`
	Checksum     uint64               `json:"checksum,omitempty"`
	CreatedAt    time.Time            `json:"created_at,omitempty"`
	ModifiedAt   time.Time            `json:"modified_at,omitempty"`
	ParentDepIDs []string             `json:"parent_ids,omitempty"`
	ParentDeps   []graph.Relationship `json:"-"`
	ParentGIDs   []int                `json:"parent_gids"`
	ChildDeps    []graph.Relationship `json:"-"`
	ChildGIDs    []int                `json:"child_gids"`
	ChildDepIDs  []string             `json:"dependency_ids,omitempty"`
	Resources    []MetaResource       `json:"resources,omitempty"`
}

// LFType describes a string representation of elements in Laforge
type LFType string

const (
	// LFTypeCompetition is a constant to define object type when serialized
	LFTypeCompetition LFType = `competition`

	// LFTypeNetwork is a constant to define object type when serialized
	LFTypeNetwork LFType = `network`

	// LFTypeHost is a constant to define object type when serialized
	LFTypeHost LFType = `host`

	// LFTypeRemoteFile is a constant to define object type when serialized
	LFTypeRemoteFile LFType = `remote_file`

	// LFTypeCommand is a constant to define object type when serialized
	LFTypeCommand LFType = `command`

	// LFTypeDNSRecord is a constant to define object type when serialized
	LFTypeDNSRecord LFType = `dns_record`

	// LFTypeScript is a constant to define object type when serialized
	LFTypeScript LFType = `script`

	// LFTypeEnvironment is a constant to define object type when serialized
	LFTypeEnvironment LFType = `environment`

	// LFTypeBuild is a constant to define object type when serialized
	LFTypeBuild LFType = `build`

	// LFTypeTeam is a constant to define object type when serialized
	LFTypeTeam LFType = `team`

	// LFTypeProvisionedNetwork is a constant to define object type when serialized
	LFTypeProvisionedNetwork LFType = `provisioned_network`

	// LFTypeProvisionedHost is a constant to define object type when serialized
	LFTypeProvisionedHost LFType = `provisioned_host`

	// LFTypeConnection is a constant to define object type when serialized
	LFTypeConnection LFType = `connection`

	// LFTypeProvisioningStep is a constant to define object type when serialized
	LFTypeProvisioningStep LFType = `provisioning_step`

	// LFTypeUnknown is totally a fucker
	LFTypeUnknown LFType = "unknown"
)

func (m *Metadata) TypeByPath() LFType {
	return TypeByPath(m.ID)
}

func (m *Metadata) IsGlobalType() bool {
	return IsGlobalType(m.ID)
}

func IsGlobalType(p string) bool {
	switch TypeByPath(p) {
	case LFTypeCompetition:
		return true
	case LFTypeCommand:
		return true
	case LFTypeNetwork:
		return true
	case LFTypeHost:
		return true
	case LFTypeDNSRecord:
		return true
	case LFTypeRemoteFile:
		return true
	case LFTypeScript:
		return true
	case LFTypeEnvironment:
		return true
	case LFTypeTeam:
		return true
	case LFTypeProvisionedNetwork:
		return true
	case LFTypeProvisionedHost:
		return true
	case LFTypeProvisioningStep:
		return true
	case LFTypeUnknown:
		return true
	default:
		return false
	}
}

func TypeByPath(p string) LFType {
	if !path.IsAbs(p) {
		return LFTypeCompetition
	}
	pelms := strings.Split(p, `/`)

	switch pelms[1] {
	case "scripts":
		return LFTypeScript
	case "networks":
		return LFTypeNetwork
	case "hosts":
		return LFTypeHost
	case "commands":
		return LFTypeCommand
	case "dns-records":
		return LFTypeDNSRecord
	case "files":
		return LFTypeRemoteFile
	}

	if path.Dir(p) == "envs" {
		return LFTypeEnvironment
	}

	if path.Dir(path.Dir(p)) == "envs" {
		return LFTypeBuild
	}

	if path.Dir(p) == "teams" {
		return LFTypeTeam
	}

	if path.Dir(p) == "networks" && pelms[1] == "envs" {
		return LFTypeProvisionedNetwork
	}

	if path.Dir(p) == "hosts" && pelms[1] == "envs" {
		return LFTypeProvisionedNetwork
	}

	if path.Base(p) == "conn" && path.Dir(path.Dir(p)) == "hosts" {
		return LFTypeConnection
	}

	if path.Dir(p) == "steps" && pelms[1] == "envs" {
		return LFTypeProvisioningStep
	}

	return LFTypeUnknown
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

func (m *Metadata) DotNode(s string, d *dag.DotOpts) *dag.DotNode {
	return &dag.DotNode{
		Name: s,
		Attrs: map[string]string{
			"checksum": fmt.Sprintf("%d", m.Checksum),
		},
	}
}

// Name implements the DotNode interface
// func (m *Metadata) Name() string {
// 	return m.ID
// }

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

// Children implements the Relationship interface
func (m *Metadata) Children() []graph.Relationship {
	return m.ChildDeps
}

// Parents implements the Relationship interface
func (m *Metadata) Parents() []graph.Relationship {
	return m.ParentDeps
}

// ParentIDs implements the relationship interface
func (m *Metadata) ParentIDs() []string {
	return m.ParentDepIDs
}

// ChildrenIDs implements the relationship interface
func (m *Metadata) ChildrenIDs() []string {
	return m.ChildDepIDs
}

// AddChild implements the relationship interface
func (m *Metadata) AddChild(r ...graph.Relationship) {
	for _, x := range r {
		if !graph.HasChild(m, x) {
			m.ChildDepIDs = append(m.ChildDepIDs, x.GetID())
			m.ChildDeps = append(m.ChildDeps, x)
			m.ChildGIDs = append(m.ChildGIDs, x.GetGID())
		}
	}
	return
}

// AddParent implements the relationship interface
func (m *Metadata) AddParent(r ...graph.Relationship) {
	for _, x := range r {
		if !graph.HasParent(m, x) {
			m.ParentDepIDs = append(m.ParentDepIDs, x.GetID())
			m.ParentDeps = append(m.ParentDeps, x)
			m.ParentGIDs = append(m.ParentGIDs, x.GetGID())
		}
	}
	return
}

// Hash implements the hasher interface
func (m *Metadata) Hash() uint64 {
	if m.Checksum == 0 {
		m.CalculateChecksum()
	}
	return m.Checksum
}

// Hashcode implements the Hashable interface
func (m *Metadata) Hashcode() interface{} {
	return m.Checksum
}

// String implements the stringer interface
func (m *Metadata) String() string {
	return m.ID
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
