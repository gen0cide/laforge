package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			StructTag(`hcl:"name,attr"`),
		field.String("uuid").
			StructTag(`hcl:"uuid,optional"`),
		field.String("email").
			StructTag(`hcl:"email,attr"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("UserToTag", Tag.Type),
		edge.From("UserToEnvironment", Environment.Type).Ref("EnvironmentToUser"),
	}
}
