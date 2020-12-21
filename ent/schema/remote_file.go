package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// RemoteFile holds the schema definition for the RemoteFile entity.
type RemoteFile struct {
	ent.Schema
}

// Fields of the RemoteFile.
func (RemoteFile) Fields() []ent.Field {
	return []ent.Field{
		field.String("source_type"),
		field.String("source"),
		field.String("destination"),
		field.JSON("vars", map[string]string{}),
		field.Bool("template"),
		field.String("perms"),
		field.Bool("disabled"),
		field.String("md5"),
		field.String("abs_path"),
		field.String("ext"),
	}
}

// Edges of the RemoteFile.
func (RemoteFile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type),
	}
}
