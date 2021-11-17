package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// FileExtract holds the schema definition for the FileExtract entity.
type FileExtract struct {
	ent.Schema
}

// Fields of the FileExtract.
func (FileExtract) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
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
		edge.From("FileExtractToEnvironment", Environment.Type).
			Ref("EnvironmentToFileExtract").
			Unique(),
	}
}
