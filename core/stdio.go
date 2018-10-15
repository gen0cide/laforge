package core

// IO represents a configuration type to define file locations for Stdin, Stdout, Stderr
type IO struct {
	Stdin  string `cty:"stdin" hcl:"stdin,optional" json:"stdin,omitempty"`
	Stdout string `cty:"stdout" hcl:"stdout,optional" json:"stdout,omitempty"`
	Stderr string `cty:"stderr" hcl:"stderr,optional" json:"stderr,omitempty"`
}
