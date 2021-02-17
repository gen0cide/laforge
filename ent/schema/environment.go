package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Environment holds the schema definition for the Environment entity.
type Environment struct {
	ent.Schema
}

// Fields of the Environment.
func (Environment) Fields() []ent.Field {
	return []ent.Field{
		field.String("competition_id"),
		field.String("name"),
		field.String("description"),
		field.String("builder"),
		field.Int("team_count"),
		field.Int("revision"),
		field.JSON("admin_cidrs", []string{}),
		field.JSON("exposed_vdi_ports", []string{}),
		field.JSON("config", map[string]string{}),
	}
}

// Edges of the Environment.
func (Environment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("EnvironmentToTag", Tag.Type),
		edge.To("EnvironmentToUser", User.Type),
		edge.To("EnvironmentToHost", Host.Type),
		edge.To("EnvironmentToCompetition", Competition.Type),
		edge.To("EnvironmentToBuild", Build.Type),
		edge.From("EnvironmentToIncludedNetwork", IncludedNetwork.Type).Ref("IncludedNetworkToEnvironment"),
		edge.From("EnvironmentToNetwork", Network.Type).Ref("NetworkToEnvironment"),
		edge.From("EnvironmentToTeam", Team.Type).Ref("TeamToEnvironment"),
	}
}
