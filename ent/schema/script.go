package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Script holds the schema definition for the Script entity.
type Script struct {
	ent.Schema
}

// Fields of the Script.
func (Script) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
		field.String("name").
			StructTag(`hcl:"name,attr"`),
		field.String("language").
			StructTag(`hcl:"language,attr"`),
		field.String("description").
			StructTag(`hcl:"description,optional"`),
		field.String("source").
			StructTag(`hcl:"source,attr"`),
		field.String("source_type").
			StructTag(`hcl:"source_type,attr"`),
		field.Int("cooldown").
			StructTag(`hcl:"cooldown,optional"`),
		field.Int("timeout").
			StructTag(`hcl:"timeout,optional"`),
		field.Bool("ignore_errors").
			StructTag(`hcl:"ignore_errors,optional"`),
		field.JSON("args", []string{}).
			StructTag(`hcl:"args,optional"`),
		field.Bool("disabled").
			StructTag(`hcl:"disabled,optional" `),
		field.JSON("vars", map[string]string{}).
			StructTag(`hcl:"vars,optional"`),
		field.String("abs_path").
			StructTag(`hcl:"abs_path,optional"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the Script.
func (Script) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ScriptToUser", User.Type).
			StructTag(`hcl:"maintainer,block"`),
		edge.To("ScriptToFinding", Finding.Type).
			StructTag(`hcl:"finding,block"`).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("ScriptToEnvironment", Environment.Type).
			Ref("EnvironmentToScript").
			Unique(),
	}
}
