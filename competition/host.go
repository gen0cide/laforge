package competition

type Host struct {
	Hostname       string   `yaml:"hostname"`
	OS             string   `yaml:"os"`
	InstanceSize   string   `yaml:"instance_size"`
	LastOctet      int      `yaml:"last_octet"`
	InternalCNAMEs []string `yaml:"internal_cnames"`
	HasPublicIP    bool     `yaml:"has_public_ip"`
	ExternalCNAMEs []string `yaml:"external_cnames"`
	SecurityGroups []string `yaml:"security_groups"`
	AdminPassword  string   `yaml:"admin_password"`
	Scripts        []string `yaml:"scripts"`
	UserDataScript string   `yaml:"userdata_script"`
	UserGroups     []string `yaml:"user_groups"`
	Vars           `yaml:"variables"`
	Network        `yaml:"-"`
}
