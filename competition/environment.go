package competition

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Vars map[string]string

type Environment struct {
	Name         string   `yaml:"name"`
	Prefix       string   `yaml:"prefix"`
	WhitelistIPs []string `yaml:"ip_whitelist"`
	Vars         `yaml:"variables"`
	AWSConfig    `yaml:"aws_config"`
	GCPConfig    `yaml:"gcp_config"`
	PodCount     int `yaml:"pod_count"`
	Pod          `yaml:"pod"`
	Competition  `yaml:"-"`
}

type AWSConfig struct {
	CIDR   string `yaml:"cidr"`
	Region string `yaml:"region"`
	Zone   string `yaml:"zone"`
}

type GCPConfig struct {
	CIDR   string `yaml:"cidr"`
	Region string `yaml:"region"`
	Zone   string `yaml:"zone"`
}

func LoadEnvironment(name string) (*Environment, error) {
	envConfigFile := filepath.Join(GetHome(), "environments", name, "env.yml")
	if !PathExists(envConfigFile) {
		return nil, errors.New("not a valid environment: no env.yml file")
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
	return &env, nil
}
