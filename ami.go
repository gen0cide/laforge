package laforge

// AMI represents a configurable object for defining custom AMIs in cloud infrastructure
type AMI struct {
	ID          string            `hcl:"id,label" json:"id,omitempty"`
	Name        string            `hcl:"name,attr" json:"name,omitempty"`
	Description string            `hcl:"description,attr" json:"description,omitempty"`
	Provider    string            `hcl:"provider,attr" json:"provider,omitempty"`
	Username    string            `hcl:"username,attr" json:"username,omitempty"`
	Config      map[string]string `hcl:"config,attr" json:"config,omitempty"`
	Tags        map[string]string `hcl:"tags,attr" json:"tags,omitempty"`
	Maintainer  *User             `hcl:"maintainer,block" json:"maintainer,omitempty"`
}
