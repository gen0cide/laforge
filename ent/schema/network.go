package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Network holds the schema definition for the Network entity.
type Network struct {
	ent.Schema
}

// Fields of the Network.
func (Network) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
		field.String("name").
			StructTag(`hcl:"name,attr"`),
		field.String("cidr").
			StructTag(`hcl:"cidr,attr"`),
		field.Bool("vdi_visible").
			StructTag(`hcl:"vdi_visible,optional"`),
		field.JSON("vars", map[string]string{}).
			StructTag(`hcl:"vars,optional"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the Network.
func (Network) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("NetworkToEnvironment", Environment.Type).
			Ref("EnvironmentToNetwork").
			Unique(),
		edge.From("NetworkToHostDependency", HostDependency.Type).
			Ref("HostDependencyToNetwork"),
		edge.From("NetworkToIncludedNetwork", IncludedNetwork.Type).
			Ref("IncludedNetworkToNetwork"),
	}
}
