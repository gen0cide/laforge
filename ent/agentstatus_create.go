// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/gen0cide/laforge/ent/agentstatus"
	"github.com/gen0cide/laforge/ent/provisionedhost"
)

// AgentStatusCreate is the builder for creating a AgentStatus entity.
type AgentStatusCreate struct {
	config
	mutation *AgentStatusMutation
	hooks    []Hook
}

// SetClientID sets the ClientID field.
func (asc *AgentStatusCreate) SetClientID(s string) *AgentStatusCreate {
	asc.mutation.SetClientID(s)
	return asc
}

// SetHostname sets the Hostname field.
func (asc *AgentStatusCreate) SetHostname(s string) *AgentStatusCreate {
	asc.mutation.SetHostname(s)
	return asc
}

// SetUpTime sets the UpTime field.
func (asc *AgentStatusCreate) SetUpTime(i int) *AgentStatusCreate {
	asc.mutation.SetUpTime(i)
	return asc
}

// SetBootTime sets the BootTime field.
func (asc *AgentStatusCreate) SetBootTime(i int) *AgentStatusCreate {
	asc.mutation.SetBootTime(i)
	return asc
}

// SetNumProcs sets the NumProcs field.
func (asc *AgentStatusCreate) SetNumProcs(i int) *AgentStatusCreate {
	asc.mutation.SetNumProcs(i)
	return asc
}

// SetOs sets the Os field.
func (asc *AgentStatusCreate) SetOs(s string) *AgentStatusCreate {
	asc.mutation.SetOs(s)
	return asc
}

// SetHostID sets the HostID field.
func (asc *AgentStatusCreate) SetHostID(s string) *AgentStatusCreate {
	asc.mutation.SetHostID(s)
	return asc
}

// SetLoad1 sets the Load1 field.
func (asc *AgentStatusCreate) SetLoad1(f float64) *AgentStatusCreate {
	asc.mutation.SetLoad1(f)
	return asc
}

// SetLoad5 sets the Load5 field.
func (asc *AgentStatusCreate) SetLoad5(f float64) *AgentStatusCreate {
	asc.mutation.SetLoad5(f)
	return asc
}

// SetLoad15 sets the Load15 field.
func (asc *AgentStatusCreate) SetLoad15(f float64) *AgentStatusCreate {
	asc.mutation.SetLoad15(f)
	return asc
}

// SetTotalMem sets the TotalMem field.
func (asc *AgentStatusCreate) SetTotalMem(i int) *AgentStatusCreate {
	asc.mutation.SetTotalMem(i)
	return asc
}

// SetFreeMem sets the FreeMem field.
func (asc *AgentStatusCreate) SetFreeMem(i int) *AgentStatusCreate {
	asc.mutation.SetFreeMem(i)
	return asc
}

// SetUsedMem sets the UsedMem field.
func (asc *AgentStatusCreate) SetUsedMem(i int) *AgentStatusCreate {
	asc.mutation.SetUsedMem(i)
	return asc
}

// AddHostIDs adds the host edge to ProvisionedHost by ids.
func (asc *AgentStatusCreate) AddHostIDs(ids ...int) *AgentStatusCreate {
	asc.mutation.AddHostIDs(ids...)
	return asc
}

// AddHost adds the host edges to ProvisionedHost.
func (asc *AgentStatusCreate) AddHost(p ...*ProvisionedHost) *AgentStatusCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return asc.AddHostIDs(ids...)
}

// Mutation returns the AgentStatusMutation object of the builder.
func (asc *AgentStatusCreate) Mutation() *AgentStatusMutation {
	return asc.mutation
}

