// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/gen0cide/laforge/ent/build"
	"github.com/gen0cide/laforge/ent/network"
	"github.com/gen0cide/laforge/ent/predicate"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	"github.com/gen0cide/laforge/ent/provisionednetwork"
	"github.com/gen0cide/laforge/ent/status"
	"github.com/gen0cide/laforge/ent/team"
)

// ProvisionedNetworkUpdate is the builder for updating ProvisionedNetwork entities.
type ProvisionedNetworkUpdate struct {
	config
	hooks    []Hook
	mutation *ProvisionedNetworkMutation
}

// Where adds a new predicate for the builder.
func (pnu *ProvisionedNetworkUpdate) Where(ps ...predicate.ProvisionedNetwork) *ProvisionedNetworkUpdate {
	pnu.mutation.predicates = append(pnu.mutation.predicates, ps...)
	return pnu
}

// SetName sets the name field.
func (pnu *ProvisionedNetworkUpdate) SetName(s string) *ProvisionedNetworkUpdate {
	pnu.mutation.SetName(s)
	return pnu
}

// SetCidr sets the cidr field.
func (pnu *ProvisionedNetworkUpdate) SetCidr(s string) *ProvisionedNetworkUpdate {
	pnu.mutation.SetCidr(s)
	return pnu
}

// AddStatuIDs adds the status edge to Status by ids.
func (pnu *ProvisionedNetworkUpdate) AddStatuIDs(ids ...int) *ProvisionedNetworkUpdate {
	pnu.mutation.AddStatuIDs(ids...)
	return pnu
}

// AddStatus adds the status edges to Status.
func (pnu *ProvisionedNetworkUpdate) AddStatus(s ...*Status) *ProvisionedNetworkUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return pnu.AddStatuIDs(ids...)
}

// AddNetworkIDs adds the network edge to Network by ids.
func (pnu *ProvisionedNetworkUpdate) AddNetworkIDs(ids ...int) *ProvisionedNetworkUpdate {
	pnu.mutation.AddNetworkIDs(ids...)
	return pnu
}

// AddNetwork adds the network edges to Network.
func (pnu *ProvisionedNetworkUpdate) AddNetwork(n ...*Network) *ProvisionedNetworkUpdate {
	ids := make([]int, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return pnu.AddNetworkIDs(ids...)
}

// AddBuildIDs adds the build edge to Build by ids.
func (pnu *ProvisionedNetworkUpdate) AddBuildIDs(ids ...int) *ProvisionedNetworkUpdate {
	pnu.mutation.AddBuildIDs(ids...)
	return pnu
}

// AddBuild adds the build edges to Build.
func (pnu *ProvisionedNetworkUpdate) AddBuild(b ...*Build) *ProvisionedNetworkUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return pnu.AddBuildIDs(ids...)
}

// AddProvisionedNetworkToTeamIDs adds the ProvisionedNetworkToTeam edge to Team by ids.
func (pnu *ProvisionedNetworkUpdate) AddProvisionedNetworkToTeamIDs(ids ...int) *ProvisionedNetworkUpdate {
	pnu.mutation.AddProvisionedNetworkToTeamIDs(ids...)
	return pnu
}

// AddProvisionedNetworkToTeam adds the ProvisionedNetworkToTeam edges to Team.
func (pnu *ProvisionedNetworkUpdate) AddProvisionedNetworkToTeam(t ...*Team) *ProvisionedNetworkUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return pnu.AddProvisionedNetworkToTeamIDs(ids...)
}

// AddProvisionedHostIDs adds the provisioned_hosts edge to ProvisionedHost by ids.
func (pnu *ProvisionedNetworkUpdate) AddProvisionedHostIDs(ids ...int) *ProvisionedNetworkUpdate {
	pnu.mutation.AddProvisionedHostIDs(ids...)
	return pnu
}

// AddProvisionedHosts adds the provisioned_hosts edges to ProvisionedHost.
func (pnu *ProvisionedNetworkUpdate) AddProvisionedHosts(p ...*ProvisionedHost) *ProvisionedNetworkUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pnu.AddProvisionedHostIDs(ids...)
}

