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
		edge.To("tag", Tag.Type),
		edge.To("user", User.Type),
		edge.To("host", Host.Type),
		edge.To("competition", Competition.Type),
		edge.To("build", Build.Type),
		edge.From("included_network", IncludedNetwork.Type).Ref("IncludedNetworkToEnvironment"),
		edge.From("network", Network.Type).Ref("NetworkToEnvironment"),
		edge.From("team", Team.Type).Ref("TeamToEnvironment"),
	}
}
