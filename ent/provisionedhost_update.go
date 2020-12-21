// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/gen0cide/laforge/ent/agentstatus"
	"github.com/gen0cide/laforge/ent/host"
	"github.com/gen0cide/laforge/ent/predicate"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	"github.com/gen0cide/laforge/ent/provisionednetwork"
	"github.com/gen0cide/laforge/ent/provisioningstep"
	"github.com/gen0cide/laforge/ent/status"
)

// ProvisionedHostUpdate is the builder for updating ProvisionedHost entities.
type ProvisionedHostUpdate struct {
	config
	hooks    []Hook
	mutation *ProvisionedHostMutation
}

// Where adds a new predicate for the builder.
func (phu *ProvisionedHostUpdate) Where(ps ...predicate.ProvisionedHost) *ProvisionedHostUpdate {
	phu.mutation.predicates = append(phu.mutation.predicates, ps...)
	return phu
}

// SetSubnetIP sets the subnet_ip field.
func (phu *ProvisionedHostUpdate) SetSubnetIP(s string) *ProvisionedHostUpdate {
	phu.mutation.SetSubnetIP(s)
	return phu
}

// AddStatuIDs adds the status edge to Status by ids.
func (phu *ProvisionedHostUpdate) AddStatuIDs(ids ...int) *ProvisionedHostUpdate {
	phu.mutation.AddStatuIDs(ids...)
	return phu
}

// AddStatus adds the status edges to Status.
func (phu *ProvisionedHostUpdate) AddStatus(s ...*Status) *ProvisionedHostUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return phu.AddStatuIDs(ids...)
}

// AddProvisionedNetworkIDs adds the provisioned_network edge to ProvisionedNetwork by ids.
func (phu *ProvisionedHostUpdate) AddProvisionedNetworkIDs(ids ...int) *ProvisionedHostUpdate {
	phu.mutation.AddProvisionedNetworkIDs(ids...)
	return phu
}

// AddProvisionedNetwork adds the provisioned_network edges to ProvisionedNetwork.
func (phu *ProvisionedHostUpdate) AddProvisionedNetwork(p ...*ProvisionedNetwork) *ProvisionedHostUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return phu.AddProvisionedNetworkIDs(ids...)
}

// AddHostIDs adds the host edge to Host by ids.
func (phu *ProvisionedHostUpdate) AddHostIDs(ids ...int) *ProvisionedHostUpdate {
	phu.mutation.AddHostIDs(ids...)
	return phu
}

// AddHost adds the host edges to Host.
func (phu *ProvisionedHostUpdate) AddHost(h ...*Host) *ProvisionedHostUpdate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return phu.AddHostIDs(ids...)
}

// AddProvisionedStepIDs adds the provisioned_steps edge to ProvisioningStep by ids.
func (phu *ProvisionedHostUpdate) AddProvisionedStepIDs(ids ...int) *ProvisionedHostUpdate {
	phu.mutation.AddProvisionedStepIDs(ids...)
	return phu
}

// AddProvisionedSteps adds the provisioned_steps edges to ProvisioningStep.
func (phu *ProvisionedHostUpdate) AddProvisionedSteps(p ...*ProvisioningStep) *ProvisionedHostUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return phu.AddProvisionedStepIDs(ids...)
}

// AddAgentStatuIDs adds the agent_status edge to AgentStatus by ids.
func (phu *ProvisionedHostUpdate) AddAgentStatuIDs(ids ...int) *ProvisionedHostUpdate {
	phu.mutation.AddAgentStatuIDs(ids...)
	return phu
}

// AddAgentStatus adds the agent_status edges to AgentStatus.
func (phu *ProvisionedHostUpdate) AddAgentStatus(a ...*AgentStatus) *ProvisionedHostUpdate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return phu.AddAgentStatuIDs(ids...)
}

// Mutation returns the ProvisionedHostMutation object of the builder.
func (phu *ProvisionedHostUpdate) Mutation() *ProvisionedHostMutation {
	return phu.mutation
}

// ClearStatus clears all "status" edges to type Status.
func (phu *ProvisionedHostUpdate) ClearStatus() *ProvisionedHostUpdate {
	phu.mutation.ClearStatus()
	return phu
}

// RemoveStatuIDs removes the status edge to Status by ids.
func (phu *ProvisionedHostUpdate) RemoveStatuIDs(ids ...int) *ProvisionedHostUpdate {
	phu.mutation.RemoveStatuIDs(ids...)
	return phu
}

// RemoveStatus removes status edges to Status.
func (phu *ProvisionedHostUpdate) RemoveStatus(s ...*Status) *ProvisionedHostUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return phu.RemoveStatuIDs(ids...)
}

