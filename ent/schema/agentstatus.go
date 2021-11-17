package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// AgentStatus holds the schema definition for the AgentStatus entity.
type AgentStatus struct {
	ent.Schema
}

// Fields of the AgentStatus.
func (AgentStatus) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("ClientID"),
		field.String("Hostname"),
		field.Int64("UpTime"),
		field.Int64("BootTime"),
		field.Int64("NumProcs"),
		field.String("Os"),
		field.String("HostID"),
		field.Float("Load1"),
		field.Float("Load5"),
		field.Float("Load15"),
		field.Int64("TotalMem"),
		field.Int64("FreeMem"),
		field.Int64("UsedMem"),
		field.Int64("Timestamp"),
	}
}

// Edges of the AgentStatus.
func (AgentStatus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("AgentStatusToProvisionedHost", ProvisionedHost.Type).
			Unique(),
		edge.To("AgentStatusToProvisionedNetwork", ProvisionedNetwork.Type).
			Unique(),
		edge.To("AgentStatusToBuild", Build.Type).
			Unique(),
	}
}
