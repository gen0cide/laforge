package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// IncludedNetwork holds the schema definition for the IncludedNetwork entity.
type IncludedNetwork struct {
	ent.Schema
}

// Fields of the IncludedNetwork.
func (IncludedNetwork) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			StructTag(`hcl:"name,label"`),
		field.JSON("hosts", []string{}).
			StructTag(`hcl:"included_hosts,attr"`),
	}
}

// Edges of the IncludedNetwork.
func (IncludedNetwork) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("IncludedNetworkToTag", Tag.Type),
		edge.To("IncludedNetworkToEnvironment", Environment.Type),
	}
}
