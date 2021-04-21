package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Plan holds the schema definition for the Plan entity.
type Plan struct {
	ent.Schema
}

// Fields of the Plan.
func (Plan) Fields() []ent.Field {
	return []ent.Field{
		field.Int("step_number"),
		field.Enum("type").
			Values(
				"start_build",
				"start_team",
				"provision_network",
				"provision_host",
				"execute_step",
			),
		field.Int("build_id"),
	}
}

// Edges of the Plan.
func (Plan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("NextPlan", Plan.Type).
			From("PrevPlan").
			Unique(),
		edge.To("PlanToBuild", Build.Type),
		edge.To("PlanToTeam", Team.Type),
		edge.To("PlanToProvisionedNetwork", ProvisionedNetwork.Type),
		edge.To("PlanToProvisionedHost", ProvisionedHost.Type),
		edge.To("PlanToProvisioningStep", ProvisioningStep.Type),
	}
}
