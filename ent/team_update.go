// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/gen0cide/laforge/ent/build"
	"github.com/gen0cide/laforge/ent/environment"
	"github.com/gen0cide/laforge/ent/predicate"
	"github.com/gen0cide/laforge/ent/provisionednetwork"
	"github.com/gen0cide/laforge/ent/tag"
	"github.com/gen0cide/laforge/ent/team"
	"github.com/gen0cide/laforge/ent/user"
)

// TeamUpdate is the builder for updating Team entities.
type TeamUpdate struct {
	config
	hooks    []Hook
	mutation *TeamMutation
}

// Where adds a new predicate for the TeamUpdate builder.
func (tu *TeamUpdate) Where(ps ...predicate.Team) *TeamUpdate {
	tu.mutation.predicates = append(tu.mutation.predicates, ps...)
	return tu
}

// SetTeamNumber sets the "team_number" field.
func (tu *TeamUpdate) SetTeamNumber(i int) *TeamUpdate {
	tu.mutation.ResetTeamNumber()
	tu.mutation.SetTeamNumber(i)
	return tu
}

// AddTeamNumber adds i to the "team_number" field.
func (tu *TeamUpdate) AddTeamNumber(i int) *TeamUpdate {
	tu.mutation.AddTeamNumber(i)
	return tu
}

// SetConfig sets the "config" field.
func (tu *TeamUpdate) SetConfig(m map[string]string) *TeamUpdate {
	tu.mutation.SetConfig(m)
	return tu
}

// SetRevision sets the "revision" field.
func (tu *TeamUpdate) SetRevision(i int64) *TeamUpdate {
	tu.mutation.ResetRevision()
	tu.mutation.SetRevision(i)
	return tu
}

// AddRevision adds i to the "revision" field.
func (tu *TeamUpdate) AddRevision(i int64) *TeamUpdate {
	tu.mutation.AddRevision(i)
	return tu
}

// AddTeamToUserIDs adds the "TeamToUser" edge to the User entity by IDs.
func (tu *TeamUpdate) AddTeamToUserIDs(ids ...int) *TeamUpdate {
	tu.mutation.AddTeamToUserIDs(ids...)
	return tu
}

// AddTeamToUser adds the "TeamToUser" edges to the User entity.
func (tu *TeamUpdate) AddTeamToUser(u ...*User) *TeamUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return tu.AddTeamToUserIDs(ids...)
}

// AddTeamToBuildIDs adds the "TeamToBuild" edge to the Build entity by IDs.
func (tu *TeamUpdate) AddTeamToBuildIDs(ids ...int) *TeamUpdate {
	tu.mutation.AddTeamToBuildIDs(ids...)
	return tu
}

// AddTeamToBuild adds the "TeamToBuild" edges to the Build entity.
func (tu *TeamUpdate) AddTeamToBuild(b ...*Build) *TeamUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return tu.AddTeamToBuildIDs(ids...)
}

// AddTeamToEnvironmentIDs adds the "TeamToEnvironment" edge to the Environment entity by IDs.
func (tu *TeamUpdate) AddTeamToEnvironmentIDs(ids ...int) *TeamUpdate {
	tu.mutation.AddTeamToEnvironmentIDs(ids...)
	return tu
}

// AddTeamToEnvironment adds the "TeamToEnvironment" edges to the Environment entity.
func (tu *TeamUpdate) AddTeamToEnvironment(e ...*Environment) *TeamUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return tu.AddTeamToEnvironmentIDs(ids...)
}

// AddTeamToTagIDs adds the "TeamToTag" edge to the Tag entity by IDs.
func (tu *TeamUpdate) AddTeamToTagIDs(ids ...int) *TeamUpdate {
	tu.mutation.AddTeamToTagIDs(ids...)
	return tu
}

// AddTeamToTag adds the "TeamToTag" edges to the Tag entity.
func (tu *TeamUpdate) AddTeamToTag(t ...*Tag) *TeamUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tu.AddTeamToTagIDs(ids...)
}

