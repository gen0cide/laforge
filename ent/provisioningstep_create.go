// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/gen0cide/laforge/ent/command"
	"github.com/gen0cide/laforge/ent/dnsrecord"
	"github.com/gen0cide/laforge/ent/filedelete"
	"github.com/gen0cide/laforge/ent/filedownload"
	"github.com/gen0cide/laforge/ent/fileextract"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	"github.com/gen0cide/laforge/ent/provisioningstep"
	"github.com/gen0cide/laforge/ent/script"
	"github.com/gen0cide/laforge/ent/status"
	"github.com/gen0cide/laforge/ent/tag"
)

// ProvisioningStepCreate is the builder for creating a ProvisioningStep entity.
type ProvisioningStepCreate struct {
	config
	mutation *ProvisioningStepMutation
	hooks    []Hook
}

// SetProvisionerType sets the "provisioner_type" field.
func (psc *ProvisioningStepCreate) SetProvisionerType(s string) *ProvisioningStepCreate {
	psc.mutation.SetProvisionerType(s)
	return psc
}

// SetStepNumber sets the "step_number" field.
func (psc *ProvisioningStepCreate) SetStepNumber(i int) *ProvisioningStepCreate {
	psc.mutation.SetStepNumber(i)
	return psc
}

// AddProvisioningStepToTagIDs adds the "ProvisioningStepToTag" edge to the Tag entity by IDs.
func (psc *ProvisioningStepCreate) AddProvisioningStepToTagIDs(ids ...int) *ProvisioningStepCreate {
	psc.mutation.AddProvisioningStepToTagIDs(ids...)
	return psc
}

// AddProvisioningStepToTag adds the "ProvisioningStepToTag" edges to the Tag entity.
func (psc *ProvisioningStepCreate) AddProvisioningStepToTag(t ...*Tag) *ProvisioningStepCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return psc.AddProvisioningStepToTagIDs(ids...)
}

// AddProvisioningStepToStatuIDs adds the "ProvisioningStepToStatus" edge to the Status entity by IDs.
func (psc *ProvisioningStepCreate) AddProvisioningStepToStatuIDs(ids ...int) *ProvisioningStepCreate {
	psc.mutation.AddProvisioningStepToStatuIDs(ids...)
	return psc
}

// AddProvisioningStepToStatus adds the "ProvisioningStepToStatus" edges to the Status entity.
func (psc *ProvisioningStepCreate) AddProvisioningStepToStatus(s ...*Status) *ProvisioningStepCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return psc.AddProvisioningStepToStatuIDs(ids...)
}

// AddProvisioningStepToProvisionedHostIDs adds the "ProvisioningStepToProvisionedHost" edge to the ProvisionedHost entity by IDs.
func (psc *ProvisioningStepCreate) AddProvisioningStepToProvisionedHostIDs(ids ...int) *ProvisioningStepCreate {
	psc.mutation.AddProvisioningStepToProvisionedHostIDs(ids...)
	return psc
}

// AddProvisioningStepToProvisionedHost adds the "ProvisioningStepToProvisionedHost" edges to the ProvisionedHost entity.
func (psc *ProvisioningStepCreate) AddProvisioningStepToProvisionedHost(p ...*ProvisionedHost) *ProvisioningStepCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return psc.AddProvisioningStepToProvisionedHostIDs(ids...)
}

// AddProvisioningStepToScriptIDs adds the "ProvisioningStepToScript" edge to the Script entity by IDs.
func (psc *ProvisioningStepCreate) AddProvisioningStepToScriptIDs(ids ...int) *ProvisioningStepCreate {
	psc.mutation.AddProvisioningStepToScriptIDs(ids...)
	return psc
}

// AddProvisioningStepToScript adds the "ProvisioningStepToScript" edges to the Script entity.
func (psc *ProvisioningStepCreate) AddProvisioningStepToScript(s ...*Script) *ProvisioningStepCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return psc.AddProvisioningStepToScriptIDs(ids...)
}

// AddProvisioningStepToCommandIDs adds the "ProvisioningStepToCommand" edge to the Command entity by IDs.
func (psc *ProvisioningStepCreate) AddProvisioningStepToCommandIDs(ids ...int) *ProvisioningStepCreate {
	psc.mutation.AddProvisioningStepToCommandIDs(ids...)
	return psc
}