// ClearProvisionedNetwork clears all "provisioned_network" edges to type ProvisionedNetwork.
func (phu *ProvisionedHostUpdate) ClearProvisionedNetwork() *ProvisionedHostUpdate {
	phu.mutation.ClearProvisionedNetwork()
	return phu
}

// RemoveProvisionedNetworkIDs removes the provisioned_network edge to ProvisionedNetwork by ids.
func (phu *ProvisionedHostUpdate) RemoveProvisionedNetworkIDs(ids ...int) *ProvisionedHostUpdate {
	phu.mutation.RemoveProvisionedNetworkIDs(ids...)
	return phu
}

// RemoveProvisionedNetwork removes provisioned_network edges to ProvisionedNetwork.
func (phu *ProvisionedHostUpdate) RemoveProvisionedNetwork(p ...*ProvisionedNetwork) *ProvisionedHostUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return phu.RemoveProvisionedNetworkIDs(ids...)
}

// ClearHost clears all "host" edges to type Host.
func (phu *ProvisionedHostUpdate) ClearHost() *ProvisionedHostUpdate {
	phu.mutation.ClearHost()
	return phu
}

// RemoveHostIDs removes the host edge to Host by ids.
func (phu *ProvisionedHostUpdate) RemoveHostIDs(ids ...int) *ProvisionedHostUpdate {
	phu.mutation.RemoveHostIDs(ids...)
	return phu
}

// RemoveHost removes host edges to Host.
func (phu *ProvisionedHostUpdate) RemoveHost(h ...*Host) *ProvisionedHostUpdate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return phu.RemoveHostIDs(ids...)
}

// ClearProvisionedSteps clears all "provisioned_steps" edges to type ProvisioningStep.
func (phu *ProvisionedHostUpdate) ClearProvisionedSteps() *ProvisionedHostUpdate {
	phu.mutation.ClearProvisionedSteps()
	return phu
}

// RemoveProvisionedStepIDs removes the provisioned_steps edge to ProvisioningStep by ids.
func (phu *ProvisionedHostUpdate) RemoveProvisionedStepIDs(ids ...int) *ProvisionedHostUpdate {
	phu.mutation.RemoveProvisionedStepIDs(ids...)
	return phu
}

// RemoveProvisionedSteps removes provisioned_steps edges to ProvisioningStep.
func (phu *ProvisionedHostUpdate) RemoveProvisionedSteps(p ...*ProvisioningStep) *ProvisionedHostUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return phu.RemoveProvisionedStepIDs(ids...)
}

// ClearAgentStatus clears all "agent_status" edges to type AgentStatus.
func (phu *ProvisionedHostUpdate) ClearAgentStatus() *ProvisionedHostUpdate {
	phu.mutation.ClearAgentStatus()
	return phu
}

// RemoveAgentStatuIDs removes the agent_status edge to AgentStatus by ids.
func (phu *ProvisionedHostUpdate) RemoveAgentStatuIDs(ids ...int) *ProvisionedHostUpdate {
	phu.mutation.RemoveAgentStatuIDs(ids...)
	return phu
}