// AddTeamToProvisionedNetworkIDs adds the "TeamToProvisionedNetwork" edge to the ProvisionedNetwork entity by IDs.
func (tu *TeamUpdate) AddTeamToProvisionedNetworkIDs(ids ...int) *TeamUpdate {
	tu.mutation.AddTeamToProvisionedNetworkIDs(ids...)
	return tu
}

// AddTeamToProvisionedNetwork adds the "TeamToProvisionedNetwork" edges to the ProvisionedNetwork entity.
func (tu *TeamUpdate) AddTeamToProvisionedNetwork(p ...*ProvisionedNetwork) *TeamUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return tu.AddTeamToProvisionedNetworkIDs(ids...)
}

// Mutation returns the TeamMutation object of the builder.
func (tu *TeamUpdate) Mutation() *TeamMutation {
	return tu.mutation
}

// ClearTeamToUser clears all "TeamToUser" edges to the User entity.
func (tu *TeamUpdate) ClearTeamToUser() *TeamUpdate {
	tu.mutation.ClearTeamToUser()
	return tu
}

// RemoveTeamToUserIDs removes the "TeamToUser" edge to User entities by IDs.
func (tu *TeamUpdate) RemoveTeamToUserIDs(ids ...int) *TeamUpdate {
	tu.mutation.RemoveTeamToUserIDs(ids...)
	return tu
}

// RemoveTeamToUser removes "TeamToUser" edges to User entities.
func (tu *TeamUpdate) RemoveTeamToUser(u ...*User) *TeamUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return tu.RemoveTeamToUserIDs(ids...)
}

// ClearTeamToBuild clears all "TeamToBuild" edges to the Build entity.
func (tu *TeamUpdate) ClearTeamToBuild() *TeamUpdate {
	tu.mutation.ClearTeamToBuild()
	return tu
}

// RemoveTeamToBuildIDs removes the "TeamToBuild" edge to Build entities by IDs.
func (tu *TeamUpdate) RemoveTeamToBuildIDs(ids ...int) *TeamUpdate {
	tu.mutation.RemoveTeamToBuildIDs(ids...)
	return tu
}

// RemoveTeamToBuild removes "TeamToBuild" edges to Build entities.
func (tu *TeamUpdate) RemoveTeamToBuild(b ...*Build) *TeamUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return tu.RemoveTeamToBuildIDs(ids...)
}

// ClearTeamToEnvironment clears all "TeamToEnvironment" edges to the Environment entity.
func (tu *TeamUpdate) ClearTeamToEnvironment() *TeamUpdate {
	tu.mutation.ClearTeamToEnvironment()
	return tu
}

// RemoveTeamToEnvironmentIDs removes the "TeamToEnvironment" edge to Environment entities by IDs.
func (tu *TeamUpdate) RemoveTeamToEnvironmentIDs(ids ...int) *TeamUpdate {
	tu.mutation.RemoveTeamToEnvironmentIDs(ids...)
	return tu
}

// RemoveTeamToEnvironment removes "TeamToEnvironment" edges to Environment entities.
func (tu *TeamUpdate) RemoveTeamToEnvironment(e ...*Environment) *TeamUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return tu.RemoveTeamToEnvironmentIDs(ids...)
}

// ClearTeamToTag clears all "TeamToTag" edges to the Tag entity.
func (tu *TeamUpdate) ClearTeamToTag() *TeamUpdate {
	tu.mutation.ClearTeamToTag()
	return tu
}

// RemoveTeamToTagIDs removes the "TeamToTag" edge to Tag entities by IDs.
func (tu *TeamUpdate) RemoveTeamToTagIDs(ids ...int) *TeamUpdate {
	tu.mutation.RemoveTeamToTagIDs(ids...)
	return tu
}

// RemoveTeamToTag removes "TeamToTag" edges to Tag entities.
func (tu *TeamUpdate) RemoveTeamToTag(t ...*Tag) *TeamUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tu.RemoveTeamToTagIDs(ids...)
}

// ClearTeamToProvisionedNetwork clears all "TeamToProvisionedNetwork" edges to the ProvisionedNetwork entity.
func (tu *TeamUpdate) ClearTeamToProvisionedNetwork() *TeamUpdate {
	tu.mutation.ClearTeamToProvisionedNetwork()
	return tu
}

