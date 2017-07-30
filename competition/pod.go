package competition

type Pod struct {
	Domain       string `yaml:"domain"`
	PodID        string `yaml:"-"`
	JumpHosts    `yaml:"jump_hosts"`
	*Environment `yaml:"-"`
	Vars         `yaml:"variables"`
	Networks     []Network `yaml:"networks"`
}

type JumpHosts struct {
	Windows struct {
		AMI   string `yaml:"ami"`
		Count int    `yaml:"count"`
		Size  string `yaml:"size"`
	} `yaml:"windows"`
	Kali struct {
		AMI   string `yaml:"ami"`
		Count int    `yaml:"count"`
		Size  string `yaml:"size"`
	} `yaml:"kali"`
	DNS struct {
		AMI  string `yaml:"ami"`
		Size string `yaml:"size"`
	} `yaml:"dns"`
}
