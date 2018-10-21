package core

import (
	"fmt"
	"sync"
	"time"
)

var (
	requiresTF = map[LFType]bool{
		LFTypeBuild:              true,
		LFTypeEnvironment:        true,
		LFTypeProvisionedHost:    true,
		LFTypeProvisionedNetwork: true,
		LFTypeTeam:               true,
	}

	requiresTFTaint = map[LFType]bool{
		LFTypeProvisionedHost: true,
	}
)

// CalculateTerraformNeeds attempts to determine what terraform steps must happen first
// for a given plan.
func CalculateTerraformNeeds(plan *Plan) (map[string][]string, error) {
	ret := map[string][]string{}
	teamsRequiringTFApply := map[string]bool{}
	for _, x := range plan.GlobalOrder {
		kind := TypeByPath(x)
		if _, ok := requiresTF[kind]; !ok {
			continue
		}
		teamID, err := GetTeamIDFromPath(x)
		if err != nil {
			continue
		}
		if ret[teamID] == nil || len(ret[teamID]) < 1 {
			ret[teamID] = []string{}
			ret[teamID] = append(ret[teamID], "refresh -no-color")
		}
		if kind == LFTypeProvisionedHost {
			host := plan.Graph.Metastore[x].Dependency
			phost, ok := host.(*ProvisionedHost)
			if !ok {
				continue
			}
			if phost.Conn == nil {
				teamsRequiringTFApply[teamID] = true
				continue
			}
			if phost.Conn.ResourceName != phost.Conn.Path() {
				ret[teamID] = append(ret[teamID], fmt.Sprintf("taint -allow-missing -no-color %s", phost.Conn.ResourceName))
			}
			teamsRequiringTFApply[teamID] = true
		}
	}

	// now to clean up
	for tid := range teamsRequiringTFApply {
		_ = tid
		ret[tid] = append(ret[tid], "apply -no-color -auto-approve -parallelism=10")
	}
	return ret, nil
}

// Plan is a type that describes how to get from one state to the next
//easyjson:json
type Plan struct {
	Checksum          uint64            `json:"checksum"`
	StartedAt         time.Time         `json:"started_at"`
	EndedAt           time.Time         `json:"ended_at"`
	Graph             *Snapshot         `json:"target,omitempty"`
	TaskTypes         map[string]string `json:"task_types"`
	Tasks             map[string]Doer   `json:"-"`
	TasksByPriority   map[int][]string  `json:"tasks_by_priority"`
	GlobalOrder       []string          `json:"global_order"`
	OrderedPriorities []int             `json:"ordered_priorities"`
	Tainted           map[string]bool   `json:"tainted"`
}

// Preflight determines what teams need terraform run on them, executing them before the plan
func (p *Plan) Preflight() error {
	tfruns, err := CalculateTerraformNeeds(p)
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	wg := new(sync.WaitGroup)

	for tid, cmds := range tfruns {
		tmeta, ok := p.Graph.Metastore[tid]
		if !ok {
			return fmt.Errorf("team %s is not in the graph", tid)
		}
		tobj, ok := tmeta.Dependency.(*Team)
		if !ok {
			return fmt.Errorf("team %s did not have a *Team dependency type", tid)
		}
		wg.Add(1)
		go tobj.RunTerraformSequence(cmds, wg, errChan)
	}

	go func() {
		wg.Wait()
		close(finChan)
	}()

	errored := false
	var exiterror error

	for {
		select {
		case err := <-errChan:
			exiterror = err
			return err
		case <-finChan:
			if errored {
				return exiterror
			}
			return nil
		}
	}
}
