// Package tfaws implements a Laforge Builder module for generating terraform configurations that target Google Compute Platform.
//go:generate fileb0x assets.toml
package tfaws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gen0cide/laforge/agent"
	"github.com/gen0cide/laforge/builder/tfaws/static"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/hashicorp/hcl/hcl/printer"

	"github.com/gen0cide/laforge/builder/buildutil/templates"
	validations "github.com/gen0cide/laforge/builder/buildutil/valdations"

	"github.com/pkg/errors"

	"github.com/gen0cide/laforge/builder/buildutil"
	"github.com/gen0cide/laforge/core"
)

// Definition of builder meta-data.
const (
	ID          = `tfaws`
	Name        = `Terraform AWS Builder`
	Description = `generates terraform configurations that isolate teams into VPCs on Google Compute Platform`
	Author      = `Alex Levinson <github.com/gen0cide>,Fred Rybin <github.com/frybin>`
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
			Name:       "DNS type not listed as bind",
			Resolution: "Make sure your dns block declaration has bind as it's type.",
			Check:      validations.FieldEquals(core.DNS{}, "Type", "bind"),
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
			Name:       "vpc CIDR not defined",
			Resolution: "define a vpc_cidr value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "vpc_cidr"),
		},
		validations.Requirement{
			Name:       "AWS Creds JSON File (aws_cred_file) not defined",
			Resolution: "define a aws_cred_file value inside your environment config = { ... } block.",
			Check:      validations.HasConfigKey(core.Environment{}, "aws_cred_file"),
		},
		validations.Requirement{
			Name:       "Root DNS Server not defined",
			Resolution: "define root dns_servers[] in the dns { ... } block within the competition configuration.",
			Check:      validations.FieldNotEmpty(core.DNS{}, "DNSServers"),
		},
		validations.Requirement{
			Name:       "no teams specified",
			Resolution: "make sure to set your team_count inside your environment config block to at least 1.",
			Check:      validations.FieldNotEmpty(core.Environment{}, "team_count"),
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
	}

	templatesToLoad = []string{
		"infra.tf.tmpl",
		"command.tf.tmpl",
		"script.tf.tmpl",
		"remote_file.tf.tmpl",
		"dns_record.tf.tmpl",
	}

	additionalTemplates = []string{
		"provisioned_host.tf.tmpl",
		"root_module.tf.tmpl",
	}

	primaryTemplate = "infra.tf.tmpl"
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
	res, ok := t.Base.CurrentBuild.Config[key]
	if ok {
		return res
	}
	r0, e0 := t.Base.CurrentEnv.Config[key]
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
	t.Base.CurrentBuild.Config[key] = fmt.Sprintf("%v", val)
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
	for _, x := range additionalTemplates {
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
	var privkey, pubkey string
	pathToPubkey := filepath.Join(t.Base.CurrentBuild.Dir, "data", "ssh.pem.pub")
	pathToPrivkey := filepath.Join(t.Base.CurrentBuild.Dir, "data", "ssh.pem")

	if _, err := os.Stat(pathToPubkey); os.IsNotExist(err) {
		privkey, pubkey, err := buildutil.GenerateSSHKeyPair(2048)
		if err != nil {
			return buildutil.Throw(err, "Could not generate a 2048-bit RSA SSH key.", nil)
		}
		err = buildutil.WriteKeyfile([]byte(privkey), pathToPrivkey)
		if err != nil {
			return buildutil.Throw(err, "Could not write the the SSH private key to the build directory", &buildutil.V{"path": pathToPrivkey})
		}
		err = buildutil.WriteKeyfile([]byte(pubkey), pathToPubkey)
		if err != nil {
			return buildutil.Throw(err, "Could not write the the SSH public key to the build directory", &buildutil.V{"path": pathToPubkey})
		}
	} else {
		pubkeyData, pubkeyErr := ioutil.ReadFile(pathToPubkey)
		if pubkeyErr != nil {
			return buildutil.Throw(pubkeyErr, "could not read already established public key", nil)
		}
		privkeyData, privkeyErr := ioutil.ReadFile(pathToPrivkey)
		if privkeyErr != nil {
			return buildutil.Throw(privkeyErr, "could not read already established private key", nil)
		}
		privkey = string(privkeyData)
		pubkey = string(pubkeyData)
	}

	t.Set("ssh_public_key_file", pathToPubkey)
	t.Set("ssh_private_key_file", pathToPrivkey)
	t.Set("rel_ssh_public_key_file", "../../data/ssh.pem.pub")
	t.Set("rel_ssh_private_key_file", "../../data/ssh.pem")
	t.Set("ssh_public_key", pubkey)
	t.Set("ssh_private_key", privkey)

	for hostid, host := range t.Base.CurrentEnv.IncludedHosts {
		uds, found := host.Vars["user_data_script_id"]
		if !found {
			return buildutil.Throw(errors.New("user_data_script_id no longer exists"), "Validation for this passed, but here we are. Likely a bug. Please report.", &buildutil.V{"host_id": hostid})
		}
		udsObj, found := t.Base.Scripts[uds]
		if !found {
			return buildutil.Throw(errors.Errorf("user_data_script_id %s not found", uds), "Host declares a user_data_script_id which was not found in the script map. Is this declared somewhere?", &buildutil.V{"host": hostid})
		}
		if _, ok := host.Scripts[uds]; ok {
			cli.Logger.Infof("UDS %s is already defined for host %s (strange?)", uds, hostid)
			continue
		}
		cli.Logger.Debugf("Adding user_data_script %s to host %s script pool", uds, hostid)
		host.Scripts[uds] = udsObj
		if _, ok := t.Library.Books[uds]; !ok {
			for _, callfile := range udsObj.Caller {
				pr, ok := t.Base.PathRegistry.DB[callfile]
				if !ok {
					continue
				}
				lfr, ok := pr.Mapping[udsObj.Source]
				if !ok {
					continue
				}
				data, err := ioutil.ReadFile(lfr.AbsPath)
				if err != nil {
					return err
				}
				_, err = t.Library.AddBook(udsObj.Path(), data)
				if err != nil {
					return err
				}
				break
			}
		}

		for _, dep := range host.Dependencies {
			depHost, ok := t.Base.CurrentEnv.IncludedHosts[dep.HostID]
			if !ok {
				return buildutil.Throw(errors.Errorf("host %s depends on host %s, which is not found in environment", host.ID, dep.HostID), "The host listed a dependency to another host which is not included in any network within the current environment.", &buildutil.V{"source_host": hostid, "depends_on_host": dep.HostID})
			}
			dep.Host = depHost

			depNet, ok := t.Base.CurrentEnv.IncludedNetworks[dep.NetworkID]
			if !ok {
				return buildutil.Throw(errors.Errorf("host %s depends on network %s, which is not found in environment", host.ID, dep.NetworkID), "The host listed a dependency to another network which is not included within the current environment.", &buildutil.V{"source_host": hostid, "depends_on_host": dep.HostID, "depends_on_network": dep.NetworkID})
			}
			dep.Network = depNet

			hostInNetwork := false
			for _, x := range t.Base.CurrentEnv.HostByNetwork[dep.NetworkID] {
				if x.ID == dep.Host.ID {
					hostInNetwork = true
					break
				}
			}
			if !hostInNetwork {
				return buildutil.Throw(errors.Errorf("host %s depends on host %s, which is not included in network %s", host.ID, dep.HostID, dep.NetworkID), "The host listed a dependency to another host, and while the network exists and is included, this host is not present within this network assignment.", &buildutil.V{"source_host": hostid, "depends_on_host": dep.HostID, "depends_on_network": dep.NetworkID})
			}

			if dep.Step != "" {
				// the Host index function within core.Host has already guarenteed that a provisioning step exists
				// don't need to check on that :)
				located := false
				for stepidx, x := range dep.Host.ProvisionSteps {
					if dep.Step == x {
						located = true
						dep.StepID = stepidx
						break
					}
				}
				if !located {
					return buildutil.Throw(errors.Errorf("host %s depends on provisioning step %s, which is not found in host %s", host.ID, dep.Step, dep.Host.ID), "The host listed a dependency to a provisioning step that is not included within the supplied host's provisioning steps.", &buildutil.V{"source_host": hostid, "depends_on_host": dep.HostID, "depends_on_network": dep.NetworkID, "depends_on_step": dep.Step})
				}
			} else {
				dep.StepID = dep.Host.FinalStepID()
			}
		}
	}

	return nil
}