// RemoveAgentStatus removes agent_status edges to AgentStatus.
func (phu *ProvisionedHostUpdate) RemoveAgentStatus(a ...*AgentStatus) *ProvisionedHostUpdate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return phu.RemoveAgentStatuIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (phu *ProvisionedHostUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(phu.hooks) == 0 {
		affected, err = phu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProvisionedHostMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			phu.mutation = mutation
			affected, err = phu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(phu.hooks) - 1; i >= 0; i-- {
			mut = phu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, phu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (phu *ProvisionedHostUpdate) SaveX(ctx context.Context) int {
	affected, err := phu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (phu *ProvisionedHostUpdate) Exec(ctx context.Context) error {
	_, err := phu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (phu *ProvisionedHostUpdate) ExecX(ctx context.Context) {
	if err := phu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (phu *ProvisionedHostUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   provisionedhost.Table,
			Columns: provisionedhost.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: provisionedhost.FieldID,
			},
		},
	}
	if ps := phu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := phu.mutation.SubnetIP(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: provisionedhost.FieldSubnetIP,
		})
	}
	if phu.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.StatusTable,
			Columns: []string{provisionedhost.StatusColumn},
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
	if nodes := phu.mutation.RemovedStatusIDs(); len(nodes) > 0 && !phu.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.StatusTable,
			Columns: []string{provisionedhost.StatusColumn},
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
	if nodes := phu.mutation.StatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.StatusTable,
			Columns: []string{provisionedhost.StatusColumn},
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
	if phu.mutation.ProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionedhost.ProvisionedNetworkTable,
			Columns: provisionedhost.ProvisionedNetworkPrimaryKey,
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
	if nodes := phu.mutation.RemovedProvisionedNetworkIDs(); len(nodes) > 0 && !phu.mutation.ProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionedhost.ProvisionedNetworkTable,
			Columns: provisionedhost.ProvisionedNetworkPrimaryKey,
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
	if nodes := phu.mutation.ProvisionedNetworkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionedhost.ProvisionedNetworkTable,
			Columns: provisionedhost.ProvisionedNetworkPrimaryKey,
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
	if phu.mutation.HostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.HostTable,
			Columns: []string{provisionedhost.HostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: host.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phu.mutation.RemovedHostIDs(); len(nodes) > 0 && !phu.mutation.HostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.HostTable,
			Columns: []string{provisionedhost.HostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: host.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phu.mutation.HostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.HostTable,
			Columns: []string{provisionedhost.HostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: host.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if phu.mutation.ProvisionedStepsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.ProvisionedStepsTable,
			Columns: provisionedhost.ProvisionedStepsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisioningstep.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phu.mutation.RemovedProvisionedStepsIDs(); len(nodes) > 0 && !phu.mutation.ProvisionedStepsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.ProvisionedStepsTable,
			Columns: provisionedhost.ProvisionedStepsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisioningstep.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phu.mutation.ProvisionedStepsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.ProvisionedStepsTable,
			Columns: provisionedhost.ProvisionedStepsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisioningstep.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if phu.mutation.AgentStatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.AgentStatusTable,
			Columns: provisionedhost.AgentStatusPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: agentstatus.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phu.mutation.RemovedAgentStatusIDs(); len(nodes) > 0 && !phu.mutation.AgentStatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.AgentStatusTable,
			Columns: provisionedhost.AgentStatusPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: agentstatus.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phu.mutation.AgentStatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.AgentStatusTable,
			Columns: provisionedhost.AgentStatusPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: agentstatus.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, phu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{provisionedhost.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ProvisionedHostUpdateOne is the builder for updating a single ProvisionedHost entity.
type ProvisionedHostUpdateOne struct {
	config
	hooks    []Hook
	mutation *ProvisionedHostMutation
}

// SetSubnetIP sets the subnet_ip field.
func (phuo *ProvisionedHostUpdateOne) SetSubnetIP(s string) *ProvisionedHostUpdateOne {
	phuo.mutation.SetSubnetIP(s)
	return phuo
}

// AddStatuIDs adds the status edge to Status by ids.
func (phuo *ProvisionedHostUpdateOne) AddStatuIDs(ids ...int) *ProvisionedHostUpdateOne {
	phuo.mutation.AddStatuIDs(ids...)
	return phuo
}

// AddStatus adds the status edges to Status.
func (phuo *ProvisionedHostUpdateOne) AddStatus(s ...*Status) *ProvisionedHostUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return phuo.AddStatuIDs(ids...)
}

// AddProvisionedNetworkIDs adds the provisioned_network edge to ProvisionedNetwork by ids.
func (phuo *ProvisionedHostUpdateOne) AddProvisionedNetworkIDs(ids ...int) *ProvisionedHostUpdateOne {
	phuo.mutation.AddProvisionedNetworkIDs(ids...)
	return phuo
}

// AddProvisionedNetwork adds the provisioned_network edges to ProvisionedNetwork.
func (phuo *ProvisionedHostUpdateOne) AddProvisionedNetwork(p ...*ProvisionedNetwork) *ProvisionedHostUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return phuo.AddProvisionedNetworkIDs(ids...)
}

// AddHostIDs adds the host edge to Host by ids.
func (phuo *ProvisionedHostUpdateOne) AddHostIDs(ids ...int) *ProvisionedHostUpdateOne {
	phuo.mutation.AddHostIDs(ids...)
	return phuo
}

// AddHost adds the host edges to Host.
func (phuo *ProvisionedHostUpdateOne) AddHost(h ...*Host) *ProvisionedHostUpdateOne {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return phuo.AddHostIDs(ids...)
}

// AddProvisionedStepIDs adds the provisioned_steps edge to ProvisioningStep by ids.
func (phuo *ProvisionedHostUpdateOne) AddProvisionedStepIDs(ids ...int) *ProvisionedHostUpdateOne {
	phuo.mutation.AddProvisionedStepIDs(ids...)
	return phuo
}

// AddProvisionedSteps adds the provisioned_steps edges to ProvisioningStep.
func (phuo *ProvisionedHostUpdateOne) AddProvisionedSteps(p ...*ProvisioningStep) *ProvisionedHostUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return phuo.AddProvisionedStepIDs(ids...)
}

// AddAgentStatuIDs adds the agent_status edge to AgentStatus by ids.
func (phuo *ProvisionedHostUpdateOne) AddAgentStatuIDs(ids ...int) *ProvisionedHostUpdateOne {
	phuo.mutation.AddAgentStatuIDs(ids...)
	return phuo
}

// AddAgentStatus adds the agent_status edges to AgentStatus.
func (phuo *ProvisionedHostUpdateOne) AddAgentStatus(a ...*AgentStatus) *ProvisionedHostUpdateOne {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return phuo.AddAgentStatuIDs(ids...)
}

// Mutation returns the ProvisionedHostMutation object of the builder.
func (phuo *ProvisionedHostUpdateOne) Mutation() *ProvisionedHostMutation {
	return phuo.mutation
}

// ClearStatus clears all "status" edges to type Status.
func (phuo *ProvisionedHostUpdateOne) ClearStatus() *ProvisionedHostUpdateOne {
	phuo.mutation.ClearStatus()
	return phuo
}

// RemoveStatuIDs removes the status edge to Status by ids.
func (phuo *ProvisionedHostUpdateOne) RemoveStatuIDs(ids ...int) *ProvisionedHostUpdateOne {
	phuo.mutation.RemoveStatuIDs(ids...)
	return phuo
}

// RemoveStatus removes status edges to Status.
func (phuo *ProvisionedHostUpdateOne) RemoveStatus(s ...*Status) *ProvisionedHostUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return phuo.RemoveStatuIDs(ids...)
}

// ClearProvisionedNetwork clears all "provisioned_network" edges to type ProvisionedNetwork.
func (phuo *ProvisionedHostUpdateOne) ClearProvisionedNetwork() *ProvisionedHostUpdateOne {
	phuo.mutation.ClearProvisionedNetwork()
	return phuo
}

// RemoveProvisionedNetworkIDs removes the provisioned_network edge to ProvisionedNetwork by ids.
func (phuo *ProvisionedHostUpdateOne) RemoveProvisionedNetworkIDs(ids ...int) *ProvisionedHostUpdateOne {
	phuo.mutation.RemoveProvisionedNetworkIDs(ids...)
	return phuo
}

// RemoveProvisionedNetwork removes provisioned_network edges to ProvisionedNetwork.
func (phuo *ProvisionedHostUpdateOne) RemoveProvisionedNetwork(p ...*ProvisionedNetwork) *ProvisionedHostUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return phuo.RemoveProvisionedNetworkIDs(ids...)
}

// ClearHost clears all "host" edges to type Host.
func (phuo *ProvisionedHostUpdateOne) ClearHost() *ProvisionedHostUpdateOne {
	phuo.mutation.ClearHost()
	return phuo
}

// RemoveHostIDs removes the host edge to Host by ids.
func (phuo *ProvisionedHostUpdateOne) RemoveHostIDs(ids ...int) *ProvisionedHostUpdateOne {
	phuo.mutation.RemoveHostIDs(ids...)
	return phuo
}

// RemoveHost removes host edges to Host.
func (phuo *ProvisionedHostUpdateOne) RemoveHost(h ...*Host) *ProvisionedHostUpdateOne {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return phuo.RemoveHostIDs(ids...)
}

// ClearProvisionedSteps clears all "provisioned_steps" edges to type ProvisioningStep.
func (phuo *ProvisionedHostUpdateOne) ClearProvisionedSteps() *ProvisionedHostUpdateOne {
	phuo.mutation.ClearProvisionedSteps()
	return phuo
}

// RemoveProvisionedStepIDs removes the provisioned_steps edge to ProvisioningStep by ids.
func (phuo *ProvisionedHostUpdateOne) RemoveProvisionedStepIDs(ids ...int) *ProvisionedHostUpdateOne {
	phuo.mutation.RemoveProvisionedStepIDs(ids...)
	return phuo
}

// RemoveProvisionedSteps removes provisioned_steps edges to ProvisioningStep.
func (phuo *ProvisionedHostUpdateOne) RemoveProvisionedSteps(p ...*ProvisioningStep) *ProvisionedHostUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return phuo.RemoveProvisionedStepIDs(ids...)
}

