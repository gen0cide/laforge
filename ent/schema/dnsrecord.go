package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// DNSRecord holds the schema definition for the DNSRecord entity.
type DNSRecord struct {
	ent.Schema
}

// Fields of the DNSRecord.
func (DNSRecord) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.JSON("values", []string{}),
		field.String("type"),
		field.String("zone"),
		field.JSON("vars", map[string]string{}),
		field.Bool("disabled"),
	}
}

// Edges of the DNSRecord.
func (DNSRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type),
	}
}
