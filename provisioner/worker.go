package provisioner

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrWorkerBusy is thrown when the work is already executing
	ErrWorkerBusy = errors.New("agent worker busy")

	// ErrLoadTimedOut is thrown when an attempt to Load(statefile) is not responsive after a few seconds
	ErrLoadTimedOut = errors.New("agent worker timed out loading state")

	// ErrStaleRevision is thrown when the config attempting to be loaded has an out of date revision
	ErrStaleRevision = errors.New("requested revision is out of date")

	// ErrRevisionMismatch is thrown when the config attempting to be loaded has a future revision
	ErrRevisionMismatch = errors.New("requested revision is ahead of currently running")

	// ErrDuplicateRevision is thrown when a config of the same revision is attempted to be loaded
	ErrDuplicateRevision = errors.New("duplicate revision already loaded")
)

// LoadRequest is a stub type that lets the Engine communicate config Load's to the worker without
// breaking async execution
type LoadRequest struct {
	Source string
	Ack    chan bool
	Err    chan error
}

// Worker is the base unit that performs the step execution of the agent
type Worker struct {
	Heartbeat  int64
	Busy       bool
	Config     *State
	ConfigFile string
	Tasks      chan LoadRequest
}

// Spawn creates the task queue and launches the worker loop in a separate goroutine
func (w *Worker) Spawn() {
	w.Tasks = make(chan LoadRequest, 1)
	go w.Do()
}

// Available is a helper function to check to see if the worker is currently working
func (w *Worker) Available() bool {
	return !w.Busy
}

// Do is the primary loop of the worker's execution
func (w *Worker) Do() {
	for s := range w.Tasks {
		w.Busy = true
		w.Handle(s)
		w.Busy = false
	}
}

// Rest simply pauses the worker function for a period of time to cooldown and handle graceful exits by scripts
func (w *Worker) Rest() {
	if runtime.GOOS == "windows" {
		time.Sleep(5 * time.Second)
		return
	}
	time.Sleep(3 * time.Second)
}

// Load will create a LoadRequest and pass it into the Worker's task queue, blocking until the Worker Ack's or an error is thrown
func (w *Worker) Load(statefile string) error {
	if w.Busy {
		return ErrWorkerBusy
	}
	lr := LoadRequest{
		Source: statefile,
		Ack:    make(chan bool, 1),
		Err:    make(chan error, 1),
	}
	w.Tasks <- lr
	select {
	case <-lr.Ack:
		return nil
	case err := <-lr.Err:
		return err
	case <-time.After(5 * time.Second):
		return ErrLoadTimedOut
	}
}

// Handle is the primary loop for managing the execution of a Config's lifecycle
func (w *Worker) Handle(lr LoadRequest) {
	state, err := LoadStateFile(lr.Source)
	if err != nil {
		lr.Err <- err
		return
	}

	if w.Config != nil {
		if state.Revision == w.Config.Revision {
			// probably need to reload and see whats happening here
			lr.Err <- ErrDuplicateRevision
			return
		} else if state.Revision < w.Config.Revision {
			lr.Err <- ErrStaleRevision
			return
		}
		lr.Err <- ErrRevisionMismatch
		return
	}

	w.Config = state
	w.ConfigFile = lr.Source
	w.Config.Source = lr.Source
	lr.Ack <- true

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for t := range ticker.C {
			atomic.StoreInt64(&w.Heartbeat, t.UTC().Unix())
		}
	}()
	defer ticker.Stop()

	for {
		err = w.Config.Normalize()
		if err != nil {
			if err == ErrAwaitingReboot {
				w.Rest()
				continue
			}
			Logger.Errorf("error attempting to normalize state: %v", err)
			return
		}

		if !w.Config.WorkExists() {
			break
		}
		err = w.Config.DoNextStep()
		if err != nil {
			return
		}
		w.Rest()
	}
	w.Config.Finalize()
}
