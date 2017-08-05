package competition

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Network struct {
	CIDR          string `yaml:"cidr"`
	Name          string `yaml:"name"`
	Subdomain     string `yaml:"subdomain"`
	Provider      string `yaml:"provider"`
	Vars          `yaml:"variables"`
	IncludedHosts []string `yaml:"included_hosts"`
	Pod           `yaml:"-"`
}

func LoadNetworkFromFile(file string) (*Network, error) {
	network := Network{}
	networkConfig, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(networkConfig, &network)
	if err != nil {
		return nil, err
	}
	return &network, nil
}

func (n *Network) ToYAML() string {
	y, err := yaml.Marshal(n)
	if err != nil {
		LogFatal("Error converting to YAML: " + err.Error())
	}
	return string(y)
}
