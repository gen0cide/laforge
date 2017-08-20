package competition

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Host struct {
	Hostname       string   `yaml:"hostname"`
	OS             string   `yaml:"os"`
	AMI            string   `yaml:"ami"`
	InstanceSize   string   `yaml:"instance_size"`
	LastOctet      int      `yaml:"last_octet"`
	InternalCNAMEs []string `yaml:"internal_cnames"`
	ExternalCNAMEs []string `yaml:"external_cnames"`
	AdminPassword  string   `yaml:"admin_password"`
	TCPPorts       []int    `yaml:"public_tcp"`
	UDPPorts       []int    `yaml:"public_udp"`
	Scripts        []string `yaml:"scripts"`
	UserDataScript string   `yaml:"userdata_script"`
	UserGroups     []string `yaml:"user_groups"`
	Vars           `yaml:"variables"`
	Network        `yaml:"-"`
}

func LoadHostFromFile(file string) (*Host, error) {
	host := Host{}
	hostConfig, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(hostConfig, &host)
	if err != nil {
		return nil, err
	}
	return &host, nil
}

func (h *Host) ToYAML() string {
	y, err := yaml.Marshal(h)
	if err != nil {
		LogFatal("Error converting to YAML: " + err.Error())
	}
	return string(y)
}
