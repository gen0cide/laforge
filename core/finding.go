package core

import "errors"

const (
	// ZeroDifficulty is a zero difficulty finding
	ZeroDifficulty FindingDifficulty = iota

	// NoviceDifficulty is a novice difficulty finding
	NoviceDifficulty

	// AdvancedDifficulty is an advanced difficulty finding
	AdvancedDifficulty

	// ExpertDifficulty is an expert difficulty finding
	ExpertDifficulty

	// ZeroSeverity is a zero severity finding
	ZeroSeverity FindingSeverity = iota

	// LowSeverity is a low severity finding
	LowSeverity

	// MediumSeverity is a medium severity finding
	MediumSeverity

	// HighSeverity is a high severity finding
	HighSeverity

	// CriticalSeverity is a critical severity finding
	CriticalSeverity
)

var (
	// ErrInvalidDifficulty is thrown if the difficulty is invalid
	ErrInvalidDifficulty = errors.New("difficulty does not fall within valid range of 0-3")

	// ErrInvalidSeverity is thrown if the difficulty is invalid
	ErrInvalidSeverity = errors.New("severity does not fall within valid range of 0-3")
)

// FindingDifficulty is a type alias to the level of difficulty the finding exists as
type FindingDifficulty int

// FindingSeverity is a type alias to the level of severity the finding exists as
type FindingSeverity int

// String implements the stringer interface
func (f FindingDifficulty) String() string {
	switch f {
	case ZeroDifficulty:
		return "ZeroDifficulty"
	case NoviceDifficulty:
		return "NoviceDifficulty"
	case AdvancedDifficulty:
		return "AdvancedDifficulty"
	case ExpertDifficulty:
		return "ExpertDifficulty"
	default:
		return "NullDifficulty"
	}
}

// String implements the stringer interface
func (f FindingSeverity) String() string {
	switch f {
	case ZeroSeverity:
		return "ZeroSeverity"
	case LowSeverity:
		return "LowSeverity"
	case MediumSeverity:
		return "MediumSeverity"
	case HighSeverity:
		return "HighSeverity"
	case CriticalSeverity:
		return "CriticalSeverity"
	default:
		return "NullSeverity"
	}
}

// Finding represents a finding to be aggregated for scoring inside a laforge scenario
type Finding struct {
	ID          string      `hcl:",id" json:"id,omitempty"`
	Name        string      `hcl:"name,attr" json:"name,omitempty"`
	Description string      `hcl:"description,attr" json:"description,omitempty"`
	Severity    int         `hcl:"severity,attr" json:"severity,omitempty"`
	Difficulty  int         `hcl:"difficulty,attr" json:"difficulty,omitempty"`
	Maintainer  *User       `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Tags        []string    `hcl:"tags,attr" json:"tags,omitempty"`
	Provisioner Provisioner `json:"-"`
	Host        *Host       `json:"-"`
}

// TotalScore returns the total score applicable to a Finding
func (f *Finding) TotalScore() int {
	return f.Severity * f.Difficulty
}
