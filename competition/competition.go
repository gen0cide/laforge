package competition

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Competition struct {
	R53ZoneID    string `yaml:"external_r53_zone_id"`
	AWSCred      `yaml:"aws_creds"`
	S3Config     `yaml:"s3_config"`
	AdminIPs     []string `yaml:"admin_ips"`
	RootPassword string   `yaml:"root_password"`
	DHCPConfig   `yaml:"dhcp"`
	UserList     map[string][]User
}

type DHCPConfig struct {
	DNSName     string   `yaml:"domain"`
	Nameservers []string `yaml:"nameservers"`
	NTPServers  []string `yaml:"ntp_servers"`
}

type NSRecord struct {
	Name        string   `yaml:"name"`
	Nameservers []string `yaml:"nameservers"`
}

type AWSCred struct {
	APIKey    string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`
	Region    string `yaml:"region"`
	Zone      string `yaml:"zone"`
}

type S3Config struct {
	Region string `yaml:"region"`
	Bucket string `yaml:"bucket"`
}

func (c *Competition) GetEnvs() map[*Environment]bool {
	ValidateHome()
	envs := make(map[*Environment]bool)
	files, _ := ioutil.ReadDir(filepath.Join(GetHome(), "environments", "."))
	for _, f := range files {
		if f.Name() == ".DS_Store" {
			continue
		}
		if f.Name() == ".gitkeep" {
			continue
		}
		e, err := LoadEnvironment(f.Name())
		if err != nil {
			LogError("Error parsing environments list: " + err.Error())
		} else {
			nowInUse := false
			if e.Name == GetEnv() {
				nowInUse = true
			}
			envs[e] = nowInUse
		}
	}

	return envs
}

func (c *Competition) EnvMap() map[string]*Environment {
	ValidateHome()
	envs := make(map[string]*Environment)
	files, _ := ioutil.ReadDir(filepath.Join(GetHome(), "environments", "."))
	for _, f := range files {
		if f.Name() == ".DS_Store" {
			continue
		}
		if f.Name() == ".gitkeep" {
			continue
		}
		e, err := LoadEnvironment(f.Name())
		if err != nil {
			LogError("Error parsing environments list: " + err.Error())
		} else {
			envs[f.Name()] = e
		}
	}

	return envs
}

func (c *Competition) GetEnvByName(name string) *Environment {
	return c.EnvMap()[name]
}

func (c *Competition) CurrentEnv() *Environment {
	return c.GetEnvByName(GetEnv())
}

func (c *Competition) CreateEnv(name, prefix string) {
	ValidateHome()
	if !ValidName(name) {
		LogFatal("The name you entered was invalid. Please keep names to lowercase alphanumeric and under 16 characters.")
	}
	if !ValidPrefix(prefix) {
		LogFatal("The prefix you entered was invalid. Please keep prefixes to lowercase alpha and under six characters.")
	}

	envs := c.GetEnvs()
	for e := range envs {
		if e.Name == name {
			LogFatal("The name you entered is already an environment!")
		}
		if e.Prefix == prefix {
			LogFatal("The prefix you entered is already in use! (" + e.Name + ")")
		}
	}
	e := Environment{
		Name:   name,
		Prefix: prefix,
	}
	os.MkdirAll(filepath.Join(GetHome(), "environments", name, "terraform", "scripts"), os.ModePerm)
	os.OpenFile(filepath.Join(GetHome(), "environments", name, "terraform", "scripts", ".gitkeep"), os.O_RDONLY|os.O_CREATE, 0644)
	os.MkdirAll(filepath.Join(GetHome(), "environments", name, "networks"), os.ModePerm)
	os.OpenFile(filepath.Join(GetHome(), "environments", name, "networks", ".gitkeep"), os.O_RDONLY|os.O_CREATE, 0644)
	os.MkdirAll(filepath.Join(GetHome(), "environments", name, "hosts"), os.ModePerm)
	os.OpenFile(filepath.Join(GetHome(), "environments", name, "hosts", ".gitkeep"), os.O_RDONLY|os.O_CREATE, 0644)
	var tpl bytes.Buffer
	tmpl, err := template.New(name).Parse(string(MustAsset("env.yml")))
	if err != nil {
		LogFatal("Fatal Error parsing environment config template: " + err.Error())
	}
	if err := tmpl.Execute(&tpl, e); err != nil {
		LogFatal("Fatal Error rendering environment config template: " + err.Error())
	}
	err = ioutil.WriteFile(filepath.Join(GetHome(), "environments", name, "env.yml"), tpl.Bytes(), 0644)
	if err != nil {
		LogFatal("Fatal Error writing environment config template: " + err.Error())
	}
	Log("Successfully created environment: " + name)
}

func (c *Competition) ChangeEnv(name string) {
	if !EnvDirExistsByName(name) {
		LogFatal("The environment you're trying to switch to doesn't exist!")
	}
	SetEnv(name)
	Log("Current environment is now set to " + name)
}

func (c *Competition) SSHPublicKey() string {
	data, err := ioutil.ReadFile(c.SSHPublicKeyPath())
	if err != nil {
		LogError("Could not read SSH public key for environment!")
		return ""
	}
	return strings.TrimSpace(string(data))

}

func (c *Competition) ParseScripts() map[string]*Script {
	scripts := make(map[string]*Script)
	scriptFiles, _ := filepath.Glob(filepath.Join(GetHome(), "scripts", "*"))
	for _, file := range scriptFiles {
		if filepath.Base(file) == ".gitkeep" {
			continue
		}
		script, err := LoadScript(file)
		if err != nil {
			continue
		}
		scripts[FileToName(file)] = script
	}

	return scripts
}

func (c *Competition) SSHPrivateKey() string {
	data, err := ioutil.ReadFile(c.SSHPrivateKeyPath())
	if err != nil {
		LogError("Could not read SSH private key for environment!")
		return ""
	}
	return string(data)
}

func (c *Competition) SSHPublicKeyPath() string {
	return filepath.Join(GetHome(), "config", "infra.pem.pub")
}

func (c *Competition) EmployeeDBPath() string {
	return filepath.Join(GetHome(), "config", "employees.json")
}

func (c *Competition) SSHPrivateKeyPath() string {
	return filepath.Join(GetHome(), "config", "infra.pem")
}

func LoadCompetition() (*Competition, error) {
	ValidateHome()
	comp := Competition{}
	LoadAMIs()
	compConfig, err := ioutil.ReadFile(filepath.Join(GetHome(), "config", "config.yml"))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(compConfig, &comp)
	if err != nil {
		return nil, err
	}
	comp.UserList = LoadUsersFromDB(comp.EmployeeDBPath())
	return &comp, nil
}

func LoadUsersFromDB(path string) map[string][]User {
	if !PathExists(path) {
		LogError("User Database does not exist at config/employees.json")
	}
	var users map[string][]User
	data, err := ioutil.ReadFile(path)
	if err != nil {
		LogFatal(err.Error())
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		LogFatal("Could not unmarshal users file: " + path)
	}
	return users
}