// Mutation returns the ProvisionedNetworkMutation object of the builder.
func (pnu *ProvisionedNetworkUpdate) Mutation() *ProvisionedNetworkMutation {
	return pnu.mutation
}

// ClearStatus clears all "status" edges to type Status.
func (pnu *ProvisionedNetworkUpdate) ClearStatus() *ProvisionedNetworkUpdate {
	pnu.mutation.ClearStatus()
	return pnu
}

// RemoveStatuIDs removes the status edge to Status by ids.
func (pnu *ProvisionedNetworkUpdate) RemoveStatuIDs(ids ...int) *ProvisionedNetworkUpdate {
	pnu.mutation.RemoveStatuIDs(ids...)
	return pnu
}

// RemoveStatus removes status edges to Status.
func (pnu *ProvisionedNetworkUpdate) RemoveStatus(s ...*Status) *ProvisionedNetworkUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return pnu.RemoveStatuIDs(ids...)
}

// ClearNetwork clears all "network" edges to type Network.
func (pnu *ProvisionedNetworkUpdate) ClearNetwork() *ProvisionedNetworkUpdate {
	pnu.mutation.ClearNetwork()
	return pnu
}

// RemoveNetworkIDs removes the network edge to Network by ids.
func (pnu *ProvisionedNetworkUpdate) RemoveNetworkIDs(ids ...int) *ProvisionedNetworkUpdate {
	pnu.mutation.RemoveNetworkIDs(ids...)
	return pnu
}

// RemoveNetwork removes network edges to Network.
func (pnu *ProvisionedNetworkUpdate) RemoveNetwork(n ...*Network) *ProvisionedNetworkUpdate {
	ids := make([]int, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return pnu.RemoveNetworkIDs(ids...)
}

// ClearBuild clears all "build" edges to type Build.
func (pnu *ProvisionedNetworkUpdate) ClearBuild() *ProvisionedNetworkUpdate {
	pnu.mutation.ClearBuild()
	return pnu
}

// RemoveBuildIDs removes the build edge to Build by ids.
func (pnu *ProvisionedNetworkUpdate) RemoveBuildIDs(ids ...int) *ProvisionedNetworkUpdate {
	pnu.mutation.RemoveBuildIDs(ids...)
	return pnu
}

// RemoveBuild removes build edges to Build.
func (pnu *ProvisionedNetworkUpdate) RemoveBuild(b ...*Build) *ProvisionedNetworkUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return pnu.RemoveBuildIDs(ids...)
}

// ClearProvisionedNetworkToTeam clears all "ProvisionedNetworkToTeam" edges to type Team.
func (pnu *ProvisionedNetworkUpdate) ClearProvisionedNetworkToTeam() *ProvisionedNetworkUpdate {
	pnu.mutation.ClearProvisionedNetworkToTeam()
	return pnu
}

// RemoveProvisionedNetworkToTeamIDs removes the ProvisionedNetworkToTeam edge to Team by ids.
func (pnu *ProvisionedNetworkUpdate) RemoveProvisionedNetworkToTeamIDs(ids ...int) *ProvisionedNetworkUpdate {
	pnu.mutation.RemoveProvisionedNetworkToTeamIDs(ids...)
	return pnu
}

// RemoveProvisionedNetworkToTeam removes ProvisionedNetworkToTeam edges to Team.
func (pnu *ProvisionedNetworkUpdate) RemoveProvisionedNetworkToTeam(t ...*Team) *ProvisionedNetworkUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return pnu.RemoveProvisionedNetworkToTeamIDs(ids...)
}

// ClearProvisionedHosts clears all "provisioned_hosts" edges to type ProvisionedHost.
func (pnu *ProvisionedNetworkUpdate) ClearProvisionedHosts() *ProvisionedNetworkUpdate {
	pnu.mutation.ClearProvisionedHosts()
	return pnu
}

