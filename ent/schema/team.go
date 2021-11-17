package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Team holds the schema definition for the Team entity.
type Team struct {
	ent.Schema
}

// Fields of the Team.
func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Int("team_number"),
		field.JSON("vars", map[string]string{}),
	}
}

// Edges of the Team.
func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("TeamToBuild", Build.Type).Unique().Required(),
		edge.To("TeamToStatus", Status.Type).Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("TeamToProvisionedNetwork", ProvisionedNetwork.Type).
			Ref("ProvisionedNetworkToTeam"),
		edge.From("TeamToPlan", Plan.Type).
			Ref("PlanToTeam").
			Unique(),
	}
}
