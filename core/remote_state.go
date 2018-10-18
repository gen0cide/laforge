package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gruntwork-io/terragrunt/util"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

// TODO: this file could be changed to use the Terraform Go code to read state files, but that code is relatively
// complicated and doesn't seem to be designed for standalone use. Fortunately, the .tfstate format is a fairly simple
// JSON format, so hopefully this simple parsing code will not be a maintenance burden.

const (
	// DefaultPathToLocalStateFile is the path to a local copy of the state file for a local state
	DefaultPathToLocalStateFile = "terraform.tfstate"

	// DefaultPathToRemoteStateFile is the path to a local copy of the state file for a remote state
	DefaultPathToRemoteStateFile = ".terraform/terraform.tfstate"
)

var (
	// PrivateIPOutputRegexp matches outputs for private IP addresses
	PrivateIPOutputRegexp = regexp.MustCompile(`\A(\w+)\.private_ip\z`)

	// PublicIPOutputRegexp matches outputs for private IP addresses
	PublicIPOutputRegexp = regexp.MustCompile(`\A(\w+)\.public_ip\z`)
)

// RemoteState is used to gather data stored in the remote etcd backend of a laforge configured environment
type RemoteState struct {
	Base         *Laforge
	LocalConfig  *TerraformState
	RemoteConfig *TerraformState
	Hosts        map[string]*HostTableInfo
}

// GetState retrieves a remote configuration of a laforge deployed terraform environment
func GetState(base *Laforge) (*RemoteState, error) {
	configFile := filepath.Join(base.TeamRoot, ".terraform", "terraform.tfstate")
	localconfig, err := ParseTerraformStateFile(configFile)
	if err != nil {
		return nil, err
	}

	endpoints := []string{}

	for _, x := range localconfig.Backend.Config["endpoints"].([]interface{}) {
		endpoints = append(endpoints, x.(string))
	}

	etcdclient, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		Username:    localconfig.Backend.Config["username"].(string),
		Password:    localconfig.Backend.Config["password"].(string),
		DialTimeout: 10 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	defer etcdclient.Close()

	key := fmt.Sprintf("%sdefault", localconfig.Backend.Config["prefix"].(string))
	resp, err := etcdclient.Get(context.TODO(), key)
	if err != nil {
		return nil, err
	}
	var data []byte
	for _, ev := range resp.Kvs {
		if string(ev.Key) == key {
			data = ev.Value
			break
		}
	}

	if len(data) == 0 {
		return nil, errors.New("remote state was empty for etcd key")
	}

	state, err := ParseTerraformStateFromBytes(data)
	if err != nil {
		return nil, err
	}

	rs := &RemoteState{
		Base:         base,
		LocalConfig:  localconfig,
		RemoteConfig: state,
		Hosts:        map[string]*HostTableInfo{},
	}

	rs.ParseHostInfo()

	return rs, nil
}

// ParseHostInfo attempts to parse host output information from the terraform configuration
func (r *RemoteState) ParseHostInfo() {
	if r.RemoteConfig == nil {
		return
	}
	for _, m := range r.RemoteConfig.Modules {
		for k, v := range m.Outputs {
			sm, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if PublicIPOutputRegexp.MatchString(k) {
				hostname := strings.Replace(k, `.public_ip`, ``, -1)
				hostinfo, found := r.Hosts[hostname]
				if !found {
					hostinfo = &HostTableInfo{Hostname: hostname}
					r.Hosts[hostname] = hostinfo
				}
				hostinfo.PublicIP = sm["value"].(string)
				continue
			}
			if PrivateIPOutputRegexp.MatchString(k) {
				hostname := strings.Replace(k, `.private_ip`, ``, -1)
				hostinfo, found := r.Hosts[hostname]
				if !found {
					hostinfo = &HostTableInfo{Hostname: hostname}
					r.Hosts[hostname] = hostinfo
				}
				hostinfo.PrivateIP = sm["value"].(string)
				continue
			}
		}
	}
}

// HostTableInfo represents host information parsed out of a remote state
type HostTableInfo struct {
	PublicIP  string `json:"public_ip"`
	PrivateIP string `json:"private_ip"`
	Hostname  string `json:"hostname"`
}

// TableInfo is a helper function for pretty pretting host information
func (h *HostTableInfo) TableInfo() []string {
	return []string{
		h.Hostname,
		h.PublicIP,
		h.PrivateIP,
	}
}

// TerraformState represents the structure of the Terraform .tfstate file
type TerraformState struct {
	Version int
	Serial  int
	Backend *TerraformBackend
	Modules []TerraformStateModule
}

// TerraformBackend represents the structure of the "backend" section of the Terraform .tfstate file
type TerraformBackend struct {
	Type   string
	Config map[string]interface{}
}

// TerraformStateModule represents the structure of a "module" section of the Terraform .tfstate file
type TerraformStateModule struct {
	Path      []string
	Outputs   map[string]interface{}
	Resources map[string]interface{}
}

// IsRemote returns true if this Terraform state is configured for remote state storage
func (state *TerraformState) IsRemote() bool {
	return state.Backend != nil && state.Backend.Type != "local"
}

// ParseTerraformStateFileFromLocation parses the Terraform .tfstate file. If a local backend is used then search the given path, or
// return nil if the file is missing. If the backend is not local then parse the Terraform .tfstate
// file from the location specified by workingDir. If no location is specified, search the current
// directory. If the file doesn't exist at any of the default locations, return nil.
func ParseTerraformStateFileFromLocation(backend string, config map[string]interface{}, workingDir string) (*TerraformState, error) {
	stateFile, ok := config["path"].(string)
	if backend == "local" && ok && util.FileExists(stateFile) {
		return ParseTerraformStateFile(stateFile)
	} else if util.FileExists(util.JoinPath(workingDir, DefaultPathToLocalStateFile)) {
		return ParseTerraformStateFile(util.JoinPath(workingDir, DefaultPathToLocalStateFile))
	} else if util.FileExists(util.JoinPath(workingDir, DefaultPathToRemoteStateFile)) {
		return ParseTerraformStateFile(util.JoinPath(workingDir, DefaultPathToRemoteStateFile))
	}
	return nil, nil
}

// ParseTerraformStateFile parses the Terraform .tfstate file at the given path
func ParseTerraformStateFile(path string) (*TerraformState, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(CantParseTerraformStateFile{Path: path, UnderlyingErr: err}, "failure to read terraform state file")
	}

	return ParseTerraformStateFromBytes(bytes)
}

// ParseTerraformStateFromBytes parses the Terraform state file data in the given byte slice
func ParseTerraformStateFromBytes(terraformStateData []byte) (*TerraformState, error) {
	terraformState := &TerraformState{}

	if err := json.Unmarshal(terraformStateData, terraformState); err != nil {
		return nil, errors.Wrap(err, "failure to parse bytes into terraform state")
	}

	return terraformState, nil
}

// CantParseTerraformStateFile is an error type used to represent a failure to parse a terraform state file
type CantParseTerraformStateFile struct {
	Path          string
	UnderlyingErr error
}

// Error implements the error interface
func (err CantParseTerraformStateFile) Error() string {
	return fmt.Sprintf("Error parsing Terraform state file %s: %s", err.Path, err.UnderlyingErr.Error())
}