// RemoveTeamToProvisionedNetworkIDs removes the "TeamToProvisionedNetwork" edge to ProvisionedNetwork entities by IDs.
func (tu *TeamUpdate) RemoveTeamToProvisionedNetworkIDs(ids ...int) *TeamUpdate {
	tu.mutation.RemoveTeamToProvisionedNetworkIDs(ids...)
	return tu
}

// RemoveTeamToProvisionedNetwork removes "TeamToProvisionedNetwork" edges to ProvisionedNetwork entities.
func (tu *TeamUpdate) RemoveTeamToProvisionedNetwork(p ...*ProvisionedNetwork) *TeamUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return tu.RemoveTeamToProvisionedNetworkIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tu *TeamUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(tu.hooks) == 0 {
		affected, err = tu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TeamMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tu.mutation = mutation
			affected, err = tu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(tu.hooks) - 1; i >= 0; i-- {
			mut = tu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TeamUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TeamUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TeamUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tu *TeamUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   team.Table,
			Columns: team.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: team.FieldID,
			},
		},
	}
	if ps := tu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.TeamNumber(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: team.FieldTeamNumber,
		})
	}
	if value, ok := tu.mutation.AddedTeamNumber(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: team.FieldTeamNumber,
		})
	}
	if value, ok := tu.mutation.Config(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: team.FieldConfig,
		})
	}
	if value, ok := tu.mutation.Revision(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: team.FieldRevision,
		})
	}
	if value, ok := tu.mutation.AddedRevision(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: team.FieldRevision,
		})
	}
	if tu.mutation.TeamToUserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToUserTable,
			Columns: []string{team.TeamToUserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedTeamToUserIDs(); len(nodes) > 0 && !tu.mutation.TeamToUserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToUserTable,
			Columns: []string{team.TeamToUserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.TeamToUserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToUserTable,
			Columns: []string{team.TeamToUserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tu.mutation.TeamToBuildCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToBuildTable,
			Columns: team.TeamToBuildPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: build.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedTeamToBuildIDs(); len(nodes) > 0 && !tu.mutation.TeamToBuildCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToBuildTable,
			Columns: team.TeamToBuildPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: build.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.TeamToBuildIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToBuildTable,
			Columns: team.TeamToBuildPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: build.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tu.mutation.TeamToEnvironmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToEnvironmentTable,
			Columns: team.TeamToEnvironmentPrimaryKey,
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
	if nodes := tu.mutation.RemovedTeamToEnvironmentIDs(); len(nodes) > 0 && !tu.mutation.TeamToEnvironmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToEnvironmentTable,
			Columns: team.TeamToEnvironmentPrimaryKey,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.TeamToEnvironmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToEnvironmentTable,
			Columns: team.TeamToEnvironmentPrimaryKey,
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
	if tu.mutation.TeamToTagCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToTagTable,
			Columns: []string{team.TeamToTagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedTeamToTagIDs(); len(nodes) > 0 && !tu.mutation.TeamToTagCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToTagTable,
			Columns: []string{team.TeamToTagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.TeamToTagIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToTagTable,
			Columns: []string{team.TeamToTagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tu.mutation.TeamToProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.TeamToProvisionedNetworkTable,
			Columns: team.TeamToProvisionedNetworkPrimaryKey,
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
	if nodes := tu.mutation.RemovedTeamToProvisionedNetworkIDs(); len(nodes) > 0 && !tu.mutation.TeamToProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.TeamToProvisionedNetworkTable,
			Columns: team.TeamToProvisionedNetworkPrimaryKey,
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
	if nodes := tu.mutation.TeamToProvisionedNetworkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.TeamToProvisionedNetworkTable,
			Columns: team.TeamToProvisionedNetworkPrimaryKey,
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
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{team.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// TeamUpdateOne is the builder for updating a single Team entity.
type TeamUpdateOne struct {
	config
	hooks    []Hook
	mutation *TeamMutation
}

// SetTeamNumber sets the "team_number" field.
func (tuo *TeamUpdateOne) SetTeamNumber(i int) *TeamUpdateOne {
	tuo.mutation.ResetTeamNumber()
	tuo.mutation.SetTeamNumber(i)
	return tuo
}

// AddTeamNumber adds i to the "team_number" field.
func (tuo *TeamUpdateOne) AddTeamNumber(i int) *TeamUpdateOne {
	tuo.mutation.AddTeamNumber(i)
	return tuo
}

// SetConfig sets the "config" field.
func (tuo *TeamUpdateOne) SetConfig(m map[string]string) *TeamUpdateOne {
	tuo.mutation.SetConfig(m)
	return tuo
}

// SetRevision sets the "revision" field.
func (tuo *TeamUpdateOne) SetRevision(i int64) *TeamUpdateOne {
	tuo.mutation.ResetRevision()
	tuo.mutation.SetRevision(i)
	return tuo
}

// AddRevision adds i to the "revision" field.
func (tuo *TeamUpdateOne) AddRevision(i int64) *TeamUpdateOne {
	tuo.mutation.AddRevision(i)
	return tuo
}

// AddTeamToUserIDs adds the "TeamToUser" edge to the User entity by IDs.
func (tuo *TeamUpdateOne) AddTeamToUserIDs(ids ...int) *TeamUpdateOne {
	tuo.mutation.AddTeamToUserIDs(ids...)
	return tuo
}

// AddTeamToUser adds the "TeamToUser" edges to the User entity.
func (tuo *TeamUpdateOne) AddTeamToUser(u ...*User) *TeamUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return tuo.AddTeamToUserIDs(ids...)
}

// AddTeamToBuildIDs adds the "TeamToBuild" edge to the Build entity by IDs.
func (tuo *TeamUpdateOne) AddTeamToBuildIDs(ids ...int) *TeamUpdateOne {
	tuo.mutation.AddTeamToBuildIDs(ids...)
	return tuo
}

// AddTeamToBuild adds the "TeamToBuild" edges to the Build entity.
func (tuo *TeamUpdateOne) AddTeamToBuild(b ...*Build) *TeamUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return tuo.AddTeamToBuildIDs(ids...)
}

// AddTeamToEnvironmentIDs adds the "TeamToEnvironment" edge to the Environment entity by IDs.
func (tuo *TeamUpdateOne) AddTeamToEnvironmentIDs(ids ...int) *TeamUpdateOne {
	tuo.mutation.AddTeamToEnvironmentIDs(ids...)
	return tuo
}

// AddTeamToEnvironment adds the "TeamToEnvironment" edges to the Environment entity.
func (tuo *TeamUpdateOne) AddTeamToEnvironment(e ...*Environment) *TeamUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return tuo.AddTeamToEnvironmentIDs(ids...)
}

// AddTeamToTagIDs adds the "TeamToTag" edge to the Tag entity by IDs.
func (tuo *TeamUpdateOne) AddTeamToTagIDs(ids ...int) *TeamUpdateOne {
	tuo.mutation.AddTeamToTagIDs(ids...)
	return tuo
}

// AddTeamToTag adds the "TeamToTag" edges to the Tag entity.
func (tuo *TeamUpdateOne) AddTeamToTag(t ...*Tag) *TeamUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tuo.AddTeamToTagIDs(ids...)
}

// AddTeamToProvisionedNetworkIDs adds the "TeamToProvisionedNetwork" edge to the ProvisionedNetwork entity by IDs.
func (tuo *TeamUpdateOne) AddTeamToProvisionedNetworkIDs(ids ...int) *TeamUpdateOne {
	tuo.mutation.AddTeamToProvisionedNetworkIDs(ids...)
	return tuo
}

// AddTeamToProvisionedNetwork adds the "TeamToProvisionedNetwork" edges to the ProvisionedNetwork entity.
func (tuo *TeamUpdateOne) AddTeamToProvisionedNetwork(p ...*ProvisionedNetwork) *TeamUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return tuo.AddTeamToProvisionedNetworkIDs(ids...)
}

// Mutation returns the TeamMutation object of the builder.
func (tuo *TeamUpdateOne) Mutation() *TeamMutation {
	return tuo.mutation
}

// ClearTeamToUser clears all "TeamToUser" edges to the User entity.
func (tuo *TeamUpdateOne) ClearTeamToUser() *TeamUpdateOne {
	tuo.mutation.ClearTeamToUser()
	return tuo
}

// RemoveTeamToUserIDs removes the "TeamToUser" edge to User entities by IDs.
func (tuo *TeamUpdateOne) RemoveTeamToUserIDs(ids ...int) *TeamUpdateOne {
	tuo.mutation.RemoveTeamToUserIDs(ids...)
	return tuo
}

// RemoveTeamToUser removes "TeamToUser" edges to User entities.
func (tuo *TeamUpdateOne) RemoveTeamToUser(u ...*User) *TeamUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return tuo.RemoveTeamToUserIDs(ids...)
}

// ClearTeamToBuild clears all "TeamToBuild" edges to the Build entity.
func (tuo *TeamUpdateOne) ClearTeamToBuild() *TeamUpdateOne {
	tuo.mutation.ClearTeamToBuild()
	return tuo
}

// RemoveTeamToBuildIDs removes the "TeamToBuild" edge to Build entities by IDs.
func (tuo *TeamUpdateOne) RemoveTeamToBuildIDs(ids ...int) *TeamUpdateOne {
	tuo.mutation.RemoveTeamToBuildIDs(ids...)
	return tuo
}

// RemoveTeamToBuild removes "TeamToBuild" edges to Build entities.
func (tuo *TeamUpdateOne) RemoveTeamToBuild(b ...*Build) *TeamUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return tuo.RemoveTeamToBuildIDs(ids...)
}

