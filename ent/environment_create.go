// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/gen0cide/laforge/ent/build"
	"github.com/gen0cide/laforge/ent/command"
	"github.com/gen0cide/laforge/ent/competition"
	"github.com/gen0cide/laforge/ent/dns"
	"github.com/gen0cide/laforge/ent/dnsrecord"
	"github.com/gen0cide/laforge/ent/environment"
	"github.com/gen0cide/laforge/ent/filedelete"
	"github.com/gen0cide/laforge/ent/filedownload"
	"github.com/gen0cide/laforge/ent/fileextract"
	"github.com/gen0cide/laforge/ent/finding"
	"github.com/gen0cide/laforge/ent/host"
	"github.com/gen0cide/laforge/ent/hostdependency"
	"github.com/gen0cide/laforge/ent/identity"
	"github.com/gen0cide/laforge/ent/includednetwork"
	"github.com/gen0cide/laforge/ent/network"
	"github.com/gen0cide/laforge/ent/script"
	"github.com/gen0cide/laforge/ent/tag"
	"github.com/gen0cide/laforge/ent/user"
)

// EnvironmentCreate is the builder for creating a Environment entity.
type EnvironmentCreate struct {
	config
	mutation *EnvironmentMutation
	hooks    []Hook
}

// SetHclID sets the "hcl_id" field.
func (ec *EnvironmentCreate) SetHclID(s string) *EnvironmentCreate {
	ec.mutation.SetHclID(s)
	return ec
}

// SetCompetitionID sets the "competition_id" field.
func (ec *EnvironmentCreate) SetCompetitionID(s string) *EnvironmentCreate {
	ec.mutation.SetCompetitionID(s)
	return ec
}

// SetName sets the "name" field.
func (ec *EnvironmentCreate) SetName(s string) *EnvironmentCreate {
	ec.mutation.SetName(s)
	return ec
}

// SetDescription sets the "description" field.
func (ec *EnvironmentCreate) SetDescription(s string) *EnvironmentCreate {
	ec.mutation.SetDescription(s)
	return ec
}

// SetBuilder sets the "builder" field.
func (ec *EnvironmentCreate) SetBuilder(s string) *EnvironmentCreate {
	ec.mutation.SetBuilder(s)
	return ec
}

// SetTeamCount sets the "team_count" field.
func (ec *EnvironmentCreate) SetTeamCount(i int) *EnvironmentCreate {
	ec.mutation.SetTeamCount(i)
	return ec
}

// SetRevision sets the "revision" field.
func (ec *EnvironmentCreate) SetRevision(i int) *EnvironmentCreate {
	ec.mutation.SetRevision(i)
	return ec
}

// SetAdminCidrs sets the "admin_cidrs" field.
func (ec *EnvironmentCreate) SetAdminCidrs(s []string) *EnvironmentCreate {
	ec.mutation.SetAdminCidrs(s)
	return ec
}

// SetExposedVdiPorts sets the "exposed_vdi_ports" field.
func (ec *EnvironmentCreate) SetExposedVdiPorts(s []string) *EnvironmentCreate {
	ec.mutation.SetExposedVdiPorts(s)
	return ec
}

// SetConfig sets the "config" field.
func (ec *EnvironmentCreate) SetConfig(m map[string]string) *EnvironmentCreate {
	ec.mutation.SetConfig(m)
	return ec
}

// SetTags sets the "tags" field.
func (ec *EnvironmentCreate) SetTags(m map[string]string) *EnvironmentCreate {
	ec.mutation.SetTags(m)
	return ec
}

// AddEnvironmentToTagIDs adds the "EnvironmentToTag" edge to the Tag entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToTagIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToTagIDs(ids...)
	return ec
}

// AddEnvironmentToTag adds the "EnvironmentToTag" edges to the Tag entity.
func (ec *EnvironmentCreate) AddEnvironmentToTag(t ...*Tag) *EnvironmentCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ec.AddEnvironmentToTagIDs(ids...)
}

