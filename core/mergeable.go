package core

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

var (
	// ErrSwapTypeMismatch is thrown when a Swap() function on the Mergeable interface cannot find it's native type
	ErrSwapTypeMismatch = errors.New("mergeable swap() call failed due to delta type mismatch")
)

// Mergeable is an interface that allows for dynamic types to be merged with coordinated strategies
type Mergeable interface {
	// GetCaller retrieves the object's Caller value
	GetCaller() Caller

	// GetID retrieves the object's ID value
	GetID() string

	// GetOnConflict retrieves the object's OnConflict value
	GetOnConflict() OnConflict

	// SetCaller sets the object's Caller to a new version
	SetCaller(ca Caller)

	// SetOnConflict replaces the object's OnConflict with a new version
	SetOnConflict(oc OnConflict)

	// Swap attempts to replace the pointer references of two mergeable objects
	Swap(ma Mergeable) error
}

// SmartMerge takes a source object (m) and a delta (diff) and attempts to merge them using settings based on the delta's OnConflict configuration.
func SmartMerge(m, diff Mergeable, appendSlices bool) (Mergeable, error) {
	strats := []func(*mergo.Config){mergo.WithOverride}
	newCaller := m.GetCaller().Stack(diff.GetCaller())
	conflict := m.GetOnConflict()
	meth := diff.GetOnConflict().Do
	switch {
	case meth == "" || meth == "default":
		if diff.GetOnConflict().Append {
			strats = append(strats, mergo.WithAppendSlice)
		}
		err := mergo.Merge(m, diff, strats...)
		m.SetCaller(newCaller)
		if err != nil {
			return m, errors.WithStack(err)
		}
		return m, nil
	case meth == "overwrite":
		swapErr := m.Swap(diff)
		m.SetOnConflict(conflict)
		return m, swapErr
	case meth == "inherit":
		if m.GetOnConflict().Append {
			strats = append(strats, mergo.WithAppendSlice)
		}
		err := mergo.Merge(diff, m, strats...)
		swapErr := m.Swap(diff)
		m.SetCaller(newCaller)
		m.SetOnConflict(conflict)
		if err != nil {
			return m, errors.WithStack(err)
		}
		if swapErr != nil {
			return m, swapErr
		}
		return m, nil
	case meth == "skip":
		return m, nil
	case meth == "panic":
		return m, NewMergeConflict(m, diff, m.GetID(), diff.GetID(), m.GetCaller().Current(), diff.GetCaller().Current())
	default:
		return m, fmt.Errorf("invalid conflict strategy %s in %s", diff.GetOnConflict().Do, diff.GetCaller().Current().CallerFile)
	}
}
