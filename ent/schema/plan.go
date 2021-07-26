package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Plan holds the schema definition for the Plan entity.
type Plan struct {
	ent.Schema
}

// Fields of the Plan.
func (Plan) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Int("step_number"),
		field.Enum("type").
			Values(
				"start_build",
				"start_team",
				"provision_network",
				"provision_host",
				"execute_step",
			),
		field.String("build_id"),
	}
}

// Edges of the Plan.
func (Plan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("NextPlan", Plan.Type).
			From("PrevPlan").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("PlanToBuild", Build.Type).Unique(),
		edge.To("PlanToTeam", Team.Type).Unique(),
		edge.To("PlanToProvisionedNetwork", ProvisionedNetwork.Type).Unique(),
		edge.To("PlanToProvisionedHost", ProvisionedHost.Type).Unique(),
		edge.To("PlanToProvisioningStep", ProvisioningStep.Type).Unique(),
		edge.To("PlanToStatus", Status.Type).Unique().Required().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("PlanToPlanDiffs", PlanDiff.Type).
			Ref("PlanDiffToPlan"),
	}
}
