package core

import (
	"context"
	"errors"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/finding"
)

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
//easyjson:json
type Finding struct {
	Name        string            `hcl:"name,attr" json:"name,omitempty"`
	Description string            `hcl:"description,optional" json:"description,omitempty"`
	Severity    FindingSeverity   `hcl:"severity,attr" json:"severity,omitempty"`
	Difficulty  FindingDifficulty `hcl:"difficulty,attr" json:"difficulty,omitempty"`
	Maintainer  *User             `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Tags        []string          `hcl:"tags,optional" json:"tags,omitempty"`
	Provisioner Provisioner       `json:"-"`
	Host        *Host             `json:"-"`
}

// TotalScore returns the total score applicable to a Finding
func (f *Finding) TotalScore() int {
	return int(f.Severity) * int(f.Difficulty)
}

// CreateFindingEntry ...
func (f *Finding) CreateFindingEntry(ph *ent.ProvisionedHost, script *ent.Script, ctx context.Context, client *ent.Client) (*ent.Finding, error) {
	// tag, err := CreateTagEntry(f.Name, f.Tags, ctx, client) // Different Type of Tag

	// if err != nil {
	// 	cli.Logger.Debugf("failed creating finding: %v", err)
	// 	return nil, err
	// }

	user, err := f.Maintainer.CreateUserEntry(ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating finding: %v", err)
		return nil, err
	}
	
	finding, err := client.Finding.
		Create().
		SetName(f.Name).
		SetDescription(f.Description).
		SetSeverity(finding.Severity(f.Severity.String())).
		SetDifficulty(finding.Difficulty(f.Difficulty.String())).
		AddUser(user).
		// AddTag(tag).
		AddHost(ph.Edges.Host...).
		AddScript(script).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating finding: %v", err)
		return nil, err
	}

	cli.Logger.Debugf("finding was created: ", finding)
	return finding, nil
}
