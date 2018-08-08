package tfaws

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/gen0cide/laforge/builder/buildutil/valdations"

	"github.com/pkg/errors"

	"github.com/gen0cide/laforge/builder/buildutil"
	"github.com/gen0cide/laforge/core"
)

const (
	_id          = `tfaws`
	_name        = `Terraform AWS Builder`
	_description = `generates terraform configurations that isolate teams into VPCs`
	_author      = `Alex Levinson <github.com/gen0cide>`
	_version     = `0.0.1`
)

var (
	rules = validations.Validations{
		validations.Requirement{
			Name:       "Environment maintainer not defined",
			Resolution: "add a maintainer block to your environment configuration",
			Check:      validations.FieldNotEmpty(core.Environment{}, "Maintainer"),
		},
		validations.Requirement{
			Name:       "terraform executable not located in path",
			Resolution: "download and ensure that terraform CLI is installed to a valid location in your PATH",
			Check:      validations.ExistsInPath("terraform"),
		},
		validations.Requirement{
			Name:       "vpc CIDR not defined",
			Resolution: "define a vpc_cidr value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "vpc_cidr"),
		},
		validations.Requirement{
			Name:       "AWS Access Key not defined",
			Resolution: "define a vpc_cidr value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "aws_access_key"),
		},
		validations.Requirement{
			Name:       "AWS Secret Key not defined",
			Resolution: "define a vpc_cidr value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "aws_secret_key"),
		},
		validations.Requirement{
			Name:       "AWS Region not defined",
			Resolution: "define a vpc_cidr value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "aws_region"),
		},
		validations.Requirement{
			Name:       "No networks have been included",
			Resolution: "Use the included_network \"$network_id\" { ... } block inside of your environment config to include networks.",
			Check:      validations.FieldNotEmpty(core.Environment{}, "IncludedNetworks"),
		},
		validations.Requirement{
			Name:       "No hosts were included",
			Resolution: "Check your included_network blocks. The field included_hosts = [ ... ] should be populated with host IDs.",
			Check:      validations.FieldNotEmpty(core.Environment{}, "IncludedHosts"),
		},
		validations.Requirement{
			Name:       "No CIDR defined for network",
			Resolution: "Check that network declarations have a cidr = ... defined in them.",
			Check:      validations.FieldNotEmpty(core.Network{}, "CIDR"),
		},
		validations.Requirement{
			Name:       "No OS defined for a host",
			Resolution: "Check that all host declarations have an os = ... attribute defined.",
			Check:      validations.FieldNotEmpty(core.Host{}, "OS"),
		},
		validations.Requirement{
			Name:       "No hostname defined for a host",
			Resolution: "Check that all host declarations have a hostname = ... attribute defined.",
			Check:      validations.FieldNotEmpty(core.Host{}, "Hostname"),
		},
		validations.Requirement{
			Name:       "No Instance Size defined for a host",
			Resolution: "Check that all host declarations have an associated instance_size = ... attribute defined.",
			Check:      validations.FieldNotEmpty(core.Host{}, "InstanceSize"),
		},
		validations.Requirement{
			Name:       "No disk defined for a host",
			Resolution: "Ensure that every host declaration has an accompanied disk { size = ... } block defined.",
			Check:      validations.FieldNotEmpty(core.Host{}, "Disk"),
		},
		validations.Requirement{
			Name:       "No disk defined for a host",
			Resolution: "Ensure that every host declaration has an accompanied disk { size = ... } block defined.",
			Check:      validations.HasVarDefined(core.Host{}, "user_data_script_id"),
		},
	}
)

// TerraformAWSBuilder implements a laforge builder that packages an environment into
// a terraform configuration targeting AWS with each team isolated into their own VPC.
type TerraformAWSBuilder struct {
	sync.RWMutex

	// Required for the Builder interface
	Base *core.Laforge
}

// Get retrieves an element from the embedded KV store
func (t *TerraformAWSBuilder) Get(key string) string {
	t.Lock()
	defer t.Unlock()
	return t.Base.Build.Config[key]
}

