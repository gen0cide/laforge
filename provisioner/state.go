package provisioner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/gen0cide/utils/uptime"

	"github.com/gen0cide/laforge/core"
)

var (
	// ErrAwaitingReboot is thrown when the state is awaiting a reboot
	ErrAwaitingReboot = errors.New("awaiting reboot for step")

	// ErrInProgress is triggered when a step is already in the process of executing
	ErrInProgress = errors.New("step currently in progress")

	// ErrStepFailure is thrown when a step has errored out
	ErrStepFailure = errors.New("current step has errored")
)

// State is the global state of this server and all required information
type State struct {
	sync.RWMutex
	Source          string                `json:"-"`
	Host            *core.Host            `json:"host,omitempty"`
	Network         *core.Network         `json:"network,omitempty"`
	Environment     *core.Environment     `json:"environment,omitempty"`
	Competition     *core.Competition     `json:"competition,omitempty"`
	Team            *core.Team            `json:"team,omitempty"`
	ProvisionedHost *core.ProvisionedHost `json:"provisioned_host,omitempty"`
	RenderedAt      time.Time             `json:"rendered_at,omitempty"`
	InitializedAt   time.Time             `json:"initialized_at,omitempty"`
	CompletedAt     time.Time             `json:"completed_at,omitempty"`
	CurrentState    string                `json:"current_state,omitempty"`
	Steps           []*Step               `json:"steps,omitempty"`
	Pending         map[int]*Step         `json:"pending_steps"`
	Completed       map[int]*Step         `json:"completed_steps"`
	CurrentStep     *Step                 `json:"current_step,omitempty"`
	Revision        int64                 `json:"revision,omitempty"`
	Errored         bool                  `json:"errored,omitempty"`
	ErrorMessage    string                `json:"error_message,omitempty"`
}

// LoadStateFile parses a JSON state file into a state object
func LoadStateFile(location string) (*State, error) {
	fdata, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	newState := &State{}
	err = json.Unmarshal(fdata, newState)
	if err != nil {
		return nil, err
	}
	return newState, nil
}

// WriteStateFile writes a JSON state file to disk
func (s *State) WriteStateFile(location string) error {
	s.Lock()
	defer s.Unlock()

	if _, err := os.Stat(location); err == nil {
		backupData, err := ioutil.ReadFile(location)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s.backup.%d", location, s.Revision), backupData, 0600)
		if err != nil {
			return err
		}
	}

	s.Revision = time.Now().UTC().Unix()

	jsonData, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(location, jsonData, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Normalize attempts to normalize the known state and prepare it for any further step cycles
func (s *State) Normalize() error {
	defer s.WriteStateFile(s.Source)
	if len(s.Pending) == 0 && len(s.Completed) == 0 && len(s.Steps) > 0 {
		for _, step := range s.Steps {
			s.Pending[step.ID] = step
		}
	}
	// did we just suffer a reboot?
	if s.CurrentStep != nil {
		currentUptime, err := uptime.Uptime()
		if err != nil {
			panic(err)
		}
		if s.CurrentStep.Status == "finished" {
			if s.CurrentStep.EndedAt.Unix() == 0 {
				s.CurrentStep.EndedAt = time.Now().UTC()
			}
			s.Completed[s.CurrentStep.ID] = s.CurrentStep
			delete(s.Pending, s.CurrentStep.ID)
			for _, x := range s.Steps {
				if x.ID == s.CurrentStep.ID {
					x = s.CurrentStep
				}
			}
			s.CurrentStep = nil
			return nil
		} else if s.CurrentStep.Status == "awaiting_reboot" && int64(time.Since(s.CurrentStep.StartedAt)*time.Second) > currentUptime {
			s.CurrentStep.Status = "finished"
			s.CurrentStep.EndedAt = time.Now().UTC()
			s.Completed[s.CurrentStep.ID] = s.CurrentStep
			delete(s.Pending, s.CurrentStep.ID)
			for _, x := range s.Steps {
				if x.ID == s.CurrentStep.ID {
					x = s.CurrentStep
				}
			}
			s.CurrentStep = nil
			return nil
		} else if s.CurrentStep.Status == "awaiting_reboot" {
			s.CurrentState = "awaiting_reboot"
			return ErrAwaitingReboot
		} else if s.CurrentStep.Status == "errored" {
			s.Errored = true
			s.ErrorMessage = s.CurrentStep.ExitError.Error()
			s.CompletedAt = time.Now().UTC()
			s.CurrentState = "errored"
			return ErrStepFailure
		} else {
			return ErrInProgress
		}
	}
	return nil
}

// WorkExists returns whether any further steps can be executed
func (s *State) WorkExists() bool {
	if len(s.Steps) == 0 {
		return false
	}

	if len(s.Completed) == len(s.Steps) {
		return false
	}

	if len(s.Pending) > 0 {
		return true
	}

	return false
}

// DoNextStep attempts to execute the subsequent step in the execution chain
func (s *State) DoNextStep() error {
	step, err := s.ResolveNextStep()
	if err != nil {
		return err
	}

	step.Prepare()
	s.CurrentStep = step
	s.Pending[step.ID] = step
	s.CurrentState = "provisioning"
	s.WriteStateFile(s.Source)
	err = step.Perform()
	if err != nil {
		step.Status = "errored"

	}
	return nil
}

// Finalize closes out all states in the provisioning config and saves it to disk
func (s *State) Finalize() error {
	return nil
}

// ResolveNextStep attempts to locate the next step in the provisioning chain
func (s *State) ResolveNextStep() (*Step, error) {
	changed := false
	defer func() {
		if changed {
			s.WriteStateFile(s.Source)
		}
	}()
	for _, x := range s.Steps {
		s, exists := s.Completed[x.ID]
		if exists {
			if x.Status != s.Status {
				changed = true
				x = s
			}
			continue
		} else if x.Status != "" {
			return x, ErrStepFailure
		}
		return x, nil
	}
	return nil, errors.New("no step was resolved")
}
