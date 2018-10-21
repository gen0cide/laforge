package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"sync"

	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
	"github.com/karrick/godirwalk"

	"github.com/gen0cide/laforge/core/cli"
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
	Base      *Laforge
	DB        *buntdb.DB
	Current   *Snapshot
	Persisted *Snapshot
	Plan      *Plan
	NewRevs   map[string]*Revision
	KnownRevs map[string]*Revision
	RevDelta  map[string]RevMod
}

// LocateRevisions attempts to load the known revision files off disk
func (s *State) LocateRevisions() error {
	if s.KnownRevs == nil {
		s.KnownRevs = map[string]*Revision{}
	}
	if err := s.Base.AssertMinContext(EnvContext); err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	errChan := make(chan error, 1)
	revChan := make(chan *Revision, 2000)
	finChan := make(chan bool, 1)

	dirname := s.Base.EnvRoot
	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if filepath.Ext(de.Name()) == `.lfrevision` {
				wg.Add(1)
				go func(fp string) {
					defer wg.Done()
					rev, err := ParseRevisionFile(fp)
					if err != nil {
						errChan <- errors.Wrapf(err, "revfile=%s", fp)
						return
					}
					revChan <- rev
					return
				}(osPathname)
			}
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})

	if err != nil {
		return err
	}

	go func() {
		wg.Wait()
		close(finChan)
	}()

	errored := false

	for {
		select {
		case rev := <-revChan:
			s.KnownRevs[rev.ID] = rev
			continue
		case err := <-errChan:
			cli.Logger.Errorf("Error reading revision file: %v", err)
			errored = true
			continue
		default:
		}

		select {
		case rev := <-revChan:
			s.KnownRevs[rev.ID] = rev
			continue
		case err := <-errChan:
			cli.Logger.Errorf("Error reading revision file: %v", err)
			errored = true
			continue
		case <-finChan:
			if errored {
				return errors.New("revision file loading has failed")
			}
			return nil
		}
	}
}

// GenerateCurrentRevs enumerates the current snapshot and generates a listing of revisions for comparison
func (s *State) GenerateCurrentRevs() error {
	if s.NewRevs == nil {
		s.NewRevs = map[string]*Revision{}
	}
	for _, x := range s.Current.Metastore {
		rev := x.ToRevision()
		s.NewRevs[rev.ID] = rev
	}
	return nil
}

// GenerateRevisionDelta compares the known verses the new revisiosn and comes up with a strategy which
// is used in the plan calculations.
func (s *State) GenerateRevisionDelta() error {
	if s.RevDelta == nil {
		s.RevDelta = map[string]RevMod{}
	}
	nrkeys := make([]string, len(s.NewRevs))
	for nrid, nrev := range s.NewRevs {
		nrkeys = append(nrkeys, nrid)
		krev, ok := s.KnownRevs[nrid]
		if !ok {
			s.RevDelta[nrid] = RevModCreate
			continue
		}
		if nrev.Checksum != krev.Checksum {
			cli.Logger.Warnf("Marking %s for a failed revision checksum comparison", nrid)
			s.RevDelta[nrid] = RevModTouch
			continue
		}
	}
	for knid := range s.KnownRevs {
		if _, ok := s.NewRevs[knid]; !ok {
			s.RevDelta[knid] = RevModDelete
		}
	}
	return nil
}

// NewRevHashes returns a hash of the new revision objects
func (s *State) NewRevHashes() uint64 {
	hashes := ChecksumList{}
	for _, x := range s.NewRevs {
		hashes = append(hashes, x.Checksum)
	}
	return hashes.Hash()
}

// KnownRevHashes returns a hash of the located revision objects
func (s *State) KnownRevHashes() uint64 {
	hashes := ChecksumList{}
	for _, x := range s.KnownRevs {
		hashes = append(hashes, x.Checksum)
	}
	return hashes.Hash()
}

