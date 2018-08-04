package laforge

// AMI represents a configurable object for defining custom AMIs in cloud infrastructure
type AMI struct {
	Name        string `hcl:"name,label" json:"name,omitempty"`
	OS          string `hcl:"os,attr" json:"os,omitempty"`
	ID          string `hcl:"id,attr" json:"id,omitempty"`
	Version     string `hcl:"version,attr" json:"version,omitempty"`
	Arch        string `hcl:"arch,attr" json:"arch,omitempty"`
	Provider    string `hcl:"provider,attr" json:"provider,omitempty"`
	Region      string `hcl:"region,attr" json:"region,omitempty"`
	Username    string `hcl:"username,attr" json:"username,omitempty"`
	Description string `hcl:"description,attr" json:"description,omitempty"`
	Maintainer  User   `hcl:"maintainer,block" json:"maintainer,omitempty"`
}
