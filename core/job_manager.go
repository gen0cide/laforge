package core

import "sync"

// Manager is a looped processor of tasks by the state machine
type Manager struct {
	Acks     chan Doer
	Errors   chan error
	Shutdown chan bool
	Inbound  chan Doer
	Laforge  *Laforge
}

// Run is the looping execution flow of the Manager
func (m *Manager) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-m.Shutdown:
			return
		case job := <-m.Inbound:
			err := job.CanProceed()
			if err != nil {
				m.Errors <- err
				continue
			}
			err = job.EnsureDependencies(m.Laforge)
			if err != nil {
				m.Errors <- err
				continue
			}
			err = job.Do()
			if err != nil {
				m.Errors <- err
				continue
			}
			err = job.Finish()
			if err != nil {
				m.Errors <- err
				continue
			}
			m.Acks <- job
			continue
		}
	}
}
