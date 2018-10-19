package core

import (
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"github.com/yourbasic/graph"
)

// Snapshot is a graph traversal structure used to create a unique fingerprint for all elements in an environment.
type Snapshot struct {
	Counter       int                   `json:"counter"`
	Graph         *graph.Mutable        `json:"-"`
	Sorted        *graph.Immutable      `json:"-"`
	OrderedGIDs   []int                 `json:"ordered_gids"`
	Metadata      map[string]*Metadata  `json:"metadata"`
	Objects       map[string]Dependency `json:"-"`
	GIDToPath     map[int]string        `json:"gid_to_path"`
	PathToGID     map[string]int        `json:"path_to_gid"`
	GIDToObject   map[int]Dependency    `json:"-"`
	ObjectToGID   map[Dependency]int    `json:"-"`
	RelationQueue chan RelateFunc       `json:"-"`
}

// DependencyList is a helper type alias to handle topological sorting of laforge state elements
type DependencyList []int

// RelateFunc is a type alias to a function that relates objects together in a promise style format
type RelateFunc func() error

// Plot writes a DOT representation of the Snapshot to the output writer, with an optional path to denote the root object.
func (s *Snapshot) Plot(output io.Writer, rootID string) error {
	if rootID != "" {
		_, ok := s.Objects[rootID]
		if !ok {
			return fmt.Errorf("the node %s could not be found in the graph", rootID)
		}
	}
	dw := NewDotWriter(output, len(s.Metadata), false)
	rootMeta := s.Metadata[rootID]
	dw.PlotGraph(rootMeta)
	return nil
}

// NewEmptySnapshot returns an empty snapshot object
func NewEmptySnapshot() *Snapshot {
	return &Snapshot{
		Counter:       0,
		GIDToObject:   map[int]Dependency{},
		ObjectToGID:   map[Dependency]int{},
		Objects:       map[string]Dependency{},
		Metadata:      map[string]*Metadata{},
		GIDToPath:     map[int]string{},
		PathToGID:     map[string]int{},
		RelationQueue: make(chan RelateFunc, 10000),
	}
}

