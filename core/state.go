package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

const (
	// DBKeySnapshot is the database key for persisting the snapshot in our local filesystem
	DBKeySnapshot = `/snapshot`
)

var (
	// ErrSnapshotsMatch is thrown when two snapshots are functionally the same during delta calculation
	ErrSnapshotsMatch = errors.New("snapshots reflected the same state")
)

// State is the primary object used to interface with the build's on disk state table
type State struct {
	DB        *buntdb.DB
	Current   *Snapshot
	Persisted *Snapshot
	Plan      *Plan
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

// CalculateDelta attempts to determine what needs to be done to bring a base in line with target
func CalculateDelta(base *Snapshot, target *Snapshot) (*Plan, error) {
	if base.Hash() == target.Hash() {
		return nil, ErrSnapshotsMatch
	}
	base.Graph.SetDebugWriter(ioutil.Discard)
	target.Graph.SetDebugWriter(ioutil.Discard)
	log.SetOutput(ioutil.Discard)

	tainted := []dag.Vertex{}
	changemap := map[string]bool{}
	changeslice := []string{}

	for k, m := range target.Metastore {
		if TypeByPath(k) != LFTypeProvisionedHost {
			continue
		}
		ph, ok := m.Dependency.(*ProvisionedHost)
		if !ok {
			cli.Logger.Errorf("Provisioned Host %s was not of proper type", k)
		}
		if ph.Conn == nil || ph.Conn.RemoteAddr == NullIP {
			changemap[k] = true
			tainted = append(tainted, k)
		}
	}

	// Find deletions or updates by walking the base and comparing against target
	base.Graph.Walk(func(v dag.Vertex) tfdiags.Diagnostics {
		basemeta := base.Metastore[v.(string)]
		targetmeta, exists := target.Metastore[v.(string)]
		if !exists {
			tainted = append(tainted, v)
			return nil
		}
		// if basemeta == nil {
		// 	tainted = append(tainted, v)
		// 	return nil
		// }
		if targetmeta.Checksum != basemeta.Checksum {
			tainted = append(tainted, v)
			// changemap[targetmeta.ID] = true
			changeslice = append(changeslice, targetmeta.ID)
		}
		return nil
	})

	target.Graph.Walk(func(v dag.Vertex) tfdiags.Diagnostics {
		targetmeta := target.Metastore[v.(string)]
		_, exists := base.Metastore[v.(string)]
		if !exists {
			tainted = append(tainted, v)
			changeslice = append(changeslice, targetmeta.ID)
			// changemap[targetmeta.ID] = true
		}
		return nil
	})

	for _, x := range changeslice {
		changemap[x] = true
	}

	// need to account for partially deployed hosts (should be rolled back and redeployed)
	if len(tainted) > 0 {
		for _, t := range tainted {
			if t == nil {
				continue
			}
			objid := t.(string)
			objtype := TypeByPath(objid)
			_ = objtype
			if TypeByPath(t.(string)) == LFTypeProvisioningStep {
				obj, ok := target.Metastore[t.(string)]
				if !ok {
					continue
				}
				badHost := obj.Dependency.ParentLaforgeID()
				changemap[badHost] = true
				tainted = append([]dag.Vertex{badHost}, tainted...)
			}
		}
	}

	// now we find all the decendents to also taint, ensuring we dont traverse the whole damn thing
	newtaints := []dag.Vertex{}
	for _, t := range tainted {
		children, err := base.Graph.Descendents(t)
		if err != nil {
			return nil, err
		}
		pruned := children.Filter(func(i interface{}) bool {
			chobj := base.Metastore[i.(string)]
			_ = chobj
			// if chobj.IsGlobalType() && !changemap[chobj.ID] {
			// 	return false
			// }
			return true
		})

		for _, x := range dag.AsVertexList(pruned) {
			newtaints = append(newtaints, x)
		}
	}

	for _, t := range tainted {
		children, err := target.Graph.Descendents(t)
		if err != nil {
			return nil, err
		}
		pruned := children.Filter(func(i interface{}) bool {
			chobj := target.Metastore[i.(string)]
			_ = chobj
			// if chobj.IsGlobalType() && !changemap[chobj.ID] {
			// 	return false
			// }
			return true
		})

		for _, x := range dag.AsVertexList(pruned) {
			if IsGlobalType(x.(string)) {
				newtaints = append(newtaints, x)
			}
		}
	}

	// now we walk the base and literally tell it to gtfo out
	edgeRemovals := []dag.Edge{}
	vertRemovals := []dag.Vertex{}
	base.Graph.DepthFirstWalk(newtaints, func(v dag.Vertex, depth int) error {
		for _, e := range base.Graph.EdgesFrom(v) {
			edgeRemovals = append(edgeRemovals, e)
		}
		vertRemovals = append(vertRemovals, v)
		return nil
	})

	for _, e := range edgeRemovals {
		base.Graph.RemoveEdge(e)
	}
	for _, v := range vertRemovals {
		base.Graph.Remove(v)
	}

	// Now we can finally make this freakin list
	tasks := map[int][]string{}
	tasktypes := map[string]string{}

	target.Graph.TransitiveReduction()

	// Now it's up to us to walk the new taint list in the *target* graph to generate the work order
	target.Graph.DepthFirstWalk(newtaints, func(v dag.Vertex, depth int) error {
		if tasks[depth] == nil {
			tasks[depth] = []string{}
		}
		bo, be := base.Metastore[v.(string)]
		to, te := target.Metastore[v.(string)]
		_ = bo
		_ = to
		if !be {
			if IsGlobalType(v.(string)) {
				return nil
			}
		}
		if !te {
			tasktypes[v.(string)] = "DESTROY"
			tasks[depth] = append(tasks[depth], v.(string))
			return nil
		}

		if TypeByPath(v.(string)) == LFTypeBuild {
			return nil
		}

		if IsGlobalType(v.(string)) {
			tasktypes[v.(string)] = "REFRESH"
			tasks[depth] = append(tasks[depth], v.(string))
			return nil
		}
		tasktypes[v.(string)] = "MODIFY"
		tasks[depth] = append(tasks[depth], v.(string))
		return nil
	})

	// lets sort these mother fuckers and get on our way
	taskkeys := []int{}
	for k := range tasks {
		taskkeys = append(taskkeys, k)
	}
	sort.Ints(taskkeys)
	taintedmap := map[string]bool{}
	globalorder := []string{}

	for _, v := range taskkeys {
		for _, l := range tasks[v] {
			taintedmap[l] = true
			globalorder = append(globalorder, l)
		}
	}

	return &Plan{
		Graph:             target,
		GlobalOrder:       globalorder,
		Tainted:           taintedmap,
		OrderedPriorities: taskkeys,
		TasksByPriority:   tasks,
		TaskTypes:         tasktypes,
		Tasks:             map[string]Doer{},
	}, nil
}

// Open attempts to create a DB connector for the state given a local file path
func (s *State) Open(dbfile string) error {
	db, err := buntdb.Open(dbfile)
	if err != nil {
		return err
	}
	s.DB = db
	return nil
}

// NewState returns an empty state
func NewState() *State {
	return &State{}
}

// SetCurrent sets the current snapshot
func (s *State) SetCurrent(snap *Snapshot) {
	s.Current = snap
}

// InitializeEmptyPersistedSnapshot sets the persisted snapshot in the state to an empty one
func (s *State) InitializeEmptyPersistedSnapshot() error {
	s.Persisted = NewEmptySnapshot()
	return nil
}

// LoadSnapshotFromDB attempts to load the last Snapshot object from the DB, assigning it to *State.Persisted and returning it if it was successful.
func (s *State) LoadSnapshotFromDB() (*Snapshot, error) {
	if s.DB == nil {
		return nil, errors.New("database driver is not initialized")
	}

	snap := NewEmptySnapshot()
	err := s.DB.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(DBKeySnapshot)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(val), snap)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	err = snap.RebuildGraph()
	if err != nil {
		return nil, err
	}

	s.Persisted = snap
	return snap, nil
}

// CreateDBSchema attempts to create the database indexes appropriately
func (s *State) CreateDBSchema() error {
	return nil
}

// PersistSnapshot will save the provided snapshot into the current snapshot entry of the database, overwriting any existing snapshot.
func (s *State) PersistSnapshot(snap *Snapshot) error {
	jsonData, err := json.Marshal(snap)
	if err != nil {
		return err
	}
	err = s.DB.Update(func(tx *buntdb.Tx) error {
		_, overwritten, err := tx.Set(DBKeySnapshot, string(jsonData), nil)
		if err != nil {
			return err
		}
		if overwritten {
			cli.Logger.Infof("Persistent Snapshot overwritten in state DB")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
