package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/karrick/godirwalk"

	"path/filepath"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
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

// Plan is a type that describes how to get from one state to the next
//easyjson:json
type Plan struct {
	TaskGroundDelay   int               `json:"ground_delay"`
	Base              *Laforge          `json:"-"`
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
	TaintedHosts      map[string]bool   `json:"tainted_hosts"`
	Walker            *dag.Walker       `json:"-"`
	Errored           bool              `json:"-"`
	FailedNodes       *dag.Set          `json:"-"`
}

// NewEmptyPlan returns an initialized, but empty plan object.
func NewEmptyPlan() *Plan {
	p := &Plan{
		TaskTypes:         map[string]string{},
		Tasks:             map[string]Doer{},
		TasksByPriority:   map[int][]string{},
		GlobalOrder:       []string{},
		OrderedPriorities: []int{},
		Tainted:           map[string]bool{},
		TaintedHosts:      map[string]bool{},
		FailedNodes:       &dag.Set{},
		TaskGroundDelay:   420,
	}
	p.Walker = &dag.Walker{
		Callback: p.Orchestrator,
	}
	return p
}

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
			alterego, found := plan.Base.Connections[path.Join(phost.Path(), "conn")]
			if found {
				if alterego.ResourceName != path.Join(phost.Path(), "conn") {
					ret[teamID] = append(ret[teamID], fmt.Sprintf("taint -allow-missing -no-color %s", alterego.ResourceName))
				}
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

// SetupTasks attempts to cull the Tasks map with Doer types to actually be performed
func (p *Plan) SetupTasks() error {
	if p.Tasks == nil {
		p.Tasks = map[string]Doer{}
	}
	for id, x := range p.GlobalOrder {
		cli.Logger.Debugf("STEP: %s", x)
		metaobj := p.Graph.Metastore[x]
		cli.Logger.Debugf("Meta: %s", metaobj.ObjectType)
		if metaobj.ObjectType == "provisioning_step" {
			pstep, ok := metaobj.Dependency.(*ProvisioningStep)
			if !ok {
				return fmt.Errorf("metadata object %s is of type %T, expected *ProvisioningStep", x, metaobj.Dependency)
			}
			var job Doer
			switch pstep.ProvisionerType {
			case "script":
				j, err := CreateScriptJob(x, id, metaobj, pstep)
				if err != nil {
					return err
				}
				job = j
			case "remote_file":
				j, err := CreateRemoteFileJob(x, id, metaobj, pstep)
				if err != nil {
					return err
				}
				job = j
			default:
				continue
			}
			job.SetTimeout(p.TaskGroundDelay)
			job.SetPlan(p)
			job.SetBase(p.Base)
			p.Tasks[x] = job
		}
	}
	return nil
}

// WriteRevisionFile writes a deng revision file
func (p *Plan) WriteRevisionFile(d Doer, status RevStatus) error {
	filename := fmt.Sprintf(".%s.pstep.lfrevision", filepath.Base(d.GetTargetID()))
	pathToRevFile := filepath.Join(p.Base.BaseDir, filepath.Dir(d.GetTargetID()), filename)
	rev := d.GetMetadata().ToRevision()
	rev.Touch()
	rev.Status = status
	err := ioutil.WriteFile(pathToRevFile, []byte(rev.ToJSONString()), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Execute walks the plan's functions against the computed dependency graph
func (p *Plan) Execute() tfdiags.Diagnostics {
	p.Walker.Update(p.Graph.AltGraph)
	err := p.Walker.Wait()
	if err.HasErrors() {
		return err
	}
	return nil
}

// Orchestrator is the walk function that is executed for each path in the dependency graph
func (p *Plan) Orchestrator(v dag.Vertex) (d tfdiags.Diagnostics) {
	if p.Errored {
		return d
	}
	id := v.(string)
	if _, ok := p.Tainted[id]; !ok {
		cli.Logger.Debugf("Node %s is unchanged. Continuing traversal.", id)
		return nil
	}
	descendents, err := p.Graph.AltGraph.Descendents(v)
	if err != nil {
		cli.Logger.Errorf("Ancestor Search Error: %v", err)
		p.Errored = true
		d.Append(tfdiags.Sourceless(tfdiags.Error, "descendent acquisition failed", tfdiags.FormatErrorPrefixed(err, id)))
		return d
	}
	if p.FailedNodes.Intersection(descendents).Len() > 0 {
		cli.Logger.Errorf("Node %s has failed lineage. Skipping execution.", id)
		d.Append(tfdiags.Sourceless(tfdiags.Error, "node has tainted lineage, skipping", id))
		return d
	}
	task, found := p.Tasks[id]
	if !found {
		cli.Logger.Errorf("Node %s did not have an associated Laforge Job! (might not be implemented yet)", id)
		// p.FailedNodes.Add(v)
		// d.Append(tfdiags.Sourceless(tfdiags.Error, "missing laforge job object for node", id))
		return d
	}
	task.SetBase(p.Base)
	err = PerformInTimeout(task.GetTimeout(), task.CanProceed)
	if err != nil {
		cli.Logger.Errorf("Task %s could not proceed: %v", id, err)
		p.FailedNodes.Add(v)
		d.Append(tfdiags.Sourceless(tfdiags.Error, "task preparation failure", tfdiags.FormatErrorPrefixed(err, id)))
		err = p.WriteRevisionFile(task, RevStatusFailed)
		if err != nil {
			d.Append(tfdiags.Sourceless(tfdiags.Error, "task cleanup failure", tfdiags.FormatErrorPrefixed(err, id)))
		}
		return d
	}
	err = PerformInTimeout(task.GetTimeout(), task.EnsureDependencies)
	if err != nil {
		cli.Logger.Errorf("Task %s failed to ensure dependencies: %v", id, err)
		p.FailedNodes.Add(v)
		d.Append(tfdiags.Sourceless(tfdiags.Error, "task dependency failure", tfdiags.FormatErrorPrefixed(err, id)))
		err = p.WriteRevisionFile(task, RevStatusFailed)
		if err != nil {
			d.Append(tfdiags.Sourceless(tfdiags.Error, "task dependency failure", tfdiags.FormatErrorPrefixed(err, id)))
		}
		return d
	}
	err = PerformInTimeout(task.GetTimeout(), task.Do)
	if err != nil {
		cli.Logger.Errorf("Task %s failed: %v", id, err)
		p.FailedNodes.Add(v)
		d.Append(tfdiags.Sourceless(tfdiags.Error, "task execution failure", tfdiags.FormatErrorPrefixed(err, id)))
		err = p.WriteRevisionFile(task, RevStatusFailed)
		if err != nil {
			d.Append(tfdiags.Sourceless(tfdiags.Error, "task execution failure", tfdiags.FormatErrorPrefixed(err, id)))
		}
		return d
	}
	err = PerformInTimeout(task.GetTimeout(), task.Finish)
	if err != nil {
		cli.Logger.Errorf("Task %s could not finish: %v", id, err)
		p.FailedNodes.Add(v)
		d.Append(tfdiags.Sourceless(tfdiags.Error, "task cleanup failure", tfdiags.FormatErrorPrefixed(err, id)))
		err = p.WriteRevisionFile(task, RevStatusFailed)
		if err != nil {
			d.Append(tfdiags.Sourceless(tfdiags.Error, "task cleanup failure", tfdiags.FormatErrorPrefixed(err, id)))
		}
		return d
	}

	err = p.WriteRevisionFile(task, RevStatusActive)
	if err != nil {
		d.Append(tfdiags.Sourceless(tfdiags.Error, "task revision file writer failure", tfdiags.FormatErrorPrefixed(err, id)))
		return d
	}
	// here is where we should do some work
	return nil
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

// BurnIt is the stub for destroying all terraform environments and their dependencies
func (p *Plan) BurnIt() error {
	errChan := make(chan error, 1)
	finChan := make(chan bool, 1)
	wg := new(sync.WaitGroup)

	sequence := []string{
		"destroy -no-color -auto-approve -parallelism=50",
	}
	for _, team := range p.Base.CurrentBuild.Teams {
		wg.Add(1)
		cli.Logger.Infof("Destroying team %s terraform environment...", team.Path())
		go team.RunTerraformSequence(sequence, wg, errChan)
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
			return p.RemoveRevisionFilesFromTeams()
		}
	}
}

// RemoveRevisionFilesFromTeams walks the team directory in the current build and removes them
func (p *Plan) RemoveRevisionFilesFromTeams() error {
	teamdirs := filepath.Join(p.Base.BaseDir, p.Base.CurrentBuild.Path(), "teams")
	err := godirwalk.Walk(teamdirs, &godirwalk.Options{
		Callback: func(ospath string, de *godirwalk.Dirent) error {
			if filepath.Ext(de.Name()) == `.lfrevision` {
				err := os.Remove(ospath)
				if err != nil {
					return err
				}
				return nil
			}
			return nil
		},
	})
	if err != nil {
		return err
	}
	return nil
}
