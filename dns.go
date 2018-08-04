package laforge

// DNS represents a configurable type for the creation of competition DNS infrastructure
type DNS struct {
	Type        string   `hcl:"type,label" json:"type,omitempty"`
	BaseFQDN    string   `hcl:"base_fqdn,attr" json:"base_fqdn,omitempty"`
	Nameservers []string `hcl:"nameservers,attr" json:"nameservers,omitempty"`
	NTPServers  []string `hcl:"ntp_servers,attr" json:"ntp_servers,omitempty"`
}
