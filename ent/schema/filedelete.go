package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// FileDelete holds the schema definition for the FileDelete entity.
type FileDelete struct {
	ent.Schema
}

// Fields of the FileDelete.
func (FileDelete) Fields() []ent.Field {
	return []ent.Field{
		field.String("path"),
	}
}

// Edges of the FileDelete.
func (FileDelete) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type),
	}
}
