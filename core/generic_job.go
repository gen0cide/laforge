package core

import (
	"time"

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
	SetStatus(s JobStatus)
	CurrentStatus() JobStatus
}

// GenericJob is an embeddable type for standardizing laforge jobs
type GenericJob struct {
	JobID      string    `json:"job_id"`
	Offset     int       `json:"offset"`
	JobType    string    `json:"job_type"`
	Metadata   *Metadata `json:"-"`
	MetadataID string    `json:"metadata_id"`
	Status     JobStatus `json:"status"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	EndedAt    time.Time `json:"ended_at,omitempty"`
}

// GetID implements the Relationship interface
func (j GenericJob) GetID() string {
	return j.JobID
}

// CurrentStatus implements the Doer interface
func (j GenericJob) CurrentStatus() JobStatus {
	return j.Status
}

// SetStatus implements the Doer interface
func (j GenericJob) SetStatus(s JobStatus) {
	j.Status = s
}

// GetTargetID implements the Doer interface
func (j GenericJob) GetTargetID() string {
	return j.MetadataID
}
