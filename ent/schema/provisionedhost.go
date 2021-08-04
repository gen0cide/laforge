package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ProvisionedHost holds the schema definition for the ProvisionedHost entity.
type ProvisionedHost struct {
	ent.Schema
}

// Fields of the ProvisionedHost.
func (ProvisionedHost) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("subnet_ip"),
		field.Enum("addon_type").Values("DNS").Nillable().Optional(),
	}
}

// Edges of the ProvisionedHost.
func (ProvisionedHost) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ProvisionedHostToStatus", Status.Type).
			Required().
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("ProvisionedHostToProvisionedNetwork", ProvisionedNetwork.Type).
			Required().
			Unique(),
		edge.To("ProvisionedHostToHost", Host.Type).
			Required().
			Unique(),
		edge.To("ProvisionedHostToEndStepPlan", Plan.Type).
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("ProvisionedHostToBuild", Build.Type).
			Unique().
			Required(),
		edge.From("ProvisionedHostToProvisioningStep", ProvisioningStep.Type).
			Ref("ProvisioningStepToProvisionedHost"),
		edge.From("ProvisionedHostToAgentStatus", AgentStatus.Type).
			Ref("AgentStatusToProvisionedHost"),
		edge.From("ProvisionedHostToAgentTask", AgentTask.Type).
			Ref("AgentTaskToProvisionedHost"),
		edge.From("ProvisionedHostToPlan", Plan.Type).
			Ref("PlanToProvisionedHost").
			Unique(),
		edge.From("ProvisionedHostToGinFileMiddleware", GinFileMiddleware.Type).
			Ref("GinFileMiddlewareToProvisionedHost").
			Unique(),
	}
}