// AddProvisioningStepToCommand adds the "ProvisioningStepToCommand" edges to the Command entity.
func (psc *ProvisioningStepCreate) AddProvisioningStepToCommand(c ...*Command) *ProvisioningStepCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return psc.AddProvisioningStepToCommandIDs(ids...)
}

// AddProvisioningStepToDNSRecordIDs adds the "ProvisioningStepToDNSRecord" edge to the DNSRecord entity by IDs.
func (psc *ProvisioningStepCreate) AddProvisioningStepToDNSRecordIDs(ids ...int) *ProvisioningStepCreate {
	psc.mutation.AddProvisioningStepToDNSRecordIDs(ids...)
	return psc
}

// AddProvisioningStepToDNSRecord adds the "ProvisioningStepToDNSRecord" edges to the DNSRecord entity.
func (psc *ProvisioningStepCreate) AddProvisioningStepToDNSRecord(d ...*DNSRecord) *ProvisioningStepCreate {
	ids := make([]int, len(d))
	for i := range d {
		ids[i] = d[i].ID
	}
	return psc.AddProvisioningStepToDNSRecordIDs(ids...)
}

// AddProvisioningStepToFileDeleteIDs adds the "ProvisioningStepToFileDelete" edge to the FileDelete entity by IDs.
func (psc *ProvisioningStepCreate) AddProvisioningStepToFileDeleteIDs(ids ...int) *ProvisioningStepCreate {
	psc.mutation.AddProvisioningStepToFileDeleteIDs(ids...)
	return psc
}

// AddProvisioningStepToFileDelete adds the "ProvisioningStepToFileDelete" edges to the FileDelete entity.
func (psc *ProvisioningStepCreate) AddProvisioningStepToFileDelete(f ...*FileDelete) *ProvisioningStepCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return psc.AddProvisioningStepToFileDeleteIDs(ids...)
}

// AddProvisioningStepToFileDownloadIDs adds the "ProvisioningStepToFileDownload" edge to the FileDownload entity by IDs.
func (psc *ProvisioningStepCreate) AddProvisioningStepToFileDownloadIDs(ids ...int) *ProvisioningStepCreate {
	psc.mutation.AddProvisioningStepToFileDownloadIDs(ids...)
	return psc
}

// AddProvisioningStepToFileDownload adds the "ProvisioningStepToFileDownload" edges to the FileDownload entity.
func (psc *ProvisioningStepCreate) AddProvisioningStepToFileDownload(f ...*FileDownload) *ProvisioningStepCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return psc.AddProvisioningStepToFileDownloadIDs(ids...)
}

// AddProvisioningStepToFileExtractIDs adds the "ProvisioningStepToFileExtract" edge to the FileExtract entity by IDs.
func (psc *ProvisioningStepCreate) AddProvisioningStepToFileExtractIDs(ids ...int) *ProvisioningStepCreate {
	psc.mutation.AddProvisioningStepToFileExtractIDs(ids...)
	return psc
}

// AddProvisioningStepToFileExtract adds the "ProvisioningStepToFileExtract" edges to the FileExtract entity.
func (psc *ProvisioningStepCreate) AddProvisioningStepToFileExtract(f ...*FileExtract) *ProvisioningStepCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return psc.AddProvisioningStepToFileExtractIDs(ids...)
}

// Mutation returns the ProvisioningStepMutation object of the builder.
func (psc *ProvisioningStepCreate) Mutation() *ProvisioningStepMutation {
	return psc.mutation
}

