package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Build holds the schema definition for the Build entity.
type Build struct {
	ent.Schema
}

// Fields of the Build.
func (Build) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Int("revision"),
		field.Int("environment_revision"),
		field.Bool("completed_plan").
			Default(false),
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
		edge.To("BuildToLatestBuildCommit", BuildCommit.Type).
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("BuildToProvisionedNetwork", ProvisionedNetwork.Type).
			Ref("ProvisionedNetworkToBuild"),
		edge.From("BuildToTeam", Team.Type).
			Ref("TeamToBuild"),
		edge.From("BuildToPlan", Plan.Type).
			Ref("PlanToBuild"),
		edge.From("BuildToBuildCommits", BuildCommit.Type).
			Ref("BuildCommitToBuild"),
		edge.From("BuildToAdhocPlans", AdhocPlan.Type).
			Ref("AdhocPlanToBuild"),
	}
}
