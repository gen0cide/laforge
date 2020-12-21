// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/gen0cide/laforge/ent/agentstatus"
	"github.com/gen0cide/laforge/ent/predicate"
	"github.com/gen0cide/laforge/ent/provisionedhost"
)

// AgentStatusUpdate is the builder for updating AgentStatus entities.
type AgentStatusUpdate struct {
	config
	hooks    []Hook
	mutation *AgentStatusMutation
}

// Where adds a new predicate for the builder.
func (asu *AgentStatusUpdate) Where(ps ...predicate.AgentStatus) *AgentStatusUpdate {
	asu.mutation.predicates = append(asu.mutation.predicates, ps...)
	return asu
}

// SetClientID sets the ClientID field.
func (asu *AgentStatusUpdate) SetClientID(s string) *AgentStatusUpdate {
	asu.mutation.SetClientID(s)
	return asu
}

// SetHostname sets the Hostname field.
func (asu *AgentStatusUpdate) SetHostname(s string) *AgentStatusUpdate {
	asu.mutation.SetHostname(s)
	return asu
}

// SetUpTime sets the UpTime field.
func (asu *AgentStatusUpdate) SetUpTime(i int) *AgentStatusUpdate {
	asu.mutation.ResetUpTime()
	asu.mutation.SetUpTime(i)
	return asu
}

// AddUpTime adds i to UpTime.
func (asu *AgentStatusUpdate) AddUpTime(i int) *AgentStatusUpdate {
	asu.mutation.AddUpTime(i)
	return asu
}

// SetBootTime sets the BootTime field.
func (asu *AgentStatusUpdate) SetBootTime(i int) *AgentStatusUpdate {
	asu.mutation.ResetBootTime()
	asu.mutation.SetBootTime(i)
	return asu
}

// AddBootTime adds i to BootTime.
func (asu *AgentStatusUpdate) AddBootTime(i int) *AgentStatusUpdate {
	asu.mutation.AddBootTime(i)
	return asu
}

// SetNumProcs sets the NumProcs field.
func (asu *AgentStatusUpdate) SetNumProcs(i int) *AgentStatusUpdate {
	asu.mutation.ResetNumProcs()
	asu.mutation.SetNumProcs(i)
	return asu
}

// AddNumProcs adds i to NumProcs.
func (asu *AgentStatusUpdate) AddNumProcs(i int) *AgentStatusUpdate {
	asu.mutation.AddNumProcs(i)
	return asu
}

// SetOs sets the Os field.
func (asu *AgentStatusUpdate) SetOs(s string) *AgentStatusUpdate {
	asu.mutation.SetOs(s)
	return asu
}

// SetHostID sets the HostID field.
func (asu *AgentStatusUpdate) SetHostID(s string) *AgentStatusUpdate {
	asu.mutation.SetHostID(s)
	return asu
}

// SetLoad1 sets the Load1 field.
func (asu *AgentStatusUpdate) SetLoad1(f float64) *AgentStatusUpdate {
	asu.mutation.ResetLoad1()
	asu.mutation.SetLoad1(f)
	return asu
}

// AddLoad1 adds f to Load1.
func (asu *AgentStatusUpdate) AddLoad1(f float64) *AgentStatusUpdate {
	asu.mutation.AddLoad1(f)
	return asu
}

// SetLoad5 sets the Load5 field.
func (asu *AgentStatusUpdate) SetLoad5(f float64) *AgentStatusUpdate {
	asu.mutation.ResetLoad5()
	asu.mutation.SetLoad5(f)
	return asu
}

// AddLoad5 adds f to Load5.
func (asu *AgentStatusUpdate) AddLoad5(f float64) *AgentStatusUpdate {
	asu.mutation.AddLoad5(f)
	return asu
}

// SetLoad15 sets the Load15 field.
func (asu *AgentStatusUpdate) SetLoad15(f float64) *AgentStatusUpdate {
	asu.mutation.ResetLoad15()
	asu.mutation.SetLoad15(f)
	return asu
}

// AddLoad15 adds f to Load15.
func (asu *AgentStatusUpdate) AddLoad15(f float64) *AgentStatusUpdate {
	asu.mutation.AddLoad15(f)
	return asu
}

