package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
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
		field.JSON("vars", []string{}),
	}
}

// Edges of the ProvisionedNetwork.
func (ProvisionedNetwork) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tag", Tag.Type),
		edge.To("status", Status.Type),
		edge.To("network", Network.Type),
		edge.To("build", Build.Type),
		edge.To("ProvisionedNetworkToTeam", Team.Type),
		edge.From("provisioned_hosts", ProvisionedHost.Type).Ref("provisioned_network"),
	}
}
