package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Finding holds the schema definition for the Finding entity.
type Finding struct {
	ent.Schema
}

// Fields of the Finding.
func (Finding) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.Enum("severity").Values("ZeroSeverity", "LowSeverity", "MediumSeverity", "HighSeverity", "CriticalSeverity", "NullSeverity"),
		field.Enum("difficulty").Values("ZeroDifficulty", "NoviceDifficulty", "AdvancedDifficulty", "ExpertDifficulty", "NullDifficulty"),
	}
}

// Edges of the Finding.
func (Finding) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type),
		edge.To("tag", Tag.Type),
		edge.To("host", Host.Type),
		edge.To("script", Script.Type),
	}
}
