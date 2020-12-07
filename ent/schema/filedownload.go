package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// FileDownload holds the schema definition for the FileDownload entity.
type FileDownload struct {
	ent.Schema
}

// Fields of the FileDownload.
func (FileDownload) Fields() []ent.Field {
	return []ent.Field{
		field.String("source_type"),
		field.String("source"),
		field.String("destination"),
		field.Bool("template"),
		field.String("mode"),
		field.Bool("disabled"),
		field.String("md5"),
		field.String("abs_path"),
	}
}

// Edges of the FileDownload.
func (FileDownload) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type),
	}
}