// SnapshotsEqual are used to test the equality of the two environments and their dependencies
func (s *State) SnapshotsEqual() bool {
	if s.Persisted.Hash() != s.Current.Hash() {
		return false
	}
	if s.KnownRevHashes() != s.NewRevHashes() {
		return false
	}
	return true
}

// CalculateDelta attempts to determine what needs to be done to bring a base in line with target
func (s *State) CalculateDelta() (*Plan, error) {
	if s.Persisted == nil {
		return nil, errors.New("the persisted state is nil and delta analysis cannot be performed")
	}
	if s.Current == nil {
		return nil, errors.New("the current state is nil and delta analysis cannot be performed")
	}

	err := s.LocateRevisions()
	if err != nil {
		return nil, err
	}

	err = s.GenerateCurrentRevs()
	if err != nil {
		return nil, err
	}

	err = s.GenerateRevisionDelta()
	if err != nil {
		return nil, err
	}

	if s.SnapshotsEqual() {
		return nil, ErrSnapshotsMatch
	}

	base := s.Persisted
	target := s.Current

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
		if _, ok := s.RevDelta[v.(string)]; ok {
			tainted = append(tainted, v)
			changeslice = append(changeslice, targetmeta.ID)
			return nil
		}
		if targetmeta.Checksum != basemeta.Checksum {
			tainted = append(tainted, v)
			changeslice = append(changeslice, targetmeta.ID)
			return nil
		}
		return nil
	})

	target.Graph.Walk(func(v dag.Vertex) tfdiags.Diagnostics {
		targetmeta := target.Metastore[v.(string)]
		_, exists := base.Metastore[v.(string)]
		if !exists {
			tainted = append(tainted, v)
			changeslice = append(changeslice, targetmeta.ID)
			return nil
		}
		if _, ok := s.RevDelta[v.(string)]; ok {
			tainted = append(tainted, v)
			changeslice = append(changeslice, targetmeta.ID)
			return nil
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

		for _, x := range dag.AsVertexList(children) {
			newtaints = append(newtaints, x)
		}
	}

	for _, t := range tainted {
		children, err := target.Graph.Descendents(t)
		if err != nil {
			return nil, err
		}

		for _, x := range dag.AsVertexList(children) {
			if IsGlobalType(x.(string)) {
				newtaints = append(newtaints, x)
			}
		}
	}

	// now we walk the base and literally tell it to gtfo out
	edgeRemovals := []dag.Edge{}
	vertRemovals := []dag.Vertex{}
	base.Graph.DepthFirstWalk(tainted, func(v dag.Vertex, depth int) error {
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
		if !be {
			if IsGlobalType(v.(string)) {
				return nil
			}
			if action, found := s.RevDelta[bo.ID]; found {
				tasktypes[bo.ID] = string(action)
			} else {
				tasktypes[to.ID] = "CREATE"
			}
			tasks[depth] = append(tasks[depth], v.(string))
			return nil
		}
		if !te {
			if action, found := s.RevDelta[to.ID]; found {
				tasktypes[to.ID] = string(action)
			} else {
				tasktypes[to.ID] = "DELETE"
			}
			tasks[depth] = append(tasks[depth], v.(string))
			return nil
		}

		brev, brevok := s.KnownRevs[bo.ID]
		trev, trevok := s.NewRevs[to.ID]
		if brevok && trevok && brev.Checksum == trev.Checksum {
			return nil
		}

		if TypeByPath(v.(string)) == LFTypeBuild {
			return nil
		}

		if IsGlobalType(v.(string)) {
			if to.Checksum != bo.Checksum {
				tasktypes[v.(string)] = "TOUCH"
				tasks[depth] = append(tasks[depth], v.(string))
			}
			return nil
		}

		if action, found := s.RevDelta[to.ID]; found {
			tasktypes[v.(string)] = string(action)
			tasks[depth] = append(tasks[depth], v.(string))
			return nil
		}

		tasktypes[v.(string)] = "UNKNOWN"
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
