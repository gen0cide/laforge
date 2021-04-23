package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProvisionedHost holds the schema definition for the ProvisionedHost entity.
type ProvisionedHost struct {
	ent.Schema
}

// Fields of the ProvisionedHost.
func (ProvisionedHost) Fields() []ent.Field {
	return []ent.Field{
		field.String("subnet_ip"),
	}
}

// Edges of the ProvisionedHost.
func (ProvisionedHost) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ProvisionedHostToStatus", Status.Type).
			Required().
			Unique(),
		edge.To("ProvisionedHostToProvisionedNetwork", ProvisionedNetwork.Type).
			Required().
			Unique(),
		edge.To("ProvisionedHostToHost", Host.Type).
			Required().
			Unique(),
		edge.From("ProvisionedHostToProvisioningStep", ProvisioningStep.Type).
			Ref("ProvisioningStepToProvisionedHost"),
		edge.From("ProvisionedHostToAgentStatus", AgentStatus.Type).
			Ref("AgentStatusToProvisionedHost"),
		edge.From("ProvisionedHostToPlan", Plan.Type).
			Ref("PlanToProvisionedHost"),
		edge.From("ProvisionedHostToGinFileMiddleware", GinFileMiddleware.Type).
			Ref("GinFileMiddlewareToProvisionedHost").
			Unique(),
	}
}
