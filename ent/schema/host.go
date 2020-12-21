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
		field.String("hostname"),
		field.String("description"),
		field.String("OS"),
		field.Int("last_octet"),
		field.Bool("allow_mac_changes"),
		field.JSON("exposed_tcp_ports", []string{}),
		field.JSON("exposed_udp_ports", []string{}),
		field.String("override_password"),
		field.JSON("vars", map[string]string{}),
		field.JSON("user_groups", []string{}),
		field.JSON("depends_on", []string{}).Optional(),
		field.JSON("scripts", []string{}).Optional(),
		field.JSON("commands", []string{}).Optional(),
		field.JSON("remote_files", []string{}).Optional(),
		field.JSON("dns_records", []string{}).Optional(),
	}
}

// Edges of the Host.
func (Host) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("disk", Disk.Type),
		edge.To("maintainer", User.Type),
		edge.To("tag", Tag.Type),
	}
}
