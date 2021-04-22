// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/gen0cide/laforge/ent/build"
	"github.com/gen0cide/laforge/ent/environment"
	"github.com/gen0cide/laforge/ent/plan"
	"github.com/gen0cide/laforge/ent/predicate"
	"github.com/gen0cide/laforge/ent/provisionednetwork"
	"github.com/gen0cide/laforge/ent/status"
	"github.com/gen0cide/laforge/ent/team"
)

// BuildUpdate is the builder for updating Build entities.
type BuildUpdate struct {
	config
	hooks    []Hook
	mutation *BuildMutation
}

// Where adds a new predicate for the BuildUpdate builder.
func (bu *BuildUpdate) Where(ps ...predicate.Build) *BuildUpdate {
	bu.mutation.predicates = append(bu.mutation.predicates, ps...)
	return bu
}

// SetRevision sets the "revision" field.
func (bu *BuildUpdate) SetRevision(i int) *BuildUpdate {
	bu.mutation.ResetRevision()
	bu.mutation.SetRevision(i)
	return bu
}

// AddRevision adds i to the "revision" field.
func (bu *BuildUpdate) AddRevision(i int) *BuildUpdate {
	bu.mutation.AddRevision(i)
	return bu
}

// SetBuildToStatusID sets the "BuildToStatus" edge to the Status entity by ID.
func (bu *BuildUpdate) SetBuildToStatusID(id int) *BuildUpdate {
	bu.mutation.SetBuildToStatusID(id)
	return bu
}

// SetNillableBuildToStatusID sets the "BuildToStatus" edge to the Status entity by ID if the given value is not nil.
func (bu *BuildUpdate) SetNillableBuildToStatusID(id *int) *BuildUpdate {
	if id != nil {
		bu = bu.SetBuildToStatusID(*id)
	}
	return bu
}

// SetBuildToStatus sets the "BuildToStatus" edge to the Status entity.
func (bu *BuildUpdate) SetBuildToStatus(s *Status) *BuildUpdate {
	return bu.SetBuildToStatusID(s.ID)
}

// SetBuildToEnvironmentID sets the "BuildToEnvironment" edge to the Environment entity by ID.
func (bu *BuildUpdate) SetBuildToEnvironmentID(id int) *BuildUpdate {
	bu.mutation.SetBuildToEnvironmentID(id)
	return bu
}

// SetBuildToEnvironment sets the "BuildToEnvironment" edge to the Environment entity.
func (bu *BuildUpdate) SetBuildToEnvironment(e *Environment) *BuildUpdate {
	return bu.SetBuildToEnvironmentID(e.ID)
}

// AddBuildToProvisionedNetworkIDs adds the "BuildToProvisionedNetwork" edge to the ProvisionedNetwork entity by IDs.
func (bu *BuildUpdate) AddBuildToProvisionedNetworkIDs(ids ...int) *BuildUpdate {
	bu.mutation.AddBuildToProvisionedNetworkIDs(ids...)
	return bu
}

// AddBuildToProvisionedNetwork adds the "BuildToProvisionedNetwork" edges to the ProvisionedNetwork entity.
func (bu *BuildUpdate) AddBuildToProvisionedNetwork(p ...*ProvisionedNetwork) *BuildUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return bu.AddBuildToProvisionedNetworkIDs(ids...)
}

// AddBuildToTeamIDs adds the "BuildToTeam" edge to the Team entity by IDs.
func (bu *BuildUpdate) AddBuildToTeamIDs(ids ...int) *BuildUpdate {
	bu.mutation.AddBuildToTeamIDs(ids...)
	return bu
}

// AddBuildToTeam adds the "BuildToTeam" edges to the Team entity.
func (bu *BuildUpdate) AddBuildToTeam(t ...*Team) *BuildUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return bu.AddBuildToTeamIDs(ids...)
}

// AddBuildToPlanIDs adds the "BuildToPlan" edge to the Plan entity by IDs.
func (bu *BuildUpdate) AddBuildToPlanIDs(ids ...int) *BuildUpdate {
	bu.mutation.AddBuildToPlanIDs(ids...)
	return bu
}

// AddBuildToPlan adds the "BuildToPlan" edges to the Plan entity.
func (bu *BuildUpdate) AddBuildToPlan(p ...*Plan) *BuildUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return bu.AddBuildToPlanIDs(ids...)
}