// RemoveProvisionedHostIDs removes the provisioned_hosts edge to ProvisionedHost by ids.
func (pnu *ProvisionedNetworkUpdate) RemoveProvisionedHostIDs(ids ...int) *ProvisionedNetworkUpdate {
	pnu.mutation.RemoveProvisionedHostIDs(ids...)
	return pnu
}

// RemoveProvisionedHosts removes provisioned_hosts edges to ProvisionedHost.
func (pnu *ProvisionedNetworkUpdate) RemoveProvisionedHosts(p ...*ProvisionedHost) *ProvisionedNetworkUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pnu.RemoveProvisionedHostIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pnu *ProvisionedNetworkUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(pnu.hooks) == 0 {
		affected, err = pnu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProvisionedNetworkMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pnu.mutation = mutation
			affected, err = pnu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(pnu.hooks) - 1; i >= 0; i-- {
			mut = pnu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pnu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (pnu *ProvisionedNetworkUpdate) SaveX(ctx context.Context) int {
	affected, err := pnu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pnu *ProvisionedNetworkUpdate) Exec(ctx context.Context) error {
	_, err := pnu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pnu *ProvisionedNetworkUpdate) ExecX(ctx context.Context) {
	if err := pnu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (pnu *ProvisionedNetworkUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   provisionednetwork.Table,
			Columns: provisionednetwork.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: provisionednetwork.FieldID,
			},
		},
	}
	if ps := pnu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pnu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: provisionednetwork.FieldName,
		})
	}
	if value, ok := pnu.mutation.Cidr(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: provisionednetwork.FieldCidr,
		})
	}
	if pnu.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.StatusTable,
			Columns: []string{provisionednetwork.StatusColumn},
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
	if nodes := pnu.mutation.RemovedStatusIDs(); len(nodes) > 0 && !pnu.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.StatusTable,
			Columns: []string{provisionednetwork.StatusColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pnu.mutation.StatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.StatusTable,
			Columns: []string{provisionednetwork.StatusColumn},
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
	if pnu.mutation.NetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.NetworkTable,
			Columns: []string{provisionednetwork.NetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: network.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pnu.mutation.RemovedNetworkIDs(); len(nodes) > 0 && !pnu.mutation.NetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.NetworkTable,
			Columns: []string{provisionednetwork.NetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: network.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pnu.mutation.NetworkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.NetworkTable,
			Columns: []string{provisionednetwork.NetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: network.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pnu.mutation.BuildCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.BuildTable,
			Columns: provisionednetwork.BuildPrimaryKey,
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
	if nodes := pnu.mutation.RemovedBuildIDs(); len(nodes) > 0 && !pnu.mutation.BuildCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.BuildTable,
			Columns: provisionednetwork.BuildPrimaryKey,
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
	if nodes := pnu.mutation.BuildIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.BuildTable,
			Columns: provisionednetwork.BuildPrimaryKey,
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
	if pnu.mutation.ProvisionedNetworkToTeamCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionednetwork.ProvisionedNetworkToTeamTable,
			Columns: provisionednetwork.ProvisionedNetworkToTeamPrimaryKey,
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
	if nodes := pnu.mutation.RemovedProvisionedNetworkToTeamIDs(); len(nodes) > 0 && !pnu.mutation.ProvisionedNetworkToTeamCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionednetwork.ProvisionedNetworkToTeamTable,
			Columns: provisionednetwork.ProvisionedNetworkToTeamPrimaryKey,
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
	if nodes := pnu.mutation.ProvisionedNetworkToTeamIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionednetwork.ProvisionedNetworkToTeamTable,
			Columns: provisionednetwork.ProvisionedNetworkToTeamPrimaryKey,
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
	if pnu.mutation.ProvisionedHostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.ProvisionedHostsTable,
			Columns: provisionednetwork.ProvisionedHostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionedhost.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pnu.mutation.RemovedProvisionedHostsIDs(); len(nodes) > 0 && !pnu.mutation.ProvisionedHostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.ProvisionedHostsTable,
			Columns: provisionednetwork.ProvisionedHostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionedhost.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pnu.mutation.ProvisionedHostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.ProvisionedHostsTable,
			Columns: provisionednetwork.ProvisionedHostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionedhost.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pnu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{provisionednetwork.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ProvisionedNetworkUpdateOne is the builder for updating a single ProvisionedNetwork entity.
type ProvisionedNetworkUpdateOne struct {
	config
	hooks    []Hook
	mutation *ProvisionedNetworkMutation
}

// SetName sets the name field.
func (pnuo *ProvisionedNetworkUpdateOne) SetName(s string) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.SetName(s)
	return pnuo
}

// SetCidr sets the cidr field.
func (pnuo *ProvisionedNetworkUpdateOne) SetCidr(s string) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.SetCidr(s)
	return pnuo
}

// AddStatuIDs adds the status edge to Status by ids.
func (pnuo *ProvisionedNetworkUpdateOne) AddStatuIDs(ids ...int) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.AddStatuIDs(ids...)
	return pnuo
}

// AddStatus adds the status edges to Status.
func (pnuo *ProvisionedNetworkUpdateOne) AddStatus(s ...*Status) *ProvisionedNetworkUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return pnuo.AddStatuIDs(ids...)
}

// AddNetworkIDs adds the network edge to Network by ids.
func (pnuo *ProvisionedNetworkUpdateOne) AddNetworkIDs(ids ...int) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.AddNetworkIDs(ids...)
	return pnuo
}

