package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Competition holds the schema definition for the Competition entity.
type Competition struct {
	ent.Schema
}

// Fields of the Competition.
func (Competition) Fields() []ent.Field {
	return []ent.Field{
		field.String("root_password").
			StructTag(`hcl:"root_password,attr"`),
		field.JSON("config", map[string]string{}).
			StructTag(`hcl:"config,optional"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,attr"`),
	}
}

// Edges of the Competition.
func (Competition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("CompetitionToTag", Tag.Type),
		edge.To("CompetitionToDNS", DNS.Type).
			StructTag(`hcl:"dns,block"`),
		edge.From("CompetitionToEnvironment", Environment.Type).Ref("EnvironmentToCompetition").
			StructTag(``),
	}
}