// AddEnvironmentToUserIDs adds the "EnvironmentToUser" edge to the User entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToUserIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToUserIDs(ids...)
	return ec
}

// AddEnvironmentToUser adds the "EnvironmentToUser" edges to the User entity.
func (ec *EnvironmentCreate) AddEnvironmentToUser(u ...*User) *EnvironmentCreate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ec.AddEnvironmentToUserIDs(ids...)
}

// AddEnvironmentToHostIDs adds the "EnvironmentToHost" edge to the Host entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToHostIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToHostIDs(ids...)
	return ec
}

// AddEnvironmentToHost adds the "EnvironmentToHost" edges to the Host entity.
func (ec *EnvironmentCreate) AddEnvironmentToHost(h ...*Host) *EnvironmentCreate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return ec.AddEnvironmentToHostIDs(ids...)
}

// AddEnvironmentToCompetitionIDs adds the "EnvironmentToCompetition" edge to the Competition entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToCompetitionIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToCompetitionIDs(ids...)
	return ec
}

// AddEnvironmentToCompetition adds the "EnvironmentToCompetition" edges to the Competition entity.
func (ec *EnvironmentCreate) AddEnvironmentToCompetition(c ...*Competition) *EnvironmentCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return ec.AddEnvironmentToCompetitionIDs(ids...)
}

// AddEnvironmentToIdentityIDs adds the "EnvironmentToIdentity" edge to the Identity entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToIdentityIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToIdentityIDs(ids...)
	return ec
}

// AddEnvironmentToIdentity adds the "EnvironmentToIdentity" edges to the Identity entity.
func (ec *EnvironmentCreate) AddEnvironmentToIdentity(i ...*Identity) *EnvironmentCreate {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return ec.AddEnvironmentToIdentityIDs(ids...)
}

// AddEnvironmentToCommandIDs adds the "EnvironmentToCommand" edge to the Command entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToCommandIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToCommandIDs(ids...)
	return ec
}

// AddEnvironmentToCommand adds the "EnvironmentToCommand" edges to the Command entity.
func (ec *EnvironmentCreate) AddEnvironmentToCommand(c ...*Command) *EnvironmentCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return ec.AddEnvironmentToCommandIDs(ids...)
}

// AddEnvironmentToScriptIDs adds the "EnvironmentToScript" edge to the Script entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToScriptIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToScriptIDs(ids...)
	return ec
}

// AddEnvironmentToScript adds the "EnvironmentToScript" edges to the Script entity.
func (ec *EnvironmentCreate) AddEnvironmentToScript(s ...*Script) *EnvironmentCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ec.AddEnvironmentToScriptIDs(ids...)
}

// AddEnvironmentToFileDownloadIDs adds the "EnvironmentToFileDownload" edge to the FileDownload entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToFileDownloadIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToFileDownloadIDs(ids...)
	return ec
}

// AddEnvironmentToFileDownload adds the "EnvironmentToFileDownload" edges to the FileDownload entity.
func (ec *EnvironmentCreate) AddEnvironmentToFileDownload(f ...*FileDownload) *EnvironmentCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return ec.AddEnvironmentToFileDownloadIDs(ids...)
}

// AddEnvironmentToFileDeleteIDs adds the "EnvironmentToFileDelete" edge to the FileDelete entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToFileDeleteIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToFileDeleteIDs(ids...)
	return ec
}

// AddEnvironmentToFileDelete adds the "EnvironmentToFileDelete" edges to the FileDelete entity.
func (ec *EnvironmentCreate) AddEnvironmentToFileDelete(f ...*FileDelete) *EnvironmentCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return ec.AddEnvironmentToFileDeleteIDs(ids...)
}

// AddEnvironmentToFileExtractIDs adds the "EnvironmentToFileExtract" edge to the FileExtract entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToFileExtractIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToFileExtractIDs(ids...)
	return ec
}

