package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// FileExtract holds the schema definition for the FileExtract entity.
type FileExtract struct {
	ent.Schema
}

// Fields of the FileExtract.
func (FileExtract) Fields() []ent.Field {
	return []ent.Field{
		field.String("source"),
		field.String("destination"),
		field.String("type"),
	}
}

// Edges of the FileExtract.
func (FileExtract) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type),
	}
}