// ClearTeamToEnvironment clears all "TeamToEnvironment" edges to the Environment entity.
func (tuo *TeamUpdateOne) ClearTeamToEnvironment() *TeamUpdateOne {
	tuo.mutation.ClearTeamToEnvironment()
	return tuo
}

// RemoveTeamToEnvironmentIDs removes the "TeamToEnvironment" edge to Environment entities by IDs.
func (tuo *TeamUpdateOne) RemoveTeamToEnvironmentIDs(ids ...int) *TeamUpdateOne {
	tuo.mutation.RemoveTeamToEnvironmentIDs(ids...)
	return tuo
}

// RemoveTeamToEnvironment removes "TeamToEnvironment" edges to Environment entities.
func (tuo *TeamUpdateOne) RemoveTeamToEnvironment(e ...*Environment) *TeamUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return tuo.RemoveTeamToEnvironmentIDs(ids...)
}

// ClearTeamToTag clears all "TeamToTag" edges to the Tag entity.
func (tuo *TeamUpdateOne) ClearTeamToTag() *TeamUpdateOne {
	tuo.mutation.ClearTeamToTag()
	return tuo
}

// RemoveTeamToTagIDs removes the "TeamToTag" edge to Tag entities by IDs.
func (tuo *TeamUpdateOne) RemoveTeamToTagIDs(ids ...int) *TeamUpdateOne {
	tuo.mutation.RemoveTeamToTagIDs(ids...)
	return tuo
}

