// Package tfibm implements a Laforge Builder module for generating terraform configurations that target IBM Softlayer.
package tfibm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/hcl/hcl/printer"

	"github.com/gen0cide/laforge/builder/tfibm/static"

	"github.com/gen0cide/laforge/builder/buildutil/templates"
	"github.com/gen0cide/laforge/builder/buildutil/valdations"

	"github.com/pkg/errors"

	"github.com/gen0cide/laforge/builder/buildutil"
	"github.com/gen0cide/laforge/core"
)

// Definition of builder meta-data.
const (
	ID          = `tfibm`
	Name        = `Terraform IBM Softlayer Builder`
	Description = `generates terraform configurations that uses IBM Softlayer`
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
			Name:       "terraform executable not located in path",
			Resolution: "download and ensure that terraform CLI is installed to a valid location in your PATH",
			Check:      validations.ExistsInPath("terraform"),
		},
		validations.Requirement{
			Name:       "etcd server password not defined",
			Resolution: "define a etcd_password attribute in the environment configuration block.",
			Check:      validations.HasConfigKey(core.Environment{}, "etcd_password"),
		},
		validations.Requirement{
			Name:       "etcd username not defined",
			Resolution: "define an etcd_username attribute in the environment configuration block.",
			Check:      validations.HasConfigKey(core.Environment{}, "etcd_username"),
		},
		validations.Requirement{
			Name:       "etcd master server not defined",
			Resolution: "define a etcd_master (host:port) attribute in the environment configuration block.",
			Check:      validations.HasConfigKey(core.Environment{}, "etcd_master"),
		},
		validations.Requirement{
			Name:       "etcd slave server not defined",
			Resolution: "define a etcd_slave (host:port) attribute in the environment configuration block.",
			Check:      validations.HasConfigKey(core.Environment{}, "etcd_slave"),
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
			Name:       "Softlayer API User not defined",
			Resolution: "define a softlayer_api_user value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "softlayer_api_user"),
		},
		validations.Requirement{
			Name:       "Softlayer API Key not defined",
			Resolution: "define a softlayer_api_key value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "softlayer_api_key"),
		},
		validations.Requirement{
			Name:       "Softlayer datacenter not defined",
			Resolution: "define a datacenter value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "datacenter"),
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
			Name:       "No user_data_script_id defined for a host",
			Resolution: "Ensure that every host declaration has a var defined for key user_data_script_id.",
			Check:      validations.HasVarDefined(core.Host{}, "user_data_script_id"),
		},
	}

	templatesToLoad = []string{
		"infra.tf.tmpl",
		"command.tf.tmpl",
		"script.tf.tmpl",
		"remote_file.tf.tmpl",
		"dns_record.tf.tmpl",
	}

	primaryTemplate = "infra.tf.tmpl"
)

// TerraformIBMBuilder implements a laforge builder that packages an environment into
// a terraform configuration targeting AWS with each team isolated into their own VPC.
type TerraformIBMBuilder struct {
	sync.RWMutex

	// Required for the Builder interface
	Base *core.Laforge

	// A place to store the templates
	Library *templates.Library
}

