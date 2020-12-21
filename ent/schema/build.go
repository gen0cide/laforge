package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Build holds the schema definition for the Build entity.
type Build struct {
	ent.Schema
}

// Fields of the Build.
func (Build) Fields() []ent.Field {
	return []ent.Field{
		field.Int("revision").Positive(),
		field.JSON("config", map[string]string{}),
	}
}

// Edges of the Build.
func (Build) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("maintainer", User.Type),
		edge.To("tag", Tag.Type),
		edge.From("team", Team.Type).Ref("build"),
		edge.To("ProvisionedNetworkToBuild", ProvisionedNetwork.Type),
	}
}