// AddEnvironmentToFileExtract adds the "EnvironmentToFileExtract" edges to the FileExtract entity.
func (ec *EnvironmentCreate) AddEnvironmentToFileExtract(f ...*FileExtract) *EnvironmentCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return ec.AddEnvironmentToFileExtractIDs(ids...)
}

// AddEnvironmentToIncludedNetworkIDs adds the "EnvironmentToIncludedNetwork" edge to the IncludedNetwork entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToIncludedNetworkIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToIncludedNetworkIDs(ids...)
	return ec
}

// AddEnvironmentToIncludedNetwork adds the "EnvironmentToIncludedNetwork" edges to the IncludedNetwork entity.
func (ec *EnvironmentCreate) AddEnvironmentToIncludedNetwork(i ...*IncludedNetwork) *EnvironmentCreate {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return ec.AddEnvironmentToIncludedNetworkIDs(ids...)
}

// AddEnvironmentToFindingIDs adds the "EnvironmentToFinding" edge to the Finding entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToFindingIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToFindingIDs(ids...)
	return ec
}

// AddEnvironmentToFinding adds the "EnvironmentToFinding" edges to the Finding entity.
func (ec *EnvironmentCreate) AddEnvironmentToFinding(f ...*Finding) *EnvironmentCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return ec.AddEnvironmentToFindingIDs(ids...)
}

// AddEnvironmentToDNSRecordIDs adds the "EnvironmentToDNSRecord" edge to the DNSRecord entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToDNSRecordIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToDNSRecordIDs(ids...)
	return ec
}

// AddEnvironmentToDNSRecord adds the "EnvironmentToDNSRecord" edges to the DNSRecord entity.
func (ec *EnvironmentCreate) AddEnvironmentToDNSRecord(d ...*DNSRecord) *EnvironmentCreate {
	ids := make([]int, len(d))
	for i := range d {
		ids[i] = d[i].ID
	}
	return ec.AddEnvironmentToDNSRecordIDs(ids...)
}

// AddEnvironmentToDNSIDs adds the "EnvironmentToDNS" edge to the DNS entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToDNSIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToDNSIDs(ids...)
	return ec
}

// AddEnvironmentToDNS adds the "EnvironmentToDNS" edges to the DNS entity.
func (ec *EnvironmentCreate) AddEnvironmentToDNS(d ...*DNS) *EnvironmentCreate {
	ids := make([]int, len(d))
	for i := range d {
		ids[i] = d[i].ID
	}
	return ec.AddEnvironmentToDNSIDs(ids...)
}

// AddEnvironmentToNetworkIDs adds the "EnvironmentToNetwork" edge to the Network entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToNetworkIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToNetworkIDs(ids...)
	return ec
}

// AddEnvironmentToNetwork adds the "EnvironmentToNetwork" edges to the Network entity.
func (ec *EnvironmentCreate) AddEnvironmentToNetwork(n ...*Network) *EnvironmentCreate {
	ids := make([]int, len(n))
	for i := range n {
		ids[i] = n[i].ID
	}
	return ec.AddEnvironmentToNetworkIDs(ids...)
}

// AddEnvironmentToHostDependencyIDs adds the "EnvironmentToHostDependency" edge to the HostDependency entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToHostDependencyIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToHostDependencyIDs(ids...)
	return ec
}

// AddEnvironmentToHostDependency adds the "EnvironmentToHostDependency" edges to the HostDependency entity.
func (ec *EnvironmentCreate) AddEnvironmentToHostDependency(h ...*HostDependency) *EnvironmentCreate {
	ids := make([]int, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return ec.AddEnvironmentToHostDependencyIDs(ids...)
}

// AddEnvironmentToBuildIDs adds the "EnvironmentToBuild" edge to the Build entity by IDs.
func (ec *EnvironmentCreate) AddEnvironmentToBuildIDs(ids ...int) *EnvironmentCreate {
	ec.mutation.AddEnvironmentToBuildIDs(ids...)
	return ec
}

// AddEnvironmentToBuild adds the "EnvironmentToBuild" edges to the Build entity.
func (ec *EnvironmentCreate) AddEnvironmentToBuild(b ...*Build) *EnvironmentCreate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return ec.AddEnvironmentToBuildIDs(ids...)
}

