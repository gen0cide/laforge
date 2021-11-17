package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// AuthUser holds the schema definition for the AuthUser entity.
type AuthUser struct {
	ent.Schema
}

// Fields of the AuthUser.
func (AuthUser) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("username"),
		field.String("password").Sensitive(),
		field.String("first_name").Default(""),
		field.String("last_name").Default(""),
		field.String("email").Default(""),
		field.String("phone").Default(""),
		field.String("company").Default(""),
		field.String("occupation").Default(""),
		field.String("private_key_path").Default(""),
		field.Enum("role").Values("USER", "ADMIN"),
		field.Enum("provider").Values("LOCAL", "GITHUB", "OPENID"),
	}
}

// Edges of the AuthUser.
func (AuthUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("AuthUserToToken", Token.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("AuthUserToServerTasks", ServerTask.Type).
			Ref("ServerTaskToAuthUser").
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
