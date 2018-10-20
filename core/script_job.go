package core

// ScriptJob attempts to upload and execute a script on the remote system
// easyjson:json
type ScriptJob struct {
	GenericJob
	Target *ProvisioningStep `json:"-"`
	Script *Script           `json:"-"`
}

// CanProceed implements the Doer interface
func (j *ScriptJob) CanProceed() error {
	return nil
}

// EnsureDependencies implements the Doer interface
func (j *ScriptJob) EnsureDependencies(l *Laforge) error {
	return nil
}

// Do implements the Doer interface
func (j *ScriptJob) Do() error {
	return nil
}

// CleanUp implements the Doer interface
func (j *ScriptJob) CleanUp() error {
	return nil
}

// Finish implements the Doer interface
func (j *ScriptJob) Finish() error {
	return nil
}
