package core

import (
	"encoding/json"
	"time"

	"github.com/cespare/xxhash"

	"github.com/gen0cide/laforge/core/graph"
)

// JobStatus defines what state the job is in
type JobStatus int

const (
	// JobStatusPlanned is the default assignment given to jobs
	JobStatusPlanned JobStatus = iota

	// JobStatusEnqueued is assigned when the job has been put on the channel for work
	JobStatusEnqueued

	// JobStatusInProgress is assigned when the job has been picked up off the channel and is being worked
	JobStatusInProgress

	// JobStatusFailed is assigned when the job status has failed to continue
	JobStatusFailed

	// JobStatusSuccessful is assigned when the job completed successfully
	JobStatusSuccessful
)

// Doer is an interface to describe types that may be executed in the flow
type Doer interface {
	graph.Hasher
	GetTargetID() string
	CanProceed() error
	EnsureDependencies(l *Laforge) error
	Do() error
	CleanUp() error
	Finish() error
	SetTimeout(t int)
	SetStatus(s JobStatus)
	SetPlan(p *Plan)
	SetBase(l *Laforge)
	CurrentStatus() JobStatus
}

// GenericJob is an embeddable type for standardizing laforge jobs
// easyjson:json
type GenericJob struct {
	JobID      string    `json:"job_id"`
	Plan       *Plan     `json:"-"`
	Base       *Laforge  `json:"-"`
	Offset     int       `json:"offset"`
	Timeout    int       `json:"timeout"`
	JobType    string    `json:"job_type"`
	Metadata   *Metadata `json:"-"`
	MetadataID string    `json:"metadata_id"`
	Status     JobStatus `json:"status"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	EndedAt    time.Time `json:"ended_at,omitempty"`
}

// Hash implements the Hasher interface
func (j *GenericJob) Hash() uint64 {
	d, err := json.Marshal(j)
	if err != nil {
		return uint64(666)
	}
	return xxhash.Sum64(d)
}

// GetID implements the Relationship interface
func (j *GenericJob) GetID() string {
	return j.JobID
}

// CurrentStatus implements the Doer interface
func (j *GenericJob) CurrentStatus() JobStatus {
	return j.Status
}

// SetStatus implements the Doer interface
func (j *GenericJob) SetStatus(s JobStatus) {
	j.Status = s
}

// GetTargetID implements the Doer interface
func (j *GenericJob) GetTargetID() string {
	return j.MetadataID
}

// SetTimeout implements the Doer interface
func (j *GenericJob) SetTimeout(t int) {
	j.Timeout = t
}

// SetPlan implements the Doer interface
func (j *GenericJob) SetPlan(p *Plan) {
	j.Plan = p
}

// SetBase implements the Doer interface
func (j *GenericJob) SetBase(l *Laforge) {
	j.Base = l
}
