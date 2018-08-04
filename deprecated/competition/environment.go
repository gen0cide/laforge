package competition

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bradfitz/iter"
	"github.com/hashicorp/terraform/terraform"

	yaml "gopkg.in/yaml.v2"
)

type Vars map[string]string

type Environment struct {
	Name             string `yaml:"name"`
	Prefix           string `yaml:"prefix"`
	Vars             `yaml:"variables"`
	PodCount         int                 `yaml:"pod_count"`
	IncludedNetworks []string            `yaml:"included_networks"`
	ResolvedNetworks map[string]*Network `yaml:"-"`
	Competition      `yaml:"-"`
	Users            []*User    `yaml:"-"`
	Networks         []*Network `yaml:"-"`
	Hosts            []*Host    `yaml:"-"`
	JumpHosts        `yaml:"jump_hosts"`
	GenesisHost      *Host `yaml:"genesis_host"`
}

type JumpHosts struct {
	CIDR    string `yaml:"cidr"`
	Windows struct {
		AMI     string   `yaml:"ami"`
		Count   int      `yaml:"count"`
		Size    string   `yaml:"size"`
		Scripts []string `yaml:"scripts"`
	} `yaml:"windows"`
	Kali struct {
		AMI     string   `yaml:"ami"`
		Count   int      `yaml:"count"`
		Size    string   `yaml:"size"`
		Scripts []string `yaml:"scripts"`
	} `yaml:"kali"`
}

type SecurityGroup struct {
	Protocol string
	FromPort int
	ToPort   int
}

type SSHConfig struct {
	SSHKey string
	Hosts  map[string]string
}

func (e *Environment) EnvRoot() string {
	return filepath.Join(GetHome(), "environments", e.Name)
}

func (e *Environment) NetworksDir() string {
	return filepath.Join(e.EnvRoot(), "networks")
}

func (e *Environment) HostsDir() string {
	return filepath.Join(e.EnvRoot(), "hosts")
}

func (e *Environment) TfDir() string {
	return filepath.Join(e.EnvRoot(), "terraform")
}

func (e *Environment) TeamIDs() []int {
	teams := []int{}
	for i := range iter.N(e.PodCount) {
		teams = append(teams, i)
	}
	return teams
}

func (e *Environment) TfDirForTeam(teamID int) string {
	return filepath.Join(e.TfDir(), strconv.Itoa(teamID))
}

func (e *Environment) TfFile(teamID int) string {
	return filepath.Join(e.TfDirForTeam(teamID), "infra.tf")
}

func (e *Environment) TfStateFile(teamID int) string {
	return filepath.Join(e.TfDirForTeam(teamID), "terraform.tfstate")
}

func (e *Environment) TfScriptsDir(teamID int) string {
	return filepath.Join(e.TfDirForTeam(teamID), "scripts")
}

func (e *Environment) SSHConfigPath() string {
	return filepath.Join(e.EnvRoot(), "ssh.conf")
}

func (e *Environment) DefaultCIDR() string {
	return "10.0.0.0/16"
}

func (e *Environment) PodPassword(podID int) string {
	return DeterminedPassword(fmt.Sprintf("%s-%d", e.Name, podID))
}