// AddNetwork adds the network edges to Network.
func (pnuo *ProvisionedNetworkUpdateOne) AddNetwork(n ...*Network) *ProvisionedNetworkUpdateOne {
	ids := make([]int, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return pnuo.AddNetworkIDs(ids...)
}

// AddBuildIDs adds the build edge to Build by ids.
func (pnuo *ProvisionedNetworkUpdateOne) AddBuildIDs(ids ...int) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.AddBuildIDs(ids...)
	return pnuo
}

// AddBuild adds the build edges to Build.
func (pnuo *ProvisionedNetworkUpdateOne) AddBuild(b ...*Build) *ProvisionedNetworkUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return pnuo.AddBuildIDs(ids...)
}

// AddProvisionedNetworkToTeamIDs adds the ProvisionedNetworkToTeam edge to Team by ids.
func (pnuo *ProvisionedNetworkUpdateOne) AddProvisionedNetworkToTeamIDs(ids ...int) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.AddProvisionedNetworkToTeamIDs(ids...)
	return pnuo
}

// AddProvisionedNetworkToTeam adds the ProvisionedNetworkToTeam edges to Team.
func (pnuo *ProvisionedNetworkUpdateOne) AddProvisionedNetworkToTeam(t ...*Team) *ProvisionedNetworkUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return pnuo.AddProvisionedNetworkToTeamIDs(ids...)
}

// AddProvisionedHostIDs adds the provisioned_hosts edge to ProvisionedHost by ids.
func (pnuo *ProvisionedNetworkUpdateOne) AddProvisionedHostIDs(ids ...int) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.AddProvisionedHostIDs(ids...)
	return pnuo
}

// AddProvisionedHosts adds the provisioned_hosts edges to ProvisionedHost.
func (pnuo *ProvisionedNetworkUpdateOne) AddProvisionedHosts(p ...*ProvisionedHost) *ProvisionedNetworkUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pnuo.AddProvisionedHostIDs(ids...)
}

// Mutation returns the ProvisionedNetworkMutation object of the builder.
func (pnuo *ProvisionedNetworkUpdateOne) Mutation() *ProvisionedNetworkMutation {
	return pnuo.mutation
}

// ClearStatus clears all "status" edges to type Status.
func (pnuo *ProvisionedNetworkUpdateOne) ClearStatus() *ProvisionedNetworkUpdateOne {
	pnuo.mutation.ClearStatus()
	return pnuo
}

// RemoveStatuIDs removes the status edge to Status by ids.
func (pnuo *ProvisionedNetworkUpdateOne) RemoveStatuIDs(ids ...int) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.RemoveStatuIDs(ids...)
	return pnuo
}