// SetTotalMem sets the TotalMem field.
func (asu *AgentStatusUpdate) SetTotalMem(i int) *AgentStatusUpdate {
	asu.mutation.ResetTotalMem()
	asu.mutation.SetTotalMem(i)
	return asu
}

// AddTotalMem adds i to TotalMem.
func (asu *AgentStatusUpdate) AddTotalMem(i int) *AgentStatusUpdate {
	asu.mutation.AddTotalMem(i)
	return asu
}

// SetFreeMem sets the FreeMem field.
func (asu *AgentStatusUpdate) SetFreeMem(i int) *AgentStatusUpdate {
	asu.mutation.ResetFreeMem()
	asu.mutation.SetFreeMem(i)
	return asu
}

// AddFreeMem adds i to FreeMem.
func (asu *AgentStatusUpdate) AddFreeMem(i int) *AgentStatusUpdate {
	asu.mutation.AddFreeMem(i)
	return asu
}

// SetUsedMem sets the UsedMem field.
func (asu *AgentStatusUpdate) SetUsedMem(i int) *AgentStatusUpdate {
	asu.mutation.ResetUsedMem()
	asu.mutation.SetUsedMem(i)
	return asu
}

// AddUsedMem adds i to UsedMem.
func (asu *AgentStatusUpdate) AddUsedMem(i int) *AgentStatusUpdate {
	asu.mutation.AddUsedMem(i)
	return asu
}

// AddHostIDs adds the host edge to ProvisionedHost by ids.
func (asu *AgentStatusUpdate) AddHostIDs(ids ...int) *AgentStatusUpdate {
	asu.mutation.AddHostIDs(ids...)
	return asu
}

// AddHost adds the host edges to ProvisionedHost.
func (asu *AgentStatusUpdate) AddHost(p ...*ProvisionedHost) *AgentStatusUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return asu.AddHostIDs(ids...)
}

// Mutation returns the AgentStatusMutation object of the builder.
func (asu *AgentStatusUpdate) Mutation() *AgentStatusMutation {
	return asu.mutation
}

// ClearHost clears all "host" edges to type ProvisionedHost.
func (asu *AgentStatusUpdate) ClearHost() *AgentStatusUpdate {
	asu.mutation.ClearHost()
	return asu
}

// RemoveHostIDs removes the host edge to ProvisionedHost by ids.
func (asu *AgentStatusUpdate) RemoveHostIDs(ids ...int) *AgentStatusUpdate {
	asu.mutation.RemoveHostIDs(ids...)
	return asu
}

