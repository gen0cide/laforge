package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// GinFileMiddleware holds the schema definition for the GinFileMiddleware entity.
type GinFileMiddleware struct {
	ent.Schema
}

// Fields of the GinFileMiddleware.
func (GinFileMiddleware) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("url_id"),
		field.String("file_path"),
		field.Bool("accessed").
			Default(false),
	}
}

// Edges of the GinFileMiddleware.
func (GinFileMiddleware) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("GinFileMiddlewareToProvisionedHost", ProvisionedHost.Type).
			Unique(),
		edge.To("GinFileMiddlewareToProvisioningStep", ProvisioningStep.Type).
			Unique(),
	}
}
