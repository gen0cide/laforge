package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// AgentStatus holds the schema definition for the AgentStatus entity.
type AgentStatus struct {
	ent.Schema
}

// Fields of the AgentStatus.
func (AgentStatus) Fields() []ent.Field {
	return []ent.Field{
		field.String("ClientID"),
		field.String("Hostname"),
		field.Int("UpTime"),
		field.Int("BootTime"),
		field.Int("NumProcs"),
		field.String("Os"),
		field.String("HostID"),
		field.Float("Load1"),
		field.Float("Load5"),
		field.Float("Load15"),
		field.Int("TotalMem"),
		field.Int("FreeMem"),
		field.Int("UsedMem"),
	}
}

// Edges of the AgentStatus.
func (AgentStatus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("host", ProvisionedHost.Type),
	}
}
