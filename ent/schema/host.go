package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Host holds the schema definition for the Host entity.
type Host struct {
	ent.Schema
}

// Fields of the Host.
func (Host) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
		field.String("hostname").
			StructTag(`hcl:"hostname,attr"`),
		field.String("description").
			StructTag(`hcl:"description,optional" `),
		field.String("OS").
			StructTag(`hcl:"os,attr"`),
		field.Int("last_octet").
			StructTag(`hcl:"last_octet,attr"`),
		field.String("instance_size").
			StructTag(`hcl:"instance_size,attr"`),
		field.Bool("allow_mac_changes").
			StructTag(`hcl:"allow_mac_changes,optional"`),
		field.JSON("exposed_tcp_ports", []string{}).
			StructTag(`hcl:"exposed_tcp_ports,optional"`),
		field.JSON("exposed_udp_ports", []string{}).
			StructTag(`hcl:"exposed_udp_ports,optional"`),
		field.String("override_password").
			StructTag(`hcl:"override_password,optional"`),
		field.JSON("vars", map[string]string{}).
			StructTag(`hcl:"vars,optional"`),
		field.JSON("user_groups", []string{}).
			StructTag(`hcl:"user_groups,optional"`),
		field.JSON("provision_steps", []string{}).Optional().
			StructTag(`hcl:"provision_steps,optional"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the Host.
func (Host) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("HostToDisk", Disk.Type).
			StructTag(`hcl:"disk,block"`).
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("HostToUser", User.Type).
			StructTag(`hcl:"maintainer,block"`),
		edge.From("HostToEnvironment", Environment.Type).
			Ref("EnvironmentToHost").
			Unique(),
		edge.From("HostToIncludedNetwork", IncludedNetwork.Type).
			Ref("IncludedNetworkToHost"),
		edge.From("DependOnHostToHostDependency", HostDependency.Type).
			Ref("HostDependencyToDependOnHost").
			StructTag(`hcl:"depends_on,block"`),
		edge.From("DependByHostToHostDependency", HostDependency.Type).
			Ref("HostDependencyToDependByHost"),
	}
}
