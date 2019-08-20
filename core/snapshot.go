package core

import (
	"sync"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/hashicorp/terraform/dag"

	mapset "github.com/deckarep/golang-set"
)

// Snapshot is a graph traversal structure used to create a unique fingerprint for all elements in an environment.
// easyjson:json
type Snapshot struct {
	Checksum  uint64               `json:"checksum"`
	AltGraph  *dag.AcyclicGraph    `json:"altgraph"`
	Metastore map[string]*Metadata `json:"metadata"`
	RootID    string               `json:"root_id"`
	Edges     interface{}          `json:"edges"`
	Metabus   chan *Metadata       `json:"-"`
	Nodebus   chan *Metadata       `json:"-"`
	Edgebus   chan Edge            `json:"-"`
	Mutex     *sync.RWMutex        `json:"-"`
}

// Edge is a type to store relationship information in the graph
// easyjson:json
type Edge struct {
	Source string `json:"source,omitempty"`
	Target string `json:"target,omitempty"`
}

// GetEdges returns a type asserted set of edges
func (s *Snapshot) GetEdges() mapset.Set {
	return s.Edges.(mapset.Set)
}

// DependencyList is a helper type alias to handle topological sorting of laforge state elements
type DependencyList []int

// RelateFunc is a type alias to a function that relates objects together in a promise style format
type RelateFunc func() error

// Hash implements the hasher interface
func (s *Snapshot) Hash() uint64 {
	hashes := ChecksumList{}
	for _, x := range s.Metastore {
		hashes = append(hashes, x.Hash())
	}
	return hashes.Hash()
}

// NewEmptySnapshot returns an empty snapshot object
func NewEmptySnapshot() *Snapshot {
	altDag := &dag.AcyclicGraph{}
	altDag.Add("root")
	return &Snapshot{
		Mutex:     &sync.RWMutex{},
		Metastore: map[string]*Metadata{},
		Edges:     mapset.NewSet(),
		AltGraph:  altDag,
		Metabus:   make(chan *Metadata, 100000),
		Nodebus:   make(chan *Metadata, 100000),
		Edgebus:   make(chan Edge, 100000),
	}
}

// NewSnapshotFromEnv creates a new snapshot from a provided environment
func NewSnapshotFromEnv(e *Environment, overwriteBuild bool) (*Snapshot, error) {
	s := NewEmptySnapshot()
	var build *Build
	if overwriteBuild || e.Build == nil {
		build = e.CreateBuild()
		e.Build = build
	} else {
		build = e.Build
	}
	err := build.CreateTeams()
	if err != nil {
		return nil, err
	}
	pgchan := make(chan struct{}, 1)
	smchan := make(chan struct{}, 1)
	rochan := make(chan struct{}, 1)
	finchan := make(chan struct{}, 3)
	wg := new(sync.WaitGroup)

	go s.PopulateGraph(pgchan, finchan)
	go s.StoreMetadata(smchan, finchan)
	go s.RelateObjects(rochan, finchan)

	wg.Add(1)
	go s.WalkEnvironment(e, wg)
	for _, t := range build.Teams {
		wg.Add(1)
		go s.WalkTeam(t, wg)
	}

	wg.Wait()

	smchan <- struct{}{}
	<-finchan
	pgchan <- struct{}{}
	<-finchan
	rochan <- struct{}{}
	<-finchan

	s.AltGraph.Remove("root")
	s.AltGraph.TransitiveReduction()

	return s, nil
}

// AddNodeV2 places items on the graph
func (s *Snapshot) AddNodeV2(id string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.AltGraph.Add(id)
	s.AltGraph.Connect(dag.BasicEdge("root", id))
}

