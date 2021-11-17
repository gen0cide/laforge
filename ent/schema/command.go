package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Command holds the schema definition for the Command entity.
type Command struct {
	ent.Schema
}

// Fields of the Command.
func (Command) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
		field.String("name").
			StructTag(`hcl:"name,attr"`),
		field.String("description").
			StructTag(`hcl:"description,attr"`),
		field.String("program").
			StructTag(`hcl:"program,attr"`),
		field.JSON("args", []string{}).
			StructTag(`hcl:"args,attr"`),
		field.Bool("ignore_errors").
			StructTag(`hcl:"ignore_errors,attr"`),
		field.Bool("disabled").
			StructTag(`hcl:"disabled,attr"`),
		field.Int("cooldown").Positive().
			StructTag(`hcl:"cooldown,attr"`),
		field.Int("timeout").Positive().
			StructTag(`hcl:"timeout,attr" `),
		field.JSON("vars", map[string]string{}).
			StructTag(`hcl:"vars,attr"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the Command.
func (Command) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("CommandToUser", User.Type).
			StructTag(`hcl:"maintainer,block"`),
		edge.From("CommandToEnvironment", Environment.Type).
			Ref("EnvironmentToCommand").
			Unique(),
	}
}
