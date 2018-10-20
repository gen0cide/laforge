package graph

// Hasher is an interface to allow types to be checksumed for potentially build breaking changes
type Hasher interface {
	Hash() uint64
}