// RelateV2 is what actually snips and splices edges into the graph
func (s *Snapshot) RelateV2(e Edge) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if !s.AltGraph.HasVertex(e.Source) {
		s.Edgebus <- e
		return
	}
	if !s.AltGraph.HasVertex(e.Target) {
		s.Edgebus <- e
		return
	}

	altoldedge := dag.BasicEdge("root", e.Source)
	if s.AltGraph.HasEdge(altoldedge) {
		s.AltGraph.RemoveEdge(altoldedge)
		// cli.Logger.Errorf("Edge root -> %s exists. Removing and connecting root -> %s", e.Source, e.Target)
		s.AltGraph.Connect(dag.BasicEdge(e.Target, "root"))
	}
	altnewedge := dag.BasicEdge(e.Source, e.Target)
	s.AltGraph.Connect(altnewedge)
	s.GetEdges().Add(e)
	return
}

// RelateObjects monitors the edgebus channel for objects that need edges defined in the graph
func (s *Snapshot) RelateObjects(end chan struct{}, fin chan struct{}) {
	ttg := false
	for {
		select {
		case e := <-s.Edgebus:
			s.RelateV2(e)
			continue
		case <-end:
			cli.Logger.Debugf("Relate objects worker termination triggered")
			ttg = true
			continue
		default:
		}
		select {
		case e := <-s.Edgebus:
			s.RelateV2(e)
			continue
		default:
			if !ttg {
				continue
			}
			if len(s.Edgebus) > 0 {
				continue
			}
			cli.Logger.Debugf("Edgebus drained")
			fin <- struct{}{}
			return
		}
	}
}

// StoreMetadata creates entries in the snapshot's Metastore
func (s *Snapshot) StoreMetadata(end chan struct{}, fin chan struct{}) {
	ttg := false
	for {
		select {
		case m := <-s.Metabus:
			s.Mutex.Lock()
			s.Metastore[m.ID] = m
			s.Mutex.Unlock()
			s.Nodebus <- m
			continue
		case <-end:
			cli.Logger.Debugf("store metadata worker termination triggered")
			ttg = true
			continue
		default:
		}
		select {
		case m := <-s.Metabus:
			s.Mutex.Lock()
			s.Metastore[m.ID] = m
			s.Mutex.Unlock()
			s.Nodebus <- m
			continue
		default:
			if !ttg {
				continue
			}
			if len(s.Metabus) > 0 {
				continue
			}
			cli.Logger.Debugf("Metabus drained")
			fin <- struct{}{}
			return
		}
	}
}

// PopulateGraph places items on the graph with a default root value
func (s *Snapshot) PopulateGraph(end chan struct{}, fin chan struct{}) {
	ttg := false
	for {
		select {
		case n := <-s.Nodebus:
			s.AddNodeV2(n.ID)
			continue
		case <-end:
			cli.Logger.Debugf("populate graph worker termination triggered")
			ttg = true
			continue
		default:
		}
		select {
		case n := <-s.Nodebus:
			s.AddNodeV2(n.ID)
			continue
		default:
			if !ttg {
				continue
			}
			if len(s.Nodebus) > 0 {
				continue
			}
			cli.Logger.Debugf("Nodebus drained")
			fin <- struct{}{}
			return
		}
	}
}

// AddObject adds a dependency to the Metastore
func (s *Snapshot) AddObject(dep Dependency) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if _, ok := s.Metastore[dep.Path()]; ok {
		return
	}

	m := &Metadata{
		ID:         dep.Path(),
		ObjectType: TypeByPath(dep.Path()),
		Dependency: dep,
		Checksum:   dep.Hash(),
	}

	s.Metabus <- m
}

// CreateEdge returns an edge creation data structure
func CreateEdge(src Dependency, target Dependency) Edge {
	return Edge{
		Source: src.Path(),
		Target: target.Path(),
	}
}

// AddRelationship places a relationship on the Edgebus
func (s *Snapshot) AddRelationship(src Dependency, target Dependency) {
	s.Edgebus <- CreateEdge(src, target)
}

