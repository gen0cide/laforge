package core

// RemoteFileJob attempts to upload a file to the remote machine
// easyjson:json
type RemoteFileJob struct {
	GenericJob
	Target *ProvisioningStep `json:"-"`
	File   *RemoteFile       `json:"-"`
}

// CanProceed implements the Doer interface
func (j *RemoteFileJob) CanProceed() error {
	return nil
}

// EnsureDependencies implements the Doer interface
func (j *RemoteFileJob) EnsureDependencies(l *Laforge) error {
	return nil
}

// Do implements the Doer interface
func (j *RemoteFileJob) Do() error {
	return nil
}

// CleanUp implements the Doer interface
func (j *RemoteFileJob) CleanUp() error {
	return nil
}

// Finish implements the Doer interface
func (j *RemoteFileJob) Finish() error {
	return nil
}
