package tf

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gen0cide/gscript/logger"

	"github.com/gruntwork-io/terragrunt/errors"
)

type Runner struct {
	Owner      string
	ID         int64
	Logger     logger.Logger
	BaseDir    string
	StdoutFile string
	StderrFile string
	Program    string
	Args       []string
	FinChan    chan bool
	Output     chan string
	Errors     chan error
	EnvVars    map[string]string
}

func NewRunner(owner, basedir string, l logger.Logger) *Runner {
	id := time.Now().UTC().Unix()
	logdir := filepath.Join(basedir, "logs")
	os.MkdirAll(logdir, 0755)
	stderrbase := fmt.Sprintf("%d_stderr.log", id)
	stdoutbase := fmt.Sprintf("%d_stdout.log", id)
	stderrfile := filepath.Join(logdir, stderrbase)
	stdoutfile := filepath.Join(logdir, stdoutbase)

	m := make(map[string]string)
	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			m[e[:i]] = e[i+1:]
		}
	}

	return &Runner{
		Owner:      owner,
		ID:         id,
		Logger:     l,
		BaseDir:    basedir,
		StdoutFile: stdoutfile,
		StderrFile: stderrfile,
		FinChan:    make(chan bool, 1),
		Output:     make(chan string, 50),
		Errors:     make(chan error, 10),
		EnvVars:    m,
	}
}

func (r *Runner) ExecuteCommand(command string, args ...string) {
	defer func() {
		close(r.FinChan)
	}()
	r.Program = command
	r.Args = args

	cmd := exec.Command(command, args...)

	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		r.Errors <- errors.WithStackTrace(err)
		return
	}

	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		r.Errors <- errors.WithStackTrace(err)
		return
	}

	stdoutFile, err := os.Create(r.StdoutFile)
	if err != nil {
		r.Errors <- errors.WithStackTrace(err)
		return
	}
	defer stdoutFile.Close()

	stderrFile, err := os.Create(r.StderrFile)
	if err != nil {
		r.Errors <- errors.WithStackTrace(err)
		return
	}
	defer stderrFile.Close()

	cmd.Env = os.Environ()
	cmd.Dir = r.BaseDir

	if err := cmd.Start(); err != nil {
		r.Errors <- errors.WithStackTrace(err)
		return
	}

	var wg sync.WaitGroup
	stdoutScanner := bufio.NewScanner(stdoutIn)
	stderrScanner := bufio.NewScanner(stderrIn)
	wg.Add(2)

	go func() {
		defer wg.Done()
		for stdoutScanner.Scan() {
			text := stdoutScanner.Text()
			r.Output <- text
			fmt.Fprintln(stdoutFile, text)
		}
	}()

	go func() {
		defer wg.Done()
		for stderrScanner.Scan() {
			text := stderrScanner.Text()
			r.Output <- text
			fmt.Fprintln(stderrFile, text)
		}
	}()

	cmdChannel := make(chan error)
	signalChannel := NewSignalsForwarder(forwardSignals, cmd, r.Logger, cmdChannel)
	defer signalChannel.Close()

	err = cmd.Wait()
	cmdChannel <- err

	wg.Wait()

	if err := stdoutScanner.Err(); err != nil {
		r.Errors <- errors.WithStackTrace(err)
	}

	if err := stderrScanner.Err(); err != nil {
		r.Errors <- errors.WithStackTrace(err)
	}

	return
}

// GetExitCode attempts to unrap an exit code into something that gives more context into why the runner exited
func GetExitCode(err error) (int, error) {
	if exiterr, ok := errors.Unwrap(err).(errors.IErrorCode); ok {
		return exiterr.ExitStatus()
	}

	if exiterr, ok := errors.Unwrap(err).(*exec.ExitError); ok {
		status := exiterr.Sys().(syscall.WaitStatus)
		return status.ExitStatus(), nil
	}

	if exiterr, ok := errors.Unwrap(err).(errors.MultiError); ok {
		for _, err := range exiterr.Errors {
			exitCode, exitCodeErr := GetExitCode(err)
			if exitCodeErr == nil {
				return exitCode, nil
			}
		}
	}

	return 0, err
}

// SignalsForwarder passes signals between the current program and the runners sub processes
type SignalsForwarder chan os.Signal

// NewSignalsForwarder creates a new signals forwarder for the runner
func NewSignalsForwarder(signals []os.Signal, c *exec.Cmd, logger logger.Logger, cmdChannel chan error) SignalsForwarder {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, signals...)

	go func() {
		for {
			select {
			case s := <-signalChannel:
				logger.Infof("Forward signal %v to terraform.", s)
				err := c.Process.Signal(s)
				if err != nil {
					logger.Errorf("Error forwarding signal: %v", err)
				}
			case <-cmdChannel:
				return
			}
		}
	}()

	return signalChannel
}

// Close closes the SignalForwardering goroutine
func (s *SignalsForwarder) Close() error {
	signal.Stop(*s)
	*s <- nil
	close(*s)
	return nil
}
