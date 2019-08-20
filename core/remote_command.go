package core

import (
	"fmt"
	"io"
	"sync"
)

// RemoteCommand represents a remote command being prepared or run.
type RemoteCommand struct {
	Command    string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	Timeout    int
	exitStatus int
	exitCh     chan struct{}
	err        error
	sync.Mutex
}

// NewRemoteCommand creates a new empty remote command object
func NewRemoteCommand() *RemoteCommand {
	r := &RemoteCommand{}
	r.Init()
	return r
}

// Init must be called before executing the command.
func (r *RemoteCommand) Init() {
	r.Lock()
	defer r.Unlock()

	r.exitCh = make(chan struct{})
}

// SetExitStatus stores the exit status of the remote command
func (r *RemoteCommand) SetExitStatus(status int, err error) {
	r.Lock()
	defer r.Unlock()

	r.exitStatus = status
	r.err = err

	close(r.exitCh)
}

// Wait waits for the remote command to complete.
func (r *RemoteCommand) Wait() error {
	<-r.exitCh

	r.Lock()
	defer r.Unlock()

	if r.err != nil || r.exitStatus != 0 {
		return &ExitError{
			Command:    r.Command,
			ExitStatus: r.exitStatus,
			Err:        r.err,
		}
	}

	return nil
}

// ExitError is returned by Wait to indicate and error executing the remote
// command, or a non-zero exit status.
type ExitError struct {
	Command    string
	ExitStatus int
	Err        error
}

// Error implements the error interface
func (e *ExitError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("error executing %q: %v", e.Command, e.Err)
	}
	return fmt.Sprintf("%q exit status: %d", e.Command, e.ExitStatus)
}
