package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
	}
}

// Edges of the Build.
func (Build) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("BuildToStatus", Status.Type).
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("BuildToEnvironment", Environment.Type).
			Unique().
			Required(),
		edge.To("BuildToCompetition", Competition.Type).
			Unique().
			Required(),
		edge.From("BuildToProvisionedNetwork", ProvisionedNetwork.Type).
			Ref("ProvisionedNetworkToBuild").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("BuildToTeam", Team.Type).
			Ref("TeamToBuild").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("BuildToPlan", Plan.Type).
			Ref("PlanToBuild").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