// Mutation returns the BuildMutation object of the builder.
func (bu *BuildUpdate) Mutation() *BuildMutation {
	return bu.mutation
}

// ClearBuildToStatus clears the "BuildToStatus" edge to the Status entity.
func (bu *BuildUpdate) ClearBuildToStatus() *BuildUpdate {
	bu.mutation.ClearBuildToStatus()
	return bu
}

// ClearBuildToEnvironment clears the "BuildToEnvironment" edge to the Environment entity.
func (bu *BuildUpdate) ClearBuildToEnvironment() *BuildUpdate {
	bu.mutation.ClearBuildToEnvironment()
	return bu
}

// ClearBuildToProvisionedNetwork clears all "BuildToProvisionedNetwork" edges to the ProvisionedNetwork entity.
func (bu *BuildUpdate) ClearBuildToProvisionedNetwork() *BuildUpdate {
	bu.mutation.ClearBuildToProvisionedNetwork()
	return bu
}

// RemoveBuildToProvisionedNetworkIDs removes the "BuildToProvisionedNetwork" edge to ProvisionedNetwork entities by IDs.
func (bu *BuildUpdate) RemoveBuildToProvisionedNetworkIDs(ids ...int) *BuildUpdate {
	bu.mutation.RemoveBuildToProvisionedNetworkIDs(ids...)
	return bu
}

// RemoveBuildToProvisionedNetwork removes "BuildToProvisionedNetwork" edges to ProvisionedNetwork entities.
func (bu *BuildUpdate) RemoveBuildToProvisionedNetwork(p ...*ProvisionedNetwork) *BuildUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return bu.RemoveBuildToProvisionedNetworkIDs(ids...)
}

// ClearBuildToTeam clears all "BuildToTeam" edges to the Team entity.
func (bu *BuildUpdate) ClearBuildToTeam() *BuildUpdate {
	bu.mutation.ClearBuildToTeam()
	return bu
}

// RemoveBuildToTeamIDs removes the "BuildToTeam" edge to Team entities by IDs.
func (bu *BuildUpdate) RemoveBuildToTeamIDs(ids ...int) *BuildUpdate {
	bu.mutation.RemoveBuildToTeamIDs(ids...)
	return bu
}

// RemoveBuildToTeam removes "BuildToTeam" edges to Team entities.
func (bu *BuildUpdate) RemoveBuildToTeam(t ...*Team) *BuildUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return bu.RemoveBuildToTeamIDs(ids...)
}

// ClearBuildToPlan clears all "BuildToPlan" edges to the Plan entity.
func (bu *BuildUpdate) ClearBuildToPlan() *BuildUpdate {
	bu.mutation.ClearBuildToPlan()
	return bu
}

// RemoveBuildToPlanIDs removes the "BuildToPlan" edge to Plan entities by IDs.
func (bu *BuildUpdate) RemoveBuildToPlanIDs(ids ...int) *BuildUpdate {
	bu.mutation.RemoveBuildToPlanIDs(ids...)
	return bu
}

