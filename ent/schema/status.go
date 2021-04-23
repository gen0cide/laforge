package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Status holds the schema definition for the Status entity.
type Status struct {
	ent.Schema
}

// Fields of the Status.
func (Status) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("state").Values("PLANNING", "AWAITING", "INPROGRESS", "FAILED", "COMPLETE", "TAINTED"),
		field.Enum("status_for").Values("Build", "Team", "ProvisionedNetwork", "ProvisionedHost"),
		field.Time("started_at").Optional(),
		field.Time("ended_at").Optional(),
		field.Bool("failed").Default(false),
		field.Bool("completed").Default(false),
		field.String("error").Optional(),
	}
}

// Edges of the Status.
func (Status) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("StatusToBuild", Build.Type).
			Ref("BuildToStatus").
			Unique(),
		edge.From("StatusToProvisionedNetwork", ProvisionedNetwork.Type).
			Ref("ProvisionedNetworkToStatus").
			Unique(),
		edge.From("StatusToProvisionedHost", ProvisionedHost.Type).
			Ref("ProvisionedHostToStatus").
			Unique(),
		edge.From("StatusToTeam", Team.Type).
			Ref("TeamToStatus").
			Unique(),
	}
}