// Get retrieves an element from the embedded KV store
func (t *TerraformIBMBuilder) Get(key string) string {
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
func (t *TerraformIBMBuilder) Set(key string, val interface{}) {
	t.Lock()
	defer t.Unlock()
	t.Base.Build.Config[key] = fmt.Sprintf("%v", val)
}

// New creates an empty TerraformIBMBuilder
func New() *TerraformIBMBuilder {
	lib := templates.NewLibrary()
	return &TerraformIBMBuilder{
		Library: lib,
	}
}

// ID implements the Builder interface (returns the ID of the builder - usually the go package name)
func (t *TerraformIBMBuilder) ID() string {
	return ID
}

// Name implements the Builder interface (returns the name of the builder - usually titleized version of the type)
func (t *TerraformIBMBuilder) Name() string {
	return Name
}

// Description implements the Builder interface (returns the builder's description)
func (t *TerraformIBMBuilder) Description() string {
	return Description
}

// Author implements the Builder interface (author's name and contact info)
func (t *TerraformIBMBuilder) Author() string {
	return Author
}

// Version implements the Builder interface (builder version)
func (t *TerraformIBMBuilder) Version() string {
	return Version
}

// Validations implements the Builder interface (builder checks)
func (t *TerraformIBMBuilder) Validations() validations.Validations {
	return rules
}

// SetLaforge implements the Builder interface
func (t *TerraformIBMBuilder) SetLaforge(base *core.Laforge) error {
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
func (t *TerraformIBMBuilder) CheckRequirements() error {
	return nil
}

// PrepareAssets implements the Builder interface
func (t *TerraformIBMBuilder) PrepareAssets() error {
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
	t.Set("rel_ssh_public_key_file", "../../data/ssh.pem.pub")
	t.Set("rel_ssh_private_key_file", "../../data/ssh.pem")
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
func (t *TerraformIBMBuilder) GenerateScripts() error {
	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	user := t.Base.User
	ctx, err := templates.NewContext(
		t.Base,
		t.Base.Build,
		t.Base.Competition,
		t.Base.Competition.DNS,
		t.Base.Environment,
		&user,
	)
	if err != nil {
		return err
	}
	for tid, teamObj := range t.Base.Build.Teams {
		wg.Add(1)
		go func(teamNum int, team *core.Team) {
			defer wg.Done()
			for netName, hosts := range t.Base.Environment.HostByNetwork {
				network := t.Base.Environment.IncludedNetworks[netName]
				for _, host := range hosts {
					for sid, script := range host.Scripts {
						wg.Add(1)
						go func(scriptID string, scriptObj *core.Script, hostObj *core.Host) {
							defer wg.Done()
							scriptCtx := ctx.Clone()
							err := scriptCtx.Attach(team, network, hostObj, scriptObj)
							if err != nil {
								errChan <- err
								return
							}
							filename := fmt.Sprintf("%s_%s", hostObj.Hostname, filepath.Base(scriptObj.Source))
							assetDir := filepath.Join(team.RelBuildPath, "assets")
							assetPath := filepath.Join(assetDir, filename)
							fileData, err := t.Library.Execute(fmt.Sprintf("script_%s", scriptID), scriptCtx)
							if err != nil {
								errChan <- err
								return
							}
							err = ioutil.WriteFile(assetPath, fileData, 0644)
							if err != nil {
								errChan <- err
								return
							}
							return
						}(sid, script, host)
					}
				}
			}

		}(tid, teamObj)
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

// StageDependencies implements the Builder interface
func (t *TerraformIBMBuilder) StageDependencies() error {
	for i := 0; i < t.Base.Environment.TeamCount; i++ {
		teamDir := filepath.Join(t.Base.EnvRoot, "build", "teams", fmt.Sprintf("%v", i))
		team := &core.Team{
			TeamNumber:    i,
			BuildID:       t.Base.Build.ID,
			Build:         t.Base.Build,
			EnvironmentID: t.Base.Environment.ID,
			Environment:   t.Base.Environment,
			Maintainer:    &t.Base.User,
			RelBuildPath:  teamDir,
		}
		t.Base.Build.Teams[i] = team
		assetDir := filepath.Join(teamDir, "assets")
		os.MkdirAll(assetDir, 0755)
		core.TouchGitKeep(assetDir)
	}

	for _, host := range t.Base.Environment.IncludedHosts {
		for sid, script := range host.Scripts {
			if _, ok := t.Library.Books[sid]; ok {
				continue
			}
			if script.Source == "" {
				continue
			}
			for _, callfile := range script.Caller {
				pr, ok := t.Base.PathRegistry.DB[callfile]
				if !ok {
					continue
				}
				lfr, ok := pr.Mapping[script.Source]
				if !ok {
					continue
				}
				data, err := ioutil.ReadFile(lfr.AbsPath)
				if err != nil {
					return err
				}
				_, err = t.Library.AddBook(fmt.Sprintf("script_%s", sid), data)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

// Render implements the Builder interface
func (t *TerraformIBMBuilder) Render() error {
	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	for i := 0; i < t.Base.Environment.TeamCount; i++ {
		wg.Add(1)
		go func(teamid int) {
			defer wg.Done()
			t.Lock()
			team, ok := t.Base.Build.Teams[teamid]
			t.Unlock()
			if !ok {
				errChan <- fmt.Errorf("team number %d not found in team index", teamid)
				return
			}
			teamDir := team.RelBuildPath
			team.ID = team.Name()
			user := t.Base.User
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
			cfgData, err := t.Library.ExecuteGroup(primaryTemplate, templatesToLoad, ctx)
			if err != nil {
				errChan <- buildutil.Throw(err, "template failed", &buildutil.V{
					"team": teamid,
					"dir":  teamDir,
				})
				return
			}
			hclPretty, err := printer.Format(cfgData)
			cfgFile := filepath.Join(teamDir, "infra.tf")
			if err != nil {
				ioutil.WriteFile(cfgFile, cfgData, 0644)
				errChan <- err
				return
			}
			err = ioutil.WriteFile(cfgFile, hclPretty, 0644)
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
