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
		field.String("path").
			StructTag(`hcl:"path,attr"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the FileDelete.
func (FileDelete) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("FileDeleteToTag", Tag.Type),
	}
}