// ClearAgentStatus clears all "agent_status" edges to type AgentStatus.
func (phuo *ProvisionedHostUpdateOne) ClearAgentStatus() *ProvisionedHostUpdateOne {
	phuo.mutation.ClearAgentStatus()
	return phuo
}

// RemoveAgentStatuIDs removes the agent_status edge to AgentStatus by ids.
func (phuo *ProvisionedHostUpdateOne) RemoveAgentStatuIDs(ids ...int) *ProvisionedHostUpdateOne {
	phuo.mutation.RemoveAgentStatuIDs(ids...)
	return phuo
}

// RemoveAgentStatus removes agent_status edges to AgentStatus.
func (phuo *ProvisionedHostUpdateOne) RemoveAgentStatus(a ...*AgentStatus) *ProvisionedHostUpdateOne {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return phuo.RemoveAgentStatuIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (phuo *ProvisionedHostUpdateOne) Save(ctx context.Context) (*ProvisionedHost, error) {
	var (
		err  error
		node *ProvisionedHost
	)
	if len(phuo.hooks) == 0 {
		node, err = phuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProvisionedHostMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			phuo.mutation = mutation
			node, err = phuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(phuo.hooks) - 1; i >= 0; i-- {
			mut = phuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, phuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (phuo *ProvisionedHostUpdateOne) SaveX(ctx context.Context) *ProvisionedHost {
	node, err := phuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (phuo *ProvisionedHostUpdateOne) Exec(ctx context.Context) error {
	_, err := phuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (phuo *ProvisionedHostUpdateOne) ExecX(ctx context.Context) {
	if err := phuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (phuo *ProvisionedHostUpdateOne) sqlSave(ctx context.Context) (_node *ProvisionedHost, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   provisionedhost.Table,
			Columns: provisionedhost.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: provisionedhost.FieldID,
			},
		},
	}
	id, ok := phuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing ProvisionedHost.ID for update")}
	}
	_spec.Node.ID.Value = id
	if value, ok := phuo.mutation.SubnetIP(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: provisionedhost.FieldSubnetIP,
		})
	}
	if phuo.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.StatusTable,
			Columns: []string{provisionedhost.StatusColumn},
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
	if nodes := phuo.mutation.RemovedStatusIDs(); len(nodes) > 0 && !phuo.mutation.StatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.StatusTable,
			Columns: []string{provisionedhost.StatusColumn},
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
	if nodes := phuo.mutation.StatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.StatusTable,
			Columns: []string{provisionedhost.StatusColumn},
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
	if phuo.mutation.ProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionedhost.ProvisionedNetworkTable,
			Columns: provisionedhost.ProvisionedNetworkPrimaryKey,
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
	if nodes := phuo.mutation.RemovedProvisionedNetworkIDs(); len(nodes) > 0 && !phuo.mutation.ProvisionedNetworkCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionedhost.ProvisionedNetworkTable,
			Columns: provisionedhost.ProvisionedNetworkPrimaryKey,
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
	if nodes := phuo.mutation.ProvisionedNetworkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisionedhost.ProvisionedNetworkTable,
			Columns: provisionedhost.ProvisionedNetworkPrimaryKey,
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
	if phuo.mutation.HostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.HostTable,
			Columns: []string{provisionedhost.HostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: host.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phuo.mutation.RemovedHostIDs(); len(nodes) > 0 && !phuo.mutation.HostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.HostTable,
			Columns: []string{provisionedhost.HostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: host.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phuo.mutation.HostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisionedhost.HostTable,
			Columns: []string{provisionedhost.HostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: host.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if phuo.mutation.ProvisionedStepsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.ProvisionedStepsTable,
			Columns: provisionedhost.ProvisionedStepsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisioningstep.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phuo.mutation.RemovedProvisionedStepsIDs(); len(nodes) > 0 && !phuo.mutation.ProvisionedStepsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.ProvisionedStepsTable,
			Columns: provisionedhost.ProvisionedStepsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisioningstep.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phuo.mutation.ProvisionedStepsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.ProvisionedStepsTable,
			Columns: provisionedhost.ProvisionedStepsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: provisioningstep.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if phuo.mutation.AgentStatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.AgentStatusTable,
			Columns: provisionedhost.AgentStatusPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: agentstatus.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phuo.mutation.RemovedAgentStatusIDs(); len(nodes) > 0 && !phuo.mutation.AgentStatusCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.AgentStatusTable,
			Columns: provisionedhost.AgentStatusPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: agentstatus.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := phuo.mutation.AgentStatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   provisionedhost.AgentStatusTable,
			Columns: provisionedhost.AgentStatusPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: agentstatus.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &ProvisionedHost{config: phuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues()
	if err = sqlgraph.UpdateNode(ctx, phuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{provisionedhost.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}
