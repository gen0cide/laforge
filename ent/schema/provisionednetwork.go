package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProvisionedNetwork holds the schema definition for the ProvisionedNetwork entity.
type ProvisionedNetwork struct {
	ent.Schema
}

// Fields of the ProvisionedNetwork.
func (ProvisionedNetwork) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("cidr"),
	}
}

// Edges of the ProvisionedNetwork.
func (ProvisionedNetwork) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ProvisionedNetworkToStatus", Status.Type).Unique(),
		edge.To("ProvisionedNetworkToNetwork", Network.Type).Unique(),
		edge.To("ProvisionedNetworkToBuild", Build.Type).Unique(),
		edge.To("ProvisionedNetworkToTeam", Team.Type).Unique(),
		edge.From("ProvisionedNetworkToProvisionedHost", ProvisionedHost.Type).
			Ref("ProvisionedHostToProvisionedNetwork"),
		edge.From("ProvisionedNetworkToPlan", Plan.Type).
			Ref("PlanToProvisionedNetwork"),
	}
}