// RemoveTeamToTag removes "TeamToTag" edges to Tag entities.
func (tuo *TeamUpdateOne) RemoveTeamToTag(t ...*Tag) *TeamUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tuo.RemoveTeamToTagIDs(ids...)
}

// ClearTeamToProvisionedNetwork clears all "TeamToProvisionedNetwork" edges to the ProvisionedNetwork entity.
func (tuo *TeamUpdateOne) ClearTeamToProvisionedNetwork() *TeamUpdateOne {
	tuo.mutation.ClearTeamToProvisionedNetwork()
	return tuo
}

// RemoveTeamToProvisionedNetworkIDs removes the "TeamToProvisionedNetwork" edge to ProvisionedNetwork entities by IDs.
func (tuo *TeamUpdateOne) RemoveTeamToProvisionedNetworkIDs(ids ...int) *TeamUpdateOne {
	tuo.mutation.RemoveTeamToProvisionedNetworkIDs(ids...)
	return tuo
}

// RemoveTeamToProvisionedNetwork removes "TeamToProvisionedNetwork" edges to ProvisionedNetwork entities.
func (tuo *TeamUpdateOne) RemoveTeamToProvisionedNetwork(p ...*ProvisionedNetwork) *TeamUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return tuo.RemoveTeamToProvisionedNetworkIDs(ids...)
}

