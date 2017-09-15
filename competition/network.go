package competition

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Network struct {
	CIDR          string `yaml:"cidr"`
	Name          string `yaml:"name"`
	Subdomain     string `yaml:"subdomain"`
	VDIVisible    bool   `yaml:"vdi_visible"`
	Vars          `yaml:"variables"`
	IncludedHosts []string         `yaml:"included_hosts"`
	ResolvedHosts map[string]*Host `yaml:"-"`
	Environment   `yaml:"-"`
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

func (n *Network) ResolveIncludedHosts() map[string]*Host {
	hosts := make(map[string]*Host)
	hostFiles, _ := filepath.Glob(filepath.Join(n.Environment.HostsDir(), "*.yml"))
	for _, file := range hostFiles {
		if !Contains(FileToName(file), n.IncludedHosts) {
			continue
		}
		host, err := LoadHostFromFile(file)
		if err != nil {
			LogError("Error reading host file: " + file)
			continue
		}
		host.Network = *n
		hosts[FileToName(file)] = host
	}

	return hosts
}

func (n *Network) ToYAML() string {
	y, err := yaml.Marshal(n)
	if err != nil {
		LogFatal("Error converting to YAML: " + err.Error())
	}
	return string(y)
}
