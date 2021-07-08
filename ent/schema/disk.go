package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Disk holds the schema definition for the Disk entity.
type Disk struct {
	ent.Schema
}

// Fields of the Disk.
func (Disk) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Int("size").Positive().
			StructTag(`hcl:"size,attr"`),
	}
}

// Edges of the Disk.
func (Disk) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("DiskToHost", Host.Type).
			Ref("HostToDisk").
			Unique(),
	}
}