func LoadEnvironment(name string) (*Environment, error) {
	envConfigFile := filepath.Join(GetHome(), "environments", name, "env.yml")
	envNetworkPath := filepath.Join(GetHome(), "environments", name, "networks")
	envHostPath := filepath.Join(GetHome(), "environments", name, "hosts")
	if !PathExists(envConfigFile) {
		return nil, fmt.Errorf("not a valid environment: no env.yml file located at %s", name)
	}
	if !PathExists(envNetworkPath) {
		return nil, errors.New("not a valid environment: no networks directory located at " + envNetworkPath)
	}
	if !PathExists(envHostPath) {
		return nil, errors.New("not a valid environment: no hosts directory located at " + envHostPath)
	}
	env := Environment{}
	envConfig, err := ioutil.ReadFile(envConfigFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(envConfig, &env)
	if err != nil {
		return nil, err
	}
	env.ResolvedNetworks = env.ResolveIncludedNetworks()
	for _, network := range env.ResolvedNetworks {
		network.ResolvedHosts = network.ResolveIncludedHosts()
	}
	return &env, nil
}

func (e *Environment) NewSSHConfig() *SSHConfig {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	teams := e.TeamIDs()
	ipMap := map[string]string{}

	for _, t := range teams {
		fileData, err := ioutil.ReadFile(e.TfStateFile(t))
		if err != nil {
			LogError("Cannot Read terraform.tfstate. You've likely not provisioned that team. team=" + strconv.Itoa(t))
			continue
		}
		tfState, err := terraform.ReadStateV3(fileData)
		if err != nil {
			LogFatal("Fatal Error Parsing Terraform State. tfstate=" + e.TfStateFile(t) + " error=" + err.Error())
		}
		for _, module := range tfState.Modules {
			for outputKey, outputState := range module.Outputs {
				if strings.Contains(outputKey, "public_ips") {
					explode := strings.Split(outputKey, ".")
					ipMap[explode[1]] = outputState.Value.(string)
				}
			}
		}
	}

	sshConf := SSHConfig{
		SSHKey: e.Competition.SSHPrivateKeyPath(),
		Hosts:  ipMap,
	}

	return &sshConf
}

func (e *Environment) GenerateSSHConfig() {
	sshConf := e.NewSSHConfig()

	tmp := template.New(RandomString(entropySize))
	newTmpl, err := tmp.Parse(string(MustAsset("ssh.conf")))
	if err != nil {
		LogFatal("Could not parse SSH Config Template.")
	}

	var tpl bytes.Buffer

	if err := newTmpl.Execute(&tpl, sshConf); err != nil {
		LogFatal("Could not render SSH Config Template.")
	}

	err = ioutil.WriteFile(e.SSHConfigPath(), tpl.Bytes(), 0644)
	if err != nil {
		LogFatal("Error Writing SSH Config: " + err.Error())
	}
}

func (e *Environment) ParseHosts() map[string]*Host {
	hosts := make(map[string]*Host)
	hostFiles, _ := filepath.Glob(filepath.Join(e.HostsDir(), "*.yml"))
	for _, file := range hostFiles {
		if filepath.Base(file) == ".gitkeep" {
			continue
		}
		host, err := LoadHostFromFile(file)
		if err != nil {
			LogError("Error reading host file: " + file + " error=" + err.Error())
			continue
		}
		hosts[FileToName(file)] = host
	}

	return hosts
}

func (e *Environment) ParseNetworks() map[string]*Network {
	networks := make(map[string]*Network)
	networkFiles, _ := filepath.Glob(filepath.Join(e.NetworksDir(), "*.yml"))
	for _, file := range networkFiles {
		if filepath.Base(file) == ".gitkeep" {
			continue
		}
		network, err := LoadNetworkFromFile(file)
		if err != nil {
			LogError("Error reading network file: " + file + " error=" + err.Error())
			continue
		}
		networks[FileToName(file)] = network
	}

	return networks
}

func (e *Environment) ResolveIncludedNetworks() map[string]*Network {
	networks := make(map[string]*Network)
	networkFiles, _ := filepath.Glob(filepath.Join(e.NetworksDir(), "*.yml"))
	for _, file := range networkFiles {
		if filepath.Base(file) == ".gitkeep" {
			continue
		}
		if !Contains(FileToName(file), e.IncludedNetworks) {
			continue
		}
		network, err := LoadNetworkFromFile(file)
		if err != nil {
			LogError("Error reading network file: " + file)
			continue
		}
		network.Environment = *e
		networks[FileToName(file)] = network
	}

	return networks
}

func (e *Environment) ResolvePublicTCP() map[string]SecurityGroup {
	portMap := map[string]SecurityGroup{}
	networks := e.ResolveIncludedNetworks()
	for _, network := range networks {
		hosts := network.ResolveIncludedHosts()
		for _, host := range hosts {
			for _, port := range host.TCPPorts {
				i, err := strconv.Atoi(port)
				if err == nil {
					portMap[port] = SecurityGroup{
						Protocol: "tcp",
						FromPort: i,
						ToPort:   i,
					}
					continue
				} else {
					matched, err := regexp.MatchString("^\\d+-\\d+$", port)
					if err == nil && matched {
						portRange := strings.Split(port, "-")
						fp, err := strconv.Atoi(portRange[0])
						if err != nil {
							LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=tcp", port, host.Hostname))
							continue
						}
						tp, err := strconv.Atoi(portRange[1])
						if err != nil {
							LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=tcp", port, host.Hostname))
							continue
						}
						portMap[port] = SecurityGroup{
							Protocol: "tcp",
							FromPort: fp,
							ToPort:   tp,
						}
						continue
					}
					LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=tcp", port, host.Hostname))
					continue
				}
			}
		}
	}
	for _, port := range e.GenesisHost.TCPPorts {
		i, err := strconv.Atoi(port)
		if err == nil {
			portMap[port] = SecurityGroup{
				Protocol: "tcp",
				FromPort: i,
				ToPort:   i,
			}
			continue
		} else {
			matched, err := regexp.MatchString("^\\d+-\\d+$", port)
			if err == nil && matched {
				portRange := strings.Split(port, "-")
				fp, err := strconv.Atoi(portRange[0])
				if err != nil {
					LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=tcp (genesis host)", port, e.GenesisHost.Hostname))
					continue
				}
				tp, err := strconv.Atoi(portRange[1])
				if err != nil {
					LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=tcp (genesis host)", port, e.GenesisHost.Hostname))
					continue
				}
				portMap[port] = SecurityGroup{
					Protocol: "tcp",
					FromPort: fp,
					ToPort:   tp,
				}
				continue
			}
			LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=tcp (genesis host)", port, e.GenesisHost.Hostname))
			continue
		}
	}
	return portMap
}

func (e *Environment) ResolvePublicUDP() map[string]SecurityGroup {
	portMap := map[string]SecurityGroup{}
	networks := e.ResolveIncludedNetworks()
	for _, network := range networks {
		hosts := network.ResolveIncludedHosts()
		for _, host := range hosts {
			for _, port := range host.UDPPorts {
				i, err := strconv.Atoi(port)
				if err == nil {
					portMap[port] = SecurityGroup{
						Protocol: "udp",
						FromPort: i,
						ToPort:   i,
					}
					continue
				} else {
					matched, err := regexp.MatchString("^\\d+-\\d+$", port)
					if err == nil && matched {
						portRange := strings.Split(port, "-")
						fp, err := strconv.Atoi(portRange[0])
						if err != nil {
							LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=udp", port, host.Hostname))
							continue
						}
						tp, err := strconv.Atoi(portRange[1])
						if err != nil {
							LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=udp", port, host.Hostname))
							continue
						}
						portMap[port] = SecurityGroup{
							Protocol: "udp",
							FromPort: fp,
							ToPort:   tp,
						}
						continue
					}
					LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=udp", port, host.Hostname))
					continue
				}
			}
		}
	}
	for _, port := range e.GenesisHost.UDPPorts {
		i, err := strconv.Atoi(port)
		if err == nil {
			portMap[port] = SecurityGroup{
				Protocol: "udp",
				FromPort: i,
				ToPort:   i,
			}
			continue
		} else {
			matched, err := regexp.MatchString("^\\d+-\\d+$", port)
			if err == nil && matched {
				portRange := strings.Split(port, "-")
				fp, err := strconv.Atoi(portRange[0])
				if err != nil {
					LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=tcp (genesis host)", port, e.GenesisHost.Hostname))
					continue
				}
				tp, err := strconv.Atoi(portRange[1])
				if err != nil {
					LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=tcp (genesis host)", port, e.GenesisHost.Hostname))
					continue
				}
				portMap[port] = SecurityGroup{
					Protocol: "udp",
					FromPort: fp,
					ToPort:   tp,
				}
				continue
			}
			LogError(fmt.Sprintf("invalid port: port=%s host=%s protocol=tcp (genesis host)", port, e.GenesisHost.Hostname))
			continue
		}
	}
	return portMap
}

func (e *Environment) Suffix(podOffset int) string {
	return fmt.Sprintf("%s%d", e.Prefix, podOffset)
}

func (e *Environment) TFName(name string, offset int) string {
	return fmt.Sprintf("%s_%s", name, e.Suffix(offset))
}

func (e *Environment) CreateNetwork(n *Network) {
	networks := e.ParseNetworks()
	if networks[n.Name] != nil {
		LogFatal("You cannot create a network that already exists!")
	}

	var tpl bytes.Buffer
	tmpl, err := template.New(n.Name).Parse(string(MustAsset("network.yml")))
	if err != nil {
		LogFatal("Fatal error parsing network config template: " + err.Error())
	}
	if err := tmpl.Execute(&tpl, n); err != nil {
		LogFatal("Fatal error rendering network config: " + err.Error())
	}
	filename := filepath.Join(e.NetworksDir(), strings.ToLower(n.Name)+".yml")
	err = ioutil.WriteFile(filename, tpl.Bytes(), 0644)
	if err != nil {
		LogFatal("Error writing network config file: " + err.Error())
	}
	Log("Network created: " + n.Name)
}

func (e *Environment) CreateHost(h *Host) {
	hosts := e.ParseHosts()
	if hosts[strings.ToLower(h.Hostname)] != nil {
		LogFatal("You cannot create a host that already exists!")
	}

	var tpl bytes.Buffer
	tmpl, err := template.New(h.Hostname).Parse(string(MustAsset("host.yml")))
	if err != nil {
		LogFatal("Fatal error parsing host config template: " + err.Error())
	}
	if err := tmpl.Execute(&tpl, h); err != nil {
		LogFatal("Fatal error rendering host config: " + err.Error())
	}
	filename := filepath.Join(e.HostsDir(), strings.ToLower(h.Hostname)+".yml")
	err = ioutil.WriteFile(filename, tpl.Bytes(), 0644)
	if err != nil {
		LogFatal("Error writing host config file: " + err.Error())
	}
	Log("Host created: " + h.Hostname)
}

func (e *Environment) KaliJumpAMI() string {
	if e.JumpHosts.Kali.AMI != "" {
		return e.JumpHosts.Kali.AMI
	}
	return AMIMap["kali"].Regions[e.AWS.Region]
}

func (e *Environment) WindowsJumpAMI() string {
	if e.JumpHosts.Windows.AMI != "" {
		return e.JumpHosts.Windows.AMI
	}
	return AMIMap["w2k16"].Regions[e.AWS.Region]
}
