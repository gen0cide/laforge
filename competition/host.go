package competition

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var (
	AMIMap map[string]*AMI
)

type AMI struct {
	OS       string            `json:"os"`
	Regions  map[string]string `json:"regions"`
	Username string            `json:"username"`
}

type DNSEntry struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	Type  string `yaml:"type"`
}

type Dependency struct {
	Host    string `yaml:"host"`
	Network string `yaml:"network"`
}

type Host struct {
	Hostname       string       `yaml:"hostname"`
	OS             string       `yaml:"os"`
	AMI            string       `yaml:"ami"`
	InstanceSize   string       `yaml:"instance_size"`
	LastOctet      int          `yaml:"last_octet"`
	InternalCNAMEs []string     `yaml:"internal_cnames"`
	ExternalCNAMEs []string     `yaml:"external_cnames"`
	TCPPorts       []int        `yaml:"public_tcp"`
	UDPPorts       []int        `yaml:"public_udp"`
	Scripts        []string     `yaml:"scripts"`
	UserGroups     []string     `yaml:"user_groups"`
	DNSEntries     []DNSEntry   `yaml:"dns_entries"`
	Dependencies   []Dependency `yaml:"dependencies"`
	Vars           `yaml:"variables"`
	Network        `yaml:"-"`
}

func (h *Host) RenderedDNSEntries(podOffset int) []DNSEntry {
	comp := h.Network.Environment.Competition
	env := h.Network.Environment
	net := h.Network
	dnsEntries := []DNSEntry{}
	for _, e := range h.DNSEntries {
		newName := StringRender(e.Name, &comp, &env, podOffset, &net, h)
		newValue := StringRender(e.Value, &comp, &env, podOffset, &net, h)
		dnsEntry := DNSEntry{
			Name:  newName,
			Value: newValue,
			Type:  e.Type,
		}
		dnsEntries = append(dnsEntries, dnsEntry)
	}
	return dnsEntries
}

func LoadAMIs() {
	var amis []*AMI
	amiMap := make(map[string]*AMI)
	json.Unmarshal(MustAsset("amis.json"), &amis)
	for _, a := range amis {
		amiMap[a.OS] = a
	}
	AMIMap = amiMap
}

func (h *Host) GetAMI() string {
	if len(h.AMI) > 0 {
		return h.AMI
	}
	if _, ok := AMIMap[h.OS]; ok {
		return AMIMap[h.OS].Regions[h.Network.Environment.AWSConfig.Region]
	}
	LogFatal(fmt.Sprintf("OS is invalid for host! host=%s os=%s", h.Hostname, h.OS))
	return ""
}

func (h *Host) RenderIP() string {
	return CustomIP(h.Network.CIDR, 0, h.LastOctet)
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