// Mutation returns the EnvironmentMutation object of the builder.
func (ec *EnvironmentCreate) Mutation() *EnvironmentMutation {
	return ec.mutation
}

// Save creates the Environment in the database.
func (ec *EnvironmentCreate) Save(ctx context.Context) (*Environment, error) {
	var (
		err  error
		node *Environment
	)
	if len(ec.hooks) == 0 {
		if err = ec.check(); err != nil {
			return nil, err
		}
		node, err = ec.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EnvironmentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ec.check(); err != nil {
				return nil, err
			}
			ec.mutation = mutation
			node, err = ec.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ec.hooks) - 1; i >= 0; i-- {
			mut = ec.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ec.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ec *EnvironmentCreate) SaveX(ctx context.Context) *Environment {
	v, err := ec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// check runs all checks and user-defined validators on the builder.
func (ec *EnvironmentCreate) check() error {
	if _, ok := ec.mutation.HclID(); !ok {
		return &ValidationError{Name: "hcl_id", err: errors.New("ent: missing required field \"hcl_id\"")}
	}
	if _, ok := ec.mutation.CompetitionID(); !ok {
		return &ValidationError{Name: "competition_id", err: errors.New("ent: missing required field \"competition_id\"")}
	}
	if _, ok := ec.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New("ent: missing required field \"name\"")}
	}
	if _, ok := ec.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New("ent: missing required field \"description\"")}
	}
	if _, ok := ec.mutation.Builder(); !ok {
		return &ValidationError{Name: "builder", err: errors.New("ent: missing required field \"builder\"")}
	}
	if _, ok := ec.mutation.TeamCount(); !ok {
		return &ValidationError{Name: "team_count", err: errors.New("ent: missing required field \"team_count\"")}
	}
	if _, ok := ec.mutation.Revision(); !ok {
		return &ValidationError{Name: "revision", err: errors.New("ent: missing required field \"revision\"")}
	}
	if _, ok := ec.mutation.AdminCidrs(); !ok {
		return &ValidationError{Name: "admin_cidrs", err: errors.New("ent: missing required field \"admin_cidrs\"")}
	}
	if _, ok := ec.mutation.ExposedVdiPorts(); !ok {
		return &ValidationError{Name: "exposed_vdi_ports", err: errors.New("ent: missing required field \"exposed_vdi_ports\"")}
	}
	if _, ok := ec.mutation.Config(); !ok {
		return &ValidationError{Name: "config", err: errors.New("ent: missing required field \"config\"")}
	}
	if _, ok := ec.mutation.Tags(); !ok {
		return &ValidationError{Name: "tags", err: errors.New("ent: missing required field \"tags\"")}
	}
	return nil
}

