package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
)

// DNS holds the schema definition for the DNS entity.
type DNS struct {
	ent.Schema
}

// Fields of the DNS.
func (DNS) Fields() []ent.Field {
	return []ent.Field{
		field.String("type"),
		field.String("root_domain"),
		field.JSON("dns_servers", []string{}),
		field.JSON("ntp_servers", []string{}),
		field.JSON("config", map[string]string{}),
	}
}

// Edges of the DNS.
func (DNS) Edges() []ent.Edge {
	return nil
}
