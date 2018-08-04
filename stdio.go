package laforge

// IO represents a configuration type to define file locations for Stdin, Stdout, Stderr
type IO struct {
	Stdin  string `hcl:"stdin,attr" json:"stdin,omitempty"`
	Stdout string `hcl:"stdout,attr" json:"stdout,omitempty"`
	Stderr string `hcl:"stderr,attr" json:"stderr,omitempty"`
}
