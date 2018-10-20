package core

// CommandJob attempts to execute remote commands on the system
// easyjson:json
type CommandJob struct {
	GenericJob
	Command          []string         `json:"command"`
	GoodOutputRegexp string           `json:"good_output_regexp"`
	Target           *ProvisionedHost `json:"-"`
}

// CanProceed implements the Doer interface
func (j *CommandJob) CanProceed() error {
	return nil
}

// EnsureDependencies implements the Doer interface
func (j *CommandJob) EnsureDependencies(l *Laforge) error {
	return nil
}

// Do implements the Doer interface
func (j *CommandJob) Do() error {
	return nil
}

// CleanUp implements the Doer interface
func (j *CommandJob) CleanUp() error {
	return nil
}

// Finish implements the Doer interface
func (j *CommandJob) Finish() error {
	return nil
}