// Save executes the query and returns the updated Team entity.
func (tuo *TeamUpdateOne) Save(ctx context.Context) (*Team, error) {
	var (
		err  error
		node *Team
	)
	if len(tuo.hooks) == 0 {
		node, err = tuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TeamMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tuo.mutation = mutation
			node, err = tuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(tuo.hooks) - 1; i >= 0; i-- {
			mut = tuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TeamUpdateOne) SaveX(ctx context.Context) *Team {
	node, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tuo *TeamUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TeamUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tuo *TeamUpdateOne) sqlSave(ctx context.Context) (_node *Team, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   team.Table,
			Columns: team.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: team.FieldID,
			},
		},
	}
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Team.ID for update")}
	}
	_spec.Node.ID.Value = id
	if ps := tuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tuo.mutation.TeamNumber(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: team.FieldTeamNumber,
		})
	}
	if value, ok := tuo.mutation.AddedTeamNumber(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: team.FieldTeamNumber,
		})
	}
	if value, ok := tuo.mutation.Config(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: team.FieldConfig,
		})
	}
	if value, ok := tuo.mutation.Revision(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: team.FieldRevision,
		})
	}
	if value, ok := tuo.mutation.AddedRevision(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: team.FieldRevision,
		})
	}
	if tuo.mutation.TeamToUserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToUserTable,
			Columns: []string{team.TeamToUserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedTeamToUserIDs(); len(nodes) > 0 && !tuo.mutation.TeamToUserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToUserTable,
			Columns: []string{team.TeamToUserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.TeamToUserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToUserTable,
			Columns: []string{team.TeamToUserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tuo.mutation.TeamToBuildCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToBuildTable,
			Columns: team.TeamToBuildPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: build.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedTeamToBuildIDs(); len(nodes) > 0 && !tuo.mutation.TeamToBuildCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToBuildTable,
			Columns: team.TeamToBuildPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: build.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.TeamToBuildIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToBuildTable,
			Columns: team.TeamToBuildPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: build.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tuo.mutation.TeamToEnvironmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToEnvironmentTable,
			Columns: team.TeamToEnvironmentPrimaryKey,
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
	if nodes := tuo.mutation.RemovedTeamToEnvironmentIDs(); len(nodes) > 0 && !tuo.mutation.TeamToEnvironmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToEnvironmentTable,
			Columns: team.TeamToEnvironmentPrimaryKey,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.TeamToEnvironmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   team.TeamToEnvironmentTable,
			Columns: team.TeamToEnvironmentPrimaryKey,
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
	if tuo.mutation.TeamToTagCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToTagTable,
			Columns: []string{team.TeamToTagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedTeamToTagIDs(); len(nodes) > 0 && !tuo.mutation.TeamToTagCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToTagTable,
			Columns: []string{team.TeamToTagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.TeamToTagIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.TeamToTagTable,
			Columns: []string{team.TeamToTagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tag.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tuo.mutation.TeamToProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.TeamToProvisionedNetworkTable,
			Columns: team.TeamToProvisionedNetworkPrimaryKey,
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
	if nodes := tuo.mutation.RemovedTeamToProvisionedNetworkIDs(); len(nodes) > 0 && !tuo.mutation.TeamToProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.TeamToProvisionedNetworkTable,
			Columns: team.TeamToProvisionedNetworkPrimaryKey,
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
	if nodes := tuo.mutation.TeamToProvisionedNetworkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   team.TeamToProvisionedNetworkTable,
			Columns: team.TeamToProvisionedNetworkPrimaryKey,
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
	_node = &Team{config: tuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{team.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}