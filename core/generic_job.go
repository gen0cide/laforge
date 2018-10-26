package core

import (
	"encoding/json"
	"time"

	"github.com/cespare/xxhash"
	"github.com/pkg/errors"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/core/graph"
)

// JobStatus defines what state the job is in
type JobStatus int

var (
	// ErrTimeoutExtensionRequested is thrown when the loop for a job should happen again
	// ErrTimeoutExtensionRequested = errors.New("job requested an extension on the timeout")

	// ErrTimeoutExceeded is thrown when the timeout has exceeded for this task
	ErrTimeoutExceeded = errors.New("timeout has exceeded for task step")
)

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
	CanProceed(e chan error)
	EnsureDependencies(e chan error)
	Do(e chan error)
	CleanUp(e chan error)
	Finish(e chan error)
	SetTimeout(t int)
	GetTimeout() int
	GetMetadata() *Metadata
	SetStatus(s JobStatus)
	StandardOutput(line string)
	StandardError(line string)
	SetPlan(p *Plan)
	SetBase(l *Laforge)
	CurrentStatus() JobStatus
}

// NewTimeoutExtension creates a wrapped error for the scheduler to retry at a later time
func NewTimeoutExtension(err error) *ErrTimeoutExtension {
	return &ErrTimeoutExtension{
		orig: err,
	}
}

// ErrTimeoutExtension is a type used to request a timeout extension for the executor
type ErrTimeoutExtension struct {
	orig error
}

// Error implements the error interface
func (e *ErrTimeoutExtension) Error() string {
	return "job requested an extension on the timeout"
}

// Cause implements the causer interface
func (e *ErrTimeoutExtension) Cause() error {
	return e.orig
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

// TimeoutFunc is a function that is retried ever half second until the interval period is hit
type TimeoutFunc func(errchan chan error)

// PerformInTimeout will perform the TimeoutFunc f every 500ms until it either returns NOT an ErrTimeoutExtensionRequested or nil
func PerformInTimeout(seconds int, f TimeoutFunc) error {
	timeout := time.After(time.Duration(seconds) * time.Second)
	tick := time.Tick(1 * time.Second)
	errchan := make(chan error, 1)
	go f(errchan)
	for {
		select {
		case <-timeout:
			return ErrTimeoutExceeded
		case err := <-errchan:
			if err == nil {
				return nil
			}
			if err != nil {
				if te, ok := err.(*ErrTimeoutExtension); ok {
					cli.Logger.Debugf("timeout extension requested: %v", te.Cause())
					<-tick
					go f(errchan)
					continue
				}
				return err
			}
		}
	}
}

// GetTimeout implements the Doer interface
func (j *GenericJob) GetTimeout() int {
	return j.Timeout
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

// GetMetadata implements the Doer interface
func (j *GenericJob) GetMetadata() *Metadata {
	return j.Metadata
}

// StandardOutput shows the standard output of a job's execution
func (j *GenericJob) StandardOutput(line string) {
	cli.Logger.Debugf("%s (STDOUT): %s", j.JobID, line)
}

// StandardError prints the standard error of a jobs execution
func (j *GenericJob) StandardError(line string) {
	cli.Logger.Debugf("%s (STDERR): %s", j.JobID, line)
}