// GenerateScripts implements the Builder interface
func (t *TerraformAWSBuilder) GenerateScripts() error {
	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	user := t.Base.User
	ctx, err := templates.NewContext(
		t.Base,
		t.Base.CurrentBuild,
		t.Base.CurrentCompetition,
		t.Base.CurrentCompetition.DNS,
		t.Base.CurrentEnv,
		user,
	)
	if err != nil {
		return err
	}
	for _, teamObj := range t.Base.CurrentBuild.Teams {
		wg.Add(1)
		go func(team *core.Team) {
			defer wg.Done()
			for netName, hosts := range t.Base.CurrentEnv.HostByNetwork {
				network := t.Base.CurrentEnv.IncludedNetworks[netName]
				for _, host := range hosts {
					for sid, script := range host.Scripts {
						wg.Add(1)
						go func(scriptID string, scriptObj *core.Script, hostObj *core.Host) {
							defer wg.Done()
							scriptCtx := ctx.Clone()
							tid := team.ID
							pnid := filepath.Join(tid, "networks", network.Base())
							pnobj := t.Base.StateManager.Current.Metastore[pnid].Dependency.(*core.ProvisionedNetwork)
							phid := filepath.Join(pnid, "hosts", hostObj.Base())
							phobj := t.Base.StateManager.Current.Metastore[phid].Dependency.(*core.ProvisionedHost)
							var pstep *core.ProvisioningStep
							for _, x := range phobj.ProvisioningSteps {
								if x.ProvisionerID == scriptID {
									pstep = x
									break
								}
							}
							conn := phobj.Conn
							err := scriptCtx.Attach(team, network, hostObj, scriptObj, pnobj, phobj, pstep, conn)
							if err != nil {
								errChan <- err
								return
							}
							filename := filepath.Base(scriptObj.Source)
							assetDir := filepath.Join(team.RelBuildPath, "networks", network.Base(), "hosts", hostObj.Base(), "assets")
							assetPath := filepath.Join(assetDir, filename)
							fileData, err := t.Library.Execute(scriptID, scriptCtx)
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
					wg.Add(1)
					go func(h *core.Host) {
						defer wg.Done()
						tid := team.ID
						pnid := filepath.Join(tid, "networks", network.Base())
						pnobj := t.Base.StateManager.Current.Metastore[pnid].Dependency.(*core.ProvisionedNetwork)
						phid := filepath.Join(pnid, "hosts", h.Base())
						phobj := t.Base.StateManager.Current.Metastore[phid].Dependency.(*core.ProvisionedHost)
						scriptCtx := ctx.Clone()
						err := scriptCtx.Attach(team, network, h, pnobj, phobj)
						if err != nil {
							errChan <- err
							return
						}
						filename := "provisioned_host.tpl"
						assetDir := filepath.Join(team.RelBuildPath, "networks", network.Base(), "hosts", h.Base(), "assets")
						assetPath := filepath.Join(assetDir, filename)
						fileData, err := t.Library.Execute("provisioned_host.tf.tmpl", scriptCtx)
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
					}(host)
				}
			}

		}(teamObj)
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
func (t *TerraformAWSBuilder) StageDependencies() error {
	if t.Base.StateManager == nil {
		return errors.New("builder cannot stage dependencies with nil state manager")
	}
	snap := t.Base.StateManager.Current
	_ = snap
	build := t.Base.CurrentBuild
	// snap, err := core.NewSnapshotFromEnv(t.Base.CurrentEnv)
	// if err != nil {
	// 	return err
	// }

	// buildnode, ok := snap.Metastore[path.Join(t.Base.CurrentEnv.Path(), t.Base.CurrentEnv.Builder)]
	// if !ok {
	// 	return errors.New("builder was not able to be resolved on the graph")
	// }
	// build, ok := buildnode.Dependency.(*core.Build)
	// if !ok {
	// 	return errors.New("build object was not of type *core.Build")
	// }

	// t.Base.CurrentBuild = build
	// // for _, x := range build.Teams {

	// }
	for _, team := range build.Teams {
		// TODO: Make team directory creation part of core
		teamDir := filepath.Join(build.Dir, "teams", fmt.Sprintf("%v", team.TeamNumber))
		team.RelBuildPath = teamDir
		os.MkdirAll(teamDir, 0755)
		core.TouchGitKeep(teamDir)
		teamCfg, err := core.RenderHCLv2Object(team)
		if err != nil {
			return err
		}
		teamCfgFile := filepath.Join(teamDir, "team.laforge")
		err = ioutil.WriteFile(teamCfgFile, teamCfg, 0644)
		if err != nil {
			return err
		}
		for _, pn := range team.ProvisionedNetworks {
			netdir := filepath.Join(teamDir, "networks", pn.Base())
			os.MkdirAll(netdir, 0755)
			core.TouchGitKeep(netdir)
			data, err := core.RenderHCLv2Object(pn)
			if err != nil {
				return err
			}
			netfile := filepath.Join(netdir, "provisioned_network.laforge")
			err = ioutil.WriteFile(netfile, data, 0644)
			if err != nil {
				return err
			}
			for _, ph := range pn.ProvisionedHosts {
				hostdir := filepath.Join(netdir, "hosts", ph.Base())
				os.Mkdir(hostdir, 0755)
				core.TouchGitKeep(hostdir)
				agentdir := filepath.Join(hostdir, "agent")
				assetdir := filepath.Join(hostdir, "assets")
				stepdir := filepath.Join(hostdir, "steps")
				os.MkdirAll(agentdir, 0755)
				os.MkdirAll(assetdir, 0755)
				os.MkdirAll(stepdir, 0755)
				core.TouchGitKeep(agentdir)
				core.TouchGitKeep(assetdir)
				core.TouchGitKeep(stepdir)
				data, err = core.RenderHCLv2Object(ph)
				if err != nil {
					return err
				}
				hostfile := filepath.Join(hostdir, "provisioned_host.laforge")
				err = ioutil.WriteFile(hostfile, data, 0644)
				if err != nil {
					return err
				}
				for _, ps := range ph.ProvisioningSteps {
					stepfile := filepath.Join(stepdir, fmt.Sprintf("%s.laforge", ps.Base()))
					data, err = core.RenderHCLv2Object(ps)
					if err != nil {
						return err
					}
					err = ioutil.WriteFile(stepfile, data, 0644)
					if err != nil {
						return err
					}
					if rfile, ok := ps.Provisioner.(*core.RemoteFile); ok {
						rfileName, err := rfile.AssetName()
						if err != nil {
							return err
						}

						dstPath := filepath.Join(t.Base.CurrentBuild.Dir, "data", rfileName)
						if _, err := os.Stat(dstPath); os.IsNotExist(err) {
							copyErr := rfile.CopyTo(dstPath)
							if copyErr != nil {
								return copyErr
							}
						}
					}
					if script, ok := ps.Provisioner.(*core.Script); ok {
						if _, ok := t.Library.Books[script.Path()]; ok {
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
							_, err = t.Library.AddBook(script.Path(), data)
							if err != nil {
								return err
							}
							break
						}
					}
				}
			}
		}

	}
	return nil
}

// Render implements the Builder interface
func (t *TerraformAWSBuilder) Render() error {
	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	for _, team := range t.Base.CurrentBuild.Teams {
		wg.Add(1)
		go func(team *core.Team) {
			defer wg.Done()

			teamDir := team.RelBuildPath
			user := t.Base.User
			ctx, err := templates.NewContext(
				t.Base,
				t.Base.CurrentBuild,
				t.Base.CurrentCompetition,
				t.Base.CurrentCompetition.DNS,
				t.Base.CurrentEnv,
				user,
				team,
			)
			if err != nil {
				errChan <- err
				return
			}
			cfgData, err := t.Library.ExecuteGroup(primaryTemplate, templatesToLoad, ctx)
			if err != nil {
				errChan <- buildutil.Throw(err, "template failed", &buildutil.V{
					"team": team.Path(),
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
			for netname, net := range t.Base.CurrentEnv.IncludedNetworks {
				for _, host := range t.Base.CurrentEnv.HostByNetwork[netname] {
					ts := time.Now()
					state := &agent.State{
						Team:         team,
						Network:      net,
						Steps:        []*agent.Step{},
						RenderedAt:   ts,
						Revision:     ts.UTC().Unix(),
						CurrentState: "pending",
					}
					for pid, prov := range host.Provisioners {
						step := &agent.Step{
							ID:       pid,
							StepType: prov.Kind(),
							Metadata: map[string]interface{}{},
						}
						state.Steps = append(state.Steps, step)
					}
					jsonData, err := json.MarshalIndent(state, "", "  ")
					if err != nil {
						errChan <- err
						return
					}
					stateFilePath := filepath.Join(teamDir, "networks", net.Base(), "hosts", host.Base(), "agent", "config.json")
					err = ioutil.WriteFile(stateFilePath, jsonData, 0644)
					if err != nil {
						errChan <- err
						return
					}
				}
			}
		}(team)
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
