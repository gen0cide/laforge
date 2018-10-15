package core

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gen0cide/laforge/tf"
	"github.com/pkg/errors"
)

var (
	whitespaceRegexp = regexp.MustCompile(`^[[:space:]]*$`)
	warningRegexp    = regexp.MustCompile(`Warning`)
	instanceRegexp   = regexp.MustCompile(`google_compute_instance`)
)

// Team represents a team specific object existing within an environment
type Team struct {
	ID            string                      `hcl:"id,label" json:"id,omitempty"`
	TeamNumber    int                         `hcl:"team_number,attr" json:"team_number,omitempty"`
	BuildID       string                      `hcl:"build_id,attr" json:"build_id,omitempty"`
	EnvironmentID string                      `hcl:"environment_id,attr" json:"environment_id,omitempty"`
	CompetitionID string                      `hcl:"competition_id,attr" json:"competition_id,omitempty"`
	Config        map[string]string           `hcl:"config,attr" json:"config,omitempty"`
	Tags          map[string]string           `hcl:"tags,attr" json:"tags,omitempty"`
	OnConflict    *OnConflict                 `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Revision      int64                       `hcl:"revision,attr" json:"revision,omitempty"`
	Maintainer    *User                       `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Hosts         map[string]*ProvisionedHost `json:"-"`
	Build         *Build                      `json:"-"`
	Environment   *Environment                `json:"-"`
	Competition   *Competition                `json:"-"`
	RelBuildPath  string                      `json:"-"`
	TeamRoot      string                      `json:"-"`
	Caller        Caller                      `json:"-"`
	Runner        *tf.Runner                  `json:"-"`
}

// GetCaller implements the Mergeable interface
func (t *Team) GetCaller() Caller {
	return t.Caller
}

// GetID implements the Mergeable interface
func (t *Team) GetID() string {
	return filepath.Join(t.CompetitionID, t.EnvironmentID, t.BuildID, t.ID)
}

// GetParentID returns the Team's parent build ID
func (t *Team) GetParentID() string {
	return filepath.Join(t.CompetitionID, t.EnvironmentID, t.BuildID)
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
	t.Revision++
	if t.BuildID == "" && t.Build != nil {
		t.BuildID = t.Build.ID
	}
	if t.EnvironmentID == "" && t.Environment != nil {
		t.EnvironmentID = t.Environment.ID
	}
	if t.CompetitionID == "" && t.Competition != nil {
		t.CompetitionID = t.Competition.ID
	}
	t.ID = fmt.Sprintf("%d", t.TeamNumber)
	return t.ID
}

// CreateRunner creates a new local command runner for the team, and returns it
func (t *Team) CreateRunner() *tf.Runner {
	runner := tf.NewRunner(t.GetID(), t.GetCaller().Current().CallerDir, Logger)
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
		Logger.Errorf("failed %s: no terraform binary located in path", t.GetID())
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
			Logger.Debugf("%s: %s", t.GetID(), i)
			continue
		case e := <-runner.Errors:
			Logger.Errorf("%s: %v", t.GetID(), e)
			execerr = e
			continue
		default:
		}
		select {
		case i := <-runner.Output:
			Logger.Debugf("%s: %s", t.GetID(), i)
			continue
		case e := <-runner.Errors:
			Logger.Errorf("%s: %v", t.GetID(), e)
			execerr = e
			continue
		case <-runner.FinChan:
			Logger.Warnf("%s command returned.", t.GetID())
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
		Logger.Errorf("%s - TF Init Error: %v", t.GetID(), err)
		return
	}

	time.Sleep(3 * time.Second)

	runner := t.CreateRunner()
	go runner.ExecuteCommand(command, args...)

	for {
		select {
		case i := <-runner.Output:
			Logger.Debugf("%s: %s", t.GetID(), i)
			continue
		case e := <-runner.Errors:
			Logger.Errorf("%s: %v", t.GetID(), e)
			continue
		default:
		}
		select {
		case i := <-runner.Output:
			Logger.Debugf("%s: %s", t.GetID(), i)
			continue
		case e := <-runner.Errors:
			Logger.Errorf("%s: %v", t.GetID(), e)
			continue
		case <-runner.FinChan:
			Logger.Warnf("%s command returned.", t.GetID())
		}
		break
	}
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
