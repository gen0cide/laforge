package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Commit holds the schema definition for the Commit entity.
type Commit struct {
	ent.Schema
}

// Fields of the Commit.
func (Commit) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Enum("type").Values("ROOT", "REBUILD", "DELETE"),
		field.Int("revision"),
		field.Enum("commit_state").Values("PLANNING", "INPROGRESS", "APPLIED"),
	}
}

// Edges of the Commit.
func (Commit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("CommitToBuild", Build.Type).
			Unique().
			Required(),
		edge.From("CommitToPlanDiffs", PlanDiff.Type).
			Ref("PlanDiffToCommit").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
