package core

import (
	"fmt"
	"io"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/hashicorp/terraform/dag"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/core/graph"
)

// Snapshot is a graph traversal structure used to create a unique fingerprint for all elements in an environment.
//easyjson:json
type Snapshot struct {
	Checksum  uint64               `json:"checksum"`
	Counter   int                  `json:"counter"`
	Graph     *dag.AcyclicGraph    `json:"graph"`
	Metastore map[string]*Metadata `json:"metadata"`
	RootID    string               `json:"root_id"`
	Edges     map[string]bool      `json:"edges"`
}

// DependencyList is a helper type alias to handle topological sorting of laforge state elements
type DependencyList []int

// RelateFunc is a type alias to a function that relates objects together in a promise style format
type RelateFunc func() error

// Plot writes a DOT representation of the Snapshot to the output writer, with an optional path to denote the root object.
func (s *Snapshot) Plot(output io.Writer, rootID string) error {
	if rootID != "" {
		_, ok := s.Metastore[rootID]
		if !ok {
			return fmt.Errorf("the node %s could not be found in the graph", rootID)
		}
	}
	dw := graph.NewDotWriter(output, len(s.Metastore), false)
	rootMeta := s.Metastore[rootID]
	dw.PlotGraph(rootMeta)
	return nil
}

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
	return &Snapshot{
		Counter:   0,
		Graph:     &dag.AcyclicGraph{},
		Metastore: map[string]*Metadata{},
		Edges:     map[string]bool{},
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
		s.Edges[strings.Join([]string{s.RootID, m.ID}, "|")] = true
	}
	return m
}

// Connect is the primary interface for associating metadata in the snapshot's graph
func (s *Snapshot) Connect(src, target *Metadata) {
	oldedge := dag.BasicEdge(s.RootID, target.ID)
	if s.Graph.HasEdge(oldedge) {
		delete(s.Edges, strings.Join([]string{s.RootID, target.ID}, "|"))
		s.Graph.RemoveEdge(oldedge)
	}
	newedge := dag.BasicEdge(src.ID, target.ID)
	s.Graph.Connect(newedge)
	s.Edges[strings.Join([]string{src.ID, target.ID}, "|")] = true
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
		s.Graph.Add(val.ID)
	}

	for x := range s.Edges {
		parts := strings.Split(x, "|")
		s.Graph.Connect(dag.BasicEdge(parts[0], parts[1]))
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
	s.AddNode(root)
	for _, l := range leafs {
		s.AddNode(l)
	}
	return nil
}
