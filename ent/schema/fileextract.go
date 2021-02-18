package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FileExtract holds the schema definition for the FileExtract entity.
type FileExtract struct {
	ent.Schema
}

// Fields of the FileExtract.
func (FileExtract) Fields() []ent.Field {
	return []ent.Field{
		field.String("source").
			StructTag(`hcl:"source,attr"`),
		field.String("destination").
			StructTag(`hcl:"destination,attr"`),
		field.String("type").
			StructTag(`hcl:"type,attr"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the FileExtract.
func (FileExtract) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("FileExtractToTag", Tag.Type),
	}
}
