package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Identity holds the schema definition for the Identity entity.
type Identity struct {
	ent.Schema
}

// Fields of the Identity.
func (Identity) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
		field.String("first_name").
			StructTag(`hcl:"firstname,attr"`),
		field.String("last_name").
			StructTag(`hcl:"lastname,attr" `),
		field.String("email").
			StructTag(`hcl:"email,attr" `),
		field.String("password").
			StructTag(`hcl:"password,attr" `),
		field.String("description").
			StructTag(`hcl:"description,optional" `),
		field.String("avatar_file").
			StructTag(`hcl:"avatar_file,optional" `),
		field.JSON("vars", map[string]string{}).
			StructTag(`hcl:"vars,optional"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the Identity.
func (Identity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("IdentityToEnvironment", Environment.Type).
			Ref("EnvironmentToIdentity").
			Unique(),
	}
}
