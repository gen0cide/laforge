package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// AgentTask holds the schema definition for the AgentTask entity.
type AgentTask struct {
	ent.Schema
}

// Fields of the AgentTask.
func (AgentTask) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Enum("command").Values(
			"DEFAULT",
			"DELETE",
			"REBOOT",
			"EXTRACT",
			"DOWNLOAD",
			"CREATEUSER",
			"CREATEUSERPASS",
			"ADDTOGROUP",
			"EXECUTE",
			"VALIDATE",
		),
		field.String("args"),
		field.Int("number"),
		field.String("output").Nillable().Optional(),
		field.Enum("state").Values("AWAITING", "INPROGRESS", "FAILED", "COMPLETE"),
	}
}

// Edges of the AgentTask.
func (AgentTask) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("AgentTaskToProvisioningStep", ProvisioningStep.Type).
			Unique(),
		edge.To("AgentTaskToProvisionedHost", ProvisionedHost.Type).
			Required().
			Unique(),
	}
}