// RemoveHost removes host edges to ProvisionedHost.
func (asu *AgentStatusUpdate) RemoveHost(p ...*ProvisionedHost) *AgentStatusUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return asu.RemoveHostIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (asu *AgentStatusUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(asu.hooks) == 0 {
		affected, err = asu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AgentStatusMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			asu.mutation = mutation
			affected, err = asu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(asu.hooks) - 1; i >= 0; i-- {
			mut = asu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, asu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (asu *AgentStatusUpdate) SaveX(ctx context.Context) int {
	affected, err := asu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (asu *AgentStatusUpdate) Exec(ctx context.Context) error {
	_, err := asu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (asu *AgentStatusUpdate) ExecX(ctx context.Context) {
	if err := asu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (asu *AgentStatusUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   agentstatus.Table,
			Columns: agentstatus.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: agentstatus.FieldID,
			},
		},
	}
	if ps := asu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := asu.mutation.ClientID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldClientID,
		})
	}
	if value, ok := asu.mutation.Hostname(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldHostname,
		})
	}
	if value, ok := asu.mutation.UpTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldUpTime,
		})
	}
	if value, ok := asu.mutation.AddedUpTime(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldUpTime,
		})
	}
	if value, ok := asu.mutation.BootTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldBootTime,
		})
	}
	if value, ok := asu.mutation.AddedBootTime(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldBootTime,
		})
	}
	if value, ok := asu.mutation.NumProcs(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldNumProcs,
		})
	}
	if value, ok := asu.mutation.AddedNumProcs(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldNumProcs,
		})
	}
	if value, ok := asu.mutation.Os(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldOs,
		})
	}
	if value, ok := asu.mutation.HostID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldHostID,
		})
	}
	if value, ok := asu.mutation.Load1(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad1,
		})
	}
	if value, ok := asu.mutation.AddedLoad1(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad1,
		})
	}
	if value, ok := asu.mutation.Load5(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad5,
		})
	}
	if value, ok := asu.mutation.AddedLoad5(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad5,
		})
	}
	if value, ok := asu.mutation.Load15(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad15,
		})
	}
	if value, ok := asu.mutation.AddedLoad15(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad15,
		})
	}
	if value, ok := asu.mutation.TotalMem(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldTotalMem,
		})
	}
	if value, ok := asu.mutation.AddedTotalMem(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldTotalMem,
		})
	}
	if value, ok := asu.mutation.FreeMem(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldFreeMem,
		})
	}
	if value, ok := asu.mutation.AddedFreeMem(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldFreeMem,
		})
	}
	if value, ok := asu.mutation.UsedMem(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldUsedMem,
		})
	}
	if value, ok := asu.mutation.AddedUsedMem(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldUsedMem,
		})
	}
	if asu.mutation.HostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   agentstatus.HostTable,
			Columns: agentstatus.HostPrimaryKey,
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
	if nodes := asu.mutation.RemovedHostIDs(); len(nodes) > 0 && !asu.mutation.HostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   agentstatus.HostTable,
			Columns: agentstatus.HostPrimaryKey,
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
	if nodes := asu.mutation.HostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   agentstatus.HostTable,
			Columns: agentstatus.HostPrimaryKey,
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
	if n, err = sqlgraph.UpdateNodes(ctx, asu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{agentstatus.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// AgentStatusUpdateOne is the builder for updating a single AgentStatus entity.
type AgentStatusUpdateOne struct {
	config
	hooks    []Hook
	mutation *AgentStatusMutation
}

// SetClientID sets the ClientID field.
func (asuo *AgentStatusUpdateOne) SetClientID(s string) *AgentStatusUpdateOne {
	asuo.mutation.SetClientID(s)
	return asuo
}

// SetHostname sets the Hostname field.
func (asuo *AgentStatusUpdateOne) SetHostname(s string) *AgentStatusUpdateOne {
	asuo.mutation.SetHostname(s)
	return asuo
}

// SetUpTime sets the UpTime field.
func (asuo *AgentStatusUpdateOne) SetUpTime(i int) *AgentStatusUpdateOne {
	asuo.mutation.ResetUpTime()
	asuo.mutation.SetUpTime(i)
	return asuo
}

// AddUpTime adds i to UpTime.
func (asuo *AgentStatusUpdateOne) AddUpTime(i int) *AgentStatusUpdateOne {
	asuo.mutation.AddUpTime(i)
	return asuo
}

// SetBootTime sets the BootTime field.
func (asuo *AgentStatusUpdateOne) SetBootTime(i int) *AgentStatusUpdateOne {
	asuo.mutation.ResetBootTime()
	asuo.mutation.SetBootTime(i)
	return asuo
}

// AddBootTime adds i to BootTime.
func (asuo *AgentStatusUpdateOne) AddBootTime(i int) *AgentStatusUpdateOne {
	asuo.mutation.AddBootTime(i)
	return asuo
}

// SetNumProcs sets the NumProcs field.
func (asuo *AgentStatusUpdateOne) SetNumProcs(i int) *AgentStatusUpdateOne {
	asuo.mutation.ResetNumProcs()
	asuo.mutation.SetNumProcs(i)
	return asuo
}

// AddNumProcs adds i to NumProcs.
func (asuo *AgentStatusUpdateOne) AddNumProcs(i int) *AgentStatusUpdateOne {
	asuo.mutation.AddNumProcs(i)
	return asuo
}

// SetOs sets the Os field.
func (asuo *AgentStatusUpdateOne) SetOs(s string) *AgentStatusUpdateOne {
	asuo.mutation.SetOs(s)
	return asuo
}

// SetHostID sets the HostID field.
func (asuo *AgentStatusUpdateOne) SetHostID(s string) *AgentStatusUpdateOne {
	asuo.mutation.SetHostID(s)
	return asuo
}

// SetLoad1 sets the Load1 field.
func (asuo *AgentStatusUpdateOne) SetLoad1(f float64) *AgentStatusUpdateOne {
	asuo.mutation.ResetLoad1()
	asuo.mutation.SetLoad1(f)
	return asuo
}

// AddLoad1 adds f to Load1.
func (asuo *AgentStatusUpdateOne) AddLoad1(f float64) *AgentStatusUpdateOne {
	asuo.mutation.AddLoad1(f)
	return asuo
}

// SetLoad5 sets the Load5 field.
func (asuo *AgentStatusUpdateOne) SetLoad5(f float64) *AgentStatusUpdateOne {
	asuo.mutation.ResetLoad5()
	asuo.mutation.SetLoad5(f)
	return asuo
}

// AddLoad5 adds f to Load5.
func (asuo *AgentStatusUpdateOne) AddLoad5(f float64) *AgentStatusUpdateOne {
	asuo.mutation.AddLoad5(f)
	return asuo
}

// SetLoad15 sets the Load15 field.
func (asuo *AgentStatusUpdateOne) SetLoad15(f float64) *AgentStatusUpdateOne {
	asuo.mutation.ResetLoad15()
	asuo.mutation.SetLoad15(f)
	return asuo
}

// AddLoad15 adds f to Load15.
func (asuo *AgentStatusUpdateOne) AddLoad15(f float64) *AgentStatusUpdateOne {
	asuo.mutation.AddLoad15(f)
	return asuo
}

// SetTotalMem sets the TotalMem field.
func (asuo *AgentStatusUpdateOne) SetTotalMem(i int) *AgentStatusUpdateOne {
	asuo.mutation.ResetTotalMem()
	asuo.mutation.SetTotalMem(i)
	return asuo
}

// AddTotalMem adds i to TotalMem.
func (asuo *AgentStatusUpdateOne) AddTotalMem(i int) *AgentStatusUpdateOne {
	asuo.mutation.AddTotalMem(i)
	return asuo
}

// SetFreeMem sets the FreeMem field.
func (asuo *AgentStatusUpdateOne) SetFreeMem(i int) *AgentStatusUpdateOne {
	asuo.mutation.ResetFreeMem()
	asuo.mutation.SetFreeMem(i)
	return asuo
}

// AddFreeMem adds i to FreeMem.
func (asuo *AgentStatusUpdateOne) AddFreeMem(i int) *AgentStatusUpdateOne {
	asuo.mutation.AddFreeMem(i)
	return asuo
}

// SetUsedMem sets the UsedMem field.
func (asuo *AgentStatusUpdateOne) SetUsedMem(i int) *AgentStatusUpdateOne {
	asuo.mutation.ResetUsedMem()
	asuo.mutation.SetUsedMem(i)
	return asuo
}

// AddUsedMem adds i to UsedMem.
func (asuo *AgentStatusUpdateOne) AddUsedMem(i int) *AgentStatusUpdateOne {
	asuo.mutation.AddUsedMem(i)
	return asuo
}

// AddHostIDs adds the host edge to ProvisionedHost by ids.
func (asuo *AgentStatusUpdateOne) AddHostIDs(ids ...int) *AgentStatusUpdateOne {
	asuo.mutation.AddHostIDs(ids...)
	return asuo
}

// AddHost adds the host edges to ProvisionedHost.
func (asuo *AgentStatusUpdateOne) AddHost(p ...*ProvisionedHost) *AgentStatusUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return asuo.AddHostIDs(ids...)
}

// Mutation returns the AgentStatusMutation object of the builder.
func (asuo *AgentStatusUpdateOne) Mutation() *AgentStatusMutation {
	return asuo.mutation
}

// ClearHost clears all "host" edges to type ProvisionedHost.
func (asuo *AgentStatusUpdateOne) ClearHost() *AgentStatusUpdateOne {
	asuo.mutation.ClearHost()
	return asuo
}

// RemoveHostIDs removes the host edge to ProvisionedHost by ids.
func (asuo *AgentStatusUpdateOne) RemoveHostIDs(ids ...int) *AgentStatusUpdateOne {
	asuo.mutation.RemoveHostIDs(ids...)
	return asuo
}

// RemoveHost removes host edges to ProvisionedHost.
func (asuo *AgentStatusUpdateOne) RemoveHost(p ...*ProvisionedHost) *AgentStatusUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return asuo.RemoveHostIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (asuo *AgentStatusUpdateOne) Save(ctx context.Context) (*AgentStatus, error) {
	var (
		err  error
		node *AgentStatus
	)
	if len(asuo.hooks) == 0 {
		node, err = asuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AgentStatusMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			asuo.mutation = mutation
			node, err = asuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(asuo.hooks) - 1; i >= 0; i-- {
			mut = asuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, asuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (asuo *AgentStatusUpdateOne) SaveX(ctx context.Context) *AgentStatus {
	node, err := asuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (asuo *AgentStatusUpdateOne) Exec(ctx context.Context) error {
	_, err := asuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (asuo *AgentStatusUpdateOne) ExecX(ctx context.Context) {
	if err := asuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (asuo *AgentStatusUpdateOne) sqlSave(ctx context.Context) (_node *AgentStatus, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   agentstatus.Table,
			Columns: agentstatus.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: agentstatus.FieldID,
			},
		},
	}
	id, ok := asuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing AgentStatus.ID for update")}
	}
	_spec.Node.ID.Value = id
	if value, ok := asuo.mutation.ClientID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldClientID,
		})
	}
	if value, ok := asuo.mutation.Hostname(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldHostname,
		})
	}
	if value, ok := asuo.mutation.UpTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldUpTime,
		})
	}
	if value, ok := asuo.mutation.AddedUpTime(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldUpTime,
		})
	}
	if value, ok := asuo.mutation.BootTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldBootTime,
		})
	}
	if value, ok := asuo.mutation.AddedBootTime(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldBootTime,
		})
	}
	if value, ok := asuo.mutation.NumProcs(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldNumProcs,
		})
	}
	if value, ok := asuo.mutation.AddedNumProcs(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldNumProcs,
		})
	}
	if value, ok := asuo.mutation.Os(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldOs,
		})
	}
	if value, ok := asuo.mutation.HostID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldHostID,
		})
	}
	if value, ok := asuo.mutation.Load1(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad1,
		})
	}
	if value, ok := asuo.mutation.AddedLoad1(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad1,
		})
	}
	if value, ok := asuo.mutation.Load5(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad5,
		})
	}
	if value, ok := asuo.mutation.AddedLoad5(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad5,
		})
	}
	if value, ok := asuo.mutation.Load15(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad15,
		})
	}
	if value, ok := asuo.mutation.AddedLoad15(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad15,
		})
	}
	if value, ok := asuo.mutation.TotalMem(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldTotalMem,
		})
	}
	if value, ok := asuo.mutation.AddedTotalMem(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldTotalMem,
		})
	}
	if value, ok := asuo.mutation.FreeMem(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldFreeMem,
		})
	}
	if value, ok := asuo.mutation.AddedFreeMem(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldFreeMem,
		})
	}
	if value, ok := asuo.mutation.UsedMem(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldUsedMem,
		})
	}
	if value, ok := asuo.mutation.AddedUsedMem(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldUsedMem,
		})
	}
	if asuo.mutation.HostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   agentstatus.HostTable,
			Columns: agentstatus.HostPrimaryKey,
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
	if nodes := asuo.mutation.RemovedHostIDs(); len(nodes) > 0 && !asuo.mutation.HostCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   agentstatus.HostTable,
			Columns: agentstatus.HostPrimaryKey,
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
	if nodes := asuo.mutation.HostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   agentstatus.HostTable,
			Columns: agentstatus.HostPrimaryKey,
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
	_node = &AgentStatus{config: asuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues()
	if err = sqlgraph.UpdateNode(ctx, asuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{agentstatus.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}
