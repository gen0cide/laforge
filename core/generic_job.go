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
	graph.Relationship
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
	JobID        string          `json:"job_id"`
	Offset       int             `json:"offset"`
	JobType      string          `json:"job_type"`
	Metadata     *Metadata       `json:"-"`
	MetadataID   string          `json:"metadata_id"`
	Status       JobStatus       `json:"status"`
	CreatedAt    time.Time       `json:"created_at,omitempty"`
	StartedAt    time.Time       `json:"started_at,omitempty"`
	EndedAt      time.Time       `json:"ended_at,omitempty"`
	ParentJobIDs map[string]bool `json:"parent_job_ids"`
	ChildJobIDs  map[string]bool `json:"child_job_ids"`
	ParentJobs   map[Doer]bool   `json:"-"`
	ChildJobs    map[Doer]bool   `json:"-"`
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

// Children implements the Relationship interface
func (j GenericJob) Children() []Doer {
	ret := make([]Doer, len(j.ChildJobs))
	for k := range j.ChildJobs {
		ret = append(ret, k)
	}
	return ret
}

// Parents implements the Relationship interface
func (j GenericJob) Parents() []Doer {
	ret := make([]Doer, len(j.ParentJobs))
	for k := range j.ParentJobs {
		ret = append(ret, k)
	}
	return ret
}

// ParentIDs implements the Relationship interface
func (j GenericJob) ParentIDs() []string {
	ret := make([]string, len(j.ParentJobIDs))
	for k := range j.ParentJobIDs {
		ret = append(ret, k)
	}
	return ret
}

// ChildrenIDs implements the Relationship interface
func (j GenericJob) ChildrenIDs() []string {
	ret := make([]string, len(j.ChildJobIDs))
	for k := range j.ChildJobIDs {
		ret = append(ret, k)
	}
	return ret
}

// AddChild implements the Relationship interface
func (j GenericJob) AddChild(d ...Doer) {
	for _, x := range d {
		j.ChildJobs[x] = true
		j.ChildJobIDs[x.GetID()] = true
	}
	return
}

// AddParent implements the Relationship interface
func (j GenericJob) AddParent(d ...Doer) {
	for _, x := range d {
		j.ChildJobs[x] = true
		j.ChildJobIDs[x.GetID()] = true
	}
	return
}
