package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Command holds the schema definition for the Command entity.
type Command struct {
	ent.Schema
}

// Fields of the Command.
func (Command) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.String("program"),
		field.JSON("args", []string{}),
		field.Bool("ignore_errors"),
		field.Bool("disabled"),
		field.Int("cooldown").Positive(),
		field.Int("timeout").Positive(),
		field.JSON("vars", map[string]string{}),
	}
}

// Edges of the Command.
func (Command) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type),
		edge.To("tag", Tag.Type),
	}
}
