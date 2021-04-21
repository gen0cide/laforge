package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Team holds the schema definition for the Team entity.
type Team struct {
	ent.Schema
}

// Fields of the Team.
func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.Int("team_number"),
		field.JSON("config", map[string]string{}),
		field.Int64("revision"),
	}
}

// Edges of the Team.
func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("TeamToUser", User.Type),
		edge.To("TeamToBuild", Build.Type),
		edge.To("TeamToEnvironment", Environment.Type),
		edge.To("TeamToTag", Tag.Type),
		edge.From("TeamToProvisionedNetwork", ProvisionedNetwork.Type).Ref("ProvisionedNetworkToTeam"),
		edge.From("TeamToPlan", Plan.Type).Ref("PlanToTeam"),
	}
}
