package core

import "github.com/gen0cide/laforge/core/graph"

// Provisioner is a meta interface to provide provisioning steps to the Builder
type Provisioner interface {
	graph.Hasher
	Pather

	Gather(s *Snapshot) error
	ParentLaforgeID() string

	// Kind denotes the type of Provisioner this is
	Kind() string
}
