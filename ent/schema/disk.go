package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Disk holds the schema definition for the Disk entity.
type Disk struct {
	ent.Schema
}

// Fields of the Disk.
func (Disk) Fields() []ent.Field {
	return []ent.Field{
		field.Int("size").Positive(),
	}
}

// Edges of the Disk.
func (Disk) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type),
	}
}
