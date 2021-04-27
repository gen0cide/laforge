package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").
			StructTag(`hcl:"name,attr"`),
		field.String("uuid").
			StructTag(`hcl:"uuid,optional"`),
		field.String("email").
			StructTag(`hcl:"email,attr"`),
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("UserToTag", Tag.Type),
		edge.From("UserToEnvironment", Environment.Type).
			Ref("EnvironmentToUser"),
	}
}