// Save creates the AgentStatus in the database.
func (asc *AgentStatusCreate) Save(ctx context.Context) (*AgentStatus, error) {
	var (
		err  error
		node *AgentStatus
	)
	if len(asc.hooks) == 0 {
		if err = asc.check(); err != nil {
			return nil, err
		}
		node, err = asc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AgentStatusMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = asc.check(); err != nil {
				return nil, err
			}
			asc.mutation = mutation
			node, err = asc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(asc.hooks) - 1; i >= 0; i-- {
			mut = asc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, asc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (asc *AgentStatusCreate) SaveX(ctx context.Context) *AgentStatus {
	v, err := asc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// check runs all checks and user-defined validators on the builder.
func (asc *AgentStatusCreate) check() error {
	if _, ok := asc.mutation.ClientID(); !ok {
		return &ValidationError{Name: "ClientID", err: errors.New("ent: missing required field \"ClientID\"")}
	}
	if _, ok := asc.mutation.Hostname(); !ok {
		return &ValidationError{Name: "Hostname", err: errors.New("ent: missing required field \"Hostname\"")}
	}
	if _, ok := asc.mutation.UpTime(); !ok {
		return &ValidationError{Name: "UpTime", err: errors.New("ent: missing required field \"UpTime\"")}
	}
	if _, ok := asc.mutation.BootTime(); !ok {
		return &ValidationError{Name: "BootTime", err: errors.New("ent: missing required field \"BootTime\"")}
	}
	if _, ok := asc.mutation.NumProcs(); !ok {
		return &ValidationError{Name: "NumProcs", err: errors.New("ent: missing required field \"NumProcs\"")}
	}
	if _, ok := asc.mutation.Os(); !ok {
		return &ValidationError{Name: "Os", err: errors.New("ent: missing required field \"Os\"")}
	}
	if _, ok := asc.mutation.HostID(); !ok {
		return &ValidationError{Name: "HostID", err: errors.New("ent: missing required field \"HostID\"")}
	}
	if _, ok := asc.mutation.Load1(); !ok {
		return &ValidationError{Name: "Load1", err: errors.New("ent: missing required field \"Load1\"")}
	}
	if _, ok := asc.mutation.Load5(); !ok {
		return &ValidationError{Name: "Load5", err: errors.New("ent: missing required field \"Load5\"")}
	}
	if _, ok := asc.mutation.Load15(); !ok {
		return &ValidationError{Name: "Load15", err: errors.New("ent: missing required field \"Load15\"")}
	}
	if _, ok := asc.mutation.TotalMem(); !ok {
		return &ValidationError{Name: "TotalMem", err: errors.New("ent: missing required field \"TotalMem\"")}
	}
	if _, ok := asc.mutation.FreeMem(); !ok {
		return &ValidationError{Name: "FreeMem", err: errors.New("ent: missing required field \"FreeMem\"")}
	}
	if _, ok := asc.mutation.UsedMem(); !ok {
		return &ValidationError{Name: "UsedMem", err: errors.New("ent: missing required field \"UsedMem\"")}
	}
	return nil
}

func (asc *AgentStatusCreate) sqlSave(ctx context.Context) (*AgentStatus, error) {
	_node, _spec := asc.createSpec()
	if err := sqlgraph.CreateNode(ctx, asc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (asc *AgentStatusCreate) createSpec() (*AgentStatus, *sqlgraph.CreateSpec) {
	var (
		_node = &AgentStatus{config: asc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: agentstatus.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: agentstatus.FieldID,
			},
		}
	)
	if value, ok := asc.mutation.ClientID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldClientID,
		})
		_node.ClientID = value
	}
	if value, ok := asc.mutation.Hostname(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldHostname,
		})
		_node.Hostname = value
	}
	if value, ok := asc.mutation.UpTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldUpTime,
		})
		_node.UpTime = value
	}
	if value, ok := asc.mutation.BootTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldBootTime,
		})
		_node.BootTime = value
	}
	if value, ok := asc.mutation.NumProcs(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldNumProcs,
		})
		_node.NumProcs = value
	}
	if value, ok := asc.mutation.Os(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldOs,
		})
		_node.Os = value
	}
	if value, ok := asc.mutation.HostID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: agentstatus.FieldHostID,
		})
		_node.HostID = value
	}
	if value, ok := asc.mutation.Load1(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad1,
		})
		_node.Load1 = value
	}
	if value, ok := asc.mutation.Load5(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad5,
		})
		_node.Load5 = value
	}
	if value, ok := asc.mutation.Load15(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: agentstatus.FieldLoad15,
		})
		_node.Load15 = value
	}
	if value, ok := asc.mutation.TotalMem(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldTotalMem,
		})
		_node.TotalMem = value
	}
	if value, ok := asc.mutation.FreeMem(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldFreeMem,
		})
		_node.FreeMem = value
	}
	if value, ok := asc.mutation.UsedMem(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: agentstatus.FieldUsedMem,
		})
		_node.UsedMem = value
	}
	if nodes := asc.mutation.HostIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// AgentStatusCreateBulk is the builder for creating a bulk of AgentStatus entities.
type AgentStatusCreateBulk struct {
	config
	builders []*AgentStatusCreate
}

// Save creates the AgentStatus entities in the database.
func (ascb *AgentStatusCreateBulk) Save(ctx context.Context) ([]*AgentStatus, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ascb.builders))
	nodes := make([]*AgentStatus, len(ascb.builders))
	mutators := make([]Mutator, len(ascb.builders))
	for i := range ascb.builders {
		func(i int, root context.Context) {
			builder := ascb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AgentStatusMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ascb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ascb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
						if cerr, ok := isSQLConstraintError(err); ok {
							err = cerr
						}
					}
				}
				mutation.done = true
				if err != nil {
					return nil, err
				}
				id := specs[i].ID.Value.(int64)
				nodes[i].ID = int(id)
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ascb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX calls Save and panics if Save returns an error.
func (ascb *AgentStatusCreateBulk) SaveX(ctx context.Context) []*AgentStatus {
	v, err := ascb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