// RemoveBuildToPlan removes "BuildToPlan" edges to Plan entities.
func (bu *BuildUpdate) RemoveBuildToPlan(p ...*Plan) *BuildUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return bu.RemoveBuildToPlanIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (bu *BuildUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(bu.hooks) == 0 {
		if err = bu.check(); err != nil {
			return 0, err
		}
		affected, err = bu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*BuildMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = bu.check(); err != nil {
				return 0, err
			}
			bu.mutation = mutation
			affected, err = bu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(bu.hooks) - 1; i >= 0; i-- {
			mut = bu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, bu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (bu *BuildUpdate) SaveX(ctx context.Context) int {
	affected, err := bu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (bu *BuildUpdate) Exec(ctx context.Context) error {
	_, err := bu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bu *BuildUpdate) ExecX(ctx context.Context) {
	if err := bu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bu *BuildUpdate) check() error {
	if _, ok := bu.mutation.BuildToEnvironmentID(); bu.mutation.BuildToEnvironmentCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"BuildToEnvironment\"")
	}
	return nil
}

func (bu *BuildUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   build.Table,
			Columns: build.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: build.FieldID,
			},
		},
	}
	if ps := bu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := bu.mutation.Revision(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: build.FieldRevision,
		})
	}
	if value, ok := bu.mutation.AddedRevision(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: build.FieldRevision,
		})
	}
	if bu.mutation.BuildToStatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   build.BuildToStatusTable,
			Columns: []string{build.BuildToStatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: status.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.BuildToStatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   build.BuildToStatusTable,
			Columns: []string{build.BuildToStatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: status.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if bu.mutation.BuildToEnvironmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   build.BuildToEnvironmentTable,
			Columns: []string{build.BuildToEnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: environment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.BuildToEnvironmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   build.BuildToEnvironmentTable,
			Columns: []string{build.BuildToEnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: environment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if bu.mutation.BuildToProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToProvisionedNetworkTable,
			Columns: []string{build.BuildToProvisionedNetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionednetwork.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.RemovedBuildToProvisionedNetworkIDs(); len(nodes) > 0 && !bu.mutation.BuildToProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToProvisionedNetworkTable,
			Columns: []string{build.BuildToProvisionedNetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionednetwork.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.BuildToProvisionedNetworkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToProvisionedNetworkTable,
			Columns: []string{build.BuildToProvisionedNetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionednetwork.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if bu.mutation.BuildToTeamCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToTeamTable,
			Columns: []string{build.BuildToTeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: team.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.RemovedBuildToTeamIDs(); len(nodes) > 0 && !bu.mutation.BuildToTeamCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToTeamTable,
			Columns: []string{build.BuildToTeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: team.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.BuildToTeamIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToTeamTable,
			Columns: []string{build.BuildToTeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: team.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if bu.mutation.BuildToPlanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToPlanTable,
			Columns: []string{build.BuildToPlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: plan.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.RemovedBuildToPlanIDs(); len(nodes) > 0 && !bu.mutation.BuildToPlanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToPlanTable,
			Columns: []string{build.BuildToPlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: plan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.BuildToPlanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToPlanTable,
			Columns: []string{build.BuildToPlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: plan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, bu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{build.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// BuildUpdateOne is the builder for updating a single Build entity.
type BuildUpdateOne struct {
	config
	hooks    []Hook
	mutation *BuildMutation
}

// SetRevision sets the "revision" field.
func (buo *BuildUpdateOne) SetRevision(i int) *BuildUpdateOne {
	buo.mutation.ResetRevision()
	buo.mutation.SetRevision(i)
	return buo
}

// AddRevision adds i to the "revision" field.
func (buo *BuildUpdateOne) AddRevision(i int) *BuildUpdateOne {
	buo.mutation.AddRevision(i)
	return buo
}

// SetBuildToStatusID sets the "BuildToStatus" edge to the Status entity by ID.
func (buo *BuildUpdateOne) SetBuildToStatusID(id int) *BuildUpdateOne {
	buo.mutation.SetBuildToStatusID(id)
	return buo
}

// SetNillableBuildToStatusID sets the "BuildToStatus" edge to the Status entity by ID if the given value is not nil.
func (buo *BuildUpdateOne) SetNillableBuildToStatusID(id *int) *BuildUpdateOne {
	if id != nil {
		buo = buo.SetBuildToStatusID(*id)
	}
	return buo
}

// SetBuildToStatus sets the "BuildToStatus" edge to the Status entity.
func (buo *BuildUpdateOne) SetBuildToStatus(s *Status) *BuildUpdateOne {
	return buo.SetBuildToStatusID(s.ID)
}

// SetBuildToEnvironmentID sets the "BuildToEnvironment" edge to the Environment entity by ID.
func (buo *BuildUpdateOne) SetBuildToEnvironmentID(id int) *BuildUpdateOne {
	buo.mutation.SetBuildToEnvironmentID(id)
	return buo
}

// SetBuildToEnvironment sets the "BuildToEnvironment" edge to the Environment entity.
func (buo *BuildUpdateOne) SetBuildToEnvironment(e *Environment) *BuildUpdateOne {
	return buo.SetBuildToEnvironmentID(e.ID)
}

// AddBuildToProvisionedNetworkIDs adds the "BuildToProvisionedNetwork" edge to the ProvisionedNetwork entity by IDs.
func (buo *BuildUpdateOne) AddBuildToProvisionedNetworkIDs(ids ...int) *BuildUpdateOne {
	buo.mutation.AddBuildToProvisionedNetworkIDs(ids...)
	return buo
}

// AddBuildToProvisionedNetwork adds the "BuildToProvisionedNetwork" edges to the ProvisionedNetwork entity.
func (buo *BuildUpdateOne) AddBuildToProvisionedNetwork(p ...*ProvisionedNetwork) *BuildUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return buo.AddBuildToProvisionedNetworkIDs(ids...)
}

// AddBuildToTeamIDs adds the "BuildToTeam" edge to the Team entity by IDs.
func (buo *BuildUpdateOne) AddBuildToTeamIDs(ids ...int) *BuildUpdateOne {
	buo.mutation.AddBuildToTeamIDs(ids...)
	return buo
}

// AddBuildToTeam adds the "BuildToTeam" edges to the Team entity.
func (buo *BuildUpdateOne) AddBuildToTeam(t ...*Team) *BuildUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return buo.AddBuildToTeamIDs(ids...)
}

// AddBuildToPlanIDs adds the "BuildToPlan" edge to the Plan entity by IDs.
func (buo *BuildUpdateOne) AddBuildToPlanIDs(ids ...int) *BuildUpdateOne {
	buo.mutation.AddBuildToPlanIDs(ids...)
	return buo
}

// AddBuildToPlan adds the "BuildToPlan" edges to the Plan entity.
func (buo *BuildUpdateOne) AddBuildToPlan(p ...*Plan) *BuildUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return buo.AddBuildToPlanIDs(ids...)
}

// Mutation returns the BuildMutation object of the builder.
func (buo *BuildUpdateOne) Mutation() *BuildMutation {
	return buo.mutation
}

// ClearBuildToStatus clears the "BuildToStatus" edge to the Status entity.
func (buo *BuildUpdateOne) ClearBuildToStatus() *BuildUpdateOne {
	buo.mutation.ClearBuildToStatus()
	return buo
}

// ClearBuildToEnvironment clears the "BuildToEnvironment" edge to the Environment entity.
func (buo *BuildUpdateOne) ClearBuildToEnvironment() *BuildUpdateOne {
	buo.mutation.ClearBuildToEnvironment()
	return buo
}

// ClearBuildToProvisionedNetwork clears all "BuildToProvisionedNetwork" edges to the ProvisionedNetwork entity.
func (buo *BuildUpdateOne) ClearBuildToProvisionedNetwork() *BuildUpdateOne {
	buo.mutation.ClearBuildToProvisionedNetwork()
	return buo
}

// RemoveBuildToProvisionedNetworkIDs removes the "BuildToProvisionedNetwork" edge to ProvisionedNetwork entities by IDs.
func (buo *BuildUpdateOne) RemoveBuildToProvisionedNetworkIDs(ids ...int) *BuildUpdateOne {
	buo.mutation.RemoveBuildToProvisionedNetworkIDs(ids...)
	return buo
}

// RemoveBuildToProvisionedNetwork removes "BuildToProvisionedNetwork" edges to ProvisionedNetwork entities.
func (buo *BuildUpdateOne) RemoveBuildToProvisionedNetwork(p ...*ProvisionedNetwork) *BuildUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return buo.RemoveBuildToProvisionedNetworkIDs(ids...)
}

// ClearBuildToTeam clears all "BuildToTeam" edges to the Team entity.
func (buo *BuildUpdateOne) ClearBuildToTeam() *BuildUpdateOne {
	buo.mutation.ClearBuildToTeam()
	return buo
}

// RemoveBuildToTeamIDs removes the "BuildToTeam" edge to Team entities by IDs.
func (buo *BuildUpdateOne) RemoveBuildToTeamIDs(ids ...int) *BuildUpdateOne {
	buo.mutation.RemoveBuildToTeamIDs(ids...)
	return buo
}

// RemoveBuildToTeam removes "BuildToTeam" edges to Team entities.
func (buo *BuildUpdateOne) RemoveBuildToTeam(t ...*Team) *BuildUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return buo.RemoveBuildToTeamIDs(ids...)
}

// ClearBuildToPlan clears all "BuildToPlan" edges to the Plan entity.
func (buo *BuildUpdateOne) ClearBuildToPlan() *BuildUpdateOne {
	buo.mutation.ClearBuildToPlan()
	return buo
}

// RemoveBuildToPlanIDs removes the "BuildToPlan" edge to Plan entities by IDs.
func (buo *BuildUpdateOne) RemoveBuildToPlanIDs(ids ...int) *BuildUpdateOne {
	buo.mutation.RemoveBuildToPlanIDs(ids...)
	return buo
}

// RemoveBuildToPlan removes "BuildToPlan" edges to Plan entities.
func (buo *BuildUpdateOne) RemoveBuildToPlan(p ...*Plan) *BuildUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return buo.RemoveBuildToPlanIDs(ids...)
}

// Save executes the query and returns the updated Build entity.
func (buo *BuildUpdateOne) Save(ctx context.Context) (*Build, error) {
	var (
		err  error
		node *Build
	)
	if len(buo.hooks) == 0 {
		if err = buo.check(); err != nil {
			return nil, err
		}
		node, err = buo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*BuildMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = buo.check(); err != nil {
				return nil, err
			}
			buo.mutation = mutation
			node, err = buo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(buo.hooks) - 1; i >= 0; i-- {
			mut = buo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, buo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (buo *BuildUpdateOne) SaveX(ctx context.Context) *Build {
	node, err := buo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (buo *BuildUpdateOne) Exec(ctx context.Context) error {
	_, err := buo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (buo *BuildUpdateOne) ExecX(ctx context.Context) {
	if err := buo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (buo *BuildUpdateOne) check() error {
	if _, ok := buo.mutation.BuildToEnvironmentID(); buo.mutation.BuildToEnvironmentCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"BuildToEnvironment\"")
	}
	return nil
}

func (buo *BuildUpdateOne) sqlSave(ctx context.Context) (_node *Build, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   build.Table,
			Columns: build.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: build.FieldID,
			},
		},
	}
	id, ok := buo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Build.ID for update")}
	}
	_spec.Node.ID.Value = id
	if ps := buo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := buo.mutation.Revision(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: build.FieldRevision,
		})
	}
	if value, ok := buo.mutation.AddedRevision(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: build.FieldRevision,
		})
	}
	if buo.mutation.BuildToStatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   build.BuildToStatusTable,
			Columns: []string{build.BuildToStatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: status.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.BuildToStatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   build.BuildToStatusTable,
			Columns: []string{build.BuildToStatusColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: status.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if buo.mutation.BuildToEnvironmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   build.BuildToEnvironmentTable,
			Columns: []string{build.BuildToEnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: environment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.BuildToEnvironmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   build.BuildToEnvironmentTable,
			Columns: []string{build.BuildToEnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: environment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if buo.mutation.BuildToProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToProvisionedNetworkTable,
			Columns: []string{build.BuildToProvisionedNetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionednetwork.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.RemovedBuildToProvisionedNetworkIDs(); len(nodes) > 0 && !buo.mutation.BuildToProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToProvisionedNetworkTable,
			Columns: []string{build.BuildToProvisionedNetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionednetwork.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.BuildToProvisionedNetworkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToProvisionedNetworkTable,
			Columns: []string{build.BuildToProvisionedNetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionednetwork.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if buo.mutation.BuildToTeamCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToTeamTable,
			Columns: []string{build.BuildToTeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: team.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.RemovedBuildToTeamIDs(); len(nodes) > 0 && !buo.mutation.BuildToTeamCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToTeamTable,
			Columns: []string{build.BuildToTeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: team.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.BuildToTeamIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToTeamTable,
			Columns: []string{build.BuildToTeamColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: team.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if buo.mutation.BuildToPlanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToPlanTable,
			Columns: []string{build.BuildToPlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: plan.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.RemovedBuildToPlanIDs(); len(nodes) > 0 && !buo.mutation.BuildToPlanCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToPlanTable,
			Columns: []string{build.BuildToPlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: plan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.BuildToPlanIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   build.BuildToPlanTable,
			Columns: []string{build.BuildToPlanColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: plan.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Build{config: buo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, buo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{build.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}
