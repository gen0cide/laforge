package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProvisioningStep holds the schema definition for the ProvisioningStep entity.
type ProvisioningStep struct {
	ent.Schema
}

// Fields of the ProvisioningStep.
func (ProvisioningStep) Fields() []ent.Field {
	return []ent.Field{
		field.String("provisioner_type"),
		field.Int("step_number"),
	}
}

// Edges of the ProvisioningStep.
func (ProvisioningStep) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ProvisioningStepToTag", Tag.Type),
		edge.To("ProvisioningStepToStatus", Status.Type),
		edge.To("ProvisioningStepToProvisionedHost", ProvisionedHost.Type),
		edge.To("ProvisioningStepToScript", Script.Type),
		edge.To("ProvisioningStepToCommand", Command.Type),
		edge.To("ProvisioningStepToDNSRecord", DNSRecord.Type),
		edge.To("ProvisioningStepToFileDelete", FileDelete.Type),
		edge.To("ProvisioningStepToFileDownload", FileDownload.Type),
		edge.To("ProvisioningStepToFileExtract", FileExtract.Type),
	}
}