// Set assigns an element to the embedded KV store
func (t *TerraformAWSBuilder) Set(key string, val interface{}) {
	t.Lock()
	defer t.Unlock()
	t.Base.Build.Config[key] = fmt.Sprintf("%v", val)
}

// New creates an empty TerraformAWSBuilder
func New() *TerraformAWSBuilder {
	return &TerraformAWSBuilder{}
}

// ID implements the Builder interface (returns the ID of the builder - usually the go package name)
func (t *TerraformAWSBuilder) ID() string {
	return _id
}

// Name implements the Builder interface (returns the name of the builder - usually titleized version of the type)
func (t *TerraformAWSBuilder) Name() string {
	return _name
}

// Description implements the Builder interface (returns the builder's description)
func (t *TerraformAWSBuilder) Description() string {
	return _description
}

// Author implements the Builder interface (author's name and contact info)
func (t *TerraformAWSBuilder) Author() string {
	return _author
}

// Version implements the Builder interface (builder version)
func (t *TerraformAWSBuilder) Version() string {
	return _version
}

// Validations implements the Builder interface (builder checks)
func (t *TerraformAWSBuilder) Validations() validations.Validations {
	return rules
}

// SetLaforge implements the Builder interface
func (t *TerraformAWSBuilder) SetLaforge(base *core.Laforge) error {
	t.Base = base
	if !base.ClearToBuild {
		return buildutil.Throw(errors.New("context is not cleared to build"), "Laforge has encountered an error and cannot continue to build. This is likely a bug in LaForge.", nil)
	}
	return nil
}

// CheckRequirements implements the Builder interface
func (t *TerraformAWSBuilder) CheckRequirements() error {
	return nil
}

// PrepareAssets implements the Builder interface
func (t *TerraformAWSBuilder) PrepareAssets() error {
	privkey, pubkey, err := buildutil.GenerateSSHKeyPair(2048)
	if err != nil {
		return buildutil.Throw(err, "Could not generate a 2048-bit RSA SSH key.", nil)
	}
	pathToPubkey := filepath.Join(t.Base.Build.Dir, "data", "ssh.pem.pub")
	pathToPrivkey := filepath.Join(t.Base.Build.Dir, "data", "ssh.pem")
	err = buildutil.WriteKeyfile([]byte(privkey), pathToPrivkey)
	if err != nil {
		return buildutil.Throw(err, "Could not write the the SSH private key to the build directory", &buildutil.V{"path": pathToPrivkey})
	}
	err = buildutil.WriteKeyfile([]byte(pubkey), pathToPubkey)
	if err != nil {
		return buildutil.Throw(err, "Could not write the the SSH public key to the build directory", &buildutil.V{"path": pathToPubkey})
	}
	t.Set("ssh_public_key", pathToPubkey)
	t.Set("ssh_private_key", pathToPrivkey)

	for hostid, host := range t.Base.Environment.IncludedHosts {
		uds, found := host.Vars["user_data_script_id"]
		if !found {
			return buildutil.Throw(errors.New("user_data_script_id no longer exists"), "Validation for this passed, but here we are. Likely a bug. Please report.", &buildutil.V{"host_id": hostid})
		}
		udsObj, found := t.Base.Scripts[uds]
		if !found {
			return buildutil.Throw(errors.Errorf("user_data_script_id %s not found", uds), "Host declares a user_data_script_id which was not found in the script map. Is this declared somewhere?", &buildutil.V{"host": hostid})
		}
		if _, ok := host.Scripts[uds]; ok {
			core.Logger.Infof("UDS %s is already defined for host %s (strange?)", uds, hostid)
			continue
		}
		core.Logger.Debugf("Adding user_data_script %s to host %s script pool", uds, hostid)
		host.Scripts[uds] = udsObj
	}

	return nil
}

// GenerateScripts implements the Builder interface
func (t *TerraformAWSBuilder) GenerateScripts() error {
	return nil
}

// StageDependencies implements the Builder interface
func (t *TerraformAWSBuilder) StageDependencies() error {
	return nil
}

// Render implements the Builder interface
func (t *TerraformAWSBuilder) Render() error {
	return nil
}
