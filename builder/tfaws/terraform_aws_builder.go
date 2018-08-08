// Package tfaws implements a Laforge Builder module for generating terraform configurations that target AWS.
package tfaws

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/gen0cide/laforge/builder/tfaws/static"

	"github.com/gen0cide/laforge/builder/buildutil/templates"
	"github.com/gen0cide/laforge/builder/buildutil/valdations"

	"github.com/pkg/errors"

	"github.com/gen0cide/laforge/builder/buildutil"
	"github.com/gen0cide/laforge/core"
)

// Definition of builder meta-data.
const (
	ID          = `tfaws`
	Name        = `Terraform AWS Builder`
	Description = `generates terraform configurations that isolate teams into VPCs`
	Author      = `Alex Levinson <github.com/gen0cide>`
	Version     = `0.0.1`
)

var (
	rules = validations.Validations{
		validations.Requirement{
			Name:       "Environment maintainer not defined",
			Resolution: "add a maintainer block to your environment configuration",
			Check:      validations.FieldNotEmpty(core.Environment{}, "Maintainer"),
		},
		validations.Requirement{
			Name:       "DNS not defined",
			Resolution: "add a DNS block to your competition configuration",
			Check:      validations.FieldNotEmpty(core.Competition{}, "DNS"),
		},
		validations.Requirement{
			Name:       "DNS type not listed as route53",
			Resolution: "Make sure your dns block declaration has route53 as it's type.",
			Check:      validations.FieldEquals(core.DNS{}, "Type", "route53"),
		},
		validations.Requirement{
			Name:       "DNS Root Domain not defined",
			Resolution: "set the root_domain parameter in your DNS config block",
			Check:      validations.FieldNotEmpty(core.DNS{}, "RootDomain"),
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
			Name:       "vdi Network CIDR not defined",
			Resolution: "define a vdi_network_cidr value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "vdi_network_cidr"),
		},
		validations.Requirement{
			Name:       "no teams specified",
			Resolution: "make sure to set your team_count inside your environment config block to at least 1.",
			Check:      validations.FieldNotEmpty(core.Environment{}, "team_count"),
		},
		validations.Requirement{
			Name:       "admin IP not defined",
			Resolution: "define an admin_ip value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "admin_ip"),
		},
		validations.Requirement{
			Name:       "AWS Access Key not defined",
			Resolution: "define a aws_access_key value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "aws_access_key"),
		},
		validations.Requirement{
			Name:       "AWS Secret Key not defined",
			Resolution: "define a aws_secret_key value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "aws_secret_key"),
		},
		validations.Requirement{
			Name:       "AWS Region not defined",
			Resolution: "define a aws_region value inside your environment config = { ... } block.",
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
			Name:       "No user_data_script_id defined for a host",
			Resolution: "Ensure that every host declaration has a var defined for key user_data_script_id.",
			Check:      validations.HasVarDefined(core.Host{}, "user_data_script_id"),
		},
		validations.Requirement{
			Name:       "No AMI defined for host",
			Resolution: "Ensure that every host declaration has an accompanied ami_id var set",
			Check:      validations.HasVarDefined(core.Host{}, "ami_id"),
		},
		validations.Requirement{
			Name:       "IP needs to be overridden",
			Resolution: "Ensure that every host declaration has an accompanied ip_override var set",
			Check:      validations.HasVarDefined(core.Host{}, "ip_override"),
		},
	}

	primaryTemplate = "demo_infra.tf.tmpl"
	templatesToLoad = []string{
		primaryTemplate,
	}
)

// TerraformAWSBuilder implements a laforge builder that packages an environment into
// a terraform configuration targeting AWS with each team isolated into their own VPC.
type TerraformAWSBuilder struct {
	sync.RWMutex

	// Required for the Builder interface
	Base *core.Laforge

	// A place to store the templates
	Library *templates.Library
}

// Get retrieves an element from the embedded KV store
func (t *TerraformAWSBuilder) Get(key string) string {
	t.Lock()
	defer t.Unlock()
	res, ok := t.Base.Build.Config[key]
	if ok {
		return res
	}
	r0, e0 := t.Base.Environment.Config[key]
	if e0 {
		defer t.Set(key, r0)
		return r0
	}
	return ""
}

