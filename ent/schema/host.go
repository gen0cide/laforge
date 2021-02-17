package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Host holds the schema definition for the Host entity.
type Host struct {
	ent.Schema
}

// Fields of the Host.
func (Host) Fields() []ent.Field {
	return []ent.Field{
		field.String("hostname").
			StructTag(`hcl:"hostname,attr"`),
		field.String("description").
			StructTag(`hcl:"description,optional" `),
		field.String("OS").
			StructTag(`hcl:"os,attr"`),
		field.Int("last_octet").
			StructTag(`hcl:"last_octet,attr"`),
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
		field.JSON("depends_on", map[string]string{}).Optional().
			StructTag(`hcl:"depends_on,optional"`),
		field.JSON("provision_steps", []string{}).Optional().
			StructTag(`hcl:"provision_steps,optional"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,attr"`),
	}
}

// Edges of the Host.
func (Host) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("HostToDisk", Disk.Type).
			StructTag(`hcl:"disk,block"`),
		edge.To("HostToUser", User.Type).
			StructTag(`hcl:"maintainer,block"`),
		edge.To("HostToTag", Tag.Type),
		edge.From("HostToEnvironment", Environment.Type).Ref("EnvironmentToHost"),
	}
}
