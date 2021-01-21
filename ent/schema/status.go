package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Status holds the schema definition for the Status entity.
type Status struct {
	ent.Schema
}

// Fields of the Status.
func (Status) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("state").Values("AWAITING", "INPROGRESS", "FAILED", "COMPLETE", "TAINTED"),
		field.Time("started_at"),
		field.Time("ended_at"),
		field.Bool("failed"),
		field.Bool("completed"),
		field.String("error"),
	}
}

// Edges of the Status.
func (Status) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type),
	}
}
