package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// FileDelete holds the schema definition for the FileDelete entity.
type FileDelete struct {
	ent.Schema
}

// Fields of the FileDelete.
func (FileDelete) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
		field.String("path").
			StructTag(`hcl:"path,attr"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the FileDelete.
func (FileDelete) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("FileDeleteToEnvironment", Environment.Type).
			Ref("EnvironmentToFileDelete").
			Unique(),
	}
}
