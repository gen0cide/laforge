package laforge

// User defines a laforge command line user and their properties
type User struct {
	Name  string `hcl:"name,attr" json:"name,omitempty"`
	ID    string `hcl:"id,attr" json:"id,omitempty"`
	Email string `hcl:"email,attr" json:"email,omitempty"`
}
