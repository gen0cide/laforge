package graph

import (
	"github.com/gen0cide/laforge/core/cli"
	"github.com/pkg/errors"
)

// Relationship is an interface to allow core objects to build relationships between them
type Relationship interface {
	Hasher
	DotNode
	GetID() string
	Children() []Relationship
	Parents() []Relationship
	ParentIDs() []string
	ChildrenIDs() []string
	AddChild(r ...Relationship)
	AddParent(r ...Relationship)
}

// AssociateChildren is a generic function to associate child dependencies with an object
func AssociateChildren(subject Relationship, children ...Relationship) {
	subject.AddChild(children...)
	return
}

// AssociateParents is a generic function to associate parent dependencies with an object
func AssociateParents(subject Relationship, parents ...Relationship) {
	subject.AddParent(parents...)
	return
}

// RelationshipWalkFunc allows for recursive relation traversal
type RelationshipWalkFunc func(rel Relationship, distance int) error

// WalkDirection allows the directon of the walk to be specified
type WalkDirection int

const (
	// InfiniteDepth is used as a flag to instruct the walker to traverse infinite depth
	InfiniteDepth int = -1

	// TraverseChildren is a flag to denote the Relationship should walk children
	TraverseChildren WalkDirection = iota

	// TraverseParents is a flag to denote the Relationship should walk parents
	TraverseParents
)

var (
	// ErrTraversalPathEnded is thrown when the walker has hit the final element in the tree
	ErrTraversalPathEnded = errors.New("traversal path has hit the final element")

	// ErrRelationshipWalkHalted is thrown when a walk is intended to discontinue walking
	ErrRelationshipWalkHalted = errors.New("relationship walk has reached a termination point")

	// ErrWalkDirectionUndefined is thrown when a direction is neither TraverseChildren or TraverseParents
	ErrWalkDirectionUndefined = errors.New("a walk direction of an unknown type has been given")
)

// HasChild determines if a direct child relationship exists between host (parent) and question (supposed child)
func HasChild(host Relationship, question Relationship) bool {
	id := question.GetID()
	for _, x := range host.ChildrenIDs() {
		if x == id {
			return true
		}
	}
	return false
}

// HasParent determines if a direct parent relationship exists between host (child) and question (supposed parent)
func HasParent(host Relationship, question Relationship) bool {
	id := question.GetID()
	for _, x := range host.ChildrenIDs() {
		if x == id {
			return true
		}
	}
	return false
}

// HasIndirectChild determines if a transitive parent relationship exists between host (parent) and question (supposed child)
func HasIndirectChild(host Relationship, question Relationship) bool {
	id := question.GetID()
	found := false
	err := WalkRelationship(host, InfiniteDepth, 0, TraverseChildren, func(rel Relationship, distance int) error {
		if rel.GetID() == id {
			found = true
			return ErrRelationshipWalkHalted
		}
		return nil
	})

	if err == nil || err == ErrRelationshipWalkHalted {
		return found
	}

	cli.Logger.Errorf("Error thrown walking %s parents looking for %s: %v", host.GetID(), id, err)
	return false
}

// HasIndirectParent determines if a transitive parent relationship exists between host (child) and question (supposed parent)
func HasIndirectParent(host Relationship, question Relationship) bool {
	id := question.GetID()
	found := false
	err := WalkRelationship(host, InfiniteDepth, 0, TraverseParents, func(rel Relationship, distance int) error {
		if rel.GetID() == id {
			found = true
			return ErrRelationshipWalkHalted
		}
		return nil
	})

	if err == nil || err == ErrRelationshipWalkHalted {
		return found
	}

	cli.Logger.Errorf("Error thrown walking %s parents looking for %s: %v", host.GetID(), id, err)
	return false
}

// WalkRelationship is a recursive walk similar to filepath.Walk except it traverses relationships
func WalkRelationship(base Relationship, maxdepth int, offset int, direction WalkDirection, walkFunc RelationshipWalkFunc) error {
	if maxdepth != InfiniteDepth && maxdepth < 1 {
		return ErrTraversalPathEnded
	}

	var subjects []Relationship
	switch direction {
	case TraverseChildren:
		subjects = base.Children()
	case TraverseParents:
		subjects = base.Parents()
	default:
		return ErrWalkDirectionUndefined
	}

	if len(subjects) == 0 {
		return ErrTraversalPathEnded
	}

	newdepth := InfiniteDepth
	if maxdepth != InfiniteDepth {
		newdepth = maxdepth - 1
	}

	for _, x := range subjects {
		if err := walkFunc(x, offset); err != nil {
			return err
		}
		if err := WalkRelationship(x, newdepth, offset+1, direction, walkFunc); err != nil {
			if err == ErrTraversalPathEnded {
				continue
			}
			return err
		}
	}

	return nil
}
