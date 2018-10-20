package core

import (
	"fmt"
	"io"

	"github.com/cespare/xxhash"
	"github.com/davecgh/go-spew/spew"

	"github.com/hashicorp/terraform/dag"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/core/graph"
)

// Snapshot is a graph traversal structure used to create a unique fingerprint for all elements in an environment.
type Snapshot struct {
	Checksum uint64            `json:"checksum"`
	Counter  int               `json:"counter"`
	Graph    *dag.AcyclicGraph `json:"graph"`
	// Graph         *graphx.Mutable               `json:"-"`
	// Sorted        *graphx.Immutable             `json:"-"`
	// OrderedGIDs   []int                         `json:"ordered_gids"`
	Objects   map[string]graph.Relationship `json:"-"`
	Metastore map[string]*Metadata          `json:"metadata"`
	RootID    string                        `json:"root_id"`
	Edges     []*Edge                       `json:"edge"`
	// GIDToPath     map[int]string                `json:"gid_to_path"`
	// PathToGID     map[string]int                `json:"path_to_gid"`
	// GIDToObject   map[int]graph.Relationship    `json:"-"`
	// ObjectToGID   map[graph.Relationship]int    `json:"-"`
	RelationQueue chan RelateFunc `json:"-"`
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
	dw := graph.NewDotWriter(output, len(s.Objects), false)
	rootMeta := s.Objects[rootID]
	dw.PlotGraph(rootMeta)
	return nil
}

type Edge struct {
	SourceID  string     `json:"source_id"`
	TargetID  string     `json:"target_id"`
	SourceV   dag.Vertex `json:"-"`
	TargetV   dag.Vertex `json:"-"`
	SourceObj *Metadata  `json:"-"`
	TargetObj *Metadata  `json:"-"`
}

func (e *Edge) Hashcode() interface{} {
	return xxhash.Sum64String(fmt.Sprintf("%x.%x", e.SourceObj.Checksum, e.TargetObj.Checksum))
}

func (e *Edge) Source() dag.Vertex {
	return e.SourceObj
}

func (e *Edge) Target() dag.Vertex {
	return e.TargetObj
}

func (e *Edge) DotNode(s string, d *dag.DotOpts) *dag.DotNode {
	return &dag.DotNode{
		Name: s,
		Attrs: map[string]string{
			"target_id:": e.TargetID,
			"source_id":  e.SourceID,
		},
	}
}

// Hash implements the hasher interface
func (s *Snapshot) Hash() uint64 {
	hashes := ChecksumList{}
	for _, x := range s.Metastore {
		hashes = append(hashes, x.Hash())
	}
	for _, x := range s.Edges {
		hashes = append(hashes, x.SourceObj.Hash())
		hashes = append(hashes, x.TargetObj.Hash())
	}
	return hashes.Hash()
}

