package core

import (
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/runner"
	"github.com/pkg/errors"
)

var (
	whitespaceRegexp = regexp.MustCompile(`^[[:space:]]*$`)
	warningRegexp    = regexp.MustCompile(`Warning`)
	instanceRegexp   = regexp.MustCompile(`google_compute_instance`)
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
			"tn=%v bid=%v config=%v",
			t.Path(),
			t.Build.Hash(),
			t.Config,
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
	runner := runner.NewRunner(t.LaforgeID(), t.GetCaller().Current().CallerDir, Logger)
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
func (t *Team) RunTerraformCommand(args []string, wg *sync.WaitGroup) {
	defer wg.Done()

	tfexe, err := FindTerraformExecutable()
	if err != nil {
		Logger.Errorf("failed %s: no terraform binary located in path", t.LaforgeID())
		return
	}

	wg.Add(1)
	t.RunLocalCommand(tfexe, args, wg)
	return
}

// TerraformInit runs terraform init for a team
func (t *Team) TerraformInit() error {
	tfexe, err := FindTerraformExecutable()
	if err != nil {
		return err
	}

	runner := t.CreateRunner()
	go runner.ExecuteCommand(tfexe, []string{"init"}...)

	var execerr error
	for {
		select {
		case i := <-runner.Output:
			Logger.Debugf("%s: %s", t.LaforgeID(), i)
			continue
		case e := <-runner.Errors:
			Logger.Errorf("%s: %v", t.LaforgeID(), e)
			execerr = e
			continue
		default:
		}
		select {
		case i := <-runner.Output:
			Logger.Debugf("%s: %s", t.LaforgeID(), i)
			continue
		case e := <-runner.Errors:
			Logger.Errorf("%s: %v", t.LaforgeID(), e)
			execerr = e
			continue
		case <-runner.FinChan:
			Logger.Warnf("%s command returned.", t.LaforgeID())
		}
		break
	}

	return execerr
}

// RunLocalCommand runs a local command in the team's local directory
func (t *Team) RunLocalCommand(command string, args []string, wg *sync.WaitGroup) {
	defer wg.Done()

	err := t.TerraformInit()
	if err != nil {
		Logger.Errorf("%s - TF Init Error: %v", t.LaforgeID(), err)
		return
	}

	time.Sleep(3 * time.Second)

	runner := t.CreateRunner()
	go runner.ExecuteCommand(command, args...)

	for {
		select {
		case i := <-runner.Output:
			Logger.Debugf("%s: %s", t.LaforgeID(), i)
			continue
		case e := <-runner.Errors:
			Logger.Errorf("%s: %v", t.LaforgeID(), e)
			continue
		default:
		}
		select {
		case i := <-runner.Output:
			Logger.Debugf("%s: %s", t.LaforgeID(), i)
			continue
		case e := <-runner.Errors:
			Logger.Errorf("%s: %v", t.LaforgeID(), e)
			continue
		case <-runner.FinChan:
			Logger.Warnf("%s command returned.", t.LaforgeID())
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
	for _, net := range t.ProvisionedNetworks {
		err := g.Relate(t, net)
		if err != nil {
			return err
		}
		err = net.Gather(g)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeEmptyLines(s string) string {
	lines := strings.Split(s, "\n")
	newLines := []string{}
	for _, x := range lines {
		newX := strings.TrimSpace(x)
		if len(newX) == 0 || whitespaceRegexp.MatchString(newX) {
			continue
		}
		newLines = append(newLines, newX)
	}
	return strings.Join(newLines, "\n")
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
