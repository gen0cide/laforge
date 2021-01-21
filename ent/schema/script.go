package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Script holds the schema definition for the Script entity.
type Script struct {
	ent.Schema
}

// Fields of the Script.
func (Script) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("language"),
		field.String("description"),
		field.String("source"),
		field.String("source_type"),
		field.Int("cooldown"),
		field.Int("timeout"),
		field.Bool("ignore_errors"),
		field.JSON("args", []string{}),
		field.Bool("disabled"),
		field.JSON("vars", map[string]string{}),
		field.String("abs_path"),
	}
}

// Edges of the Script.
func (Script) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type),
		edge.To("maintainer", User.Type),
		edge.From("finding", Finding.Type).Ref("script"),
	}
}
