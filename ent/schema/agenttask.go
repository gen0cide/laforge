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
			"CHANGEPERMS",
			"APPENDFILE",
		),
		field.String("args"),
		field.Int("number"),
		field.String("output").Default(""),
		field.Enum("state").Values("AWAITING", "INPROGRESS", "FAILED", "COMPLETE"),
		field.String("error_message").Default(""),
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
		edge.From("AgentTaskToAdhocPlan", AdhocPlan.Type).
			Ref("AdhocPlanToAgentTask"),
	}
}