// NewEmptySnapshot returns an empty snapshot object
func NewEmptySnapshot() *Snapshot {
	return &Snapshot{
		Counter: 0,
		Graph:   &dag.AcyclicGraph{},
		// GIDToObject:   map[int]graph.Relationship{},
		Metastore: map[string]*Metadata{},
		Edges:     []*Edge{},
		// ObjectToGID:   map[graph.Relationship]int{},
		// Objects: map[string]graph.Relationship{},
		// GIDToPath:     map[int]string{},
		// PathToGID:     map[string]int{},
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
func (s *Snapshot) AddNode(dep Dependency) *Metadata {
	if mdep, ok := s.Metastore[dep.Path()]; ok {
		return mdep
	}
	// if taken, bad := s.GIDToObject[s.Counter]; bad {
	// 	return fmt.Errorf("unable to add %s to graph: GID %d is already taken by %s", dep.Path(), s.Counter, taken.GetID())
	// }
	m := &Metadata{
		Name:       dep.Path(),
		ID:         dep.Path(),
		GID:        s.Counter,
		ObjectType: DependencyType(dep),
		Dependency: dep,
		Checksum:   dep.Hash(),
	}

	if len(s.Metastore) == 0 {
		s.RootID = m.ID
	}

	s.Metastore[dep.Path()] = m
	s.Graph.Add(m.ID)

	if s.RootID != m.ID {
		s.Graph.Connect(dag.BasicEdge(s.RootID, m.ID))
	}
	return m
}

func (s *Snapshot) Connect(src, target *Metadata) {
	oldedge := dag.BasicEdge(s.RootID, target.ID)
	if s.Graph.HasEdge(oldedge) {
		s.Graph.RemoveEdge(oldedge)
	}
	newedge := dag.BasicEdge(src.ID, target.ID)
	s.Graph.Connect(newedge)
	if err := s.Graph.Validate(); err != nil {
		cli.Logger.Errorf("FAILED TO CONNECT %s => %s", src.ID, target.ID)
		spew.Dump(src)
		spew.Dump(target)
		panic(err)
	}
}

// RebuildGraph will attempt to recreate the dependency tree based on it's persisted metadata instead of building from environment.
func (s *Snapshot) RebuildGraph() error {

	s.Graph = &dag.AcyclicGraph{}
	for _, val := range s.Metastore {
		s.Graph.Add(val)
	}

	for _, x := range s.Edges {
		x.SourceObj = s.Metastore[x.SourceID]
		x.TargetObj = s.Metastore[x.TargetID]
		x.SourceV = s.Graph.Add(x.SourceObj)
		x.TargetV = s.Graph.Add(x.TargetObj)
		s.Graph.Connect(x)
	}
	// if s.Graph != nil {
	// 	return errors.New("snapshot already has built dependency graph")
	// }
	// s.Graph = graphx.New(len(s.Metastore))
	// wg := new(sync.WaitGroup)
	// echan := make(chan edge, len(s.Metastore)*100)
	// finChan := make(chan bool, 1)

	// for key, mref := range s.Metastore {
	// 	wg.Add(1)
	// 	s.Objects[key] = mref
	// 	go func(mobj *Metadata) {
	// 		defer wg.Done()
	// 		mobj.ParentDeps = make([]graph.Relationship, len(mobj.ParentGIDs))
	// 		for offset := 0; offset < len(mobj.ParentGIDs); offset++ {
	// 			ppath := mobj.ParentDepIDs[offset]
	// 			pgid := mobj.ParentGIDs[offset]
	// 			par, ok := s.Metastore[ppath]
	// 			if !ok {
	// 				cli.Logger.Errorf("could not find metastore object %s", ppath)
	// 				continue
	// 			}
	// 			mobj.ParentDeps[offset] = par
	// 			echan <- edge{v: pgid, w: mobj.GetGID(), c: mobj.GCost}
	// 		}
	// 		mobj.ChildDeps = make([]graph.Relationship, len(mobj.ChildGIDs))
	// 		for offset := 0; offset < len(mobj.ChildGIDs); offset++ {
	// 			cpath := mobj.ChildDepIDs[offset]
	// 			cgid := mobj.ChildGIDs[offset]
	// 			chi, ok := s.Metastore[cpath]
	// 			if !ok {
	// 				cli.Logger.Errorf("could not find metastore object %s as child of %s", cpath, mobj.GetID())
	// 				continue
	// 			}
	// 			mobj.ChildDeps[offset] = chi
	// 			echan <- edge{v: mobj.GetGID(), w: cgid, c: mobj.GCost}
	// 		}
	// 	}(mref)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(finChan)
	// }()

	// for {
	// 	select {
	// 	case e := <-echan:
	// 		s.Graph.AddCost(e.v, e.w, e.c)
	// 		continue
	// 	case <-finChan:
	// 		return nil
	// 	}
	// }
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
	// deps := []Dependency{root}
	// rootnode := s.AddNode(root)
	// source := s.Graph.Add(rootnode)

	// rootmeta := s.Metastore[root.Path()]

	s.AddNode(root)

	for _, l := range leafs {
		leafnode := s.AddNode(l)
		_ = leafnode

		// vert := s.Graph.Add(leafnode)
		// edge := &Edge{
		// 	SourceV:   source,
		// 	TargetV:   vert,
		// 	SourceObj: rootnode,
		// 	TargetObj: leafnode,
		// 	SourceID:  rootnode.ID,
		// 	TargetID:  leafnode.ID,
		// }
		// if leafnode.ID == "/envs/lfdev/tfgcp/teams/0/networks/corp/hosts/ad-00/steps/1-set-bad-local-secpol" {
		// 	pp.Println(leafnode)
		// 	pp.Println(vert)
		// }
		// s.Edges = append(s.Edges, edge)
		// s.Graph.Connect(dag.BasicEdge(source, vert))
	}
	// s.RelationQueue <- func() error {
	// 	rootrel, ok := s.Objects[root.Path()]
	// 	if !ok {
	// 		return fmt.Errorf("could not find GID for root %s", root.Path())
	// 	}

	// 	for _, x := range leafs {
	// 		leafrel, ok := s.Objects[x.Path()]
	// 		if !ok {
	// 			return fmt.Errorf("could not find GID for leaf %s", x.Path())
	// 		}
	// 		cost := DependencyCost(x)
	// 		s.Graph.AddCost(rootrel.GetGID(), leafrel.GetGID(), cost)
	// 	}
	// 	return nil
	// }
	return nil
}

// MapRelations enumerates the relation queue and actually mints the graph for the current objects
func (s *Snapshot) MapRelations() error {
	// size := s.Counter
	// s.Graph = graphx.New(size)
	// close(s.RelationQueue)

	// for {
	// 	relateFunc, valid := <-s.RelationQueue
	// 	if !valid {
	// 		break
	// 	}
	// 	err := relateFunc()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// s.Sorted = graphx.Sort(s.Graph)
	// gids, ok := graphx.TopSort(s.Sorted)
	// if !ok {
	// 	return errors.New("dependency graph contains a cyclical dependency")
	// }
	// s.OrderedGIDs = gids
	return nil
}

// WalkChildren performs reverse traversal looking for all parent objects of an element.
func (s *Snapshot) WalkChildren() error {
	return nil
	// drev := graphx.Transpose(s.Graph)
	// gids, ok := graphx.TopSort(drev)
	// if !ok {
	// 	return errors.New("dependency graph contained cyclical dependencies during inversion")
	// }
	// wg := new(sync.WaitGroup)
	// retChan := make(chan map[graph.Relationship]DependencyList, s.Counter)
	// errChan := make(chan error, 1)

	// for _, x := range gids {
	// 	wg.Add(1)
	// 	o := s.GIDToObject[x]
	// 	go func(obj graph.Relationship, i int) {
	// 		defer wg.Done()
	// 		newDepList := map[graph.Relationship]DependencyList{
	// 			obj: DependencyList{},
	// 		}

	// 		depfin := []int{}
	// 		graphx.BFS(drev, i, func(v, w int, c int64) {
	// 			if s.Graph.Cost(w, i) > 0 {
	// 				depfin = append(depfin, w)
	// 			}
	// 		})
	// 		sort.Ints(depfin)
	// 		newDepList[obj] = depfin
	// 		retChan <- newDepList
	// 	}(o, x)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(retChan)
	// }()

	// for {
	// 	select {
	// 	case resp, valid := <-retChan:
	// 		if !valid {
	// 			return nil
	// 		}
	// 		for k, depgroup := range resp {
	// 			for _, dep := range depgroup {
	// 				d, ok := s.GIDToObject[dep]
	// 				if !ok {
	// 					cli.Logger.Errorf("could find object ID %d as a dependency for %s", dep, k.GetID())
	// 					continue
	// 				}
	// 				k.AddParent(d)
	// 			}
	// 		}
	// 		continue
	// 	case err := <-errChan:
	// 		return err
	// 	}
	// }
}

// WalkParents performs breadth-first search semantics to elements in the graph, to map objects that depend on them.
func (s *Snapshot) WalkParents() error {
	return nil
	// drev := graphx.Sort(s.Sorted)
	// wg := new(sync.WaitGroup)
	// retChan := make(chan map[graph.Relationship]DependencyList, s.Counter)
	// errChan := make(chan error, 1)

	// for _, x := range s.OrderedGIDs {
	// 	o := s.GIDToObject[x]
	// 	wg.Add(1)
	// 	go func(obj graph.Relationship, i int) {
	// 		defer wg.Done()
	// 		newDepList := map[graph.Relationship]DependencyList{
	// 			obj: DependencyList{},
	// 		}
	// 		depmap := []int{}
	// 		graphx.BFS(drev, i, func(v, w int, c int64) {
	// 			if s.Graph.Cost(i, w) > 0 {
	// 				depmap = append(depmap, w)
	// 			}
	// 			// depmap = append(depmap, w)
	// 		})
	// 		sort.Ints(depmap)
	// 		newDepList[obj] = depmap
	// 		retChan <- newDepList
	// 	}(o, x)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(retChan)
	// }()

	// for {
	// 	select {
	// 	case resp, valid := <-retChan:
	// 		if !valid {
	// 			return nil
	// 		}
	// 		for k, depgroup := range resp {
	// 			for _, dep := range depgroup {
	// 				realobj, ok := s.GIDToObject[dep]
	// 				if !ok {
	// 					cli.Logger.Errorf("could not find object for %d in snapshot", dep)
	// 					continue
	// 				}
	// 				k.AddChild(realobj)
	// 			}
	// 		}
	// 		continue
	// 	case err := <-errChan:
	// 		return err
	// 	}
	// }
}

// Sort is the meta-function to kick off the mapping and creation of the snapshot's graph.
func (s *Snapshot) Sort() error {
	// err := s.MapRelations()
	// if err != nil {
	// 	return err
	// }
	// wg := new(sync.WaitGroup)
	// finChan := make(chan bool, 1)
	// errChan := make(chan error, 1)
	// wg.Add(2)
	// go func() {
	// 	defer wg.Done()
	// 	e := s.WalkChildren()
	// 	if e != nil {
	// 		errChan <- e
	// 	}
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	e := s.WalkParents()
	// 	if e != nil {
	// 		errChan <- e
	// 	}
	// }()
	// go func() {
	// 	wg.Wait()
	// 	close(finChan)
	// }()
	// select {
	// case <-finChan:
	// 	return nil
	// case finErr := <-errChan:
	// 	return finErr
	// }
	return nil
}
