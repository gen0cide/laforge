package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ProvisioningStep holds the schema definition for the ProvisioningStep entity.
type ProvisioningStep struct {
	ent.Schema
}

// Fields of the ProvisioningStep.
func (ProvisioningStep) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Enum("type").
			Values(
				"Script",
				"Command",
				"DNSRecord",
				"FileDelete",
				"FileDownload",
				"FileExtract",
			),
		field.Int("step_number"),
	}
}

// Edges of the ProvisioningStep.
func (ProvisioningStep) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ProvisioningStepToStatus", Status.Type).
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("ProvisioningStepToProvisionedHost", ProvisionedHost.Type).
			Unique(),
		edge.To("ProvisioningStepToScript", Script.Type).
			Unique(),
		edge.To("ProvisioningStepToCommand", Command.Type).
			Unique(),
		edge.To("ProvisioningStepToDNSRecord", DNSRecord.Type).
			Unique(),
		edge.To("ProvisioningStepToFileDelete", FileDelete.Type).
			Unique(),
		edge.To("ProvisioningStepToFileDownload", FileDownload.Type).
			Unique(),
		edge.To("ProvisioningStepToFileExtract", FileExtract.Type).
			Unique(),
		edge.From("ProvisioningStepToPlan", Plan.Type).
			Ref("PlanToProvisioningStep").
			Unique(),
		edge.From("ProvisioningStepToAgentTask", AgentTask.Type).
			Ref("AgentTaskToProvisioningStep"),
		edge.From("ProvisioningStepToGinFileMiddleware", GinFileMiddleware.Type).
			Ref("GinFileMiddlewareToProvisioningStep").
			Unique(),
	}
}