// RemoveStatus removes status edges to Status.
func (pnuo *ProvisionedNetworkUpdateOne) RemoveStatus(s ...*Status) *ProvisionedNetworkUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return pnuo.RemoveStatuIDs(ids...)
}

// ClearNetwork clears all "network" edges to type Network.
func (pnuo *ProvisionedNetworkUpdateOne) ClearNetwork() *ProvisionedNetworkUpdateOne {
	pnuo.mutation.ClearNetwork()
	return pnuo
}

// RemoveNetworkIDs removes the network edge to Network by ids.
func (pnuo *ProvisionedNetworkUpdateOne) RemoveNetworkIDs(ids ...int) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.RemoveNetworkIDs(ids...)
	return pnuo
}

// RemoveNetwork removes network edges to Network.
func (pnuo *ProvisionedNetworkUpdateOne) RemoveNetwork(n ...*Network) *ProvisionedNetworkUpdateOne {
	ids := make([]int, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return pnuo.RemoveNetworkIDs(ids...)
}

// ClearBuild clears all "build" edges to type Build.
func (pnuo *ProvisionedNetworkUpdateOne) ClearBuild() *ProvisionedNetworkUpdateOne {
	pnuo.mutation.ClearBuild()
	return pnuo
}

// RemoveBuildIDs removes the build edge to Build by ids.
func (pnuo *ProvisionedNetworkUpdateOne) RemoveBuildIDs(ids ...int) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.RemoveBuildIDs(ids...)
	return pnuo
}

// RemoveBuild removes build edges to Build.
func (pnuo *ProvisionedNetworkUpdateOne) RemoveBuild(b ...*Build) *ProvisionedNetworkUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return pnuo.RemoveBuildIDs(ids...)
}

// ClearProvisionedNetworkToTeam clears all "ProvisionedNetworkToTeam" edges to type Team.
func (pnuo *ProvisionedNetworkUpdateOne) ClearProvisionedNetworkToTeam() *ProvisionedNetworkUpdateOne {
	pnuo.mutation.ClearProvisionedNetworkToTeam()
	return pnuo
}

// RemoveProvisionedNetworkToTeamIDs removes the ProvisionedNetworkToTeam edge to Team by ids.
func (pnuo *ProvisionedNetworkUpdateOne) RemoveProvisionedNetworkToTeamIDs(ids ...int) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.RemoveProvisionedNetworkToTeamIDs(ids...)
	return pnuo
}

// RemoveProvisionedNetworkToTeam removes ProvisionedNetworkToTeam edges to Team.
func (pnuo *ProvisionedNetworkUpdateOne) RemoveProvisionedNetworkToTeam(t ...*Team) *ProvisionedNetworkUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return pnuo.RemoveProvisionedNetworkToTeamIDs(ids...)
}

// ClearProvisionedHosts clears all "provisioned_hosts" edges to type ProvisionedHost.
func (pnuo *ProvisionedNetworkUpdateOne) ClearProvisionedHosts() *ProvisionedNetworkUpdateOne {
	pnuo.mutation.ClearProvisionedHosts()
	return pnuo
}

// RemoveProvisionedHostIDs removes the provisioned_hosts edge to ProvisionedHost by ids.
func (pnuo *ProvisionedNetworkUpdateOne) RemoveProvisionedHostIDs(ids ...int) *ProvisionedNetworkUpdateOne {
	pnuo.mutation.RemoveProvisionedHostIDs(ids...)
	return pnuo
}

