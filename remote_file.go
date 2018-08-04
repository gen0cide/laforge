package laforge

// RemoteFile is a configurable type that defines a static file that will be placed on a configured target host.
type RemoteFile struct {
	ID          string `hcl:"id,attr" json:"id,omitempty"`
	SourceType  string `hcl:"source_type,attr" json:"source_type,omitempty"`
	SourcePath  string `hcl:"source_path,attr" json:"source_path,omitempty"`
	Destination string `hcl:"destination,label" json:"destination,omitempty"`
	Perms       string `hcl:"perms,attr" json:"perms,omitempty"`
	Disabled    bool   `hcl:"disabled,attr" json:"disabled,omitempty"`
}
