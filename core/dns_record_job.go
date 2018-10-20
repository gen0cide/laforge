package core

// DNSRecordJob attempts to set a DNS record via NSUpdate
// easyjson:json
type DNSRecordJob struct {
	GenericJob
	Target    *ProvisioningStep `json:"-"`
	DNSRecord *DNSRecord        `json:"-"`
}

// CanProceed implements the Doer interface
func (j *DNSRecordJob) CanProceed() error {
	return nil
}

// EnsureDependencies implements the Doer interface
func (j *DNSRecordJob) EnsureDependencies(l *Laforge) error {
	return nil
}

// Do implements the Doer interface
func (j *DNSRecordJob) Do() error {
	return nil
}

// CleanUp implements the Doer interface
func (j *DNSRecordJob) CleanUp() error {
	return nil
}

// Finish implements the Doer interface
func (j *DNSRecordJob) Finish() error {
	return nil
}
