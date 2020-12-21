package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// IncludedNetwork holds the schema definition for the IncludedNetwork entity.
type IncludedNetwork struct {
	ent.Schema
}

// Fields of the IncludedNetwork.
func (IncludedNetwork) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.JSON("hosts", []string{}),
	}
}

// Edges of the IncludedNetwork.
func (IncludedNetwork) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type),
		edge.To("IncludedNetworkToEnvironment", Environment.Type),
	}
}
