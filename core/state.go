package core

// Remote defines a configuration object that keeps terraform and remote files synchronized
type Remote struct {
	ID     string            `hcl:"id,label" json:"id,omitempty"`
	Type   string            `hcl:"type,attr" json:"type,omitempty"`
	Config map[string]string `hcl:"config,attr" json:"config,omitempty"`
}
