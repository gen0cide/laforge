package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// BuildCommit holds the schema definition for the BuildCommit entity.
type BuildCommit struct {
	ent.Schema
}

// Fields of the BuildCommit.
func (BuildCommit) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Enum("type").Values("ROOT", "REBUILD", "DELETE"),
		field.Int("revision"),
		field.Enum("state").Values("PLANNING", "INPROGRESS", "APPLIED", "CANCELLED", "APPROVED"),
	}
}

// Edges of the BuildCommit.
func (BuildCommit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("BuildCommitToBuild", Build.Type).
			Unique().
			Required(),
		edge.From("BuildCommitToPlanDiffs", PlanDiff.Type).
			Ref("PlanDiffToBuildCommit").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
