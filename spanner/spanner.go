package spanner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gen0cide/laforge/core"
	"github.com/gen0cide/laforge/core/cli"
)

// Spanner is a multiplexer for teams in a laforge environment
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
	tn := 0
	for tid, team := range s.Laforge.CurrentBuild.Teams {
		teamDir := filepath.Join(s.BuildDir, "teams", fmt.Sprintf("%d", tn))

		if _, err := os.Stat(teamDir); os.IsNotExist(err) {
			return fmt.Errorf("no team directory exists for team number %d", tn)
		}
		s.Workers[tn] = &Worker{
			TeamDir:    teamDir,
			Job:        s.Job,
			ID:         s.ExecTime,
			TeamID:     tid,
			TeamNumber: tn,
			Team:       team,
			Parent:     s,
			LogFile:    filepath.Join(s.LogDir, fmt.Sprintf("team_%d.log", tn)),
		}
		notValid := s.Workers[tn].Verify()
		if notValid != nil {
			return notValid
		}
	}
	return nil
}

// Do does stuff
func (s *Spanner) Do() error {
	cli.SetLogLevel("info")
	cli.Logger.Infof("Logs found at: %s", s.LogDir)

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
			cli.Logger.Errorf("Worker for team %d errored: %v", worker.TeamID, worker.ExitError)
		} else {
			cli.Logger.Infof("Worker for team %d has completed successfully.", worker.TeamID)
		}
	}

	return nil
}
