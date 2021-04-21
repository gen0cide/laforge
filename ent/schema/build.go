package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Build holds the schema definition for the Build entity.
type Build struct {
	ent.Schema
}

// Fields of the Build.
func (Build) Fields() []ent.Field {
	return []ent.Field{
		field.Int("revision"),
		field.JSON("config", map[string]string{}),
	}
}

// Edges of the Build.
func (Build) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("BuildToUser", User.Type),
		edge.To("BuildToTag", Tag.Type),
		edge.To("BuildToProvisionedNetwork", ProvisionedNetwork.Type),
		edge.From("BuildToTeam", Team.Type).Ref("TeamToBuild"),
		edge.From("BuildToEnvironment", Environment.Type).Ref("EnvironmentToBuild"),
		edge.From("BuildToPlan", Plan.Type).Ref("PlanToBuild"),
	}
}
