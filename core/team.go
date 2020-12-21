package core

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/runner"
	"github.com/pkg/errors"
)

// Team represents a team specific object existing within an environment
//easyjson:json
type Team struct {
	ID                  string                         `hcl:"id,label" json:"id,omitempty"`
	TeamNumber          int                            `hcl:"team_number,attr" json:"team_number,omitempty"`
	Config              map[string]string              `hcl:"config,attr" json:"config,omitempty"`
	Tags                map[string]string              `hcl:"tags,attr" json:"tags,omitempty"`
	OnConflict          *OnConflict                    `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Revision            int64                          `hcl:"revision,attr" json:"revision,omitempty"`
	Maintainer          *User                          `hcl:"maintainer,block" json:"maintainer,omitempty"`
	ProvisionedNetworks map[string]*ProvisionedNetwork `json:"provisioned_networks"`
	ProvisionedHosts    map[string]*ProvisionedHost    `json:"provisioned_hosts"`
	Build               *Build                         `json:"-"`
	Environment         *Environment                   `json:"-"`
	Competition         *Competition                   `json:"-"`
	RelBuildPath        string                         `json:"-"`
	TeamRoot            string                         `json:"-"`
	Dir                 string                         `json:"-"`
	Caller              Caller                         `json:"-"`
	Runner              *runner.Runner                 `json:"-"`
}

// Hash implements the Hasher interface
func (t *Team) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"tn=%v config=%v",
			t.Path(),
			HashConfigMap(t.Config),
		),
	)
}

// Path implements the Pather interface
func (t *Team) Path() string {
	return t.ID
}

// Base implements the Pather interface
func (t *Team) Base() string {
	return path.Base(t.ID)
}

// ValidatePath implements the Pather interface
func (t *Team) ValidatePath() error {
	if err := ValidateGenericPath(t.Path()); err != nil {
		return err
	}
	return nil
}

// GetCaller implements the Mergeable interface
func (t *Team) GetCaller() Caller {
	return t.Caller
}

// LaforgeID implements the Mergeable interface
func (t *Team) LaforgeID() string {
	return t.ID
}

// ParentLaforgeID returns the Team's parent build ID
func (t *Team) ParentLaforgeID() string {
	return path.Dir(path.Dir(t.LaforgeID()))
}

// GetOnConflict implements the Mergeable interface
func (t *Team) GetOnConflict() OnConflict {
	if t.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *t.OnConflict
}

// SetCaller implements the Mergeable interface
func (t *Team) SetCaller(ca Caller) {
	t.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (t *Team) SetOnConflict(o OnConflict) {
	t.OnConflict = &o
}

// Swap implements the Mergeable interface
func (t *Team) Swap(m Mergeable) error {
	rawVal, ok := m.(*Team)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", t, m)
	}
	*t = *rawVal
	return nil
}

// SetID increments the revision and sets the team ID if needed
func (t *Team) SetID() string {
	if t.ID == "" {
		t.ID = path.Join(t.Build.Path(), "teams", fmt.Sprintf("%d", t.TeamNumber))
	}
	if t.Environment == nil {
		t.Environment = t.Build.Environment
	}
	return t.ID
}

// CreateRunner creates a new local command runner for the team, and returns it
func (t *Team) CreateRunner() *runner.Runner {
	if t.Dir == "" {
		cfg, err := LocateBaseConfig()
		if err != nil {
			panic(err)
		}
		t.Dir = filepath.Join(filepath.Dir(cfg), t.ID)
	}
	runner := runner.NewRunner(t.LaforgeID(), t.Dir)
	return runner
}

// FindTerraformExecutable attempts to find the terraform executable in the given path
func FindTerraformExecutable() (string, error) {
	binname := "terraform"
	if runtime.GOOS == "windows" {
		binname = "terraform.exe"
	}
	tfexepath, err := exec.LookPath(binname)
	if err != nil {
		return "", err
	}
	return tfexepath, nil
}

// RunTerraformCommand runs terraform subcommands inside a team's local directory
func (t *Team) RunTerraformCommand(args []string, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	tfexe, err := FindTerraformExecutable()
	if err != nil {
		cli.Logger.Errorf("failed %s: no terraform binary located in path", t.LaforgeID())
		errChan <- errors.New("terraform binary was not located in the path")
		return
	}

	wg.Add(1)
	t.RunLocalCommand(tfexe, args, wg, errChan)
}

// RunTerraformSequence attempts to run a series of commands on a team
func (t *Team) RunTerraformSequence(cmds []string, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	tfexe, err := FindTerraformExecutable()
	if err != nil {
		errChan <- err
		return
	}

	err = t.TerraformInit()
	if err != nil {
		cli.Logger.Errorf("%s - TF Init Error: %v", t.LaforgeID(), err)
		errChan <- err
		return
	}

	for _, tfcmd := range cmds {
		time.Sleep(3 * time.Second)

		errors := []error{}
		runner := t.CreateRunner()

		go runner.ExecuteCommand(tfexe, strings.Split(tfcmd, " ")...)

		for {
			select {
			case i := <-runner.Output:
				cli.Logger.Debugf("%s: %s", t.LaforgeID(), i)
				continue
			case e := <-runner.Errors:
				cli.Logger.Errorf("%s: %v", t.LaforgeID(), e)
				errChan <- e
				errors = append(errors, e)
			default:
			}

			select {
			case i := <-runner.Output:
				cli.Logger.Debugf("%s: %s", t.LaforgeID(), i)
				continue
			case e := <-runner.Errors:
				cli.Logger.Errorf("%s: %v", t.LaforgeID(), e)
				errChan <- e
				//nolint:ineffassign
				errors = append(errors, e)
			case <-runner.FinChan:
				if len(errors) > 0 {
					cli.Logger.Errorf("Runner for %s has failed performing %s %s:", t.LaforgeID(), tfexe, tfcmd)
					for idx, x := range errors {
						cli.Logger.Errorf("Error #%d: %v", idx, x)
					}
					return
				}
				cli.Logger.Warnf("Team %s executed \"terraform %s\" successfully. Logs located at %s", t.LaforgeID(), tfcmd, runner.StdoutFile)
			}
			break
		}
	}
}

// TerraformInit runs terraform init for a team
func (t *Team) TerraformInit() error {
	tfexe, err := FindTerraformExecutable()
	if err != nil {
		return err
	}

	runner := t.CreateRunner()
	go runner.ExecuteCommand(tfexe, []string{"init", "-force-copy", "-no-color"}...)

	var execerr error
	for {
		select {
		case i := <-runner.Output:
			cli.Logger.Debugf("%s: %s", t.LaforgeID(), i)
			continue
		case e := <-runner.Errors:
			cli.Logger.Errorf("%s: %v", t.LaforgeID(), e)
			execerr = e
			continue
		default:
		}
		select {
		case i := <-runner.Output:
			cli.Logger.Debugf("%s: %s", t.LaforgeID(), i)
			continue
		case e := <-runner.Errors:
			cli.Logger.Errorf("%s: %v", t.LaforgeID(), e)
			execerr = e
			continue
		case <-runner.FinChan:
			cli.Logger.Warnf("%s command returned. Logs located at %s", t.LaforgeID(), runner.StdoutFile)
		}
		break
	}

	return execerr
}

// RunLocalCommand runs a local command in the team's local directory
func (t *Team) RunLocalCommand(command string, args []string, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	err := t.TerraformInit()
	if err != nil {
		cli.Logger.Errorf("%s - TF Init Error: %v", t.LaforgeID(), err)
		return
	}

	time.Sleep(3 * time.Second)

	var runner *runner.Runner
	if t.Runner == nil {
		runner = t.CreateRunner()
	} else {
		runner = t.Runner
	}

	go runner.ExecuteCommand(command, args...)

	for {
		select {
		case i := <-runner.Output:
			cli.Logger.Debugf("%s: %s", t.LaforgeID(), i)
			continue
		case e := <-runner.Errors:
			cli.Logger.Errorf("%s: %v", t.LaforgeID(), e)
			errChan <- e
			continue
		default:
		}
		select {
		case i := <-runner.Output:
			cli.Logger.Debugf("%s: %s", t.LaforgeID(), i)
			continue
		case e := <-runner.Errors:
			cli.Logger.Errorf("%s: %v", t.LaforgeID(), e)
			errChan <- e
			continue
		case <-runner.FinChan:
			cli.Logger.Warnf("%s command returned.", t.LaforgeID())
		}
		break
	}
}

// CreateProvisionedNetwork actually creates the provisioned network object and assigns parent pointers accordingly.
func (t *Team) CreateProvisionedNetwork(net *Network) *ProvisionedNetwork {
	p := &ProvisionedNetwork{
		Name:             net.Name,
		CIDR:             net.CIDR,
		Network:          net,
		Team:             t,
		Build:            t.Build,
		Environment:      t.Environment,
		Competition:      t.Competition,
		ProvisionedHosts: map[string]*ProvisionedHost{},
	}

	t.ProvisionedNetworks[p.SetID()] = p
	return p
}

// CreateProvisionResources enumerates the environment's included networks and creates provisioned network objects
func (t *Team) CreateProvisionResources() error {
	for _, n := range t.Environment.IncludedNetworks {
		pn := t.CreateProvisionedNetwork(n)
		err := pn.CreateProvisionedHosts()
		if err != nil {
			return err
		}
	}
	return nil
}

// Gather implements the Dependency interface
func (t *Team) Gather(g *Snapshot) error {
	// for _, net := range t.ProvisionedNetworks {
	// 	g.AddNode(net)
	// 	err := net.Gather(g)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// Associate attempts to actually draw the relationships on the graph between dependencies
func (t *Team) Associate(g *Snapshot) error {
	// depwalker := map[string]*ProvisionedHost{}
	// for nid, net := range t.ProvisionedNetworks {
	// 	nmeta := g.Metastore[nid]
	// 	for hid, host := range net.ProvisionedHosts {
	// 		hmeta := g.Metastore[hid]
	// 		connmeta := g.Metastore[host.Conn.Path()]
	// 		g.Connect(hmeta, connmeta)
	// 		if len(host.Host.Dependencies) > 0 {
	// 			depwalker[hid] = host
	// 		} else {
	// 			for idx, s := range host.StepsByOffset {
	// 				smeta := g.Metastore[s.Path()]
	// 				if s.StepNumber != 0 {
	// 					prevIdx := idx - 1
	// 					prevStep := host.StepsByOffset[prevIdx]
	// 					prevmeta := g.Metastore[prevStep.Path()]
	// 					g.Connect(prevmeta, smeta)
	// 					err := g.Graph.Validate()
	// 					if err != nil {
	// 						cli.Logger.Errorf("Graph trust has failed connecting %s to %s", prevmeta.ID, smeta.ID)
	// 						fmt.Printf("%s\n", g.Graph.StringWithNodeTypes())
	// 						panic(err)
	// 					}
	// 				} else {
	// 					g.Connect(connmeta, smeta)
	// 					err := g.Graph.Validate()
	// 					if err != nil {
	// 						cli.Logger.Errorf("Graph trust has failed connecting %s to %s", hmeta.ID, smeta.ID)
	// 						fmt.Printf("%s\n", g.Graph.StringWithNodeTypes())
	// 						panic(err)
	// 					}
	// 				}
	// 			}
	// 			g.Connect(nmeta, hmeta)
	// 			err := g.Graph.Validate()
	// 			if err != nil {
	// 				cli.Logger.Errorf("Graph trust has failed connecting %s to %s", nmeta.ID, hmeta.ID)
	// 				fmt.Printf("%s\n", g.Graph.StringWithNodeTypes())
	// 				panic(err)
	// 			}
	// 		}
	// 	}
	// }

	// for {
	// 	if len(depwalker) == 0 {
	// 		break
	// 	}
	// 	maxVal := 0
	// 	maxHost := ""
	// 	for hid, host := range depwalker {
	// 		if host.Host.DependencyCount(t.Environment) >= maxVal {
	// 			maxHost = hid
	// 			maxVal = host.Host.DependencyCount(t.Environment)
	// 		}
	// 	}
	// 	hmeta := g.Metastore[maxHost]

	// 	if maxVal == 1 {
	// 		dep := depwalker[maxHost].Host.Dependencies[0]
	// 		dh, err := t.LocateProvisionedHost(dep.NetworkID, dep.HostID)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fsid := dh.Host.FinalStepID()
	// 		if fsid != -1 {
	// 			smeta := g.Metastore[dh.StepsByOffset[fsid].Path()]
	// 			g.Connect(smeta, hmeta)
	// 			err := g.Graph.Validate()
	// 			if err != nil {
	// 				cli.Logger.Errorf("Graph trust has failed connecting %s to %s", smeta.ID, hmeta.ID)
	// 				fmt.Printf("%s\n", g.Graph.StringWithNodeTypes())
	// 				panic(err)
	// 			}
	// 			delete(depwalker, maxHost)
	// 			continue
	// 		} else {
	// 			// dhmeta := g.Metastore[dh.Path()]
	// 			dhconnmeta := g.Metastore[dh.Conn.Path()]
	// 			g.Connect(dhconnmeta, hmeta)
	// 			err := g.Graph.Validate()
	// 			if err != nil {
	// 				cli.Logger.Errorf("Graph trust has failed connecting %s to %s", dhconnmeta.ID, hmeta.ID)
	// 				fmt.Printf("%s\n", g.Graph.StringWithNodeTypes())
	// 				panic(err)
	// 			}
	// 			delete(depwalker, maxHost)
	// 			continue
	// 		}
	// 	} else {
	// 		maxDepVal := 0
	// 		maxDepOffset := 0
	// 		for do, dep := range depwalker[maxHost].Host.Dependencies {
	// 			depHost, ok := t.Environment.IncludedHosts[dep.HostID]
	// 			if !ok {
	// 				continue
	// 			}
	// 			if depHost.DependencyCount(t.Environment) >= maxDepVal {
	// 				maxDepOffset = do
	// 				maxDepVal = depHost.DependencyCount(t.Environment)
	// 			}
	// 		}
	// 		dep := depwalker[maxHost].Host.Dependencies[maxDepOffset]
	// 		dh, err := t.LocateProvisionedHost(dep.NetworkID, dep.HostID)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fsid := dh.Host.FinalStepID()
	// 		if fsid != -1 {
	// 			smeta := g.Metastore[dh.StepsByOffset[fsid].Path()]
	// 			g.Connect(smeta, hmeta)
	// 			err := g.Graph.Validate()
	// 			if err != nil {
	// 				cli.Logger.Errorf("Graph trust has failed connecting %s to %s", smeta.ID, hmeta.ID)
	// 				fmt.Printf("%s\n", g.Graph.StringWithNodeTypes())
	// 				panic(err)
	// 			}
	// 			delete(depwalker, maxHost)
	// 			continue
	// 		} else {
	// 			// dhmeta := g.Metastore[dh.Path()]
	// 			dhconnmeta := g.Metastore[dh.Conn.Path()]
	// 			g.Connect(dhconnmeta, hmeta)
	// 			err := g.Graph.Validate()
	// 			if err != nil {
	// 				cli.Logger.Errorf("Graph trust has failed connecting %s to %s", dhconnmeta.ID, hmeta.ID)
	// 				fmt.Printf("%s\n", g.Graph.StringWithNodeTypes())
	// 				panic(err)
	// 			}
	// 			delete(depwalker, maxHost)
	// 			continue
	// 		}
	// 	}
	// }
	// for _, net := range t.ProvisionedNetworks {
	// 	for _, host := range net.ProvisionedHosts {
	// 		for _, s := range host.StepsByOffset {
	// 			smeta := g.Metastore[s.Path()]
	// 			if s.StepNumber != 0 {
	// 				prevIdx := s.StepNumber - 1
	// 				prevStep := host.StepsByOffset[prevIdx]
	// 				prevmeta := g.Metastore[prevStep.Path()]
	// 				g.Connect(prevmeta, smeta)
	// 				err := g.Graph.Validate()
	// 				if err != nil {
	// 					cli.Logger.Errorf("Graph trust has failed connecting %s to %s", prevmeta.ID, smeta.ID)
	// 					fmt.Printf("%s\n", g.Graph.StringWithNodeTypes())
	// 					panic(err)
	// 				}
	// 			} else {
	// 				hconnmeta := g.Metastore[host.Conn.Path()]
	// 				g.Connect(hconnmeta, smeta)
	// 				err := g.Graph.Validate()
	// 				if err != nil {
	// 					cli.Logger.Errorf("Graph trust has failed connecting %s to %s", hconnmeta.ID, smeta.ID)
	// 					fmt.Printf("%s\n", g.Graph.StringWithNodeTypes())
	// 					panic(err)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// // fmt.Printf("%s\n\n\n", g.Graph.StringWithNodeTypes())
	return nil
}

// LocateProvisionedHost is used to locate the provisioned host object by specifying a global host and network ID. (useed in dependency traversal)
func (t *Team) LocateProvisionedHost(netid, hostid string) (*ProvisionedHost, error) {
	for _, x := range t.ProvisionedNetworks {
		if x.Network.Path() != netid {
			continue
		}
		for _, h := range x.ProvisionedHosts {
			if h.Host.Path() != hostid {
				continue
			}
			return h, nil
		}
	}
	return nil, fmt.Errorf("host %s was not located within network %s for the given build", hostid, netid)
}

// CreateTeamEntry ...
func (t *Team) CreateTeamEntry(env *ent.Environment, build *ent.Build, ctx context.Context, client *ent.Client) (*ent.Team, error) {
	// user, err := t.Maintainer.CreateUserEntry(ctx, client)

	// if err != nil {
	// 	cli.Logger.Debugf("failed creating team: %v", err)
	// 	return nil, err
	// }

	tag, err := CreateTagEntry(t.ID, t.Tags, ctx, client)

	if err != nil {
		cli.Logger.Debugf("failed creating team: %v", err)
		return nil, err
	}

	team, err := client.Team.
		Create().
		SetTeamNumber(t.TeamNumber).
		SetConfig(t.Config).
		SetRevision(t.Revision).
		// AddMaintainer(user).
		AddBuild(build).
		AddTeamToEnvironment(env).
		AddTag(tag).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating team: %v", err)
		return nil, err
	}

	for _, v := range t.ProvisionedNetworks {
		_, err := v.CreateProvisionedNetworkEntry(ctx,build,team,client)

		if err != nil {
			cli.Logger.Debugf("failed creating team: %v", err)
			return nil, err
		}
	}

	cli.Logger.Debugf("team was created: ", team)
	return team, nil
}
