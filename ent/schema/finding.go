package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Finding holds the schema definition for the Finding entity.
type Finding struct {
	ent.Schema
}

// Fields of the Finding.
func (Finding) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").
			StructTag(`hcl:"name,attr"`),
		field.String("description").
			StructTag(`hcl:"description,optional"`),
		field.Enum("severity").Values("ZeroSeverity", "LowSeverity", "MediumSeverity", "HighSeverity", "CriticalSeverity", "NullSeverity").
			StructTag(`hcl:"severity,attr"`),
		field.Enum("difficulty").Values("ZeroDifficulty", "NoviceDifficulty", "AdvancedDifficulty", "ExpertDifficulty", "NullDifficulty").
			StructTag(`hcl:"difficulty,attr"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the Finding.
func (Finding) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("FindingToUser", User.Type).
			StructTag(`hcl:"maintainer,block"`),
		edge.To("FindingToHost", Host.Type).
			Unique(),
		edge.From("FindingToScript", Script.Type).
			Ref("ScriptToFinding").
			Unique(),
		edge.From("FindingToEnvironment", Environment.Type).
			Ref("EnvironmentToFinding").
			Unique(),
	}
}
