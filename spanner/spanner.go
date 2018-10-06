package spanner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gen0cide/laforge/core"
)

type Spanner struct {
	Job
	Laforge  *core.Laforge
	TeamDirs map[int]string
	Workers  map[int]*Worker
	LogDir   string
	BuildDir string
	Result   chan *Worker
	ExecTime int64
}

type Job struct {
	Command    []string
	SaveOutput bool
	Silent     bool
	ExecType   string
	HostID     string
}

// New returns a new Spanner
func New(base *core.Laforge, command []string, exectype, hostid string, saveOutput, silent bool) (*Spanner, error) {
	var err error
	if base == nil {
		base, err = core.Bootstrap()
		if err != nil {
			return nil, err
		}
	}

	err = base.AssertExactContext(core.BuildContext)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now().Unix()
	logDir := filepath.Join(base.BuildRoot, "spanner_logs", fmt.Sprintf("%d", timeNow))

	os.MkdirAll(logDir, 0755)

	return &Spanner{
		Laforge:  base,
		BuildDir: base.BuildRoot,
		TeamDirs: map[int]string{},
		Workers:  map[int]*Worker{},
		LogDir:   logDir,
		Result:   make(chan *Worker),
		ExecTime: timeNow,
		Job: Job{
			Command:    command,
			ExecType:   exectype,
			HostID:     hostid,
			SaveOutput: saveOutput,
			Silent:     silent,
		},
	}, nil
}

// CreateWorkerPool creates the individual workers for the spanner
func (s *Spanner) CreateWorkerPool() error {
	for i := 0; i < s.Laforge.Environment.TeamCount; i++ {
		teamDir := filepath.Join(s.BuildDir, "teams", fmt.Sprintf("%d", i))

		if _, err := os.Stat(teamDir); os.IsNotExist(err) {
			return fmt.Errorf("no team directory exists for team number %d", i)
		}
		s.Workers[i] = &Worker{
			TeamDir: teamDir,
			Job:     s.Job,
			ID:      s.ExecTime,
			TeamID:  i,
			Parent:  s,
			LogFile: filepath.Join(s.LogDir, fmt.Sprintf("team_%d.log", i)),
		}
		notValid := s.Workers[i].Verify()
		if notValid != nil {
			return notValid
		}
	}
	return nil
}

// Verify attempts to validate the constructs of the spanner
func (w *Worker) Verify() error {
	if w.ExecType == "remote-exec" {
		provisionedHostFile := filepath.Join(w.TeamDir, "provisioned_hosts", fmt.Sprintf("%s.laforge", w.HostID))
		if _, err := os.Stat(provisionedHostFile); os.IsNotExist(err) {
			return fmt.Errorf("team %d does not have an active host %s", w.TeamID, w.HostID)
		}
		err := w.ResolveProvisionedHost()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Spanner) Do() error {
	core.SetLogLevel("info")
	core.Logger.Infof("Logs found at: %s", s.LogDir)

	switch s.ExecType {
	case "local-exec":
		for _, w := range s.Workers {
			go w.RunLocalCommand(s.Result)
		}
	case "remote-exec":
		for _, w := range s.Workers {
			go w.RunRemoteCommand(s.Result)
		}
	default:
		return errors.New("invalid exectype passed into spanner")
	}

	for i := 0; i < len(s.Workers); i++ {
		worker := <-s.Result
		if worker.ExitError != nil {
			core.Logger.Errorf("Worker for team %d errored: %v", worker.TeamID, worker.ExitError)
		} else {
			core.Logger.Infof("Worker for team %d has completed successfully.", worker.TeamID)
		}
	}

	return nil
}
