package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// PlanDiff holds the schema definition for the PlanDiff entity.
type PlanDiff struct {
	ent.Schema
}

// Fields of the PlanDiff.
func (PlanDiff) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Int("revision"),
		field.Enum("new_state").Values("PLANNING", "AWAITING", "INPROGRESS", "FAILED", "COMPLETE", "TAINTED", "TODELETE", "DELETEINPROGRESS", "DELETED", "TOREBUILD"),
	}
}

// Edges of the PlanDiff.
func (PlanDiff) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("PlanDiffToBuildCommit", BuildCommit.Type).
			Unique().
			Required(),
		edge.To("PlanDiffToPlan", Plan.Type).
			Unique().
			Required(),
	}
}