// RemoveProvisionedHosts removes provisioned_hosts edges to ProvisionedHost.
func (pnuo *ProvisionedNetworkUpdateOne) RemoveProvisionedHosts(p ...*ProvisionedHost) *ProvisionedNetworkUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pnuo.RemoveProvisionedHostIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (pnuo *ProvisionedNetworkUpdateOne) Save(ctx context.Context) (*ProvisionedNetwork, error) {
	var (
		err  error
		node *ProvisionedNetwork
	)
	if len(pnuo.hooks) == 0 {
		node, err = pnuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProvisionedNetworkMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pnuo.mutation = mutation
			node, err = pnuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(pnuo.hooks) - 1; i >= 0; i-- {
			mut = pnuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pnuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (pnuo *ProvisionedNetworkUpdateOne) SaveX(ctx context.Context) *ProvisionedNetwork {
	node, err := pnuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (pnuo *ProvisionedNetworkUpdateOne) Exec(ctx context.Context) error {
	_, err := pnuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pnuo *ProvisionedNetworkUpdateOne) ExecX(ctx context.Context) {
	if err := pnuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (pnuo *ProvisionedNetworkUpdateOne) sqlSave(ctx context.Context) (_node *ProvisionedNetwork, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   provisionednetwork.Table,
			Columns: provisionednetwork.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: provisionednetwork.FieldID,
			},
		},
	}
	id, ok := pnuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing ProvisionedNetwork.ID for update")}
	}
	_spec.Node.ID.Value = id
	if value, ok := pnuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: provisionednetwork.FieldName,
		})
	}
	if value, ok := pnuo.mutation.Cidr(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: provisionednetwork.FieldCidr,
		})
	}
	if pnuo.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.StatusTable,
			Columns: []string{provisionednetwork.StatusColumn},
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
	if nodes := pnuo.mutation.RemovedStatusIDs(); len(nodes) > 0 && !pnuo.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.StatusTable,
			Columns: []string{provisionednetwork.StatusColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pnuo.mutation.StatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.StatusTable,
			Columns: []string{provisionednetwork.StatusColumn},
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
	if pnuo.mutation.NetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.NetworkTable,
			Columns: []string{provisionednetwork.NetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: network.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pnuo.mutation.RemovedNetworkIDs(); len(nodes) > 0 && !pnuo.mutation.NetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.NetworkTable,
			Columns: []string{provisionednetwork.NetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: network.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pnuo.mutation.NetworkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionednetwork.NetworkTable,
			Columns: []string{provisionednetwork.NetworkColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: network.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pnuo.mutation.BuildCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.BuildTable,
			Columns: provisionednetwork.BuildPrimaryKey,
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
	if nodes := pnuo.mutation.RemovedBuildIDs(); len(nodes) > 0 && !pnuo.mutation.BuildCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.BuildTable,
			Columns: provisionednetwork.BuildPrimaryKey,
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
	if nodes := pnuo.mutation.BuildIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.BuildTable,
			Columns: provisionednetwork.BuildPrimaryKey,
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
	if pnuo.mutation.ProvisionedNetworkToTeamCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionednetwork.ProvisionedNetworkToTeamTable,
			Columns: provisionednetwork.ProvisionedNetworkToTeamPrimaryKey,
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
	if nodes := pnuo.mutation.RemovedProvisionedNetworkToTeamIDs(); len(nodes) > 0 && !pnuo.mutation.ProvisionedNetworkToTeamCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionednetwork.ProvisionedNetworkToTeamTable,
			Columns: provisionednetwork.ProvisionedNetworkToTeamPrimaryKey,
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
	if nodes := pnuo.mutation.ProvisionedNetworkToTeamIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionednetwork.ProvisionedNetworkToTeamTable,
			Columns: provisionednetwork.ProvisionedNetworkToTeamPrimaryKey,
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
	if pnuo.mutation.ProvisionedHostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.ProvisionedHostsTable,
			Columns: provisionednetwork.ProvisionedHostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionedhost.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pnuo.mutation.RemovedProvisionedHostsIDs(); len(nodes) > 0 && !pnuo.mutation.ProvisionedHostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.ProvisionedHostsTable,
			Columns: provisionednetwork.ProvisionedHostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionedhost.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pnuo.mutation.ProvisionedHostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionednetwork.ProvisionedHostsTable,
			Columns: provisionednetwork.ProvisionedHostsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisionedhost.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &ProvisionedNetwork{config: pnuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues()
	if err = sqlgraph.UpdateNode(ctx, pnuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{provisionednetwork.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}
