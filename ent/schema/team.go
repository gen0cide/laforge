package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
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
		edge.To("maintainer", User.Type),
		edge.To("build", Build.Type),
		edge.To("TeamToEnvironment", Environment.Type),
		edge.To("tag", Tag.Type),
		edge.From("provisioned_networks", ProvisionedNetwork.Type).Ref("ProvisionedNetworkToTeam"),
	}
}
