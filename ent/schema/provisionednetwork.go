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
		edge.To("ProvisionedNetworkToTag", Tag.Type),
		edge.To("ProvisionedNetworkToStatus", Status.Type),
		edge.To("ProvisionedNetworkToNetwork", Network.Type),
		edge.From("ProvisionedNetworkToBuild", Build.Type).
			Ref("BuildToProvisionedNetwork"),
		edge.To("ProvisionedNetworkToTeam", Team.Type),
		edge.From("ProvisionedNetworkToProvisionedHost", ProvisionedHost.Type).
			Ref("ProvisionedHostToProvisionedNetwork"),
		edge.From("ProvisionedNetworkToPlan", Plan.Type).
			Ref("PlanToProvisionedNetwork"),
	}
}