// Save creates the ProvisioningStep in the database.
func (psc *ProvisioningStepCreate) Save(ctx context.Context) (*ProvisioningStep, error) {
	var (
		err  error
		node *ProvisioningStep
	)
	if len(psc.hooks) == 0 {
		if err = psc.check(); err != nil {
			return nil, err
		}
		node, err = psc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ProvisioningStepMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = psc.check(); err != nil {
				return nil, err
			}
			psc.mutation = mutation
			node, err = psc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(psc.hooks) - 1; i >= 0; i-- {
			mut = psc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, psc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (psc *ProvisioningStepCreate) SaveX(ctx context.Context) *ProvisioningStep {
	v, err := psc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// check runs all checks and user-defined validators on the builder.
func (psc *ProvisioningStepCreate) check() error {
	if _, ok := psc.mutation.ProvisionerType(); !ok {
		return &ValidationError{Name: "provisioner_type", err: errors.New("ent: missing required field \"provisioner_type\"")}
	}
	if _, ok := psc.mutation.StepNumber(); !ok {
		return &ValidationError{Name: "step_number", err: errors.New("ent: missing required field \"step_number\"")}
	}
	return nil
}

func (psc *ProvisioningStepCreate) sqlSave(ctx context.Context) (*ProvisioningStep, error) {
	_node, _spec := psc.createSpec()
	if err := sqlgraph.CreateNode(ctx, psc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (psc *ProvisioningStepCreate) createSpec() (*ProvisioningStep, *sqlgraph.CreateSpec) {
	var (
		_node = &ProvisioningStep{config: psc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: provisioningstep.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: provisioningstep.FieldID,
			},
		}
	)
	if value, ok := psc.mutation.ProvisionerType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: provisioningstep.FieldProvisionerType,
		})
		_node.ProvisionerType = value
	}
	if value, ok := psc.mutation.StepNumber(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: provisioningstep.FieldStepNumber,
		})
		_node.StepNumber = value
	}
	if nodes := psc.mutation.ProvisioningStepToTagIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisioningstep.ProvisioningStepToTagTable,
			Columns: []string{provisioningstep.ProvisioningStepToTagColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := psc.mutation.ProvisioningStepToStatusIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisioningstep.ProvisioningStepToStatusTable,
			Columns: []string{provisioningstep.ProvisioningStepToStatusColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := psc.mutation.ProvisioningStepToProvisionedHostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   provisioningstep.ProvisioningStepToProvisionedHostTable,
			Columns: provisioningstep.ProvisioningStepToProvisionedHostPrimaryKey,
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
	if nodes := psc.mutation.ProvisioningStepToScriptIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisioningstep.ProvisioningStepToScriptTable,
			Columns: []string{provisioningstep.ProvisioningStepToScriptColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: script.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := psc.mutation.ProvisioningStepToCommandIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisioningstep.ProvisioningStepToCommandTable,
			Columns: []string{provisioningstep.ProvisioningStepToCommandColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: command.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := psc.mutation.ProvisioningStepToDNSRecordIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisioningstep.ProvisioningStepToDNSRecordTable,
			Columns: []string{provisioningstep.ProvisioningStepToDNSRecordColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: dnsrecord.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := psc.mutation.ProvisioningStepToFileDeleteIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisioningstep.ProvisioningStepToFileDeleteTable,
			Columns: []string{provisioningstep.ProvisioningStepToFileDeleteColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: filedelete.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := psc.mutation.ProvisioningStepToFileDownloadIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisioningstep.ProvisioningStepToFileDownloadTable,
			Columns: []string{provisioningstep.ProvisioningStepToFileDownloadColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: filedownload.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := psc.mutation.ProvisioningStepToFileExtractIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   provisioningstep.ProvisioningStepToFileExtractTable,
			Columns: []string{provisioningstep.ProvisioningStepToFileExtractColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: fileextract.FieldID,
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

// ProvisioningStepCreateBulk is the builder for creating many ProvisioningStep entities in bulk.
type ProvisioningStepCreateBulk struct {
	config
	builders []*ProvisioningStepCreate
}

// Save creates the ProvisioningStep entities in the database.
func (pscb *ProvisioningStepCreateBulk) Save(ctx context.Context) ([]*ProvisioningStep, error) {
	specs := make([]*sqlgraph.CreateSpec, len(pscb.builders))
	nodes := make([]*ProvisioningStep, len(pscb.builders))
	mutators := make([]Mutator, len(pscb.builders))
	for i := range pscb.builders {
		func(i int, root context.Context) {
			builder := pscb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ProvisioningStepMutation)
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
					_, err = mutators[i+1].Mutate(root, pscb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pscb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, pscb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (pscb *ProvisioningStepCreateBulk) SaveX(ctx context.Context) []*ProvisioningStep {
	v, err := pscb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}