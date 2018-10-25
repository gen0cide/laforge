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
			err = job.EnsureDependencies()
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

// sleepcalc := map[string]int{}

// for pri, tasks := range plan.TasksByPriority {
// 	for _, x := range tasks {
// 		sleepcalc[x] = pri
// 	}
// }

// _ = tfcmds

// root, err := plan.Graph.AltGraph.Root()
// if err != nil {
// 	return err
// }

// _ = root

// walker := &dag.Walker{
// 	Callback: func(v dag.Vertex) tfdiags.Diagnostics {
// 		id := v.(string)
// 		sleeptimer, found := sleepcalc[id]
// 		if !found {
// 			return nil
// 		}
// 		cliLogger.Infof("Performing Task: %s", id)
// 		time.Sleep(time.Duration(sleeptimer) * time.Second)
// 		return nil
// 	},
// }

// walker.Update(plan.Graph.AltGraph)
// walker.Wait()
