package core

// Remote defines a configuration object that keeps terraform and remote files synchronized
type Remote struct {
	ID            string `hcl:"id,label" json:"id,omitempty"`
	Type          string `hcl:"type,attr" json:"type,omitempty"`
	Region        string `hcl:"region,attr" json:"region,omitempty"`
	Key           string `hcl:"key,attr" json:"key,omitempty"`
	Secret        string `hcl:"secret,attr" json:"secret,omitempty"`
	StateBucket   string `hcl:"state_bucket,attr" json:"state_bucket,omitempty"`
	StorageBucket string `hcl:"storage_bucket,attr" json:"storage_bucket,omitempty"`
}