// Set assigns an element to the embedded KV store
func (t *TerraformAWSBuilder) Set(key string, val interface{}) {
	t.Lock()
	defer t.Unlock()
	t.Base.Build.Config[key] = fmt.Sprintf("%v", val)
}

// New creates an empty TerraformAWSBuilder
func New() *TerraformAWSBuilder {
	lib := templates.NewLibrary()
	return &TerraformAWSBuilder{
		Library: lib,
	}
}

// ID implements the Builder interface (returns the ID of the builder - usually the go package name)
func (t *TerraformAWSBuilder) ID() string {
	return ID
}

// Name implements the Builder interface (returns the name of the builder - usually titleized version of the type)
func (t *TerraformAWSBuilder) Name() string {
	return Name
}

// Description implements the Builder interface (returns the builder's description)
func (t *TerraformAWSBuilder) Description() string {
	return Description
}

// Author implements the Builder interface (author's name and contact info)
func (t *TerraformAWSBuilder) Author() string {
	return Author
}

// Version implements the Builder interface (builder version)
func (t *TerraformAWSBuilder) Version() string {
	return Version
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
	for _, x := range templatesToLoad {
		d, err := static.ReadFile(x)
		if err != nil {
			return buildutil.Throw(err, "could not read template", &buildutil.V{"template_name": x})
		}
		_, err = t.Library.AddBook(x, d)
		if err != nil {
			return buildutil.Throw(err, "could not parse template", &buildutil.V{"template_name": x})
		}
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
	t.Set("ssh_public_key_file", pathToPubkey)
	t.Set("ssh_private_key_file", pathToPrivkey)
	t.Set("ssh_public_key", pubkey)
	t.Set("ssh_private_key", privkey)

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
	for i := 0; i < t.Base.Environment.TeamCount; i++ {
		teamDir := filepath.Join(t.Base.EnvRoot, "build", "teams", fmt.Sprintf("%v", i))
		os.MkdirAll(teamDir, 0755)
		core.TouchGitKeep(teamDir)
	}
	return nil
}

// Render implements the Builder interface
func (t *TerraformAWSBuilder) Render() error {
	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	for i := 0; i < t.Base.Environment.TeamCount; i++ {
		wg.Add(1)
		go func(teamid int) {
			defer wg.Done()
			teamDir := filepath.Join(t.Base.EnvRoot, "build", "teams", fmt.Sprintf("%v", teamid))
			team := &core.Team{
				TeamNumber:    teamid,
				BuildID:       t.Base.Build.ID,
				Build:         t.Base.Build,
				EnvironmentID: t.Base.Environment.ID,
				Environment:   t.Base.Environment,
				Maintainer:    &t.Base.User,
				RelBuildPath:  teamDir,
			}
			// t.Base.Build.Teams[teamid] = team
			team.ID = team.Name()
			user := t.Base.User
			t.Base.Team = team
			ctx, err := templates.NewContext(
				t.Base,
				t.Base.Build,
				t.Base.Competition,
				t.Base.Competition.DNS,
				t.Base.Environment,
				&user,
				team,
			)
			if err != nil {
				errChan <- err
				return
			}
			cfgData, err := t.Library.Execute(primaryTemplate, ctx)
			if err != nil {
				errChan <- buildutil.Throw(err, "template failed", &buildutil.V{
					"team": teamid,
					"dir":  teamDir,
				})
				return
			}
			cfgFile := filepath.Join(teamDir, "infra.tf")
			err = ioutil.WriteFile(cfgFile, cfgData, 0644)
			if err != nil {
				errChan <- err
				return
			}
			teamCfg, err := core.RenderHCLv2Object(team)
			if err != nil {
				errChan <- err
				return
			}
			teamCfgFile := filepath.Join(teamDir, "team.laforge")
			err = ioutil.WriteFile(teamCfgFile, teamCfg, 0644)
			if err != nil {
				errChan <- err
				return
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(finChan)
	}()

	select {
	case <-finChan:
		return nil
	case err := <-errChan:
		return err
	}
}