func (ec *EnvironmentCreate) sqlSave(ctx context.Context) (*Environment, error) {
	_node, _spec := ec.createSpec()
	if err := sqlgraph.CreateNode(ctx, ec.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (ec *EnvironmentCreate) createSpec() (*Environment, *sqlgraph.CreateSpec) {
	var (
		_node = &Environment{config: ec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: environment.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: environment.FieldID,
			},
		}
	)
	if value, ok := ec.mutation.HclID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: environment.FieldHclID,
		})
		_node.HclID = value
	}
	if value, ok := ec.mutation.CompetitionID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: environment.FieldCompetitionID,
		})
		_node.CompetitionID = value
	}
	if value, ok := ec.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: environment.FieldName,
		})
		_node.Name = value
	}
	if value, ok := ec.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: environment.FieldDescription,
		})
		_node.Description = value
	}
	if value, ok := ec.mutation.Builder(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: environment.FieldBuilder,
		})
		_node.Builder = value
	}
	if value, ok := ec.mutation.TeamCount(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: environment.FieldTeamCount,
		})
		_node.TeamCount = value
	}
	if value, ok := ec.mutation.Revision(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: environment.FieldRevision,
		})
		_node.Revision = value
	}
	if value, ok := ec.mutation.AdminCidrs(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: environment.FieldAdminCidrs,
		})
		_node.AdminCidrs = value
	}
	if value, ok := ec.mutation.ExposedVdiPorts(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: environment.FieldExposedVdiPorts,
		})
		_node.ExposedVdiPorts = value
	}
	if value, ok := ec.mutation.Config(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: environment.FieldConfig,
		})
		_node.Config = value
	}
	if value, ok := ec.mutation.Tags(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: environment.FieldTags,
		})
		_node.Tags = value
	}
	if nodes := ec.mutation.EnvironmentToTagIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.EnvironmentToTagTable,
			Columns: []string{environment.EnvironmentToTagColumn},
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
	if nodes := ec.mutation.EnvironmentToUserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToUserTable,
			Columns: environment.EnvironmentToUserPrimaryKey,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.EnvironmentToHostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToHostTable,
			Columns: environment.EnvironmentToHostPrimaryKey,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.EnvironmentToCompetitionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToCompetitionTable,
			Columns: environment.EnvironmentToCompetitionPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: competition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.EnvironmentToIdentityIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToIdentityTable,
			Columns: environment.EnvironmentToIdentityPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: identity.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.EnvironmentToCommandIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToCommandTable,
			Columns: environment.EnvironmentToCommandPrimaryKey,
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
	if nodes := ec.mutation.EnvironmentToScriptIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToScriptTable,
			Columns: environment.EnvironmentToScriptPrimaryKey,
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
	if nodes := ec.mutation.EnvironmentToFileDownloadIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToFileDownloadTable,
			Columns: environment.EnvironmentToFileDownloadPrimaryKey,
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
	if nodes := ec.mutation.EnvironmentToFileDeleteIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToFileDeleteTable,
			Columns: environment.EnvironmentToFileDeletePrimaryKey,
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
	if nodes := ec.mutation.EnvironmentToFileExtractIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToFileExtractTable,
			Columns: environment.EnvironmentToFileExtractPrimaryKey,
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
	if nodes := ec.mutation.EnvironmentToIncludedNetworkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToIncludedNetworkTable,
			Columns: environment.EnvironmentToIncludedNetworkPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: includednetwork.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.EnvironmentToFindingIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToFindingTable,
			Columns: environment.EnvironmentToFindingPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: finding.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.EnvironmentToDNSRecordIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToDNSRecordTable,
			Columns: environment.EnvironmentToDNSRecordPrimaryKey,
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
	if nodes := ec.mutation.EnvironmentToDNSIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToDNSTable,
			Columns: environment.EnvironmentToDNSPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: dns.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.EnvironmentToNetworkIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToNetworkTable,
			Columns: environment.EnvironmentToNetworkPrimaryKey,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.EnvironmentToHostDependencyIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.EnvironmentToHostDependencyTable,
			Columns: environment.EnvironmentToHostDependencyPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: hostdependency.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.EnvironmentToBuildIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   environment.EnvironmentToBuildTable,
			Columns: []string{environment.EnvironmentToBuildColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// EnvironmentCreateBulk is the builder for creating many Environment entities in bulk.
type EnvironmentCreateBulk struct {
	config
	builders []*EnvironmentCreate
}

// Save creates the Environment entities in the database.
func (ecb *EnvironmentCreateBulk) Save(ctx context.Context) ([]*Environment, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ecb.builders))
	nodes := make([]*Environment, len(ecb.builders))
	mutators := make([]Mutator, len(ecb.builders))
	for i := range ecb.builders {
		func(i int, root context.Context) {
			builder := ecb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EnvironmentMutation)
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
					_, err = mutators[i+1].Mutate(root, ecb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ecb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ecb *EnvironmentCreateBulk) SaveX(ctx context.Context) []*Environment {
	v, err := ecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
