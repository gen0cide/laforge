// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/gen0cide/laforge/ent/adhocplan"
	"github.com/gen0cide/laforge/ent/predicate"
)

// AdhocPlanDelete is the builder for deleting a AdhocPlan entity.
type AdhocPlanDelete struct {
	config
	hooks    []Hook
	mutation *AdhocPlanMutation
}

// Where adds a new predicate to the AdhocPlanDelete builder.
func (apd *AdhocPlanDelete) Where(ps ...predicate.AdhocPlan) *AdhocPlanDelete {
	apd.mutation.predicates = append(apd.mutation.predicates, ps...)
	return apd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (apd *AdhocPlanDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(apd.hooks) == 0 {
		affected, err = apd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AdhocPlanMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			apd.mutation = mutation
			affected, err = apd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(apd.hooks) - 1; i >= 0; i-- {
			mut = apd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, apd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (apd *AdhocPlanDelete) ExecX(ctx context.Context) int {
	n, err := apd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (apd *AdhocPlanDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: adhocplan.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: adhocplan.FieldID,
			},
		},
	}
	if ps := apd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, apd.driver, _spec)
}

// AdhocPlanDeleteOne is the builder for deleting a single AdhocPlan entity.
type AdhocPlanDeleteOne struct {
	apd *AdhocPlanDelete
}

// Exec executes the deletion query.
func (apdo *AdhocPlanDeleteOne) Exec(ctx context.Context) error {
	n, err := apdo.apd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{adhocplan.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (apdo *AdhocPlanDeleteOne) ExecX(ctx context.Context) {
	apdo.apd.ExecX(ctx)
}
