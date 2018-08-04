package laforge

// Competition is a configurable type that holds competition wide settings
type Competition struct {
	BaseDir      string `hcl:"base_dir,attr" json:"base_dir,omitempty"`
	RootPassword string `hcl:"root_password,attr" json:"root_password,omitempty"`
	DNS          DNS    `hcl:"dns,block" json:"domain,omitempty"`
	Remote       Remote `hcl:"remote,block" json:"remote,omitempty"`
}
