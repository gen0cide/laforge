package provisioner

import "time"

// Step is a unique action to be taken on the local machine
type Step struct {
	ID          string                 `json:"id,omitempty"`
	Revision    string                 `json:"revision,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	StepType    string                 `json:"step_type,omitempty"`
	Source      string                 `json:"source,omitempty"`
	Destination string                 `json:"destination,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Status      string                 `json:"status,omitempty"`
	StartedAt   time.Time              `json:"started_at,omitempty"`
	EndedAt     time.Time              `json:"ended_at,omitempty"`
	StdoutFile  string                 `json:"stdout_file,omitempty"`
	StderrFile  string                 `json:"stderr_file,omitempty"`
	ExitStatus  int                    `json:"exit_status,omitempty"`
	ExitError   error                  `json:"exit_error,omitempty"`
}
