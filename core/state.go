package core

import (
	"fmt"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
	"github.com/yourbasic/graph"
)

// State is the primary object used to interface with the build's on disk state table
type State struct {
	*sync.RWMutex
	DB      *buntdb.DB
	Tainted map[string]bool
}

// Snapshot is a graph traversal structure used to create a unique fingerprint for all elements in an environment.
type Snapshot struct {
	sync.RWMutex
	Counter       int
	Graph         *graph.Mutable
	Sorted        *graph.Immutable
	OrderedGIDs   []int
	Metadata      map[string]*Metadata
	Objects       map[string]Dependency
	GIDToObject   map[int]Dependency
	ObjectToGID   map[Dependency]int
	RelationQueue chan RelateFunc
}

// DependencyList is a helper type alias to handle topological sorting of laforge state elements
type DependencyList []int

// RelateFunc is a type alias to a function that relates objects together in a promise style format
type RelateFunc func() error

// NewEmptySnapshot returns an empty snapshot object
func NewEmptySnapshot() *Snapshot {
	return &Snapshot{
		Counter:       0,
		GIDToObject:   map[int]Dependency{},
		ObjectToGID:   map[Dependency]int{},
		Objects:       map[string]Dependency{},
		Metadata:      map[string]*Metadata{},
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
				ID:          dep.Path(),
				ObjectType:  DependencyType(dep),
				Dependency:  dep,
				ParentIDs:   []string{},
				ChildrenIDs: []string{},
			}
			s.Metadata[dep.Path()] = m
			go m.CalculateChecksum()
			s.Objects[dep.Path()] = dep
			s.GIDToObject[s.Counter] = dep
			s.ObjectToGID[dep] = s.Counter
			s.Counter++
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
					s.Metadata[k.Path()].ParentIDs = append(s.Metadata[k.Path()].ParentIDs, d.Path())

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
				depmap = append(depmap, w)
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
					s.Metadata[k.Path()].ChildrenIDs = append(s.Metadata[k.Path()].ChildrenIDs, d.Path())
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
