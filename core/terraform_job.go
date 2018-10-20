package core

// OutputCallback is used to monitor the output of terraform as it runs
type OutputCallback func(state *State, job Doer, line string) error

// TerraformJob attempts to run one or more terraform commands on the system
// easyjson:json
type TerraformJob struct {
	GenericJob
	Commands       [][]string     `json:"commands"`
	Target         *Team          `json:"-"`
	OutputCallback OutputCallback `json:"-"`
}

// CanProceed implements the Doer interface
func (j *TerraformJob) CanProceed() error {
	return nil
}

// EnsureDependencies implements the Doer interface
func (j *TerraformJob) EnsureDependencies(l *Laforge) error {
	return nil
}

// Do implements the Doer interface
func (j *TerraformJob) Do() error {
	return nil
}

// CleanUp implements the Doer interface
func (j *TerraformJob) CleanUp() error {
	return nil
}

// Finish implements the Doer interface
func (j *TerraformJob) Finish() error {
	return nil
}
