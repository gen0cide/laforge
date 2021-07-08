package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// FileDownload holds the schema definition for the FileDownload entity.
type FileDownload struct {
	ent.Schema
}

// Fields of the FileDownload.
func (FileDownload) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
		field.String("source_type").
			StructTag(`hcl:"source_type,attr"`),
		field.String("source").
			StructTag(`hcl:"source,attr"`),
		field.String("destination").
			StructTag(`hcl:"destination,attr"`),
		field.Bool("template").
			StructTag(`hcl:"template,optional"`),
		field.String("perms").
			StructTag(`hcl:"perms,optional"`),
		field.Bool("disabled").
			StructTag(`hcl:"disabled,optional"`),
		field.String("md5").
			StructTag(`hcl:"md5,optional"`),
		field.String("abs_path").
			StructTag(`hcl:"abs_path,optional"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the FileDownload.
func (FileDownload) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("FileDownloadToEnvironment", Environment.Type).
			Ref("EnvironmentToFileDownload").
			Unique(),
	}
}