// NewSnapshotFromEnv creates a new snapshot from a provided environment
func NewSnapshotFromEnv(e *Environment) (*Snapshot, error) {
	s := NewEmptySnapshot()
	build := e.CreateBuild()
	err := build.CreateTeams()
	if err != nil {
		return nil, err
	}
	err = build.Gather(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// AddNode is used to generate verticies on the graph for laforge state elements
func (s *Snapshot) AddNode(deps ...Dependency) error {
	for _, dep := range deps {
		if _, ok := s.ObjectToGID[dep]; !ok {
			if taken, bad := s.GIDToObject[s.Counter]; bad {
				return fmt.Errorf("unable to add %s to graph: GID %d is already taken by %s", dep.Path(), s.Counter, taken.Path())
			}
			m := &Metadata{
				ID:         dep.Path(),
				GID:        s.Counter,
				ObjectType: DependencyType(dep),
				GCost:      DependencyCost(dep),
				Dependency: dep,
				ParentIDs:  []string{},
				ChildIDs:   []string{},
				ParentDeps: []DotNode{},
				ChildDeps:  []DotNode{},
				ParentGIDs: []int{},
				ChildGIDs:  []int{},
			}
			s.Metadata[dep.Path()] = m
			m.CalculateChecksum()
			s.Objects[dep.Path()] = dep
			s.GIDToObject[s.Counter] = dep
			s.ObjectToGID[dep] = s.Counter
			s.GIDToPath[s.Counter] = dep.Path()
			s.PathToGID[dep.Path()] = s.Counter
			s.Counter++
		}
	}
	return nil
}

// RebuildGraph will attempt to recreate the dependency tree based on it's persisted metadata instead of building from environment.
func (s *Snapshot) RebuildGraph() error {
	if s.Graph != nil {
		return errors.New("snapshot already has built dependency graph")
	}
	s.Graph = graph.New(len(s.Metadata))

	for _, mobj := range s.Metadata {
		mobj.ParentDeps = make([]DotNode, len(mobj.ParentGIDs))
		for offset := 0; offset < len(mobj.ParentGIDs); offset++ {
			ppath := mobj.ParentIDs[offset]
			pgid := mobj.ParentGIDs[offset]
			par := s.Metadata[ppath]
			mobj.ParentDeps[offset] = par
			s.Graph.AddCost(pgid, mobj.GID, mobj.GCost)
		}
		mobj.ChildDeps = make([]DotNode, len(mobj.ChildGIDs))
		for offset := 0; offset < len(mobj.ChildGIDs); offset++ {
			cpath := mobj.ChildIDs[offset]
			cgid := mobj.ChildGIDs[offset]
			chi := s.Metadata[cpath]
			mobj.ChildDeps[offset] = chi
			s.Graph.AddCost(mobj.GID, cgid, chi.GCost)
		}
	}
	return nil
}

// DependencyType is used to return a string representation of the ObjectType of a dependency
func DependencyType(d Dependency) string {
	switch d.(type) {
	case *Script:
		return "script"
	case *RemoteFile:
		return "remote_file"
	case *Command:
		return "command"
	case *DNSRecord:
		return "dns_record"
	case *Host:
		return "host"
	case *Network:
		return "network"
	case *Environment:
		return "environment"
	case *Build:
		return "build"
	case *Team:
		return "team"
	case *ProvisionedNetwork:
		return "provisioned_network"
	case *ProvisionedHost:
		return "provisioned_host"
	case *ProvisioningStep:
		return "provisioning_step"
	default:
		return "unknown"
	}
}

// DependencyCost returns a pre-specified cost for paths within the graph traversal
func DependencyCost(d Dependency) int64 {
	switch v := d.(type) {
	case *ProvisioningStep:
		return int64(20 + (v.StepNumber + 1))
	case *Script:
		return int64(999999)
	case *RemoteFile:
		return int64(999999)
	case *Command:
		return int64(999999)
	case *DNSRecord:
		return int64(999999)
	case *Host:
		return int64(0)
	case *Network:
		return int64(0)
	case *Environment:
		return int64(1)
	case *Build:
		return int64(2)
	case *Team:
		return int64(5)
	case *ProvisionedNetwork:
		return int64(10)
	case *ProvisionedHost:
		return int64(15)
	default:
		return int64(-55)
	}
}

// Relate takes a root dependency, and creates associations in the graph from it to the leafs.
func (s *Snapshot) Relate(root Dependency, leafs ...Dependency) error {
	deps := []Dependency{root}
	deps = append(deps, leafs...)
	s.AddNode(deps...)

	s.RelationQueue <- func() error {
		rootGID, ok := s.ObjectToGID[root]
		if !ok {
			return fmt.Errorf("could not find GID for root %s", root.Path())
		}
		for _, x := range leafs {
			leafGID, ok := s.ObjectToGID[x]
			if !ok {
				return fmt.Errorf("could not find GID for leaf %s", x.Path())
			}
			cost := DependencyCost(x)
			s.Graph.AddCost(rootGID, leafGID, cost)
		}
		return nil
	}
	return nil
}

// MapRelations enumerates the relation queue and actually mints the graph for the current objects
func (s *Snapshot) MapRelations() error {
	size := s.Counter
	s.Graph = graph.New(size)
	close(s.RelationQueue)

	for {
		relateFunc, valid := <-s.RelationQueue
		if !valid {
			break
		}
		err := relateFunc()
		if err != nil {
			return err
		}
	}
	s.Sorted = graph.Sort(s.Graph)
	gids, ok := graph.TopSort(s.Sorted)
	if !ok {
		return errors.New("dependency graph contains a cyclical dependency")
	}
	s.OrderedGIDs = gids
	return nil
}

// WalkChildren performs reverse traversal looking for all parent objects of an element.
func (s *Snapshot) WalkChildren() error {
	drev := graph.Transpose(s.Graph)
	gids, ok := graph.TopSort(drev)
	if !ok {
		return errors.New("dependency graph contained cyclical dependencies during inversion")
	}
	wg := new(sync.WaitGroup)
	retChan := make(chan map[Dependency]DependencyList, s.Counter)
	errChan := make(chan error, 1)

	for _, x := range gids {
		wg.Add(1)
		o := s.GIDToObject[x]
		go func(obj Dependency, i int) {
			defer wg.Done()
			newDepList := map[Dependency]DependencyList{
				obj: DependencyList{},
			}

			depfin := []int{}
			graph.BFS(drev, i, func(v, w int, c int64) {
				if s.Graph.Cost(w, i) > 0 {
					depfin = append(depfin, w)
				}
			})
			sort.Ints(depfin)
			newDepList[obj] = depfin
			retChan <- newDepList
		}(o, x)
	}

	go func() {
		wg.Wait()
		close(retChan)
	}()

	for {
		select {
		case resp, valid := <-retChan:
			if !valid {
				return nil
			}
			for k, depgroup := range resp {
				for _, dep := range depgroup {
					d := s.GIDToObject[dep]
					mdobj := s.Metadata[k.Path()]
					mdobj.ParentDeps = append(mdobj.ParentDeps, s.Metadata[d.Path()])
					mdobj.ParentIDs = append(mdobj.ParentIDs, d.Path())
					mdobj.ParentGIDs = append(mdobj.ParentGIDs, dep)
				}
			}
			continue
		case err := <-errChan:
			return err
		}
	}
}

// WalkParents performs breadth-first search semantics to elements in the graph, to map objects that depend on them.
func (s *Snapshot) WalkParents() error {
	drev := graph.Sort(s.Sorted)
	wg := new(sync.WaitGroup)
	retChan := make(chan map[Dependency]DependencyList, s.Counter)
	errChan := make(chan error, 1)

	for _, x := range s.OrderedGIDs {
		o := s.GIDToObject[x]
		wg.Add(1)
		go func(obj Dependency, i int) {
			defer wg.Done()
			newDepList := map[Dependency]DependencyList{
				obj: DependencyList{},
			}
			depmap := []int{}
			graph.BFS(drev, i, func(v, w int, c int64) {
				if s.Graph.Cost(i, w) > 0 {
					depmap = append(depmap, w)
				}
				// depmap = append(depmap, w)
			})
			sort.Ints(depmap)
			newDepList[obj] = depmap
			retChan <- newDepList
		}(o, x)
	}

	go func() {
		wg.Wait()
		close(retChan)
	}()

	for {
		select {
		case resp, valid := <-retChan:
			if !valid {
				return nil
			}
			for k, depgroup := range resp {
				for _, dep := range depgroup {
					d := s.GIDToObject[dep]
					mdobj := s.Metadata[k.Path()]
					mdobj.ChildDeps = append(mdobj.ChildDeps, s.Metadata[d.Path()])
					mdobj.ChildIDs = append(mdobj.ChildIDs, d.Path())
					mdobj.ChildGIDs = append(mdobj.ChildGIDs, dep)
				}
			}
			continue
		case err := <-errChan:
			return err
		}
	}
}

// Sort is the meta-function to kick off the mapping and creation of the snapshot's graph.
func (s *Snapshot) Sort() error {
	err := s.MapRelations()
	if err != nil {
		return err
	}
	wg := new(sync.WaitGroup)
	finChan := make(chan bool, 1)
	errChan := make(chan error, 1)
	wg.Add(2)
	go func() {
		defer wg.Done()
		e := s.WalkChildren()
		if e != nil {
			errChan <- e
		}
	}()
	go func() {
		defer wg.Done()
		e := s.WalkParents()
		if e != nil {
			errChan <- e
		}
	}()
	go func() {
		wg.Wait()
		close(finChan)
	}()
	select {
	case <-finChan:
		return nil
	case finErr := <-errChan:
		return finErr
	}
}