// WalkEnvironment walks the environment looking for dependencies to graph
func (s *Snapshot) WalkEnvironment(e *Environment, wg *sync.WaitGroup) {
	defer wg.Done()
	s.AddObject(e)
	s.AddObject(e.Build)
	s.AddRelationship(e, e.Build)
	for _, net := range e.IncludedNetworks {
		s.AddObject(net)
		s.AddRelationship(e.Build, net)
		for _, host := range e.HostByNetwork[net.Path()] {
			s.AddObject(host)
			s.AddRelationship(net, host)
			wg.Add(1)
			go s.WalkHost(host, wg)
		}
	}
}

// WalkHost is used to identify global host dependencies
func (s *Snapshot) WalkHost(h *Host, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, x := range h.Scripts {
		s.AddObject(x)
		s.AddRelationship(h, x)
	}
	for _, x := range h.DNSRecords {
		s.AddObject(x)
		s.AddRelationship(h, x)
	}
	for _, x := range h.Commands {
		s.AddObject(x)
		s.AddRelationship(h, x)
	}
	for _, x := range h.RemoteFiles {
		s.AddObject(x)
		s.AddRelationship(h, x)
	}
}

// WalkTeam is used to enumerate the resources of a team
func (s *Snapshot) WalkTeam(t *Team, wg *sync.WaitGroup) {
	defer wg.Done()
	s.AddObject(t)
	s.AddRelationship(t.Build, t)
	for _, pn := range t.ProvisionedNetworks {
		s.AddObject(pn)
		s.AddRelationship(t, pn)
		s.AddRelationship(pn, pn.Network)
		for _, ph := range pn.ProvisionedHosts {
			wg.Add(1)
			s.AddObject(ph)
			s.AddRelationship(pn, ph)
			go s.WalkProvisionedHost(ph, wg)
		}
	}
}

// WalkProvisionedHost is used to walk all the elements of a provisioned host
func (s *Snapshot) WalkProvisionedHost(ph *ProvisionedHost, wg *sync.WaitGroup) {
	defer wg.Done()
	s.AddObject(ph.Conn)
	s.AddRelationship(ph, ph.Conn)
	// s.AddRelationship(ph.Host, ph)
	s.AddRelationship(ph, ph.Host)
	for psidx, ps := range ph.StepsByOffset {
		s.AddObject(ps)
		switch v := ps.Provisioner.(type) {
		case *Command:
			s.AddObject(v)
			s.AddRelationship(ph, v)
			s.AddRelationship(ps, v)
		case *DNSRecord:
			s.AddObject(v)
			s.AddRelationship(ph, v)
			s.AddRelationship(ps, v)
		case *RemoteFile:
			s.AddObject(v)
			s.AddRelationship(ph, v)
			s.AddRelationship(ps, v)
		case *Script:
			s.AddObject(v)
			s.AddRelationship(ph, v)
			s.AddRelationship(ps, v)
		}
		if psidx == 0 {
			s.AddRelationship(ph.Conn, ps)
			continue
		}
		s.AddRelationship(ph.Conn, ps)
		prevstep := ph.StepsByOffset[psidx-1]
		s.AddRelationship(prevstep, ps)
	}
	for _, dep := range ph.Host.Dependencies {
		dh, err := ph.Team.LocateProvisionedHost(dep.NetworkID, dep.HostID)
		if err != nil {
			panic(err)
		}

		fsid := dh.Host.FinalStepID()
		if fsid != -1 {
			fs := dh.StepsByOffset[fsid]
			s.AddRelationship(fs, ph)
			continue
		}
		s.AddRelationship(dh, ph)
		s.AddRelationship(dh.Conn, ph)
	}
}

// RebuildGraph will attempt to recreate the dependency tree based on it's persisted metadata instead of building from environment.
func (s *Snapshot) RebuildGraph() error {

	s.AltGraph = &dag.AcyclicGraph{}
	for _, val := range s.Metastore {
		s.AltGraph.Add(val.ID)
	}

	for x := range s.GetEdges().Iter() {
		edge := x.(Edge)
		s.AltGraph.Connect(dag.BasicEdge(edge.Source, edge.Target))
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
// DEPRECATED
func (s *Snapshot) Relate(root Dependency, leafs ...Dependency) error {
	return nil
}
