package graph

import (
	"fmt"
	"sync"
)

// Node a single node that composes the tree
type Node struct {
	value Relationship
}

// NodeQueue the queue of Nodes
type NodeQueue struct {
	items []Node
	lock  sync.RWMutex
}

// ItemGraph the Items graph
type ItemGraph struct {
	nodes  []*Node
	unique map[string]*Node
	edges  map[Node][]*Node
	lock   sync.RWMutex
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.value)
}

// ToNode converts a relationship to a node
func ToNode(rel Relationship) *Node {
	return &Node{
		value: rel,
	}
}

// AddNode adds a node to the graph
func (g *ItemGraph) AddNode(n *Node) {
	if _, found := g.unique[n.value.GetID()]; found {
		return
	}
	g.lock.Lock()
	g.nodes = append(g.nodes, n)
	g.unique[n.value.GetID()] = n
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *ItemGraph) AddEdge(n1, n2 *Node) {
	g.lock.Lock()
	if g.edges == nil {
		g.edges = make(map[Node][]*Node)
	}
	g.edges[*n1] = append(g.edges[*n1], n2)
	g.edges[*n2] = append(g.edges[*n2], n1)
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *ItemGraph) String() {
	g.lock.RLock()
	s := ""
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.lock.RUnlock()
}

// Traverse implements the BFS traversing algorithm
func (g *ItemGraph) Traverse(f func(*Node)) {
	g.lock.RLock()
	q := NodeQueue{}
	q.New()
	n := g.nodes[0]
	q.Enqueue(*n)
	visited := make(map[*Node]bool)
	for {
		if q.IsEmpty() {
			break
		}
		node := q.Dequeue()
		visited[node] = true
		near := g.edges[*node]

		for i := 0; i < len(near); i++ {
			j := near[i]
			if !visited[j] {
				q.Enqueue(*j)
				visited[j] = true
			}
		}
		if f != nil {
			f(node)
		}
	}
	g.lock.RUnlock()
}

// New creates a new NodeQueue
func (s *NodeQueue) New() *NodeQueue {
	s.lock.Lock()
	s.items = []Node{}
	s.lock.Unlock()
	return s
}

// Enqueue adds an Node to the end of the queue
func (s *NodeQueue) Enqueue(t Node) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes an Node from the start of the queue
func (s *NodeQueue) Dequeue() *Node {
	s.lock.Lock()
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	s.lock.Unlock()
	return &item
}

// Front returns the item next in the queue, without removing it
func (s *NodeQueue) Front() *Node {
	s.lock.RLock()
	item := s.items[0]
	s.lock.RUnlock()
	return &item
}

// IsEmpty returns true if the queue is empty
func (s *NodeQueue) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items) == 0
}

// Size returns the number of Nodes in the queue
func (s *NodeQueue) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}
